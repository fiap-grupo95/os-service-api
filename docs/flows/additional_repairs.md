# Fluxo de Reparos Adicionais (Additional Repairs)

## O que são Reparos Adicionais?

Reparos Adicionais são serviços ou peças extras identificados como necessários após a aprovação inicial de uma Ordem de Serviço (OS). Isso geralmente ocorre quando, durante a execução do serviço principal, o mecânico descobre novos problemas que não foram mapeados no diagnóstico inicial.

Diferente de editar a OS original (que já pode estar aprovada ou em execução), cria-se uma entidade separada de **Reparo Adicional** vinculada à OS. Isso permite um controle granular de aprovação, onde o cliente pode aprovar ou rejeitar apenas o extra, sem cancelar o serviço principal que já está em andamento.

## Relação com a Ordem de Serviço

- **Vínculo**: Cada Reparo Adicional pertence a uma única OS (`service_order_id`).
- **Independência de Status**: O status do Reparo Adicional (`IN_ANALYSIS`, `APPROVED`, `DENIED`) evolui independentemente do status da OS principal.
- **Impacto Financeiro**: 
  - O orçamento do Reparo Adicional é calculado separadamente.
  - Ao ser **Aprovado**, o valor do orçamento da OS principal é substituído pelo valor do orçamento do Reparo Adicional (`repoOS.UpdateEstimate`).

---

## Detalhamento do Fluxo

### 1. Solicitação (Criação)
O mecânico identifica a necessidade de novos itens.
- **Pré-requisito**: Existir uma OS principal (geralmente em status `APROVADO` ou `EM_EXECUCAO`).
- **Inputs Necessários**: `service_order_id`, `description`, lista de `services`, lista de `parts_supplies`.
- **Validações**:
  - Validação da existência da OS.
  - Validação de existência de cada Serviço.
  - Validação de disponibilidade de estoque para as Peças.
- **Ações**:
  - Reserva das peças solicitadas no estoque.
  - Cálculo do orçamento adicional.
- **Status Inicial**: `IN_ANALYSIS` (Em Análise/Recebido).

### 2. Edição (Opcional)
Enquanto estiver `IN_ANALYSIS`, é possível adicionar ou remover itens.
- **Métodos**: `AddPartSupplyAndService`, `RemovePartSupplyAndService`.
- **Comportamento**: Recalcula o orçamento e ajusta as reservas de estoque conforme necessário.

### 3. Aprovação ou Rejeição (Pelo Cliente)
O cliente recebe a notificação do novo orçamento e decide.

#### Opção A: Aprovação (`APPROVED`)
- **Ação**: Cliente aprova o orçamento adicional.
- **Processamento**:
  1. **Baixa de Estoque**: Confirma o uso das peças (transforma reserva em baixa definitiva).
  2. **Atualização da OS**: O valor do reparo é somado ao `Estimate` da OS principal.
  3. **Atualização de Status**: Define status como `APPROVED`.
- **Resultado**: O mecânico está autorizado a executar os serviços extras.

#### Opção B: Rejeição (`DENIED`)
- **Ação**: Cliente nega o orçamento adicional.
- **Processamento**:
  1. **Estorno de Estoque**: As peças reservadas para este reparo são liberadas de volta para o estoque.
  2. **Atualização de Status**: Define status como `DENIED`.
- **Resultado**: A OS principal segue apenas com o escopo original. O reparo adicional é descartado.

---

## Alterações na OS Principal
A principal alteração refletida na OS original é financeira.
- O campo `Estimate` da OS é atualizado para refletir `Valor Original + Valor do Reparo Adicional Aprovado`.
- Isso garante que, no momento do pagamento (Billing), o valor total cobrado inclua todos os serviços autorizados.

---

## Dependências
- **OS Service**: Para vínculo e atualização do orçamento global.
- **Billing Service**: Para cálculo do orçamento adicional.
- **Entity API**: Para gestão de estoque (Reserva/Baixa/Estorno).
- **Notification Service**: Para alertar o cliente sobre a nova solicitação de aprovação.

---

## Referência Visual
Para visualizar a interação entre os sistemas e a ordem cronológica das mensagens, consulte o diagrama de sequência específico para este fluxo:

[Diagrama de Sequência de Reparos Adicionais](../diagrams/diagrama_sequencia_ar.md)
