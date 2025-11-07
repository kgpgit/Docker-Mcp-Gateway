# Exemplos Práticos - Docker MCP Gateway Modo Nativo

Este documento fornece exemplos práticos e cenários de uso real para o Docker MCP Gateway em modo nativo.

## Exemplo 1: Configuração Básica de Desenvolvimento

### Cenário
Desenvolvedor quer usar o MCP Gateway com Docker Engine nativo para desenvolvimento local.

### Passos

```bash
#!/bin/bash
# setup-dev-env.sh

# 1. Criar diretórios necessários
mkdir -p ~/.docker/mcp/{catalogs,servers,logs}

# 2. Configurar modo nativo
export DOCKER_MCP_NATIVE_MODE=1

# 3. Criar arquivo de segredos de desenvolvimento
cat > ~/.docker/mcp/dev-secrets.env << 'EOF'
# Segredos de desenvolvimento
GITHUB_TOKEN=ghp_dev_token_xxxxxxxxxxxx
OPENAI_API_KEY=sk_dev_key_xxxxxxxxxxxx
DATABASE_URL=postgresql://dev_user:dev_pass@localhost:5432/dev_db
REDIS_URL=redis://localhost:6379/0
EOF

# 4. Configurar permissões
chmod 600 ~/.docker/mcp/dev-secrets.env

# 5. Configurar variável de ambiente
export DOCKER_MCP_SECRETS_FILE=$HOME/.docker/mcp/dev-secrets.env

# 6. Criar catálogo de desenvolvimento
cat > ~/.docker/mcp/catalogs/dev-catalog.yaml << 'EOF'
servers:
  github-mcp:
    image: docker/github-mcp:latest
    description: "GitHub integration for development"
    environment:
      - GITHUB_TOKEN
    secrets:
      - GITHUB_TOKEN
      
  database-mcp:
    image: docker/db-mcp:latest
    description: "Database operations for development"
    environment:
      - DATABASE_URL
    secrets:
      - DATABASE_URL
      
  redis-mcp:
    image: docker/redis-mcp:latest
    description: "Redis operations for development"
    environment:
      - REDIS_URL
    secrets:
      - REDIS_URL
EOF

echo "Ambiente de desenvolvimento configurado!"
echo "Execute: docker mcp gateway run --catalog ~/.docker/mcp/catalogs/dev-catalog.yaml"
```

### Uso

```bash
# Carregar configuração
source setup-dev-env.sh

# Executar gateway
docker mcp gateway run \
  --catalog ~/.docker/mcp/catalogs/dev-catalog.yaml \
  --verbose \
  --log-calls \
  --log ~/.docker/mcp/logs/dev-gateway.log
```

## Exemplo 2: Ambiente de Produção com Systemd

### Cenário
Configurar o MCP Gateway como serviço systemd em ambiente de produção.

### Arquivo de Serviço

```bash
# /etc/systemd/system/docker-mcp.service
[Unit]
Description=Docker MCP Gateway Service
After=docker.service
Requires=docker.service

[Service]
Type=simple
User=docker-mcp
Group=docker
Environment=DOCKER_MCP_NATIVE_MODE=1
Environment=DOCKER_MCP_SECRETS_FILE=/opt/docker-mcp/secrets.env
ExecStart=/usr/local/bin/docker-mcp gateway run \
  --working-set production \
  --transport sse \
  --port 8080 \
  --cpus 2 \
  --memory 4Gb \
  --verify-signatures \
  --block-secrets \
  --log /var/log/docker-mcp/gateway.log
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
```

### Script de Instalação

```bash
#!/bin/bash
# install-production.sh

# 1. Criar usuário de serviço
sudo useradd -r -s /bin/false docker-mcp
sudo usermod -a -G docker docker-mcp

# 2. Criar diretórios
sudo mkdir -p /opt/docker-mcp
sudo mkdir -p /var/log/docker-mcp
sudo chown docker-mcp:docker-mcp /opt/docker-mcp
sudo chown docker-mcp:docker-mcp /var/log/docker-mcp

# 3. Instalar binário
sudo cp bin/docker-mcp /usr/local/bin/
sudo chmod +x /usr/local/bin/docker-mcp

# 4. Criar arquivo de segredos
sudo tee /opt/docker-mcp/secrets.env > /dev/null << 'EOF'
# Segredos de produção
GITHUB_TOKEN=ghp_prod_token_xxxxxxxxxxxx
OPENAI_API_KEY=sk_prod_key_xxxxxxxxxxxx
DATABASE_URL=postgresql://prod_user:secure_pass@db.example.com:5432/prod_db
REDIS_URL=redis://redis.example.com:6379/0
EOF

# 5. Configurar permissões
sudo chmod 600 /opt/docker-mcp/secrets.env
sudo chown docker-mcp:docker-mcp /opt/docker-mcp/secrets.env

# 6. Instalar serviço systemd
sudo cp docker-mcp.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable docker-mcp.service

# 7. Iniciar serviço
sudo systemctl start docker-mcp.service

echo "Instalação concluída!"
echo "Status: sudo systemctl status docker-mcp"
echo "Logs: sudo journalctl -u docker-mcp -f"
```

## Exemplo 3: Configuração com Múltiplos Ambientes

### Cenário
Gerenciar múltiplos ambientes (desenvolvimento, staging, produção) com working sets.

### Estrutura de Diretórios

```
~/.docker/mcp/
├── working-sets/
│   ├── development/
│   │   ├── registry.yaml
│   │   ├── config.yaml
│   │   └── tools.yaml
│   ├── staging/
│   │   ├── registry.yaml
│   │   ├── config.yaml
│   │   └── tools.yaml
│   └── production/
│       ├── registry.yaml
│       ├── config.yaml
│       └── tools.yaml
├── secrets/
│   ├── dev-secrets.env
│   ├── staging-secrets.env
│   └── prod-secrets.env
└── catalogs/
    ├── dev-catalog.yaml
    ├── staging-catalog.yaml
    └── prod-catalog.yaml
```

### Script de Gerenciamento

```bash
#!/bin/bash
# env-manager.sh

ENVIRONMENT=$1
ACTION=$2

case $ENVIRONMENT in
  dev|development)
    export DOCKER_MCP_SECRETS_FILE=$HOME/.docker/mcp/secrets/dev-secrets.env
    WORKING_SET="development"
    ;;
  staging)
    export DOCKER_MCP_SECRETS_FILE=$HOME/.docker/mcp/secrets/staging-secrets.env
    WORKING_SET="staging"
    ;;
  prod|production)
    export DOCKER_MCP_SECRETS_FILE=$HOME/.docker/mcp/secrets/prod-secrets.env
    WORKING_SET="production"
    ;;
  *)
    echo "Uso: $0 {dev|staging|prod} {start|stop|status}"
    exit 1
    ;;
esac

export DOCKER_MCP_NATIVE_MODE=1

case $ACTION in
  start)
    echo "Iniciando MCP Gateway para ambiente: $ENVIRONMENT"
    docker mcp gateway run --working-set $WORKING_SET
    ;;
  stop)
    echo "Parando MCP Gateway para ambiente: $ENVIRONMENT"
    pkill -f "docker-mcp.*--working-set $WORKING_SET"
    ;;
  status)
    echo "Status do MCP Gateway para ambiente: $ENVIRONMENT"
    ps aux | grep "docker-mcp.*--working-set $WORKING_SET" || echo "Não está rodando"
    ;;
  *)
    echo "Uso: $0 {dev|staging|prod} {start|stop|status}"
    exit 1
    ;;
esac
```

### Uso

```bash
# Iniciar ambiente de desenvolvimento
./env-manager.sh dev start

# Verificar status do ambiente de produção
./env-manager.sh prod status

# Parar ambiente de staging
./env-manager.sh staging stop
```

## Exemplo 4: Configuração com Interceptors

### Cenário
Implementar interceptors para logging e auditoria de operações.

### Script de Interceptor

```bash
#!/bin/bash
# ~/.docker/mcp/interceptors/audit-log.sh

# Interceptor para logging de auditoria
OPERATION=$1
TIMESTAMP=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
USER=$(whoami)
HOSTNAME=$(hostname)

# Log de auditoria
echo "[$TIMESTAMP] User: $USER, Host: $HOSTNAME, Operation: $OPERATION" >> ~/.docker/mcp/logs/audit.log

# Continuar operação
exit 0
```

### Configuração com Interceptors

```bash
#!/bin/bash
# run-with-interceptors.sh

export DOCKER_MCP_NATIVE_MODE=1
export DOCKER_MCP_SECRETS_FILE=$HOME/.docker/mcp/secrets.env

# Criar diretório de logs
mkdir -p ~/.docker/mcp/logs

# Executar com interceptors
docker mcp gateway run \
  --interceptor before:exec:$HOME/.docker/mcp/interceptors/audit-log.sh \
  --interceptor after:exec:$HOME/.docker/mcp/interceptors/audit-log.sh \
  --verbose \
  --log ~/.docker/mcp/logs/gateway.log
```

## Exemplo 5: Configuração com OCI References

### Cenário
Usar imagens OCI diretamente de registries personalizados.

### Script de Configuração

```bash
#!/bin/bash
# oci-setup.sh

export DOCKER_MCP_NATIVE_MODE=1
export DOCKER_MCP_SECRETS_FILE=$HOME/.docker/mcp/secrets.env

# Configurar registry credentials
export REGISTRY_USER=myuser
export REGISTRY_PASS=mypass

# Executar com referências OCI
docker mcp gateway run \
  --oci-ref docker.io/mcp/github-server:latest \
  --oci-ref docker.io/mcp/database-server:v1.2.0 \
  --oci-ref registry.example.com/custom/mcp-server:custom-tag \
  --verbose
```

## Exemplo 6: Configuração com Rede Personalizada

### Cenário
Configurar o MCP Gateway para operar em redes Docker personalizadas.

### Script de Configuração de Rede

```bash
#!/bin/bash
# network-setup.sh

# 1. Criar rede personalizada
docker network create mcp-network --driver bridge

# 2. Configurar variáveis de ambiente
export DOCKER_MCP_NATIVE_MODE=1
export DOCKER_MCP_SECRETS_FILE=$HOME/.docker/mcp/secrets.env

# 3. Executar gateway com rede específica
docker mcp gateway run \
  --network mcp-network \
  --verbose \
  --log ~/.docker/mcp/logs/network-gateway.log
```

## Exemplo 7: Monitoramento e Métricas

### Cenário
Configurar monitoramento do MCP Gateway com métricas detalhadas.

### Script de Monitoramento

```bash
#!/bin/bash
# monitoring-setup.sh

# Configurar telemetria
export DOCKER_MCP_TELEMETRY=1
export DOCKER_MCP_METRICS_INTERVAL=30s
export DOCKER_MCP_TELEMETRY_DEBUG=1

# Configurar modo nativo
export DOCKER_MCP_NATIVE_MODE=1
export DOCKER_MCP_SECRETS_FILE=$HOME/.docker/mcp/secrets.env

# Criar diretório de métricas
mkdir -p ~/.docker/mcp/metrics

# Executar com monitoramento
docker mcp gateway run \
  --verbose \
  --log ~/.docker/mcp/logs/monitored-gateway.log \
  --cpus 2 \
  --memory 4Gb &

# Monitorar processos
echo "Monitorando MCP Gateway..."
while true; do
  echo "$(date): Verificando processo..."
  ps aux | grep "docker-mcp" | grep -v grep
  sleep 60
done
```

## Exemplo 8: Backup e Restauração

### Script de Backup

```bash
#!/bin/bash
# backup-config.sh

BACKUP_DIR="$HOME/docker-mcp-backup-$(date +%Y%m%d-%H%M%S)"
mkdir -p "$BACKUP_DIR"

echo "Criando backup em: $BACKUP_DIR"

# Backup de configurações
cp -r ~/.docker/mcp "$BACKUP_DIR/"

# Backup de variáveis de ambiente
env | grep DOCKER_MCP > "$BACKUP_DIR/env-vars"

# Backup de scripts
cp ~/.local/bin/docker-mcp-* "$BACKUP_DIR/" 2>/dev/null

# Comprimir backup
tar -czf "$BACKUP_DIR.tar.gz" -C "$HOME" "$(basename $BACKUP_DIR)"
rm -rf "$BACKUP_DIR"

echo "Backup criado: $BACKUP_DIR.tar.gz"
```

### Script de Restauração

```bash
#!/bin/bash
# restore-config.sh

BACKUP_FILE=$1

if [ -z "$BACKUP_FILE" ]; then
    echo "Uso: $0 <arquivo_de_backup.tar.gz>"
    exit 1
fi

# Extrair backup
TEMP_DIR=$(mktemp -d)
tar -xzf "$BACKUP_FILE" -C "$TEMP_DIR"

# Restaurar configurações
BACKUP_NAME=$(basename "$BACKUP_FILE" .tar.gz)
cp -r "$TEMP_DIR/$BACKUP_NAME/mcp" ~/.docker/

# Restaurar variáveis de ambiente
source "$TEMP_DIR/$BACKUP_NAME/env-vars"

# Restaurar scripts
cp "$TEMP_DIR/$BACKUP_NAME"/docker-mcp-* ~/.local/bin/ 2>/dev/null

# Limpar
rm -rf "$TEMP_DIR"

echo "Restauração concluída!"
```

## Exemplo 9: Teste de Carga

### Script de Teste de Carga

```bash
#!/bin/bash
# load-test.sh

export DOCKER_MCP_NATIVE_MODE=1
export DOCKER_MCP_SECRETS_FILE=$HOME/.docker/mcp/secrets.env

# Iniciar gateway em background
docker mcp gateway run --transport sse --port 8080 &
GATEWAY_PID=$!

# Esperar inicialização
sleep 10

# Executar teste de carga
echo "Iniciando teste de carga..."
for i in {1..100}; do
  curl -s http://localhost:8080/sse > /dev/null &
  if [ $((i % 10)) -eq 0 ]; then
    echo "Enviadas $i requisições..."
  fi
done

# Esperar conclusão
wait

# Parar gateway
kill $GATEWAY_PID

echo "Teste de carga concluído!"
```

## Exemplo 10: Configuração com SSL/TLS

### Script de Configuração SSL

```bash
#!/bin/bash
# ssl-setup.sh

# Gerar certificados autoassinados
mkdir -p ~/.docker/mcp/ssl
cd ~/.docker/mcp/ssl

# Gerar chave privada
openssl genrsa -out server.key 2048

# Gerar CSR
openssl req -new -key server.key -out server.csr -subj "/C=BR/ST=State/L=City/O=Organization/CN=localhost"

# Gerar certificado
openssl x509 -req -days 365 -in server.csr -signkey server.key -out server.crt

# Configurar modo nativo
export DOCKER_MCP_NATIVE_MODE=1
export DOCKER_MCP_SECRETS_FILE=$HOME/.docker/mcp/secrets.env

# Executar com SSL
docker mcp gateway run \
  --transport sse \
  --port 8443 \
  --tls-cert ~/.docker/mcp/ssl/server.crt \
  --tls-key ~/.docker/mcp/ssl/server.key \
  --verbose
```

## Dicas Adicionais

### Performance
- Use `--long-lived` para servidores stateful
- Configure `--cpus` e `--memory` adequadamente
- Desabilite `--verify-signatures` em ambiente de desenvolvimento

### Segurança
- Mantenha permissões `600` nos arquivos de segredos
- Use variáveis de ambiente para segredos sensíveis
- Configure firewall para modo SSE/HTTP

### Debugging
- Use `--verbose` e `--log-calls` para detalhes
- Configure `DOCKER_MCP_TELEMETRY_DEBUG=1`
- Use `--dry-run` para testar configuração

## Conclusão

Estes exemplos demonstram a flexibilidade do Docker MCP Gateway em modo nativo. Adapte-os conforme suas necessidades específicas e consulte a documentação principal para detalhes adicionais.