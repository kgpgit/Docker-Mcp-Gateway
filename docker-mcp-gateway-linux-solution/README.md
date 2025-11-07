# Docker MCP Gateway - Modo Nativo para Docker Engine

## Visão Geral

Esta documentação descreve uma solução implementada para permitir que o **Docker MCP Gateway** funcione com o **Docker Engine nativo do Linux**, sem depender do Docker Desktop. A solução foi desenvolvida para resolver problemas de compatibilidade e permitir que usuários de Linux possam utilizar o MCP Gateway com sua instalação nativa do Docker.

## Problema Original

O Docker MCP Gateway foi originalmente projetado para funcionar com o Docker Desktop, dependendo de vários componentes específicos:

- Verificação de recursos do Docker Desktop
- Gerenciamento de segredos via sockets específicos do Docker Desktop
- Monitoramento OAuth via backend socket
- Autenticação via credential helper do Docker Desktop

Essas dependências impediam o uso do MCP Gateway em ambientes Linux com Docker Engine nativo, resultando em erros de conexão e funcionalidades limitadas.

## Solução Implementada

A solução implementa um **modo nativo** que permite operação completa do Docker MCP Gateway sem dependências do Docker Desktop. As principais modificações incluem:

### 1. Detecção Automática de Ambiente Nativo
- **Detecção automática**: O sistema agora detecta automaticamente quando está rodando em Linux nativo
- **Override manual**: A variável `DOCKER_MCP_NATIVE_MODE` ainda pode ser usada para forçar um modo específico
- **Bypass inteligente**: Verificações do Docker Desktop são automaticamente ignoradas em modo nativo

### 2. Gerenciamento Alternativo de Segredos
- Leitura de segredos via arquivos `.env` locais
- Suporte a variáveis de ambiente como fallback
- Configuração flexível através da variável `DOCKER_MCP_SECRETS_FILE`

### 3. Configuração do RegistryPath
- **Correção do problema "no server is enabled"**: Adicionado `RegistryPath: []string{"registry.yaml"}` na configuração padrão
- Permite que o gateway carregue corretamente os servidores MCP habilitados
- Essencial para o funcionamento adequado do modo nativo

### 4. Desabilitamento de Funcionalidades Específicas
- Desabilitamento do monitor OAuth em modo nativo
- Remoção da dependência de autenticação específica do Docker Desktop
- Configurações otimizadas para ambiente nativo

### 5. Compatibilidade Mantida
- Funcionalidade completa preservada para usuários do Docker Desktop
- Transparência na operação entre modos
- Migração simples entre ambientes

## Pré-requisitos

### Sistema Operacional
- Linux (qualquer distribuição moderna)
- Docker Engine instalado e configurado
- Go 1.19+ para compilação (se necessário)

### Software Necessário
- Docker Engine (versão 20.10+ recomendada)
- Docker CLI
- Git (para clonar o repositório)
- Make (para build automatizado)

### Permissões
- Acesso ao socket do Docker (`/var/run/docker.sock`)
- Permissões de execução para binários gerados

## Início Rápido

### 1. Configurar Arquivo de Segredos (Opcional)
```bash
echo "SEU_SEGREDO=valor_secreto" > .env
export DOCKER_MCP_SECRETS_FILE=./.env
```

### 2. Executar o Gateway
```bash
# O sistema detectará automaticamente o ambiente Linux nativo
docker mcp gateway run --dry-run
```

### 3. Forçar Modo Específico (Opcional)
```bash
# Para forçar modo nativo manualmente
export DOCKER_MCP_NATIVE_MODE=1
docker mcp gateway run --dry-run

# Para forçar modo Docker Desktop manualmente
export DOCKER_MCP_NATIVE_MODE=0
docker mcp gateway run --dry-run
```

## Estrutura da Documentação

Esta documentação está organizada nos seguintes documentos:

- **[Guia de Instalação](guia-instalacao.md)**: Passos detalhados para compilar e configurar o ambiente
- **[Guia de Configuração](guia-configuracao.md)**: Explicação detalhada das variáveis de ambiente e opções avançadas
- **[Soluções de Problemas](solucoes-problemas.md)**: Problemas comuns, FAQ e dicas de troubleshooting
- **[Referência Técnica](referencia-tecnica.md)**: Detalhes técnicos das modificações e arquitetura da solução

## Benefícios da Solução

### Para Usuários
- **Compatibilidade**: Funciona com Docker Engine nativo do Linux
- **Simplicidade**: Configuração mínima necessária
- **Performance**: Melhor performance sem camadas adicionais
- **Flexibilidade**: Suporte a diferentes métodos de gerenciamento de segredos

### Para Desenvolvedores
- **Código Limpo**: Modificações mínimas e focadas
- **Manutenibilidade**: Código bem documentado e estruturado
- **Extensibilidade**: Arquitetura que permite futuras expansões
- **Testabilidade**: Facilita testes em diferentes ambientes

## Exemplos de Uso

### Básico (Detecção Automática)
```bash
# O sistema detecta automaticamente o ambiente Linux nativo
docker mcp gateway run
```

### Forçando Modo Nativo
```bash
# Forçar modo nativo manualmente
export DOCKER_MCP_NATIVE_MODE=1
docker mcp gateway run
```

### Com RegistryPath Personalizado
```bash
# Criar arquivo registry.yaml com servidores habilitados
cat > registry.yaml << 'EOF'
servers:
  filesystem:
    enabled: true
  docker:
    enabled: true
EOF

# Executar gateway com registry personalizado (detecção automática)
docker mcp gateway run --registry ./registry.yaml
```

### Com Segredos Personalizados
```bash
# Configurar arquivo de segredos
export DOCKER_MCP_SECRETS_FILE=./meus-segredos.env

# Executar gateway (detecção automática)
docker mcp gateway run --secrets file:./meus-segredos.env
```

### Exemplo Prático: Carregando Servidores
```bash
# 1. Criar arquivo de configuração de servidores
cat > registry.yaml << 'EOF'
servers:
  filesystem:
    enabled: true
    image: docker/mcp-server-filesystem:latest
  docker:
    enabled: true
    image: docker/mcp-server-docker:latest
  git:
    enabled: false
    image: docker/mcp-server-git:latest
EOF

# 2. Executar gateway em modo verbose para ver servidores sendo carregados
docker mcp gateway run --registry ./registry.yaml --verbose

# Saída esperada:
# - Detecção automática: Linux nativo detectado, ativando modo nativo
# ✓ Loading registry from: ./registry.yaml
# ✓ Server filesystem: enabled
# ✓ Server docker: enabled
# ✓ Server git: disabled
# ✓ Starting MCP Gateway with 2 enabled servers
```

### Modo SSE com Autenticação
```bash
# Executar em modo SSE (detecção automática)
docker mcp gateway run --transport sse --port 8080
```

## Suporte e Contribuições

Esta solução foi desenvolvida pela comunidade e está disponível para uso e contribuição. Para报告 problemas ou sugerir melhorias:

1. Verifique a documentação de [Soluções de Problemas](solucoes-problemas.md)
2. Consulte a [Referência Técnica](referencia-tecnica.md) para detalhes de implementação
3. Abra uma issue no repositório original com detalhes do problema

## Licença

Esta solução mantém a mesma licença do projeto Docker MCP Gateway original.