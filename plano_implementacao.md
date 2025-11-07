# Plano Detalhado de Implementação

## Resumo Executivo

Este documento descreve o plano detalhado para modificar o Docker MCP Gateway a fim de permitir sua operação com Docker Engine nativo do Linux, sem dependências do Docker Desktop.

## Objetivos

1. Permitir execução com Docker Engine nativo do Linux
2. Implementar gerenciamento de segredos alternativo
3. Desabilitar monitor OAuth quando não disponível
4. Manter compatibilidade com instalações existentes

## Fases de Implementação

### Fase 1: Modificações Básicas de Detecção

#### 1.1 Modificar pkg/desktop/sockets_linux.go

**Arquivo:** `pkg/desktop/sockets_linux.go`
**Função:** `getDockerDesktopPaths()`
**Linhas:** 13-33

**Implementação:**
```go
func getDockerDesktopPaths() (DockerDesktopPaths, error) {
    // Se estiver em modo contêiner ou modo nativo, retornar paths vazios
    if os.Getenv("DOCKER_MCP_IN_CONTAINER") == "1" || os.Getenv("DOCKER_MCP_NATIVE_MODE") == "1" {
        return DockerDesktopPaths{
            AdminSettingPath:     "",
            BackendSocket:        "",
            RawDockerSocket:      "",
            JFSSocket:            "",
            ToolsSocket:          "",
            CredentialHelperPath: func() string { return "" },
        }, nil
    }
    
    // Código existente permanece para compatibilidade com Docker Desktop
    _, err := os.Stat("/run/host-services/backend.sock")
    if err != nil {
        if !errors.Is(err, os.ErrNotExist) {
            return DockerDesktopPaths{}, err
        }
        // ... resto do código existente
    }
    // ... resto do código existente
}
```

#### 1.2 Modificar cmd/docker-mcp/commands/root.go

**Arquivo:** `cmd/docker-mcp/commands/root.go`
**Função:** `PersistentPreRunE`
**Linhas:** 48-57

**Implementação:**
```go
PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
    cmd.SetContext(ctx)
    if err := plugin.PersistentPreRunE(cmd, args); err != nil {
        return err
    }

    // Pular verificação do Docker Desktop em modo contêiner ou modo nativo
    if os.Getenv("DOCKER_MCP_IN_CONTAINER") != "1" && os.Getenv("DOCKER_MCP_NATIVE_MODE") != "1" {
        runningInDockerCE, err := docker.RunningInDockerCE(ctx, dockerCli)
        if err != nil {
            return err
        }

        if !runningInDockerCE {
            return desktop.CheckFeatureIsEnabled(ctx, "enableDockerMCPToolkit", "Docker MCP Toolkit")
        }
    }

    return nil
},
```

### Fase 2: Gerenciamento Alternativo de Segredos

#### 2.1 Criar pkg/docker/secrets_native.go

**Arquivo novo:** `pkg/docker/secrets_native.go`

**Implementação:**
```go
package docker

import (
    "bufio"
    "context"
    "os"
    "strings"
)

// readSecretsFromFile lê segredos de um arquivo .env
func (c *dockerClient) readSecretsFromFile(ctx context.Context, names []string, filePath string) (map[string]string, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    values := map[string]string{}
    scanner := bufio.NewScanner(file)
    
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line == "" || strings.HasPrefix(line, "#") {
            continue
        }
        
        parts := strings.SplitN(line, "=", 2)
        if len(parts) == 2 {
            key := strings.TrimSpace(parts[0])
            value := strings.TrimSpace(parts[1])
            
            for _, name := range names {
                if key == name {
                    values[name] = value
                    break
                }
            }
        }
    }
    
    return values, scanner.Err()
}

// readSecretsFromEnv lê segredos de variáveis de ambiente
func (c *dockerClient) readSecretsFromEnv(ctx context.Context, names []string) (map[string]string, error) {
    values := map[string]string{}
    for _, name := range names {
        if envValue := os.Getenv(name); envValue != "" {
            values[name] = envValue
        }
    }
    return values, nil
}
```

#### 2.2 Modificar pkg/docker/secrets.go

**Arquivo:** `pkg/docker/secrets.go`
**Função:** `readSecrets()`
**Linhas:** 45-94

**Implementação:**
```go
func (c *dockerClient) readSecrets(ctx context.Context, names []string) (map[string]string, error) {
    // Se estiver em modo contêiner ou modo nativo, usar método alternativo
    if os.Getenv("DOCKER_MCP_IN_CONTAINER") == "1" || os.Getenv("DOCKER_MCP_NATIVE_MODE") == "1" {
        return c.readSecretsAlternative(ctx, names)
    }
    
    // Código existente para Docker Desktop
    flags := []string{"--network=none", "--pull=never"}
    // ... resto do código existente
}

// Adicionar novo método
func (c *dockerClient) readSecretsAlternative(ctx context.Context, names []string) (map[string]string, error) {
    // Tentar ler de arquivo .env primeiro
    if secretsFile := os.Getenv("DOCKER_MCP_SECRETS_FILE"); secretsFile != "" {
        return c.readSecretsFromFile(ctx, names, secretsFile)
    }
    
    // Tentar ler de variáveis de ambiente
    return c.readSecretsFromEnv(ctx, names)
}
```

### Fase 3: Modificações no Gateway

#### 3.1 Modificar pkg/gateway/run.go

**Arquivo:** `pkg/gateway/run.go`
**Função:** `Run()`
**Linhas:** 272-296

**Implementação:**
```go
// Linha 272-273: Modificar lógica de detecção
inContainer := os.Getenv("DOCKER_MCP_IN_CONTAINER") == "1"
nativeMode := os.Getenv("DOCKER_MCP_NATIVE_MODE") == "1"
isDockerDesktop := !inContainer && !nativeMode

// Linha 275: Modificar condição OAuth
if g.McpOAuthDcrEnabled && isDockerDesktop {
    // Start OAuth notification monitor to receive OAuth related events from Docker Desktop
    log.Log("- Starting OAuth notification monitor")
    monitor := oauth.NewNotificationMonitor()
    monitor.OnOAuthEvent = func(event oauth.Event) {
        // Route event to specific provider
        g.routeEventToProvider(event)
    }
    monitor.Start(ctx)
    
    // ... resto do código OAuth existente
} else if g.McpOAuthDcrEnabled && (inContainer || nativeMode) {
    log.Log("- OAuth disabled in container/native mode")
}
```

**Linhas:** 334-341

**Implementação:**
```go
// Linha 334: Modificar condição de autenticação
if (transport == "sse" || transport == "http" || transport == "streamable" || transport == "streaming" || transport == "streamable-http") && !inContainer && !nativeMode {
    token, wasGenerated, err := getOrGenerateAuthToken()
    if err != nil {
        return fmt.Errorf("failed to initialize auth token: %w", err)
    }
    g.authToken = token
    g.authTokenWasGenerated = wasGenerated
} else if (transport == "sse" || transport == "http" || transport == "streamable" || transport == "streaming" || transport == "streamable-http") && (inContainer || nativeMode) {
    log.Log("- Authentication disabled in container/native mode")
    g.authToken = ""
    g.authTokenWasGenerated = false
}
```

### Fase 4: Configurações Padrão

#### 4.1 Modificar cmd/docker-mcp/commands/gateway.go

**Arquivo:** `cmd/docker-mcp/commands/gateway.go`
**Função:** `gatewayCommand()`
**Linhas:** 32-46

**Implementação:**
```go
// Linha 32: Modificar condição
if os.Getenv("DOCKER_MCP_IN_CONTAINER") == "1" || os.Getenv("DOCKER_MCP_NATIVE_MODE") == "1" {
    // In-container or native mode.
    options = gateway.Config{
        CatalogPath: []string{catalog.DockerCatalogURLV2},
        SecretsPath: "file:./.env", // Usar arquivo local em vez de docker-desktop
        Options: gateway.Options{
            Cpus:         1,
            Memory:       "2Gb",
            Transport:    "stdio",
            LogCalls:     true,
            BlockSecrets: true,
            Verbose:      true,
        },
    }
} else {
    // On-host (Docker Desktop)
    // ... código existente permanece
}
```

**Linha:** 184

**Implementação:**
```go
runCmd.Flags().StringVar(&options.SecretsPath, "secrets", options.SecretsPath, "Colon separated paths to search for secrets. Can be `docker-desktop`, `file:./.env` for local files, or environment variables")
```

### Fase 5: Suporte a Modo Nativo

#### 5.1 Modificar pkg/docker/client.go

**Arquivo:** `pkg/docker/client.go`
**Após linha 70**

**Implementação:**
```go
// IsNativeDockerEngine verifica se está rodando em Docker Engine nativo
func IsNativeDockerEngine(ctx context.Context, dockerCli command.Cli) (bool, error) {
    if runtime.GOOS != "linux" {
        return false, nil
    }
    
    info, err := dockerCli.Client().Info(ctx)
    if err != nil {
        return false, fmt.Errorf("failed to get Docker info: %w", err)
    }
    
    return info.OperatingSystem != "Docker Desktop", nil
}
```

## Variáveis de Ambiente

### Novas Variáveis
- `DOCKER_MCP_NATIVE_MODE=1`: Ativa modo nativo para Docker Engine
- `DOCKER_MCP_SECRETS_FILE=./.env`: Caminho para arquivo de segredos

### Variáveis Existentes (Mantidas)
- `DOCKER_MCP_IN_CONTAINER=1`: Mantém comportamento atual para contêineres

## Testes

### Teste 1: Modo Nativo Básico
```bash
export DOCKER_MCP_NATIVE_MODE=1
docker mcp gateway run --dry-run
```

### Teste 2: Segredos via Arquivo
```bash
export DOCKER_MCP_NATIVE_MODE=1
export DOCKER_MCP_SECRETS_FILE=./test-secrets.env
echo "TEST_SECRET=test_value" > test-secrets.env
docker mcp gateway run --dry-run
```

### Teste 3: Segredos via Ambiente
```bash
export DOCKER_MCP_NATIVE_MODE=1
export TEST_SECRET=test_value
docker mcp gateway run --dry-run
```

## Compatibilidade

### Mantida
- Compatibilidade total com Docker Desktop
- Funcionalidade existente em contêineres
- Configurações personalizadas existentes

### Adicionada
- Suporte a Docker Engine nativo no Linux
- Gerenciamento de segredos via arquivos locais
- Operação sem dependências do Docker Desktop

## Rollback

Caso necessário, as modificações podem ser revertidas:
1. Remover verificações de `DOCKER_MCP_NATIVE_MODE`
2. Restaurar código original de gerenciamento de segredos
3. Reverter modificações no gateway

## Documentação

### Atualizações Necessárias
1. README.md: Adicionar instruções para modo nativo
2. docs/troubleshooting.md: Adicionar problemas comuns do modo nativo
3. docs/self-configured.md: Documentar configuração de segredos locais

## Implantação

### Passos
1. Aplicar modificações fase por fase
2. Testar cada fase individualmente
3. Executar testes de integração completos
4. Atualizar documentação
5. Release com notas de versão

### Migração
Para usuários existentes:
1. Nenhuma ação necessária (compatibilidade mantida)
2. Opcional: migrar para modo nativo se desejado
3. Configurar variáveis de ambiente conforme necessário