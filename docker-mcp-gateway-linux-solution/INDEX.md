# Índice da Documentação - Docker MCP Gateway Modo Nativo

## Visão Rápida

Esta documentação descreve a solução implementada para permitir que o **Docker MCP Gateway** funcione com **Docker Engine nativo do Linux**, sem depender do Docker Desktop.

## Documentos Disponíveis

### 📖 [README.md](README.md)
**Visão geral da solução**
- Problema original e solução implementada
- Pré-requisitos e início rápido
- Benefícios e exemplos de uso
- Estrutura da documentação

### 🛠️ [Guia de Instalação](guia-instalacao.md)
**Passos detalhados para compilar e configurar**
- Pré-requisitos do sistema
- Obtenção e compilação do código fonte
- Configuração do ambiente
- Verificação da instalação
- Automação e configuração avançada

### ⚙️ [Guia de Configuração](guia-configuracao.md)
**Detalhes de configuração e opções avançadas**
- Variáveis de ambiente principais
- Configuração de segredos
- Opções de linha de comando
- Arquivos de configuração
- Configuração avançada (autenticação, interceptors, etc.)

### 🔧 [Soluções de Problemas](solucoes-problemas.md)
**Troubleshooting e FAQ**
- Problemas comuns e soluções
- Perguntas frequentes
- Dicas de troubleshooting
- Scripts de diagnóstico
- Recuperação de desastres

### 🔍 [Detecção Automática](deteccao-automatica.md)
**Como funciona a detecção automática de ambiente**
- Fluxo de decisão da detecção
- Regras e critérios utilizados
- Cenários de uso e comportamentos esperados
- Solução de problemas de detecção
- Testes e validação

### 📚 [Referência Técnica](referencia-tecnica.md)
**Detalhes técnicos das modificações**
- Arquitetura da solução
- Modificações implementadas arquivo por arquivo
- Fluxo de dados e decisões
- Compatibilidade e performance
- Segurança e manutenção

### 💡 [Exemplos Práticos](exemplos.md)
**Scripts e cenários de uso real**
- Configuração básica de desenvolvimento
- Ambiente de produção com Systemd
- Múltiplos ambientes com working sets
- Configuração com interceptors
- Monitoramento e métricas

### 📝 [Atualizações](atualizacoes.md)
**Histórico de mudanças e melhorias**
- Novas funcionalidades implementadas
- Correções de bugs e problemas
- Melhorias de performance e usabilidade
- Roadmap e planejamento futuro

## Mapa de Navegação

### Para Usuários Iniciantes
1. Comece com o [README.md](README.md) para entender a solução
2. Siga o [Guia de Instalação](guia-instalacao.md) passo a passo
3. Consulte o [Guia de Configuração](guia-configuracao.md) para personalizar seu ambiente
4. Leia sobre [Detecção Automática](deteccao-automatica.md) para entender como o sistema funciona

### Para Usuários Avançados
1. Revise a [Referência Técnica](referencia-tecnica.md) para detalhes de implementação
2. Consulte o [Guia de Configuração](guia-configuracao.md) para opções avançadas
3. Use [Soluções de Problemas](solucoes-problemas.md) para troubleshooting específico
4. Entenda a [Detecção Automática](deteccao-automatica.md) para casos especiais

### Para Desenvolvedores
1. Estude a [Referência Técnica](referencia-tecnica.md) para entender as modificações
2. Revise os exemplos de código e testes
3. Consulte o [Guia de Configuração](guia-configuracao.md) para entender as opções
4. Analise a [Detecção Automática](deteccao-automatica.md) para implementação

## Resumo Rápido dos Comandos

### Configuração Básica (Detecção Automática)
```bash
# O sistema detecta automaticamente o ambiente Linux nativo
# Não é necessário configurar variáveis

# Configurar arquivo de segredos (opcional)
echo "API_KEY=seu_valor" > ~/.docker/mcp/secrets.env
export DOCKER_MCP_SECRETS_FILE=$HOME/.docker/mcp/secrets.env

# Executar gateway
docker mcp gateway run
```

### Forçar Modo Específico (Opcional)
```bash
# Forçar modo nativo
export DOCKER_MCP_NATIVE_MODE=1
docker mcp gateway run

# Forçar modo Docker Desktop
export DOCKER_MCP_NATIVE_MODE=0
docker mcp gateway run
```

### Modo SSE
```bash
# Executar em modo SSE
docker mcp gateway run --transport sse --port 8080
```

### Verificação
```bash
# Testar configuração
docker mcp gateway run --dry-run

# Verificar versão
docker mcp version
```

## Variáveis de Ambiente Principais

| Variável | Valor | Descrição |
|----------|-------|-----------|
| `DOCKER_MCP_NATIVE_MODE` | `1`/`0`/não definido | Força modo nativo/Docker Desktop (opcional, detecção automática padrão) |
| `DOCKER_MCP_SECRETS_FILE` | `/path/to/.env` | Caminho para arquivo de segredos |
| `DOCKER_MCP_IN_CONTAINER` | `1` | Detectado automaticamente quando em contêiner |

## Arquivos de Configuração Importantes

| Arquivo | Descrição |
|---------|-----------|
| `~/.docker/mcp/secrets.env` | Segredos do MCP Gateway |
| `~/.docker/mcp/catalogs/` | Catálogos de servidores MCP |
| `~/.docker/mcp/registry.yaml` | Configuração de registry |
| `~/.docker/mcp/config.yaml` | Configuração geral |

## Problemas Comuns e Soluções Rápidas

| Problema | Solução |
|----------|---------|
| Detecção automática não funciona | Use `export DOCKER_MCP_NATIVE_MODE=1` para forçar modo nativo |
| Segredos não são lidos | Configure `DOCKER_MCP_SECRETS_FILE` |
| Permissão negada no Docker socket | Adicione usuário ao grupo docker |
| Erro de autenticação em modo SSE | Verifique se está em modo nativo (detecção automática ou `DOCKER_MCP_NATIVE_MODE=1`) |

## Recursos Adicionais

### Scripts Úteis
- [Script de diagnóstico](solucoes-problemas.md#diagnóstico-rápido)
- [Script de validação de configuração](guia-instalacao.md#verificação-final)
- [Script de backup](solucoes-problemas.md#backup-de-configuração)

### Exemplos de Configuração
- [Ambiente de desenvolvimento](guia-configuracao.md#exemplos-de-configuração-completa)
- [Ambiente de produção](guia-configuracao.md#exemplos-de-configuração-completa)
- [Systemd service](guia-instalacao.md#configurar-systemd-service-opcional)

## Contribuição e Suporte

### Reportar Problemas
1. Verifique [Soluções de Problemas](solucoes-problemas.md)
2. Execute o [script de diagnóstico](solucoes-problemas.md#script-de-coleta-de-informações)
3. Abra uma issue no repositório original

### Contribuir com a Documentação
1. Fork do repositório
2. Faça as melhorias na documentação
3. Abra um pull request

## Histórico de Versões

### v1.1.0 - Detecção Automática (07/11/2025)
- Implementação de detecção automática de ambiente Linux nativo
- Eliminação da necessidade de configurar `DOCKER_MCP_NATIVE_MODE`
- Centralização da lógica de detecção em `pkg/desktop/detection.go`
- Manutenção de override manual para casos especiais
- Simplificação do processo de instalação e configuração

### v1.0.0 - Modo Nativo
- Implementação inicial do modo nativo
- Suporte a Docker Engine nativo do Linux
- Gerenciamento alternativo de segredos
- Desabilitamento de OAuth em modo nativo

## Licença

Esta documentação segue a mesma licença do projeto Docker MCP Gateway original.

---

**Dica**: Use a função de busca do seu editor para encontrar rapidamente informações específicas nesta documentação.