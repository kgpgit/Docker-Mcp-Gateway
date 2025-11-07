# Soluções de Problemas - Docker MCP Gateway Modo Nativo

Este documento aborda problemas comuns, soluções e dicas de troubleshooting para o Docker MCP Gateway em modo nativo.

## Problemas Comuns

### 1. Erro de Conexão com Docker Desktop

**Sintoma**:
```
Error: Docker Desktop is not running or not accessible
```

**Causa**: O gateway está tentando se conectar ao Docker Desktop em vez do Docker Engine nativo.

**Solução**:
```bash
# Verificar se o modo nativo está ativo
echo $DOCKER_MCP_NATIVE_MODE

# Ativar modo nativo
export DOCKER_MCP_NATIVE_MODE=1

# Verificar se o Docker Engine está rodando
docker info
```

**Verificação**:
```bash
# Testar configuração
docker mcp gateway run --dry-run
```

### 2. Segredos Não São Lidos

**Sintoma**:
```
Warning: Could not read secret SECRET_NAME
```

**Causa**: O arquivo de segredos não está configurado corretamente ou não existe.

**Solução**:
```bash
# Verificar variável de ambiente
echo $DOCKER_MCP_SECRETS_FILE

# Criar arquivo de segredos
cat > ~/.docker/mcp/secrets.env << 'EOF'
API_KEY=seu_valor_aqui
DATABASE_URL=postgresql://user:pass@localhost/db
EOF

# Configurar permissões corretas
chmod 600 ~/.docker/mcp/secrets.env

# Configurar variável de ambiente
export DOCKER_MCP_SECRETS_FILE=$HOME/.docker/mcp/secrets.env
```

**Verificação**:
```bash
# Testar leitura de segredos
docker mcp gateway run --dry-run --secrets file:$HOME/.docker/mcp/secrets.env
```

### 3. Permissão Negada ao Acessar Docker Socket

**Sintoma**:
```
Error: permission denied while trying to connect to the Docker daemon socket
```

**Causa**: O usuário não tem permissão para acessar o socket do Docker.

**Solução**:
```bash
# Adicionar usuário ao grupo docker
sudo usermod -aG docker $USER

# Aplicar alterações (requer logout/login)
newgrp docker

# Verificar permissões
docker ps
```

**Alternativa (temporária)**:
```bash
# Usar sudo (não recomendado para produção)
sudo docker mcp gateway run
```

### 4. Erro de Autenticação em Modo SSE/HTTP

**Sintoma**:
```
Error: Authentication required for SSE/HTTP transport
```

**Causa**: O modo nativo deveria desabilitar autenticação, mas a configuração não foi aplicada.

**Solução**:
```bash
# Verificar modo nativo
export DOCKER_MCP_NATIVE_MODE=1

# Executar novamente
docker mcp gateway run --transport sse --port 8080
```

**Verificação**:
```bash
# Deve mostrar "Authentication disabled" nos logs
docker mcp gateway run --transport sse --port 8080 --verbose
```

### 5. Servidores MCP Não Iniciam

**Sintoma**:
```
Error: Failed to start MCP server: container not found
```

**Causa**: Imagens do Docker não estão disponíveis ou não podem ser baixadas.

**Solução**:
```bash
# Verificar conexão com registry
docker pull docker/mcp-server:latest

# Verificar imagens disponíveis
docker images | grep mcp

# Forçar pull de imagens
docker mcp gateway run --pull-always
```

### 6. Problema "No Server is Enabled"

**Sintoma**:
```
Error: no server is enabled
Warning: No MCP servers found in registry
```

**Causa Raiz**: O gateway não consegue localizar o arquivo de registro (`registry.yaml`) que contém a configuração dos servidores MCP habilitados. Em modo nativo, este arquivo não era configurado por padrão na configuração inicial.

**Solução Implementada**:
Foi adicionada a linha `RegistryPath: []string{"registry.yaml"}` na configuração padrão do modo nativo no arquivo [`cmd/docker-mcp/commands/gateway.go`](cmd/docker-mcp/commands/gateway.go:37).

**Soluções Imediatas**:

1. **Verificar se o arquivo registry.yaml existe**:
```bash
# Verificar arquivo padrão
ls -la ~/.docker/mcp/registry.yaml

# Criar arquivo se não existir
mkdir -p ~/.docker/mcp
cat > ~/.docker/mcp/registry.yaml << 'EOF'
servers:
  filesystem:
    enabled: true
  docker:
    enabled: true
EOF
```

2. **Especificar caminho do registry explicitamente**:
```bash
# Usar registry personalizado
docker mcp gateway run --registry ./meu-registry.yaml

# Verificar conteúdo do registry
cat ./meu-registry.yaml
```

3. **Habilitar todos os servidores disponíveis**:
```bash
# Habilitar todos os servidores do catálogo
docker mcp gateway run --enable-all-servers

# Verificar quais servidores estão disponíveis
docker mcp catalog list
```

**Comandos para Verificação**:
```bash
# 1. Verificar configuração atual
docker mcp gateway run --dry-run --verbose

# 2. Verificar se o registry está sendo lido
docker mcp gateway run --dry-run --verbose 2>&1 | grep -i registry

# 3. Listar servidores disponíveis no catálogo
docker mcp catalog list

# 4. Verificar conteúdo do registry atual
cat ~/.docker/mcp/registry.yaml

# 5. Testar com registry de exemplo
cat > test-registry.yaml << 'EOF'
servers:
  filesystem:
    enabled: true
    image: docker/mcp-server-filesystem:latest
  docker:
    enabled: true
    image: docker/mcp-server-docker:latest
EOF

docker mcp gateway run --registry ./test-registry.yaml --dry-run
```

**Verificação da Solução**:
```bash
# Aplicar a correção e testar
export DOCKER_MCP_NATIVE_MODE=1

# Criar registry de teste
echo 'servers:
  filesystem:
    enabled: true' > registry.yaml

# Executar gateway
docker mcp gateway run --registry ./registry.yaml --verbose

# Saída esperada:
# ✓ Loading registry from: ./registry.yaml
# ✓ Server filesystem: enabled
# ✓ Starting MCP Gateway with 1 enabled servers
```

### 7. Erro de Rede

**Sintoma**:
```
Error: Network connection failed
```

**Causa**: Configuração de rede incorreta ou conflito de portas.

**Solução**:
```bash
# Verificar portas em uso
netstat -tlnp | grep :8080

# Usar porta diferente
docker mcp gateway run --transport sse --port 8081

# Verificar redes Docker
docker network ls
```

## FAQ

### Perguntas Gerais

**Q: O modo nativo funciona em outras distribuições Linux além de Ubuntu?**
R: Sim, o modo nativo funciona em qualquer distribuição Linux com Docker Engine instalado, incluindo Debian, CentOS, Fedora, Arch Linux, etc.

**Q: Posso usar o modo nativo e o Docker Desktop simultaneamente?**
R: Não é recomendado. O modo nativo é projetado para substituir completamente o Docker Desktop.

**Q: Como migro do Docker Desktop para o modo nativo?**
R: Simplesmente defina `DOCKER_MCP_NATIVE_MODE=1` e configure o arquivo de segredos. Não é necessário reinstalar ou reconfigurar outros componentes.

**Q: O modo nativo afeta o desempenho?**
R: Na verdade, o modo nativo geralmente oferece melhor desempenho, pois elimina a camada adicional do Docker Desktop.

### Configuração

**Q: Onde devo armazenar meus segredos?**
R: Use `~/.docker/mcp/secrets.env` com permissões 600. Para ambientes de produção, considere usar um gerenciador de segredos como HashiCorp Vault.

**Q: Como configuro múltiplos ambientes (dev, staging, prod)?**
R: Use working sets ou arquivos de configuração separados para cada ambiente:
```bash
docker mcp gateway run --working-set development
docker mcp gateway run --working-set production
```

**Q: Posso usar catálogos personalizados no modo nativo?**
R: Sim, os catálogos funcionam da mesma forma no modo nativo:
```bash
docker mcp gateway run --catalog ~/.docker/mcp/catalogs/custom.yaml
```

### Troubleshooting

**Q: Como habilito logs detalhados para debugging?**
R: Use as opções `--verbose` e `--log-calls`:
```bash
docker mcp gateway run --verbose --log-calls
```

**Q: Como verifico se o modo nativo está realmente ativo?**
R: Verifique as variáveis de ambiente e os logs:
```bash
echo $DOCKER_MCP_NATIVE_MODE
docker mcp gateway run --dry-run --verbose | grep -i native
```

**Q: O que fazer se o gateway não iniciar?**
R: Siga estes passos:
1. Verifique se o Docker Engine está rodando: `docker info`
2. Verifique as variáveis de ambiente: `env | grep DOCKER_MCP`
3. Tente modo dry-run: `docker mcp gateway run --dry-run`
4. Verifique os logs com `--verbose`

## Dicas de Troubleshooting

### Diagnóstico Rápido

```bash
#!/bin/bash
# Script de diagnóstico rápido

echo "=== Diagnóstico do Docker MCP Gateway ==="

# 1. Verificar Docker
echo "1. Status do Docker:"
docker info > /dev/null 2>&1
if [ $? -eq 0 ]; then
    echo "   ✓ Docker Engine está rodando"
else
    echo "   ✗ Docker Engine não está acessível"
fi

# 2. Verificar variáveis de ambiente
echo "2. Variáveis de ambiente:"
echo "   DOCKER_MCP_NATIVE_MODE: ${DOCKER_MCP_NATIVE_MODE:-não definido}"
echo "   DOCKER_MCP_SECRETS_FILE: ${DOCKER_MCP_SECRETS_FILE:-não definido}"

# 3. Verificar arquivo de segredos
if [ -n "$DOCKER_MCP_SECRETS_FILE" ]; then
    if [ -f "$DOCKER_MCP_SECRETS_FILE" ]; then
        echo "   ✓ Arquivo de segredos existe"
        echo "   Permissões: $(ls -l "$DOCKER_MCP_SECRETS_FILE" | cut -d' ' -f1)"
    else
        echo "   ✗ Arquivo de segredos não encontrado"
    fi
fi

# 4. Testar configuração
echo "4. Teste de configuração:"
if docker mcp gateway run --dry-run > /dev/null 2>&1; then
    echo "   ✓ Configuração válida"
else
    echo "   ✗ Erro na configuração"
    echo "   Execute: docker mcp gateway run --dry-run --verbose"
fi

echo "=== Fim do diagnóstico ==="
```

### Logs Detalhados

```bash
# Habilitar logs máximos
export DOCKER_MCP_TELEMETRY_DEBUG=1
export DOCKER_MCP_SHOW_HIDDEN=1

# Executar com verbose
docker mcp gateway run --verbose --log-calls --log debug.log

# Monitorar logs em tempo real
tail -f debug.log
```

### Verificação de Rede

```bash
# Verificar conectividade
docker network ls
docker network inspect bridge

# Testar resolução DNS
docker run --rm alpine nslookup google.com

# Verificar portas
netstat -tlnp | grep :8080
```

### Depuração de Segredos

```bash
# Verificar conteúdo do arquivo de segredos
cat $DOCKER_MCP_SECRETS_FILE

# Testar leitura de segredos
docker mcp gateway run --dry-run --secrets file:$DOCKER_MCP_SECRETS_FILE --verbose

# Verificar variáveis de ambiente
env | grep -E "(API|TOKEN|KEY|SECRET)"
```

## Problemas Específicos

### Problemas com Systemd

**Sintoma**: O serviço systemd não inicia corretamente.

**Solução**:
```bash
# Verificar status do serviço
sudo systemctl status docker-mcp

# Verificar logs
sudo journalctl -u docker-mcp -f

# Recarregar configuração
sudo systemctl daemon-reload
sudo systemctl restart docker-mcp
```

### Problemas com Firewalls

**Sintoma**: Conexões recusadas em modo SSE/HTTP.

**Solução**:
```bash
# Verificar firewall (ufw)
sudo ufw status
sudo ufw allow 8080/tcp

# Verificar firewall (firewalld)
sudo firewall-cmd --list-all
sudo firewall-cmd --add-port=8080/tcp --permanent
sudo firewall-cmd --reload
```

### Problemas com SELinux

**Sintoma**: Erros de permissão mesmo com configuração correta.

**Solução**:
```bash
# Verificar status SELinux
sestatus

# Verificar contextos
ls -Z ~/.docker/mcp/

# Ajustar contextos se necessário
sudo semanage fcontext -a -t container_file_t "/home/user/.docker/mcp(/.*)?"
sudo restorecon -R ~/.docker/mcp/
```

## Recuperação de Desastres

### Backup de Configuração

```bash
#!/bin/bash
# Script de backup

BACKUP_DIR="$HOME/docker-mcp-backup-$(date +%Y%m%d)"
mkdir -p "$BACKUP_DIR"

# Backup de configurações
cp -r ~/.docker/mcp "$BACKUP_DIR/"

# Backup de variáveis de ambiente
env | grep DOCKER_MCP > "$BACKUP_DIR/env-vars"

# Backup de scripts de inicialização
cp ~/.local/bin/docker-mcp-* "$BACKUP_DIR/" 2>/dev/null

echo "Backup criado em: $BACKUP_DIR"
```

### Restauração

```bash
#!/bin/bash
# Script de restauração

BACKUP_DIR="$1"

if [ -z "$BACKUP_DIR" ]; then
    echo "Uso: $0 <diretorio_do_backup>"
    exit 1
fi

# Restaurar configurações
cp -r "$BACKUP_DIR/mcp" ~/.docker/

# Restaurar variáveis de ambiente
source "$BACKUP_DIR/env-vars"

# Restaurar scripts
cp "$BACKUP_DIR"/docker-mcp-* ~/.local/bin/ 2>/dev/null

echo "Restauração concluída"
```

## Performance e Otimização

### Monitoramento de Recursos

```bash
# Monitorar uso de CPU e memória
docker stats

# Monitorar logs de performance
docker mcp gateway run --verbose 2>&1 | grep -E "(performance|memory|cpu)"

# Verificar uso de disco
du -sh ~/.docker/mcp/
```

### Otimizações

```bash
# Limpar caches não utilizados
docker system prune -f

# Usar servidores de longa duração para melhor performance
docker mcp gateway run --long-lived

# Configurar recursos adequados
docker mcp gateway run --cpus 2 --memory 4Gb
```

## Contato e Suporte

### Recursos de Ajuda

1. **Documentação Oficial**: Consulte os outros documentos desta série
2. **GitHub Issues**: Reporte problemas no repositório original
3. **Comunidade**: Participe das discussões no Docker Community Slack

### Informações para Reportar Problemas

Ao reportar um problema, inclua:

1. **Versão do Docker MCP Gateway**:
   ```bash
   docker mcp version
   ```

2. **Versão do Docker Engine**:
   ```bash
   docker version
   ```

3. **Variáveis de Ambiente**:
   ```bash
   env | grep DOCKER_MCP
   ```

4. **Logs Detalhados**:
   ```bash
   docker mcp gateway run --dry-run --verbose
   ```

5. **Sistema Operacional**:
   ```bash
   uname -a
   cat /etc/os-release
   ```

### Script de Coleta de Informações

```bash
#!/bin/bash
# Coletar informações para suporte

INFO_FILE="docker-mcp-support-info.txt"

{
    echo "=== Informações do Sistema ==="
    uname -a
    cat /etc/os-release
    
    echo -e "\n=== Versões ==="
    docker version
    docker mcp version
    
    echo -e "\n=== Variáveis de Ambiente ==="
    env | grep DOCKER_MCP
    
    echo -e "\n=== Teste de Configuração ==="
    docker mcp gateway run --dry-run --verbose
    
    echo -e "\n=== Status do Docker ==="
    docker info
} > "$INFO_FILE"

echo "Informações coletadas em: $INFO_FILE"
```

## Próximos Passos

Se você não encontrou uma solução para seu problema:

1. Consulte a [Referência Técnica](referencia-tecnica.md) para detalhes de implementação
2. Execute o script de diagnóstico acima
3. Abra uma issue no repositório com as informações coletadas
4. Participe da comunidade Docker para obter ajuda adicional