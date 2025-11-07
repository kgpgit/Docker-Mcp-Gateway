# Guia de Instalação - Docker MCP Gateway Modo Nativo

Este guia fornece instruções detalhadas para compilar e configurar o Docker MCP Gateway para operação com Docker Engine nativo do Linux.

## Pré-requisitos

### Sistema Operacional
- Linux (Ubuntu 20.04+, Debian 11+, CentOS 8+, Fedora 35+ ou similar)
- Acesso a um terminal com privilégios de usuário

### Software Necessário
- **Docker Engine** (versão 20.10+ recomendada)
- **Docker CLI** (compatível com a versão do Engine)
- **Go** (versão 1.19 ou superior)
- **Git** (para clonar o repositório)
- **Make** (para build automatizado)

### Verificação de Pré-requisitos

#### Verificar Docker Engine
```bash
docker --version
# Saída esperada: Docker version 20.10.x ou superior

docker info
# Deve mostrar informações do Docker Engine nativo
```

#### Verificar Go
```bash
go version
# Saída esperada: go version go1.19+ linux/amd64
```

#### Verificar Git e Make
```bash
git --version
make --version
```

## Passo 1: Obter o Código Fonte

### Clonar o Repositório
```bash
# Clonar o repositório original
git clone https://github.com/docker/mcp-gateway.git
cd mcp-gateway

# OU, se já tiver o repositório, atualizar
git pull origin main
```

### Verificar Branch
```bash
git branch
# Certifique-se de estar na branch correta (main ou develop)
```

## Passo 2: Compilar o Código

### Método 1: Usando Make (Recomendado)
```bash
# Compilar o binário
make build

# Verificar se o binário foi criado
ls -la bin/docker-mcp
```

### Método 2: Compilação Manual
```bash
# Baixar dependências
go mod download

# Compilar
go build -o bin/docker-mcp ./cmd/docker-mcp

# Tornar executável
chmod +x bin/docker-mcp
```

### Método 3: Instalação Global
```bash
# Compilar e instalar globalmente
go install ./cmd/docker-mcp

# Verificar instalação
which docker-mcp
docker-mcp --version
```

## Passo 3: Configurar o Ambiente

### 3.1 Configuração de Segredos (Opcional)

O Docker MCP Gateway agora detecta automaticamente o ambiente Linux nativo, não sendo mais necessário configurar a variável `DOCKER_MCP_NATIVE_MODE`.

```bash
# Criar arquivo de segredos (opcional)
echo "API_KEY=sua_chave_api" > ~/.docker/mcp/secrets.env
echo "DATABASE_URL=postgresql://user:pass@localhost/db" >> ~/.docker/mcp/secrets.env

# Configurar variável de ambiente para segredos (opcional)
export DOCKER_MCP_SECRETS_FILE=$HOME/.docker/mcp/secrets.env
```

### 3.2 Forçar Modo Específico (Opcional)

Se precisar forçar um modo específico, você ainda pode usar a variável de ambiente:

```bash
# Forçar modo nativo
export DOCKER_MCP_NATIVE_MODE=1

# Forçar modo Docker Desktop
export DOCKER_MCP_NATIVE_MODE=0
```

### 3.3 Criar Diretórios Necessários
```bash
# Criar diretórios de configuração
mkdir -p ~/.docker/mcp/catalogs
mkdir -p ~/.docker/mcp/servers
mkdir -p ~/.docker/mcp/logs
```

### 3.4 Configurar Permissões do Docker
```bash
# Adicionar usuário ao grupo docker (se necessário)
sudo usermod -aG docker $USER

# Aplicar alterações (requer logout/login)
newgrp docker

# Verificar permissões
docker ps
```

## Passo 4: Verificar Instalação

### 4.1 Teste Básico
```bash
# Verificar versão
docker mcp version

# Testar ajuda
docker mcp --help

# Listar comandos disponíveis
docker mcp gateway --help
```

### 4.2 Teste de Configuração
```bash
# Testar configuração em modo dry-run
docker mcp gateway run --dry-run

# Saída esperada: Configuração carregada sem erros
```

### 4.3 Teste de Segredos
```bash
# Criar arquivo de teste
echo "TEST_SECRET=test_value" > test-secrets.env

# Testar leitura de segredos
export DOCKER_MCP_SECRETS_FILE=./test-secrets.env
docker mcp gateway run --dry-run --secrets file:./test-secrets.env

# Limpar
rm test-secrets.env
```

## Passo 5: Configurar Automação (Opcional)

### 5.1 Criar Script de Inicialização
```bash
# Criar script (detecção automática)
cat > ~/.local/bin/docker-mcp << 'EOF'
#!/bin/bash
export DOCKER_MCP_SECRETS_FILE=$HOME/.docker/mcp/secrets.env
exec docker mcp "$@"
EOF

# Tornar executável
chmod +x ~/.local/bin/docker-mcp

# Adicionar ao PATH (se necessário)
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc
```

### 5.2 Configurar Systemd Service (Opcional)
```bash
# Criar arquivo de serviço (detecção automática)
sudo tee /etc/systemd/system/docker-mcp.service > /dev/null << 'EOF'
[Unit]
Description=Docker MCP Gateway
After=docker.service
Requires=docker.service

[Service]
Type=simple
User=docker-mcp
Group=docker
Environment=DOCKER_MCP_SECRETS_FILE=/opt/docker-mcp/secrets.env
ExecStart=/usr/local/bin/docker-mcp gateway run
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF

# Habilitar serviço
sudo systemctl daemon-reload
sudo systemctl enable docker-mcp.service
```

## Passo 6: Configuração Avançada

### 6.1 Configurar Catalogos Personalizados
```bash
# Criar catálogo personalizado
cat > ~/.docker/mcp/catalogs/custom.yaml << 'EOF'
servers:
  my-server:
    image: docker/mcp-server:latest
    command: ["--option1", "value1"]
    environment:
      - VAR1=value1
      - VAR2=value2
EOF
```

### 6.2 Configurar Registry
```bash
# Criar arquivo de registry
cat > ~/.docker/mcp/registry.yaml << 'EOF'
servers:
  custom-registry:
    url: https://my-registry.example.com
    auth:
      username: myuser
      password: mypass
EOF
```

### 6.3 Configurar Logs
```bash
# Configurar arquivo de log
mkdir -p ~/.docker/mcp/logs

# Testar com arquivo de log
docker mcp gateway run --log ~/.docker/mcp/logs/gateway.log
```

## Solução de Problemas na Instalação

### Problema: Permissão Negada ao Acessar Docker Socket
```bash
# Solução 1: Adicionar usuário ao grupo docker
sudo usermod -aG docker $USER
newgrp docker

# Solução 2: Usar sudo temporariamente
sudo docker mcp gateway run
```

### Problema: Go Não Encontrado
```bash
# Instalar Go (Ubuntu/Debian)
sudo apt update
sudo apt install golang-go

# Instalar Go (CentOS/RHEL)
sudo yum install golang

# Instalar Go (Fedora)
sudo dnf install golang
```

### Problema: Erro de Compilação
```bash
# Limpar cache de módulos
go clean -modcache
go mod download

# Recompilar
make clean
make build
```

### Problema: Detecção Automática Não Funcionando
```bash
# Verificar se o sistema detectou corretamente
docker mcp gateway run --dry-run --verbose

# Verificar variáveis (se definidas manualmente)
echo $DOCKER_MCP_NATIVE_MODE
echo $DOCKER_MCP_SECRETS_FILE

# Forçar modo nativo se necessário
export DOCKER_MCP_NATIVE_MODE=1
docker mcp gateway run --dry-run
```

## Verificação Final

### Script de Verificação Completo
```bash
#!/bin/bash
echo "=== Verificação da Instalação do Docker MCP Gateway ==="

# Verificar Docker
echo "1. Verificando Docker..."
docker --version || { echo "Docker não encontrado"; exit 1; }

# Verificar Go
echo "2. Verificando Go..."
go version || { echo "Go não encontrado"; exit 1; }

# Verificar binário
echo "3. Verificando binário..."
which docker-mcp || { echo "docker-mcp não encontrado no PATH"; exit 1; }

# Verificar detecção automática
echo "4. Verificando detecção automática..."
docker mcp gateway run --dry-run --verbose 2>&1 | grep -i "detecção\|detection" || echo "Mensagem de detecção não encontrada"

# Verificar variáveis (se definidas manualmente)
echo "DOCKER_MCP_NATIVE_MODE: ${DOCKER_MCP_NATIVE_MODE:-não definido (detecção automática)}"
echo "DOCKER_MCP_SECRETS_FILE: ${DOCKER_MCP_SECRETS_FILE:-não definido}"

# Testar funcionamento
echo "5. Testando funcionamento..."
docker mcp version || { echo "Erro ao executar docker-mcp"; exit 1; }

echo "=== Instalação verificada com sucesso! ==="
```

## Próximos Passos

Após a instalação bem-sucedida:

1. Leia o [Guia de Configuração](guia-configuracao.md) para detalhes sobre opções avançadas
2. Consulte [Soluções de Problemas](solucoes-problemas.md) para problemas comuns
3. Veja a [Referência Técnica](referencia-tecnica.md) para detalhes de implementação

## Atualização

Para atualizar para uma versão mais recente:

```bash
# Atualizar código fonte
git pull origin main

# Recompilar
make clean
make build

# Verificar versão
docker mcp version