# Modificações Necessárias para Docker MCP Gateway com Docker Engine Nativo

## Análise Geral

O Docker MCP Gateway atualmente depende fortemente do Docker Desktop para várias funcionalidades, incluindo:
- Verificação de recursos do Docker Desktop
- Gerenciamento de segredos via sockets específicos
- Monitoramento OAuth via backend socket
- Autenticação via credential helper do Docker Desktop

## Arquivos e Linhas que Precisam ser Modificados

### 1. pkg/desktop/sockets_linux.go

**Problema:** Verifica a existência do socket do Docker Desktop em `/run/host-services/backend.sock`

**Linhas 13-33:** Modificar para ignorar verificação do Docker Desktop quando `DOCKER_MCP_IN_CONTAINER=1`

```go
// Linha 13-33: Modificar getDockerDesktopPaths()
func getDockerDesktopPaths() (DockerDesktopPaths, error) {
    // Se estiver em modo contêiner, retornar paths vazios para ignorar Docker Desktop
    if os.Getenv("DOCKER_MCP_IN_CONTAINER") == "1" {
        return DockerDesktopPaths{
            AdminSettingPath:     "",
            BackendSocket:        "",
            RawDockerSocket:      "",
            JFSSocket:            "",
            ToolsSocket:          "",
            CredentialHelperPath: func() string { return "" },
        }, nil
    }
    
    _, err := os.Stat("/run/host-services/backend.sock")
    // ... resto do código existente
}
```

### 2. pkg/docker/secrets.go

**Problema:** Tenta usar o socket do Docker Desktop para ler segredos quando não está em contêiner

**Linhas 59-73:** Modificar para sempre usar método alternativo quando `DOCKER_MCP_IN_CONTAINER=1`

```go
// Linha 59: Modificar condição
if os.Getenv("DOCKER_MCP_IN_CONTAINER") != "1" {
    // Modificar para sempre usar método alternativo
    return c.readSecretsAlternative(ctx, names)
}

// Adicionar novo método:
func (c *dockerClient) readSecretsAlternative(ctx context.Context, names []string) (map[string]string, error) {
    // Implementar leitura de segredos via arquivo .env ou variáveis de ambiente
    // Este método deve ler segredos de arquivos locais ou variáveis de ambiente
    // em vez de depender do socket do Docker Desktop
    values := map[string]string{}
    for _, name := range names {
        if envValue := os.Getenv(name); envValue != "" {
            values[name] = envValue
        }
    }
    return values, nil
}
```

### 3. pkg/gateway/run.go

**Problema:** Desabilita OAuth apenas quando está em contêiner, mas deveria ser o contrário

**Linhas 272-296:** Modificar para desabilitar OAuth quando não está em contêiner Docker Desktop

```go
// Linha 272-273: Modificar lógica
inContainer := os.Getenv("DOCKER_MCP_IN_CONTAINER") == "1"
isDockerDesktop := os.Getenv("DOCKER_MCP_IN_CONTAINER") != "1" && !isNativeDockerEngine()

// Linha 275: Modificar condição
if g.McpOAuthDcrEnabled && isDockerDesktop {
    // ... código OAuth existente
}

// Adicionar função auxiliar:
func isNativeDockerEngine() bool {
    // Verificar se está rodando em Docker Engine nativo (não Docker Desktop)
    info, err := dockerCli.Client().Info(context.Background())
    if err != nil {
        return false
    }
    return info.OperatingSystem != "Docker Desktop"
}
```

**Linhas 334-341:** Modificar para sempre desabilitar autenticação quando `DOCKER_MCP_IN_CONTAINER=1`

```go
// Linha 334: Modificar condição
if (transport == "sse" || transport == "http" || transport == "streamable" || transport == "streaming" || transport == "streamable-http") && !inContainer {
    // Modificar para sempre desabilitar autenticação em modo nativo
    if os.Getenv("DOCKER_MCP_IN_CONTAINER") == "1" {
        // Pular autenticação
        g.authToken = ""
        g.authTokenWasGenerated = false
    } else {
        // Código existente de autenticação
        token, wasGenerated, err := getOrGenerateAuthToken()
        // ...
    }
}
```

### 4. cmd/docker-mcp/commands/root.go

**Problema:** Verifica se está rodando em Docker CE e exige Docker Desktop

**Linhas 48-57:** Modificar para pular verificação do Docker Desktop quando `DOCKER_MCP_IN_CONTAINER=1`

```go
// Linha 48: Modificar condição
if os.Getenv("DOCKER_MCP_IN_CONTAINER") != "1" {
    // Adicionar verificação para modo nativo
    if os.Getenv("DOCKER_MCP_NATIVE_MODE") == "1" {
        // Pular verificação do Docker Desktop
        return nil
    }
    
    runningInDockerCE, err := docker.RunningInDockerCE(ctx, dockerCli)
    if err != nil {
        return err
    }

    if !runningInDockerCE {
        return desktop.CheckFeatureIsEnabled(ctx, "enableDockerMCPToolkit", "Docker MCP Toolkit")
    }
}
```

### 5. cmd/docker-mcp/commands/gateway.go

**Problema:** Configurações padrão dependem do Docker Desktop

**Linhas 32-46:** Modificar configurações padrão para modo nativo

```go
// Linha 32: Modificar condição
if os.Getenv("DOCKER_MCP_IN_CONTAINER") == "1" || os.Getenv("DOCKER_MCP_NATIVE_MODE") == "1" {
    // Modificar configurações para modo nativo
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
}
```

**Linhas 184:** Modificar descrição do parâmetro secrets

```go
// Linha 184: Modificar texto de ajuda
runCmd.Flags().StringVar(&options.SecretsPath, "secrets", options.SecretsPath, "Colon separated paths to search for secrets. Can be `docker-desktop`, `file:./.env` for local files, or environment variables")
```

## Novas Funcionalidades Necessárias

### 1. Gerenciamento Alternativo de Segredos

Criar novo arquivo `pkg/docker/secrets_native.go`:

```go
package docker

import (
    "bufio"
    "context"
    "os"
    "strings"
)

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
```

### 2. Modificação do Cliente Docker

Modificar `pkg/docker/client.go` para adicionar suporte a modo nativo:

```go
// Adicionar após linha 70
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

## Plano de Implementação

### Fase 1: Modificações Básicas
1. Modificar `pkg/desktop/sockets_linux.go` para ignorar verificação do Docker Desktop
2. Modificar `cmd/docker-mcp/commands/root.go` para pular verificação em modo nativo
3. Adicionar variável de ambiente `DOCKER_MCP_NATIVE_MODE=1`

### Fase 2: Gerenciamento de Segredos
1. Modificar `pkg/docker/secrets.go` para usar método alternativo
2. Implementar `pkg/docker/secrets_native.go` com leitura de arquivos .env
3. Modificar `cmd/docker-mcp/commands/gateway.go` para usar segredos locais

### Fase 3: OAuth e Autenticação
1. Modificar `pkg/gateway/run.go` para desabilitar OAuth em modo nativo
2. Implementar autenticação alternativa se necessária
3. Modificar fluxos de autenticação para modo nativo

### Fase 4: Configurações Padrão
1. Modificar configurações padrão em `cmd/docker-mcp/commands/gateway.go`
2. Atualizar documentação e textos de ajuda
3. Testar integração completa

## Variáveis de Ambiente

- `DOCKER_MCP_IN_CONTAINER=1`: Mantém comportamento atual para contêineres
- `DOCKER_MCP_NATIVE_MODE=1`: Novo modo para Docker Engine nativo
- `DOCKER_MCP_SECRETS_FILE`: Caminho para arquivo de segredos alternativo

## Compatibilidade

As modificações mantêm compatibilidade com:
- Instalações existentes do Docker Desktop
- Modo contêiner atual
- Configurações personalizadas existentes

As novas funcionalidades adicionam suporte para:
- Docker Engine nativo no Linux
- Gerenciamento de segredos via arquivos locais
- Operação sem dependências do Docker Desktop