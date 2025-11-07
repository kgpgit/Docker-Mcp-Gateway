# Guia de Configuração - Docker MCP Gateway Modo Nativo

Este guia explica em detalhes as variáveis de ambiente, opções de configuração e recursos avançados do Docker MCP Gateway em modo nativo.

## Variáveis de Ambiente Principais

### DOCKER_MCP_NATIVE_MODE (Opcional)

**Descrição**: Força o modo de operação nativo para Docker Engine do Linux. Se não definido, o sistema detectará automaticamente o ambiente.

**Valores**:
- `1` - Força modo nativo
- `0` - Força modo Docker Desktop
- Não definido - **Detecção automática** (recomendado)

**Exemplo**:
```bash
# Forçar modo nativo (geralmente não necessário)
export DOCKER_MCP_NATIVE_MODE=1

# Forçar modo Docker Desktop (se a detecção automática falhar)
export DOCKER_MCP_NATIVE_MODE=0
```

**Impacto**:
- Desabilita verificações do Docker Desktop (modo nativo)
- Ativa gerenciamento alternativo de segredos (modo nativo)
- Desabilita monitor OAuth (modo nativo)
- Configura paths padrão para ambiente nativo

### Detecção Automática

O Docker MCP Gateway agora detecta automaticamente o ambiente Linux nativo seguindo estas regras:

1. **Override manual**: Se `DOCKER_MCP_NATIVE_MODE` estiver definida, respeita o valor
2. **Sistema operacional**: Apenas Linux pode ter modo nativo
3. **Modo contêiner**: Se `DOCKER_MCP_IN_CONTAINER=1`, não usa modo nativo
4. **Socket Docker Desktop**: Verifica se `/run/host-services/backend.sock` não existe
5. **Docker Engine**: Verifica se `/var/run/docker.sock` existe

**Exemplo de saída**:
```
- Detecção automática: Linux nativo detectado, ativando modo nativo
```

### DOCKER_MCP_SECRETS_FILE

**Descrição**: Especifica o caminho para o arquivo de segredos em formato `.env`.

**Valores**: Caminho absoluto ou relativo para arquivo de segredos

**Exemplo**:
```bash
export DOCKER_MCP_SECRETS_FILE=$HOME/.docker/mcp/secrets.env
```

**Formato do Arquivo**:
```bash
# Comentários são ignorados
API_KEY=valor_da_chave_api
DATABASE_URL=postgresql://user:pass@localhost/db
SECRET_TOKEN=token_secreto_aqui
```

### DOCKER_MCP_IN_CONTAINER

**Descrição**: Indica se o gateway está rodando dentro de um contêiner Docker.

**Valores**:
- `1` - Modo contêiner
- Não definido ou `0` - Modo host

**Nota**: Esta variável é geralmente definida automaticamente pelo Docker e não precisa ser configurada manualmente.

## Configuração de Segredos

### Métodos Suportados

#### 1. Arquivo .env (Recomendado)
```bash
# Criar arquivo de segredos
cat > ~/.docker/mcp/secrets.env << 'EOF'
# Segredos do MCP Gateway
GITHUB_TOKEN=ghp_xxxxxxxxxxxxxxxxxxxx
OPENAI_API_KEY=sk-xxxxxxxxxxxxxxxxxxxx
DATABASE_URL=postgresql://user:password@localhost:5432/mcp
REDIS_URL=redis://localhost:6379/0
EOF

# Configurar variável de ambiente
export DOCKER_MCP_SECRETS_FILE=$HOME/.docker/mcp/secrets.env
```

#### 2. Variáveis de Ambiente
```bash
# Definir segredos diretamente
export GITHUB_TOKEN=ghp_xxxxxxxxxxxxxxxxxxxx
export OPENAI_API_KEY=sk-xxxxxxxxxxxxxxxxxxxx

# O gateway irá ler automaticamente estas variáveis
```

#### 3. Combinação de Métodos
```bash
# Configurar arquivo base
export DOCKER_MCP_SECRETS_FILE=$HOME/.docker/mcp/secrets.env

# Adicionar ou sobrescrever segredos específicos
export TEMP_SECRET=temporary_value
```

### Boas Práticas de Segredos

#### 1. Permissões de Arquivo
```bash
# Restringir permissões do arquivo de segredos
chmod 600 ~/.docker/mcp/secrets.env
```

#### 2. Segredos Sensíveis
```bash
# Usar ferramentas de gerenciamento de segredos
# Exemplo com HashiCorp Vault
vault kv get -field=token secret/mcp/github > ~/.docker/mcp/github.token
```

#### 3. Rotação de Segredos
```bash
# Script para rotação de segredos
#!/bin/bash
NEW_TOKEN=$(generate_new_token)
sed -i "s/OPENAI_API_KEY=.*/OPENAI_API_KEY=$NEW_TOKEN/" ~/.docker/mcp/secrets.env
```

## Configuração do Gateway

### Opções de Linha de Comando

#### Transporte
```bash
# STDIO (padrão)
docker mcp gateway run --transport stdio

# Server-Sent Events
docker mcp gateway run --transport sse --port 8080

# Streaming HTTP
docker mcp gateway run --transport streaming --port 8080
```

#### Catalogos
```bash
# Usar catálogo padrão
docker mcp gateway run --catalog docker-mcp.yaml

# Múltiplos catálogos
docker mcp gateway run --catalog catalog1.yaml --catalog catalog2.yaml

# Catálogo remoto
docker mcp gateway run --catalog https://example.com/catalog.yaml
```

#### Servidores
```bash
# Servidores específicos
docker mcp gateway run --servers server1,server2

# Habilitar todos os servidores
docker mcp gateway run --enable-all-servers
```

#### Recursos
```bash
# Configurar CPU e memória
docker mcp gateway run --cpus 2 --memory 4Gb

# Habilitar logs detalhados
docker mcp gateway run --verbose --log-calls

# Bloquear acesso à rede
docker mcp gateway run --block-network
```

### Arquivos de Configuração

#### Estrutura de Diretórios
```
~/.docker/mcp/
├── catalogs/
│   ├── docker-mcp.yaml
│   └── custom.yaml
├── servers/
│   └── server-configs/
├── secrets.env
├── registry.yaml
├── config.yaml
└── tools.yaml
```

#### Catalog Configuration
```yaml
# ~/.docker/mcp/catalogs/custom.yaml
servers:
  github-mcp:
    image: docker/github-mcp:latest
    description: "GitHub integration MCP server"
    environment:
      - GITHUB_TOKEN
    secrets:
      - GITHUB_TOKEN
      
  database-mcp:
    image: docker/db-mcp:latest
    description: "Database operations MCP server"
    environment:
      - DATABASE_URL
    secrets:
      - DATABASE_URL
```

#### Registry Configuration
```yaml
# ~/.docker/mcp/registry.yaml
servers:
  custom-registry:
    url: https://registry.example.com
    auth:
      username: ${REGISTRY_USER}
      password: ${REGISTRY_PASS}
    tls:
      verify: true
      cert_file: /path/to/cert.pem
```

#### Tools Configuration
```yaml
# ~/.docker/mcp/tools.yaml
tools:
  custom-tool:
    name: "Custom Operations"
    description: "Custom tool for specific operations"
    image: docker/custom-tool:latest
    command: ["--option1", "value1"]
    environment:
      - CUSTOM_VAR=value
```

## Configuração Avançada

### Autenticação

#### Modo STDIO (Sem Autenticação)
```bash
# Padrão para modo nativo (detecção automática)
docker mcp gateway run --transport stdio
```

#### Modo SSE/HTTP (Sem Autenticação em Modo Nativo)
```bash
# Em modo nativo, autenticação é desabilitada automaticamente
docker mcp gateway run --transport sse --port 8080
```

#### Modo SSE/HTTP (Com Autenticação Manual)
```bash
# Forçar autenticação (não recomendado em modo nativo)
export MCP_GATEWAY_AUTH_TOKEN=seu_token_secreto
docker mcp gateway run --transport sse --port 8080
```

#### Forçar Modo Específico
```bash
# Forçar modo Docker Desktop mesmo em Linux
export DOCKER_MCP_NATIVE_MODE=0
docker mcp gateway run --transport sse --port 8080

# Forçar modo nativo mesmo com Docker Desktop instalado
export DOCKER_MCP_NATIVE_MODE=1
docker mcp gateway run --transport sse --port 8080
```

### Interceptors

#### Configuração de Interceptors
```bash
# Antes da execução
docker mcp gateway run --interceptor before:exec:/path/to/pre-hook

# Depois da execução
docker mcp gateway run --interceptor after:exec:/path/to/post-hook

# Múltiplos interceptors
docker mcp gateway run \
  --interceptor before:exec:/path/to/pre-hook \
  --interceptor after:exec:/path/to/post-hook
```

#### Exemplo de Interceptor
```bash
#!/bin/bash
# ~/.docker/mcp/interceptors/log-request.sh
echo "[$(date)] Request: $1" >> ~/.docker/mcp/logs/requests.log
```

### Working Sets

#### Configurar Working Set
```bash
# Usar working set específico
docker mcp gateway run --working-set development

# Listar working sets disponíveis
docker mcp working-set ls
```

#### Working Set Configuration
```yaml
# ~/.docker/mcp/working-sets/development.yaml
servers:
  - github-mcp
  - database-mcp
  - redis-mcp
  
config:
  log_level: debug
  verbose: true
  
tools:
  - custom-dev-tool
```

### OCI References

#### Usar Imagens OCI
```bash
# Referência direta à imagem OCI
docker mcp gateway run --oci-ref docker.io/user/mcp-server:latest

# Múltiplas referências OCI
docker mcp gateway run \
  --oci-ref docker.io/user/server1:latest \
  --oci-ref docker.io/user/server2:latest
```

## Configuração de Rede

### Modo Nativo
```bash
# Detectar automaticamente redes disponíveis
docker mcp gateway run

# Especificar redes manualmente
docker mcp gateway run --network bridge --network host
```

### Modo Contêiner
```bash
# O gateway detecta automaticamente as redes do contêiner
export DOCKER_MCP_IN_CONTAINER=1
docker mcp gateway run
```

## Configuração de Logs

### Níveis de Log
```bash
# Verboso
docker mcp gateway run --verbose

# Log de chamadas
docker mcp gateway run --log-calls

# Arquivo de log
docker mcp gateway run --log ~/.docker/mcp/logs/gateway.log
```

### Configuração Avançada de Logs
```bash
# Configurar rotação de logs
docker mcp gateway run --log ~/.docker/mcp/logs/gateway.log

# Com logrotate
cat > /etc/logrotate.d/docker-mcp << 'EOF'
~/.docker/mcp/logs/*.log {
    daily
    rotate 7
    compress
    delaycompress
    missingok
    notifempty
    create 644 docker docker
}
EOF
```

## Configuração de Performance

### Recursos do Servidor
```bash
# Configurar recursos para cada servidor MCP
docker mcp gateway run --cpus 2 --memory 4Gb

# Servidores de longa duração
docker mcp gateway run --long-lived

# Modo estático (servidores pré-iniciados)
docker mcp gateway run --static
```

### Otimizações
```bash
# Desabilitar verificação de assinaturas (mais rápido)
docker mcp gateway run --verify-signatures=false

# Bloquear segredos em logs
docker mcp gateway run --block-secrets

# Bloquear acesso à rede
docker mcp gateway run --block-network
```

## Configuração de Sessão

### Sessões Nomeadas
```bash
# Criar sessão específica
docker mcp gateway run --session development

# Usar configuração de sessão
docker mcp gateway run --session production
```

### Estrutura de Sessão
```
~/.docker/mcp/
├── development/
│   ├── registry.yaml
│   ├── config.yaml
│   └── tools.yaml
├── production/
│   ├── registry.yaml
│   ├── config.yaml
│   └── tools.yaml
└── staging/
    ├── registry.yaml
    ├── config.yaml
    └── tools.yaml
```

## Variáveis de Ambiente Adicionais

### Telemetria
```bash
# Habilitar telemetria
export DOCKER_MCP_TELEMETRY=1

# Intervalo de métricas
export DOCKER_MCP_METRICS_INTERVAL=60s

# Debug de telemetria
export DOCKER_MCP_TELEMETRY_DEBUG=1
```

### Debug
```bash
# Mostrar comandos ocultos
export DOCKER_MCP_SHOW_HIDDEN=1

# Debug de DNS
docker mcp gateway run --debug-dns
```

## Exemplos de Configuração Completa

### Ambiente de Desenvolvimento
```bash
#!/bin/bash
# setup-dev.sh

# Configurar segredos de desenvolvimento
export DOCKER_MCP_SECRETS_FILE=$HOME/.docker/mcp/dev-secrets.env

# Executar com configuração de desenvolvimento (detecção automática)
docker mcp gateway run \
  --working-set development \
  --verbose \
  --log-calls \
  --log ~/.docker/mcp/logs/dev-gateway.log \
  --enable-all-servers
```

### Ambiente de Produção
```bash
#!/bin/bash
# setup-prod.sh

# Configurar segredos de produção
export DOCKER_MCP_SECRETS_FILE=/opt/docker-mcp/secrets.env

# Executar com configuração de produção (detecção automática)
docker mcp gateway run \
  --working-set production \
  --transport sse \
  --port 8080 \
  --servers github-mcp,database-mcp \
  --cpus 2 \
  --memory 4Gb \
  --verify-signatures \
  --block-secrets \
  --log /var/log/docker-mcp/gateway.log
```

### Forçar Modo Docker Desktop em Ambiente Linux
```bash
#!/bin/bash
# setup-docker-desktop.sh

# Forçar modo Docker Desktop mesmo em Linux
export DOCKER_MCP_NATIVE_MODE=0

# Executar com configuração Docker Desktop
docker mcp gateway run \
  --working-set development \
  --verbose \
  --log-calls
```

## Validação de Configuração

### Testar Configuração
```bash
# Teste seco (dry-run)
docker mcp gateway run --dry-run

# Verificar configuração
docker mcp config dump

# Listar servidores disponíveis
docker mcp server ls
```

### Script de Validação
```bash
#!/bin/bash
# validate-config.sh

echo "Validando configuração do Docker MCP Gateway..."

# Verificar detecção automática
echo "Verificando detecção automática..."
docker mcp gateway run --dry-run --verbose 2>&1 | grep -i "detecção\|detection" || echo "Mensagem de detecção não encontrada"

# Verificar variáveis de ambiente (se definidas manualmente)
if [ -n "$DOCKER_MCP_NATIVE_MODE" ]; then
    echo "Modo forçado via DOCKER_MCP_NATIVE_MODE: $DOCKER_MCP_NATIVE_MODE"
else
    echo "Usando detecção automática (DOCKER_MCP_NATIVE_MODE não definido)"
fi

# Verificar arquivo de segredos
if [ -n "$DOCKER_MCP_SECRETS_FILE" ]; then
    if [ -f "$DOCKER_MCP_SECRETS_FILE" ]; then
        echo "✓ Arquivo de segredos encontrado: $DOCKER_MCP_SECRETS_FILE"
    else
        echo "✗ Arquivo de segredos não encontrado: $DOCKER_MCP_SECRETS_FILE"
    fi
else
    echo "Arquivo de segredos não especificado (usando padrão)"
fi

# Testar configuração
if docker mcp gateway run --dry-run; then
    echo "✓ Configuração válida"
else
    echo "✗ Erro na configuração"
fi
```

## Próximos Passos

Após configurar o ambiente:

1. Consulte [Soluções de Problemas](solucoes-problemas.md) para problemas comuns
2. Veja a [Referência Técnica](referencia-tecnica.md) para detalhes de implementação
3. Explore os exemplos no diretório `examples/` do projeto original