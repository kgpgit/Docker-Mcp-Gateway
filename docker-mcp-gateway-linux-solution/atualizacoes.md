# Atualizações - Docker MCP Gateway Modo Nativo

Este documento mantém um histórico de todas as atualizações e correções implementadas no Docker MCP Gateway para suporte ao modo nativo.

## 07 de Novembro de 2025

### Implementação de Detecção Automática de Ambiente Linux Nativo

**Descrição**: Foi implementado um sistema de detecção automática que identifica quando o Docker MCP Gateway está rodando em ambiente Linux nativo, eliminando a necessidade de configurar manualmente a variável `DOCKER_MCP_NATIVE_MODE`.

**Motivação**: A configuração manual era um ponto de fricção para usuários e podia levar a erros. A detecção automática simplifica o uso e melhora a experiência do usuário.

**Implementação**:
- Criado novo arquivo [`pkg/desktop/detection.go`](pkg/desktop/detection.go) com funções centralizadas de detecção
- Implementada função `IsNativeMode()` com lógica inteligente de detecção
- Mantida compatibilidade com override manual via variável de ambiente
- Atualizadas todas as verificações espalhadas pelo código para usar as novas funções

**Lógica de Detecção**:
1. **Override manual**: Se `DOCKER_MCP_NATIVE_MODE` estiver definida, respeita o valor
2. **Sistema operacional**: Apenas Linux pode ter modo nativo
3. **Modo contêiner**: Se `DOCKER_MCP_IN_CONTAINER=1`, não usa modo nativo
4. **Socket Docker Desktop**: Verifica se `/run/host-services/backend.sock` não existe
5. **Docker Engine**: Verifica se `/var/run/docker.sock` existe

**Arquivos Modificados**:
1. **pkg/desktop/detection.go** (novo)
   - Implementadas funções `IsNativeMode()`, `IsContainerMode()` e `IsDockerDesktop()`
   - Centralizada toda lógica de detecção em um único arquivo

2. **pkg/desktop/sockets_linux.go**
   - Atualizada para usar `IsNativeMode()` e `IsContainerMode()`
   - Simplificada a lógica de verificação

3. **cmd/docker-mcp/commands/root.go**
   - Modificada função `PersistentPreRunE` para usar funções de detecção
   - Simplificada a verificação do Docker Desktop

4. **pkg/docker/secrets.go**
   - Atualizada função `readSecrets()` para usar `IsNativeMode()`
   - Mantida compatibilidade com método alternativo de segredos

5. **pkg/gateway/run.go**
   - Modificada função `Run()` para usar `IsDockerDesktop()`
   - Simplificada a lógica de configuração OAuth

6. **cmd/docker-mcp/commands/gateway.go**
   - Atualizada função `gatewayCommand()` para usar funções de detecção
   - Simplificada a configuração padrão

**Documentação Atualizada**:
1. **README.md**
   - Atualizada seção de início rápido para remover necessidade de configuração manual
   - Adicionados exemplos com detecção automática
   - Mantida documentação sobre override manual

2. **guia-instalacao.md**
   - Simplificados passos de instalação
   - Removida necessidade de configurar variável de ambiente
   - Adicionada seção sobre override manual (opcional)

3. **guia-configuracao.md**
   - Documentada detecção automática e suas regras
   - Atualizados exemplos para usar detecção automática
   - Mantida documentação sobre override manual

4. **referencia-tecnica.md**
   - Documentado novo arquivo detection.go
   - Explicada implementação das funções de detecção
   - Atualizados fluxos de dados e testes

**Benefícios**:
- **Simplicidade**: Usuários não precisam mais configurar variáveis manualmente
- **Robustez**: Detecção baseada em múltiplos critérios confiáveis
- **Flexibilidade**: Override manual ainda disponível quando necessário
- **Manutenibilidade**: Lógica centralizada facilita manutenção futura

**Testes Recomendados**:
```bash
# Teste 1: Detecção automática
docker mcp gateway run --dry-run --verbose

# Teste 2: Override manual para modo nativo
export DOCKER_MCP_NATIVE_MODE=1
docker mcp gateway run --dry-run

# Teste 3: Override manual para modo Docker Desktop
export DOCKER_MCP_NATIVE_MODE=0
docker mcp gateway run --dry-run
```

### Correção do Problema "No Server is Enabled"

**Descrição**: Foi identificada e corrigida uma issue crítica onde o gateway não conseguia carregar servidores MCP em modo nativo, resultando na mensagem de erro "no server is enabled".

**Causa Raiz**: O arquivo [`cmd/docker-mcp/commands/gateway.go`](cmd/docker-mcp/commands/gateway.go:37) não configurava o `RegistryPath` na configuração padrão do modo nativo, impedindo que o gateway localizasse o arquivo `registry.yaml` que contém a configuração dos servidores habilitados.

**Solução Implementada**:
- Adicionada a linha `RegistryPath: []string{"registry.yaml"}` na configuração padrão do modo nativo
- Esta modificação permite que o gateway encontre e carregue automaticamente os servidores configurados

**Arquivos Modificados**:
1. **cmd/docker-mcp/commands/gateway.go** (linha 37)
   - Adicionado: `RegistryPath: []string{"registry.yaml"}`

**Documentação Atualizada**:
1. **README.md**
   - Adicionada seção sobre configuração do RegistryPath
   - Incluído exemplo prático mostrando servidores sendo carregados
   - Mencionada resolução do problema "no server is enabled"

2. **solucoes-problemas.md**
   - Adicionada seção completa sobre o problema "no server is enabled"
   - Documentada causa raiz e solução implementada
   - Incluídos comandos para verificação e diagnóstico

3. **referencia-tecnica.md**
   - Documentada modificação específica no gateway.go
   - Explicado por que a modificação foi necessária
   - Incluída linha exata que foi adicionada

4. **atualizacoes.md** (novo arquivo)
   - Criado para manter histórico das mudanças
   - Documentadas todas as atualizações com datas e descrições

**Impacto**: Esta correção é essencial para o funcionamento adequado do modo nativo, permitindo que usuários configurem e utilizem servidores MCP sem dependências do Docker Desktop.

**Testes Recomendados**:
```bash
# Testar configuração básica
export DOCKER_MCP_NATIVE_MODE=1
echo 'servers:
  filesystem:
    enabled: true' > registry.yaml

docker mcp gateway run --registry ./registry.yaml --dry-run --verbose
```

## Histórico de Modificações Anteriores

### Implementação Inicial do Modo Nativo

**Data**: Implementação inicial (data não especificada)

**Descrição**: Implementação do modo nativo para permitir operação do Docker MCP Gateway com Docker Engine nativo do Linux, sem dependências do Docker Desktop.

**Principais Modificações**:
1. **pkg/desktop/sockets_linux.go**
   - Modificada função `getDockerDesktopPaths()` para ignorar verificação do Docker Desktop em modo nativo

2. **cmd/docker-mcp/commands/root.go**
   - Modificada função `PersistentPreRunE` para pular verificação do Docker Desktop

3. **pkg/docker/secrets.go**
   - Implementado gerenciamento alternativo de segredos via arquivos .env

4. **pkg/gateway/run.go**
   - Desabilitado monitor OAuth em modo nativo
   - Removida dependência de autenticação específica do Docker Desktop

5. **cmd/docker-mcp/commands/gateway.go**
   - Configuradas opções padrão para modo nativo
   - Definido uso de arquivo .env local para segredos

**Variáveis de Ambiente Introduzidas**:
- `DOCKER_MCP_NATIVE_MODE=1`: Ativa modo de operação nativo
- `DOCKER_MCP_SECRETS_FILE`: Especifica arquivo alternativo de segredos

## Próximas Atualizações Planejadas

### Melhorias Futuras

1. **Autenticação Nativa**
   - Implementar sistema de autenticação específico para modo nativo
   - Permitir configuração de tokens de API personalizados

2. **Gerenciamento Avançado de Segredos**
   - Suporte a múltiplos formatos de arquivo de segredos
   - Integração com gerenciadores de segredos externos

3. **Monitoramento e Métricas**
   - Implementar métricas específicas para modo nativo
   - Dashboard de monitoramento de servidores MCP

4. **Clustering**
   - Suporte a múltiplas instâncias em modo nativo
   - Balanceamento de carga entre servidores

## Diretrizes de Atualização

### Processo de Documentação

1. **Identificar Mudança**: Qualquer modificação que afete o comportamento do sistema
2. **Documentar Código**: Adicionar comentários explicando o propósito da mudança
3. **Atualizar Docs**: Modificar documentos relevantes (README, soluções-problemas, etc.)
4. **Registrar Aqui**: Adicionar entrada neste arquivo com data e descrição
5. **Testar**: Verificar que a documentação está correta e completa

### Padrão de Entrada

```markdown
## DD de MMMM de YYYY

### Título da Atualização

**Descrição**: Breve descrição da mudança e seu propósito.

**Motivação**: Por que a mudança foi necessária.

**Implementação**: Detalhes técnicos da implementação.

**Arquivos Modificados**:
1. **arquivo.md** (linha X)
   - Descrição da modificação

**Documentação Atualizada**:
1. **documento.md**
   - Descrição das atualizações na documentação

**Impacto**: Efeito da mudança no sistema.

**Testes**: Comandos ou procedimentos para testar a mudança.
```

## Contribuições

Este documento é mantido pela comunidade e contribuições são bem-vindas. Para adicionar uma nova entrada:

1. Siga o padrão acima
2. Seja claro e conciso
3. Inclua datas precisas
4. Referencie os arquivos modificados
5. Descreva o impacto da mudança

## Contato

Para dúvidas ou sugestões sobre este documento:
1. Abra uma issue no repositório
2. Participe das discussões da comunidade
3. Consulte os outros documentos desta série para contexto adicional