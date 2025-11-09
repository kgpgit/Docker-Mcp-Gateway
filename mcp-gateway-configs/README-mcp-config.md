# Docker MCP Gateway - Arquivo de Configuração Completo

Este documento explica como usar o arquivo `mcp-gateway-config.json` para configurar o Docker MCP Gateway para diferentes agentes de IA.

## Visão Geral

O arquivo `mcp-gateway-config.json` contém configurações pré-definidas para facilitar a integração do Docker MCP Gateway com diversos clientes MCP, incluindo:

* Claude Desktop
* VS Code
* Cursor
* Continue
* Zed
* LM Studio
* Outros agentes de IA

## Configuração Rápida

### 1. Claude Desktop

Adicione a seguinte configuração ao seu arquivo `claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "docker-mcp-gateway": {
      "command": "docker",
      "args": [
        "mcp",
        "gateway",
        "run",
        "--servers=duckduckgo,github,wikipedia,filesystem"
      ]
    }
  }
}
```

### 2. VS Code

Adicione a seguinte configuração ao seu `settings.json`:

```json
{
  "mcp.servers": {
    "docker-mcp-gateway": {
      "command": "docker",
      "args": [
        "mcp",
        "gateway",
        "run",
        "--servers=duckduckgo,github,wikipedia,filesystem"
      ]
    }
  }
}
```

### 3. Outros Agentes

Use a configuração apropriada da seção `other_agents` no arquivo `mcp-gateway-config.json`.

## Servidores Pré-configurados

O arquivo de configuração inclui os seguintes servidores MCP pré-configurados:

### DuckDuckGo Search

* **Nome**: `duckduckgo`
* **Descrição**: Pesquisas web via DuckDuckGo
* **Ferramentas**: `search`, `fetch_content`

### GitHub Integration

* **Nome**: `github`
* **Descrição**: Integração com GitHub
* **Ferramentas**: `create_repository`, `list_repositories`, `get_repository`, `create_issue`, `list_issues`
* **Requer**: `GITHUB_PERSONAL_ACCESS_TOKEN`

### Wikipedia

* **Nome**: `wikipedia`
* **Descrição**: Consultas enciclopédicas
* **Ferramentas**: `search`, `get_page`, `get_summary`

### Filesystem Operations

* **Nome**: `filesystem`
* **Descrição**: Operações de sistema de arquivos locais
* **Ferramentas**: `read_file`, `write_file`, `list_directory`, `create_directory`, `delete_file`, `move_file`

### PostgreSQL Database

* **Nome**: `postgres`
* **Descrição**: Operações com banco PostgreSQL
* **Ferramentas**: `execute_query`, `list_tables`, `describe_table`, `insert_data`, `update_data`, `delete_data`
* **Requer**: `POSTGRES_CONNECTION_STRING`

### Neo4j Graph Database

* **Nome**: `neo4j`
* **Descrição**: Operações com banco Neo4j
* **Ferramentas**: `execute_cypher`, `read_graph`, `write_graph`, `list_nodes`, `list_relationships`
* **Requer**: `NEO4J_URI`, `NEO4J_USERNAME`, `NEO4J_PASSWORD`

## Configuração de Variáveis de Ambiente

### Variáveis Obrigatórias

```bash
# Socket do Docker (geralmente já configurado)
export DOCKER_HOST=unix:///var/run/docker.sock
```

### Variáveis Opcionais

```bash
# Nível de log
export DOCKER_MCP_LOG_LEVEL=info

# Modo nativo Linux (sem Docker Desktop)
export DOCKER_MCP_NATIVE_MODE=1

# Caminhos personalizados
export DOCKER_MCP_CONFIG_PATH=~/.docker/mcp/config.yaml
export DOCKER_MCP_CATALOG_PATH=~/.docker/mcp/catalogs/docker-mcp.yaml
export DOCKER_MCP_REGISTRY_PATH=~/.docker/mcp/registry.yaml
```

### Secrets (Variáveis Sensíveis)

Crie um arquivo `.env` na raiz do seu projeto:

```bash
# GitHub Token
GITHUB_PERSONAL_ACCESS_TOKEN=seu_token_aqui

# PostgreSQL
POSTGRES_CONNECTION_STRING=postgresql://usuario:senha@host:port/database

# Neo4j
NEO4J_PASSWORD=sua_senha_aqui
```

## Exemplos de Uso

### Execução Básica (stdio)

```bash
docker mcp gateway run --servers=duckduckgo,github,wikipedia,filesystem
```

### Execução com Transporte Streaming

```bash
docker mcp gateway run --transport=streaming --port=8080 --servers=duckduckgo,github,wikipedia,filesystem
```

### Execução com Transporte SSE

```bash
docker mcp gateway run --transport=sse --port=8081 --servers=duckduckgo,github,wikipedia,filesystem
```

### Execução com Secrets

```bash
docker mcp gateway run --servers=github,postgres --secrets=docker-desktop:./.env
```

### Execução via Docker Compose

Crie um arquivo `docker-compose.yml`:

```yaml
version: '3.8'
services:
  gateway:
    image: docker/mcp-gateway
    command:
      - --servers=duckduckgo,github,wikipedia,filesystem
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
    ports:
      - "8080:8080"
```

Execute com:

```bash
docker-compose up
```

## Troubleshooting

### Verificar Configuração

```bash
# Modo de teste para verificar configuração
docker mcp gateway run --verbose --dry-run
```

### Listar Ferramentas

```bash
# Listar todas as ferramentas disponíveis
docker mcp tools ls --verbose
```

### Testar Ferramenta Específica

```bash
# Testar uma ferramenta específica
docker mcp tools call --gateway-arg="--servers=duckduckgo" --verbose search query=Docker
```

### Inspecionar Servidor

```bash
# Inspecionar um servidor específico
docker mcp server inspect duckduckgo
```

### Problemas Comuns


1. **Permissão negada ao acessar o socket do Docker**

   ```bash
   sudo usermod -aG docker $USER
   # Faça logout e login novamente
   ```
2. **Servidor MCP não encontrado**

   ```bash
   docker mcp gateway run --verbose --dry-run --servers=nome_do_servidor
   ```
3. **Conexão recusada**

   ```bash
   netstat -tlnp | grep :8080
   ```
4. **Falha de autenticação**

   ```bash
   docker mcp gateway run --secrets=docker-desktop:./.env
   ```

## Configuração Avançada

### Limites de Recursos

```bash
# Configurar CPUs e memória para servidores MCP
docker mcp gateway run --cpus=2 --memory=4Gb --servers=duckduckgo
```

### Opções de Segurança

```bash
# Bloquear acesso à rede
docker mcp gateway run --block-network --servers=duckduckgo

# Verificar assinaturas de imagens
docker mcp gateway run --verify-signatures --servers=duckduckgo
```

### Interceptores

```bash
# Adicionar interceptores
docker mcp gateway run --interceptor=before:exec:/bin/path --servers=duckduckgo
```

### Working Sets

```bash
# Habilitar feature de working sets
docker mcp feature enable working-sets

# Usar working set específico
docker mcp gateway run --working-set my-working-set
```

## Estrutura do Arquivo de Configuração

O arquivo `mcp-gateway-config.json` está organizado nas seguintes seções:


1. **configurations**: Configurações para diferentes clientes MCP
2. **preconfigured_servers**: Definições dos servidores MCP disponíveis
3. **environment_variables**: Variáveis de ambiente necessárias
4. **usage_examples**: Exemplos de uso do gateway
5. **troubleshooting**: Dicas para resolver problemas
6. **advanced_configuration**: Opções avançadas de configuração

## Contribuição

Para adicionar novos servidores ou configurações:


1. Adicione o servidor à seção `preconfigured_servers`
2. Inclua as variáveis de ambiente necessárias em `environment_variables`
3. Adicione exemplos de uso em `usage_examples`
4. Documente problemas comuns em `troubleshooting`

## Recursos Adicionais

* [Documentação oficial do Docker MCP Gateway](docs/mcp-gateway.md)
* [Guia de troubleshooting](docs/troubleshooting.md)
* [Exemplos de configuração](examples/README.md)
* [Fluxo de configuração](docs/config-flow.md)


