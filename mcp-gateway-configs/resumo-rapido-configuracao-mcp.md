# Resumo Rápido: Configuração MCP para Agentes de IA

## 📍 Localizações Principais

### Claude Desktop
- **Windows**: `%APPDATA%\Claude\claude_desktop_config.json`
- **macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
- **Linux**: `~/.config/claude/claude_desktop_config.json`

### VS Code
- **Windows**: `%APPDATA%\Code\User\settings.json`
- **macOS**: `~/Library/Application Support/Code/User/settings.json`
- **Linux**: `~/.config/Code/User/settings.json`

### Kilo Code
- **Windows**: `%APPDATA%\KiloCode\settings.json`
- **macOS**: `~/Library/Application Support/KiloCode/settings.json`
- **Linux**: `~/.config/kilocode/settings.json`

---

## 🚀 Configuração Rápida

### 1. Instalar Dependências
```bash
# Instalar Docker e Docker MCP Gateway
curl -fsSL https://get.docker.com | sh
docker plugin install docker/mcp-gateway

# Instalar ferramentas úteis
# macOS: brew install jq
# Linux: sudo apt install jq
```

### 2. Criar Arquivo de Configuração

#### Claude Desktop
```json
{
  "mcpServers": {
    "docker-mcp-gateway": {
      "command": "docker",
      "args": [
        "mcp",
        "gateway",
        "run",
        "--servers=duckduckgo,wikipedia-mcp,github,filesystem,postgres"
      ]
    }
  }
}
```

#### VS Code
```json
{
  "mcp.servers": {
    "docker-mcp-gateway": {
      "command": "docker",
      "args": [
        "mcp",
        "gateway",
        "run",
        "--servers=github,filesystem,kubernetes,postgres,redis"
      ],
      "env": {
        "DOCKER_MCP_LOG_LEVEL": "info",
        "DOCKER_MCP_NATIVE_MODE": "1"
      }
    }
  }
}
```

### 3. Configurar Secrets
```bash
# Criar diretório de secrets
mkdir -p ~/.docker/mcp/secrets
chmod 700 ~/.docker/mcp/secrets

# Criar arquivo de secrets
cat > ~/.docker/mcp/secrets/.env << 'EOF'
GITHUB_PERSONAL_ACCESS_TOKEN=seu_token_aqui
POSTGRES_CONNECTION_STRING=postgresql://user:pass@host:port/db
REDIS_CONNECTION_STRING=redis://user:pass@host:port/db
EOF

chmod 600 ~/.docker/mcp/secrets/.env
```

### 4. Testar Configuração
```bash
# Verificar configuração
docker mcp gateway run --verbose --dry-run

# Listar ferramentas disponíveis
docker mcp tools ls --verbose

# Testar ferramenta específica
docker mcp tools call --gateway-arg="--servers=duckduckgo" search query=Docker
```

---

## 🔧 Scripts Automatizados

### Configuração Automática
```bash
# Baixar e executar script de configuração
curl -O https://raw.githubusercontent.com/docker/mcp-gateway/main/scripts/auto-config-mcp.sh
chmod +x auto-config-mcp.sh
./auto-config-mcp.sh
```

### Backup de Configurações
```bash
# Criar backup
./backup-restore-mcp.sh backup

# Restaurar backup
./backup-restore-mcp.sh restore /path/to/backup
```

### Validação de Configurações
```bash
# Validar todas as configurações
./validate-mcp-configs.sh
```

---

## 🛠️ Comandos Úteis

### Docker MCP Gateway
```bash
# Verificar versão
docker mcp version

# Listar servidores disponíveis
docker mcp server ls

# Inspecionar servidor específico
docker mcp server inspect duckduckgo

# Listar ferramentas
docker mcp tools ls

# Executar com perfil de otimização
docker mcp gateway run --optimization=performance --servers=github,filesystem
```

### Troubleshooting
```bash
# Verificar permissões do Docker
sudo usermod -aG docker $USER

# Verificar portas em uso
netstat -tlnp | grep :8080

# Verificar logs de containers
docker logs $(docker ps -q --filter "name=mcp")
```

---

## 📋 Servidores Principais

| Categoria | Servidores | Recursos |
|-----------|------------|----------|
| **Pesquisa** | duckduckgo, wikipedia-mcp, brave | 25 |
| **Desenvolvimento** | github, filesystem, kubernetes, context7 | 35 |
| **Bancos de Dados** | postgres, redis, clickhouse, sqlite-mcp-server | 35 |
| **Comunicação** | slack, mcp-discord | 12 |
| **Produtividade** | notion | 7 |

**Total**: 13 servidores, 105 recursos

---

## 🎯 Casos de Uso Rápidos

### Desenvolvedor Full Stack
```bash
docker mcp gateway run --servers=github,filesystem,postgres,redis,duckduckgo,wikipedia-mcp
```

### Analista de Dados
```bash
docker mcp gateway run --servers=clickhouse,postgres,sqlite-mcp-server,wikipedia-mcp,notion
```

### DevOps Engineer
```bash
docker mcp gateway run --servers=kubernetes,slack,mcp-discord,filesystem,github
```

### Pesquisador
```bash
docker mcp gateway run --servers=duckduckgo,wikipedia-mcp,github,notion,brave
```

---

## 🔗 Links Úteis

- [Guia Completo](guia-completo-configuracao-mcp-agentes-ia.md)
- [Scripts Automatizados](scripts-automatizados-configuracao-mcp.md)
- [Documentação Oficial](https://docs.docker.com/mcp-gateway/)
- [Catálogo Docker MCP](https://hub.docker.com/mcp)
- [Repositório GitHub](https://github.com/docker/mcp-gateway)

---

## ⚠️ Importante

1. **Reinicie** seus aplicativos após alterar configurações
2. **Verifique** as permissões do arquivo de secrets (600)
3. **Faça backup** antes de modificar configurações existentes
4. **Teste** com `--dry-run` antes de aplicar configurações

---

**Última atualização**: 7 de novembro de 2025  
**Versão**: 1.0.0  
**Compatível com**: Docker MCP Gateway 2.0.0+