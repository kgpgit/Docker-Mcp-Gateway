# Detecção Automática de Ambiente Linux Nativo

Este documento explica em detalhes como funciona o sistema de detecção automática do Docker MCP Gateway para identificar ambientes Linux nativos.

## Visão Geral

O Docker MCP Gateway agora inclui um sistema inteligente de detecção automática que identifica quando está rodando em um ambiente Linux nativo (sem Docker Desktop), eliminando a necessidade de configuração manual.

## Como Funciona a Detecção

### Fluxo de Decisão

A detecção segue uma sequência de verificações em ordem de prioridade:

```
Iniciar Detecção
        ↓
DOCKER_MCP_NATIVE_MODE está definido?
     Sim                Não
      ↓                  ↓
   Usar valor         Verificar SO = Linux?
   manual               Sim      Não
   (1/0)                ↓        ↓
      ↓            Modo      Modo
   Retornar      Contêiner? Docker Desktop
   resultado       Sim   Não      ↓
      ↓            ↓     ↓      Verificar
   Modo          Não    Sim    Socket
   Específico     ↓      ↓      Docker
      ↓         Modo    Modo   Desktop
   Retornar     Nativo  Contêiner   ↓
   resultado     ↓        ↓      Sim Não
                ↓        ↓       ↓   ↓
            Verificar  Retornar  Modo  Modo
            Socket     Falso     Docker Docker
            Docker               Desktop
            Engine
                ↓
            Sim Não
             ↓   ↓
        Modo  Retornar
        Nativo Falso
```

### Regras de Detecção

#### 1. Override Manual (Prioridade Máxima)

Se a variável `DOCKER_MCP_NATIVE_MODE` estiver definida, o sistema respeita o valor:

```bash
# Forçar modo nativo
export DOCKER_MCP_NATIVE_MODE=1

# Forçar modo Docker Desktop
export DOCKER_MCP_NATIVE_MODE=0
```

**Quando usar**: Em casos especiais onde a detecção automática não funciona corretamente.

#### 2. Verificação de Sistema Operacional

Apenas sistemas Linux podem ter modo nativo:

```go
if runtime.GOOS != "linux" {
    return false
}
```

**Racional**: O modo nativo foi projetado especificamente para Docker Engine em Linux.

#### 3. Verificação de Modo Contêiner

Se estiver rodando dentro de um contêiner Docker, não usa modo nativo:

```bash
# Detectado automaticamente pelo Docker
export DOCKER_MCP_IN_CONTAINER=1
```

**Racional**: Contêineres têm necessidades diferentes de configuração.

#### 4. Verificação do Socket do Docker Desktop

Verifica se o socket do Docker Desktop existe:

```go
if _, err := os.Stat("/run/host-services/backend.sock"); err == nil {
    return false  // Docker Desktop detectado
}
```

**Racional**: A presença deste socket indica que o Docker Desktop está instalado e rodando.

#### 5. Verificação do Docker Engine

Verifica se o socket do Docker Engine nativo existe:

```go
if _, err := os.Stat("/var/run/docker.sock"); err != nil {
    return false  // Docker Engine não encontrado
}
```

**Racional**: Confirma que o Docker Engine nativo está disponível.

## Implementação Técnica

### Arquivo: pkg/desktop/detection.go

O sistema de detecção está centralizado no arquivo [`pkg/desktop/detection.go`](../pkg/desktop/detection.go):

```go
// IsNativeMode detecta automaticamente se o sistema está rodando em Linux nativo
func IsNativeMode() bool {
    // 1. Respeita override manual via variável de ambiente
    if env := os.Getenv("DOCKER_MCP_NATIVE_MODE"); env != "" {
        return env == "1"
    }

    // 2. Apenas Linux pode ter modo nativo
    if runtime.GOOS != "linux" {
        return false
    }

    // 3. Se estiver em contêiner, não é modo nativo
    if os.Getenv("DOCKER_MCP_IN_CONTAINER") == "1" {
        return false
    }

    // 4. Verifica se o socket do Docker Desktop não existe
    if _, err := os.Stat("/run/host-services/backend.sock"); err == nil {
        return false
    }

    // 5. Verifica se o Docker Engine está disponível
    if _, err := os.Stat("/var/run/docker.sock"); err != nil {
        return false
    }

    // Se chegou aqui, é modo nativo
    log.Log("- Detecção automática: Linux nativo detectado, ativando modo nativo")
    return true
}
```

### Funções Auxiliares

#### IsContainerMode()

Detecta se está rodando em contêiner:

```go
func IsContainerMode() bool {
    return os.Getenv("DOCKER_MCP_IN_CONTAINER") == "1"
}
```

#### IsDockerDesktop()

Detecta se está rodando com Docker Desktop:

```go
func IsDockerDesktop() bool {
    // Se estiver em modo contêiner ou modo nativo, não é Docker Desktop
    if IsContainerMode() || IsNativeMode() {
        return false
    }

    // Verifica se o socket do Docker Desktop existe
    if _, err := os.Stat("/run/host-services/backend.sock"); err == nil {
        return true
    }

    // Verifica paths alternativos
    paths := Paths()
    if paths.BackendSocket != "" {
        if _, err := os.Stat(paths.BackendSocket); err == nil {
            return true
        }
    }

    return false
}
```

## Cenários de Uso

### Cenário 1: Linux Nativo Puro

**Ambiente**: Linux sem Docker Desktop instalado

```
SO: Linux ✓
DOCKER_MCP_IN_CONTAINER: não definido ✓
/run/host-services/backend.sock: não existe ✓
/var/run/docker.sock: existe ✓
Resultado: Modo nativo ativado
```

**Saída esperada**:
```
- Detecção automática: Linux nativo detectado, ativando modo nativo
```

### Cenário 2: Linux com Docker Desktop

**Ambiente**: Linux com Docker Desktop instalado e rodando

```
SO: Linux ✓
DOCKER_MCP_IN_CONTAINER: não definido ✓
/run/host-services/backend.sock: existe ✗
Resultado: Modo Docker Desktop
```

### Cenário 3: Contêiner Docker

**Ambiente**: Contêiner Docker rodando em qualquer sistema

```
DOCKER_MCP_IN_CONTAINER=1
Resultado: Modo contêiner
```

### Cenário 4: Override Manual

**Ambiente**: Linux nativo, mas forçando modo Docker Desktop

```bash
export DOCKER_MCP_NATIVE_MODE=0
Resultado: Modo Docker Desktop (override manual)
```

### Cenário 5: Windows/macOS

**Ambiente**: Windows ou macOS

```
SO: Windows/macOS ✗
Resultado: Modo Docker Desktop
```

## Comportamento Esperado

### Modo Nativo Detectado

Quando o modo nativo é detectado:

1. **Log informativo**: `- Detecção automática: Linux nativo detectado, ativando modo nativo`
2. **Configuração automática**: Paths configurados para ambiente nativo
3. **Segredos**: Usa sistema alternativo de gerenciamento de segredos
4. **OAuth**: Desabilitado automaticamente
5. **Autenticação**: Não requerida para transportes HTTP/SSE

### Modo Docker Desktop Detectado

Quando o Docker Desktop é detectado:

1. **Configuração padrão**: Usa configurações originais
2. **Segredos**: Usa sockets do Docker Desktop
3. **OAuth**: Habilitado se configurado
4. **Autenticação**: Requerida para transportes HTTP/SSE

## Solução de Problemas

### Problema: Detecção Incorreta

**Sintoma**: O sistema detecta modo incorreto

**Solução**: Use override manual:

```bash
# Forçar modo nativo
export DOCKER_MCP_NATIVE_MODE=1

# Forçar modo Docker Desktop
export DOCKER_MCP_NATIVE_MODE=0
```

### Problema: Socket Docker Engine Não Encontrado

**Sintoma**: Erro ao acessar `/var/run/docker.sock`

**Solução**: Verifique se o Docker Engine está rodando:

```bash
sudo systemctl status docker
sudo systemctl start docker
```

### Problema: Permissão Negada

**Sintoma**: Erro de permissão ao acessar sockets

**Solução**: Adicione usuário ao grupo docker:

```bash
sudo usermod -aG docker $USER
newgrp docker
```

## Testes de Validação

### Teste Básico de Detecção

```bash
# Limpar variáveis
unset DOCKER_MCP_NATIVE_MODE

# Testar detecção automática
docker mcp gateway run --dry-run --verbose 2>&1 | grep -i "detecção"
```

### Teste de Override

```bash
# Testar override para modo nativo
export DOCKER_MCP_NATIVE_MODE=1
docker mcp gateway run --dry-run

# Testar override para modo Docker Desktop
export DOCKER_MCP_NATIVE_MODE=0
docker mcp gateway run --dry-run
```

### Script Completo de Teste

```bash
#!/bin/bash
# test-detection.sh

echo "=== Teste de Detecção Automática ==="

# Teste 1: Detecção automática
echo "Teste 1: Detecção automática"
unset DOCKER_MCP_NATIVE_MODE
docker mcp gateway run --dry-run --verbose 2>&1 | grep -q "Linux nativo detectado"
if [ $? -eq 0 ]; then
    echo "✓ Detecção automática funcionando"
else
    echo "✗ Detecção automática não funcionando"
fi

# Teste 2: Override para modo nativo
echo "Teste 2: Override para modo nativo"
export DOCKER_MCP_NATIVE_MODE=1
docker mcp gateway run --dry-run
if [ $? -eq 0 ]; then
    echo "✓ Override para modo nativo funcionando"
else
    echo "✗ Override para modo nativo não funcionando"
fi

# Teste 3: Override para modo Docker Desktop
echo "Teste 3: Override para modo Docker Desktop"
export DOCKER_MCP_NATIVE_MODE=0
docker mcp gateway run --dry-run
if [ $? -eq 0 ]; then
    echo "✓ Override para modo Docker Desktop funcionando"
else
    echo "✗ Override para modo Docker Desktop não funcionando"
fi

# Limpar
unset DOCKER_MCP_NATIVE_MODE

echo "=== Testes concluídos ==="
```

## Perguntas Frequentes

### P: Preciso configurar alguma variável de ambiente?

R: Não. A detecção automática funciona sem configuração manual. As variáveis são opcionais apenas para casos especiais.

### P: Como sei qual modo foi detectado?

R: Execute com `--verbose` e procure pela mensagem de detecção:

```bash
docker mcp gateway run --dry-run --verbose
```

### P: A detecção automática funciona em contêineres?

R: Não. Em contêineres, o sistema usa modo contêiner automaticamente.

### P: Posso forçar um modo específico?

R: Sim. Use `DOCKER_MCP_NATIVE_MODE=1` para forçar modo nativo ou `DOCKER_MCP_NATIVE_MODE=0` para forçar modo Docker Desktop.

### P: O que acontece se a detecção falhar?

R: Use override manual ou verifique se os sockets estão acessíveis e com permissões corretas.

## Evolução Futura

### Melhorias Planejadas

1. **Detecção mais robusta**: Verificação de múltiplos indicadores
2. **Logs detalhados**: Mais informações sobre o processo de detecção
3. **Configuração persistente**: Opção de salvar preferência de modo
4. **Detecção de distribuições**: Suporte específico para distribuições Linux

### Contribuições

Para contribuir com melhorias no sistema de detecção:

1. Teste em diferentes ambientes
2. Reporte casos de detecção incorreta
3. Sugira novos critérios de detecção
4. Contribua com código para melhorias

## Conclusão

O sistema de detecção automática representa um avanço significativo na usabilidade do Docker MCP Gateway, eliminando a necessidade de configuração manual e proporcionando uma experiência mais fluida para usuários de Linux nativo.

A implementação é robusta, flexível e mantém compatibilidade total com casos especiais através do sistema de override manual.