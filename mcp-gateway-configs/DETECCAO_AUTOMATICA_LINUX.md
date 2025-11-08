# Detecção Automática do Ambiente Linux Nativo

## Overview

Esta implementação adiciona detecção automática do ambiente Linux nativo ao Docker MCP Gateway, eliminando a necessidade de definir manualmente a variável de ambiente `DOCKER_MCP_NATIVE_MODE`.

## Funcionalidades Implementadas

### 1. Funções de Detecção (`pkg/desktop/detection.go`)

#### `IsNativeMode()`

Detecta automaticamente se o sistema está rodando em Linux nativo (sem Docker Desktop):

**Lógica de detecção:**


1. **Override manual**: Se `DOCKER_MCP_NATIVE_MODE` estiver definida, respeita o valor manual
2. **Verificação de SO**: Apenas Linux pode ter modo nativo
3. **Modo contêiner**: Se `DOCKER_MCP_IN_CONTAINER=1`, não é modo nativo
4. **Socket Docker Desktop**: Verifica se `/run/host-services/backend.sock` não existe
5. **Docker Engine**: Verifica se `/var/run/docker.sock` existe

**Retorno:** `true` se todas as condições forem satisfeitas

#### `IsContainerMode()`

Detecta se está rodando em contêiner baseado na variável `DOCKER_MCP_IN_CONTAINER`.

#### `IsDockerDesktop()`

Detecta se está rodando com Docker Desktop verificando a existência dos sockets específicos.

### 2. Modificações nos Arquivos Existentes

#### `pkg/desktop/sockets_linux.go`

* **Antes:** Verificava manualmente `os.Getenv("DOCKER_MCP_NATIVE_MODE")`
* **Depois:** Usa `desktop.IsContainerMode() || desktop.IsNativeMode()`

#### `cmd/docker-mcp/commands/root.go`

* **Antes:** Verificava manualmente as variáveis de ambiente
* **Depois:** Usa `desktop.IsContainerMode() && !desktop.IsNativeMode()`

#### `pkg/docker/secrets.go`

* **Antes:** Verificava manualmente `DOCKER_MCP_IN_CONTAINER` e `DOCKER_MCP_NATIVE_MODE`
* **Depois:** Usa `desktop.IsContainerMode() || desktop.IsNativeMode()`

#### `pkg/gateway/run.go`

* **Antes:** Verificação manual das variáveis para OAuth e autenticação
* **Depois:** Usa as funções de detecção automática

#### `cmd/docker-mcp/commands/gateway.go`

* **Antes:** Verificação manual para configurações padrão
* **Depois:** Usa `desktop.IsContainerMode() || desktop.IsNativeMode()`

## Comportamento Esperado

### Detecção Automática

Em Linux nativo sem Docker Desktop:

```
- Detecção automática: Linux nativo detectado, ativando modo nativo
```

### Override Manual

A variável `DOCKER_MCP_NATIVE_MODE` ainda funciona como override:

```bash
# Forçar modo nativo
export DOCKER_MCP_NATIVE_MODE=1

# Forçar modo Docker Desktop
export DOCKER_MCP_NATIVE_MODE=0
```

### Configurações Automáticas

#### Modo Nativo (Linux sem Docker Desktop)

* **Catálogo:** URL remoto v3
* **Segredos:** `file:./.env`
* **OAuth:** Desabilitado
* **Autenticação:** Desabilitada para SSE/streaming

#### Modo Docker Desktop

* **Catálogo:** Arquivo local `docker-mcp.yaml`
* **Segredos:** `docker-desktop`
* **OAuth:** Habilitado se feature flag ativa
* **Autenticação:** Requerida para SSE/streaming

#### Modo Contêiner

* **Catálogo:** URL remoto v3
* **Segredos:** `file:./.env`
* **OAuth:** Desabilitado
* **Autenticação:** Desabilitada

## Testes Realizados

### 1. Detecção Automática

```bash
./docker-mcp gateway run --dry-run --verbose
```

**Resultado:** ✅ Modo nativo detectado automaticamente

### 2. Override Manual

```bash
DOCKER_MCP_NATIVE_MODE=0 ./docker-mcp gateway run --dry-run --verbose
```

**Resultado:** ✅ Override manual respeitado

### 3. Compilação

```bash
go build ./cmd/docker-mcp
```

**Resultado:** ✅ Sem erros de compilação

## Benefícios


1. **Experiência do usuário:** Não há mais necessidade de configuração manual
2. **Detecção inteligente:** Sistema identifica corretamente o ambiente
3. **Compatibilidade:** Override manual ainda funciona para casos especiais
4. **Logging:** Usuário é informado quando modo nativo é ativado
5. **Manutenibilidade:** Lógica centralizada em um único arquivo

## Compatibilidade

* ✅ **Linux nativo:** Detecção automática funciona
* ✅ **Linux com Docker Desktop:** Funciona normalmente
* ✅ **Contêineres:** Funciona normalmente
* ✅ **Override manual:** Continua funcionando
* ✅ **Outros SOs:** Sem impacto (funções retornam `false`)

## Arquivos Modificados


1. **Novo:** `pkg/desktop/detection.go` - Funções de detecção
2. **Modificado:** `pkg/desktop/sockets_linux.go`
3. **Modificado:** `cmd/docker-mcp/commands/root.go`
4. **Modificado:** `pkg/docker/secrets.go`
5. **Modificado:** `pkg/gateway/run.go`
6. **Modificado:** `cmd/docker-mcp/commands/gateway.go`

## Implementação Futura

Possíveis melhorias:

* Detecção mais refinada para diferentes distribuições Linux
* Suporte para outros engines Docker (Podman, etc.)
* Logging mais detalhado do processo de detecção


