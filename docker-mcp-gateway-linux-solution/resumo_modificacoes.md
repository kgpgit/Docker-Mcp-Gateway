# Resumo das Modificações Críticas

## Modificações Imediatas Necessárias

### 1. Ignorar Verificação do Docker Desktop

**Arquivo:** `pkg/desktop/sockets_linux.go`
**Função:** `getDockerDesktopPaths()`
**Modificação:** Adicionar verificação para `DOCKER_MCP_NATIVE_MODE` no início da função

```go
// Adicionar após linha 12
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
```

### 2. Pular Verificação de Recursos do Docker Desktop

**Arquivo:** `cmd/docker-mcp/commands/root.go`
**Função:** `PersistentPreRunE`
**Modificação:** Modificar condição na linha 48

```go
// Modificar linha 48
if os.Getenv("DOCKER_MCP_IN_CONTAINER") != "1" && os.Getenv("DOCKER_MCP_NATIVE_MODE") != "1" {
```

### 3. Implementar Gerenciamento de Segredos Alternativo

**Arquivo:** `pkg/docker/secrets.go`
**Função:** `readSecrets()`
**Modificação:** Adicionar verificação no início da função (linha 59)

```go
// Modificar linha 59
if os.Getenv("DOCKER_MCP_IN_CONTAINER") == "1" || os.Getenv("DOCKER_MCP_NATIVE_MODE") == "1" {
    return c.readSecretsAlternative(ctx, names)
}
```

**Adicionar novo método no final do arquivo:**
```go
func (c *dockerClient) readSecretsAlternative(ctx context.Context, names []string) (map[string]string, error) {
    values := map[string]string{}
    
    // Tentar ler de arquivo se especificado
    if secretsFile := os.Getenv("DOCKER_MCP_SECRETS_FILE"); secretsFile != "" {
        file, err := os.Open(secretsFile)
        if err == nil {
            defer file.Close()
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
        }
    }
    
    // Complementar com variáveis de ambiente
    for _, name := range names {
        if _, exists := values[name]; !exists {
            if envValue := os.Getenv(name); envValue != "" {
                values[name] = envValue
            }
        }
    }
    
    return values, nil
}
```

### 4. Desabilitar OAuth em Modo Nativo

**Arquivo:** `pkg/gateway/run.go`
**Função:** `Run()`
**Modificação:** Modificar linhas 272-275

```go
// Modificar linha 272-275
inContainer := os.Getenv("DOCKER_MCP_IN_CONTAINER") == "1"
nativeMode := os.Getenv("DOCKER_MCP_NATIVE_MODE") == "1"
isDockerDesktop := !inContainer && !nativeMode

// Modificar linha 275
if g.McpOAuthDcrEnabled && isDockerDesktop {
```

### 5. Configurar Segredos Padrão para Modo Nativo

**Arquivo:** `cmd/docker-mcp/commands/gateway.go`
**Função:** `gatewayCommand()`
**Modificação:** Modificar linha 32 e SecretsPath na linha 37

```go
// Modificar linha 32
if os.Getenv("DOCKER_MCP_IN_CONTAINER") == "1" || os.Getenv("DOCKER_MCP_NATIVE_MODE") == "1" {

// Modificar linha 37
SecretsPath: "file:./.env",
```

## Variáveis de Ambiente Chave

### Para Ativar Modo Nativo
```bash
export DOCKER_MCP_NATIVE_MODE=1
```

### Para Configurar Arquivo de Segredos
```bash
export DOCKER_MCP_SECRETS_FILE=./.env
```

## Teste Rápido

Aplicar as modificações acima e testar com:

```bash
# Criar arquivo de segredos de teste
echo "TEST_SECRET=test_value" > .env

# Ativar modo nativo
export DOCKER_MCP_NATIVE_MODE=1
export DOCKER_MCP_SECRETS_FILE=./.env

# Testar gateway
docker mcp gateway run --dry-run
```

## Importações Necessárias

Adicionar em `pkg/docker/secrets.go`:
```go
import (
    "bufio"
    "os"
    "strings"
)
```

## Resumo da Execução

1. **Modificação 1:** `pkg/desktop/sockets_linux.go` - Ignorar Docker Desktop
2. **Modificação 2:** `cmd/docker-mcp/commands/root.go` - Pular verificação
3. **Modificação 3:** `pkg/docker/secrets.go` - Segredos alternativos
4. **Modificação 4:** `pkg/gateway/run.go` - Desabilitar OAuth
5. **Modificação 5:** `cmd/docker-mcp/commands/gateway.go` - Configuração padrão

Estas 5 modificações são suficientes para permitir operação básica do Docker MCP Gateway com Docker Engine nativo do Linux.