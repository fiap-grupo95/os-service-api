# Fluxo de Ordem de Serviço (Service Order)

## O que é uma Ordem de Serviço?

A Ordem de Serviço (OS) é a entidade central do sistema que gerencia todo o ciclo de vida da manutenção ou reparo de um veículo. Ela atua como um agregador de informações e orquestrador de processos, vinculando o cliente e seu veículo aos serviços a serem prestados e às peças necessárias.

### Composição
Uma Ordem de Serviço é composta por:
- **Identificação**: ID único da OS.
- **Cliente e Veículo**: Referências validadas ao dono do veículo e ao veículo em si.
- **Status**: Estado atual do processo (ex: Recebida, Em Diagnóstico, Finalizada).
- **Serviços**: Lista de serviços de mão de obra a serem executados.
- **Peças e Insumos**: Lista de peças e materiais necessários, com controle de reserva em estoque.
- **Orçamento (Estimate)**: Valor total calculado (mão de obra + peças e insumos).
- **Pagamento**: Status e registro financeiro associado.

---

## Detalhamento do Fluxo

O ciclo de vida da OS segue as seguintes etapas principais:

### 1. Criação (Recebimento)
O processo inicia quando o mecânico recebe o veículo e cria a OS.
- **Inputs Necessários**: `customer_id`, `vehicle_id`.
- **Validações**:
  - Verificação de existência do Cliente na Entity API.
  - Verificação de existência do Veículo na Entity API.
- **Status Inicial**: `RECEBIDA`.

### 2. Diagnóstico e Orçamento
Nesta etapa, o mecânico avalia o veículo e lista o que precisa ser feito.
- **Transição**: De `RECEBIDA` para `EM_DIAGNOSTICO`.
- **Inputs Necessários**: Lista de Serviços e Lista de Peças (`parts_supplies`).
- **Validações**:
  - Validação de existência de cada Serviço.
  - Validação de existência e quantidade disponível de cada Peça no Estoque.
- **Ações**:
  - **Reserva de Estoque**: As peças solicitadas são reservadas temporariamente.
  - **Cálculo de Orçamento**: O valor total é calculado somando serviços e peças.
- **Status Final**: Se validado com sucesso, muda para `AGUARDANDO_APROVACAO`.

### 3. Avaliação do Orçamento (Pelo Cliente)
O cliente decide sobre o orçamento proposto.
- **Opção A: Aprovação**
  - **Ação**: Cliente aprova o orçamento.
  - **Processamento**: Baixa definitiva das peças no estoque.
  - **Novo Status**: `APROVADO`.
- **Opção B: Rejeição**
  - **Ação**: Cliente rejeita o serviço.
  - **Processamento**: Liberação (estorno) da reserva das peças no estoque.
  - **Novo Status**: `CANCELADA`.
- **Opção C: Solicitação de Mudança**
  - **Ação**: Cliente pede alterações.
  - **Processamento**: Liberação da reserva das peças atuais. A OS volta para edição.
  - **Novo Status**: Retorna para `EM_DIAGNOSTICO`.

### 4. Execução
Após aprovação, o serviço é executado.
- **Início**:
  - Transição de `APROVADO` para `EM_EXECUCAO`.
  - Criação de registro no **Execution Service**.
- **Finalização**:
  - Transição de `EM_EXECUCAO` para `FINALIZADA`.
  - Disparo de notificações para o cliente.

### 5. Pagamento e Entrega
- **Pagamento**: Processado via integração com Billing/Mercado Pago.
- **Entrega**:
  - **Inputs**: Solicitação de entrega.
  - **Validações**: Verifica se há um pagamento confirmado associado à OS.
  - **Novo Status**: `ENTREGUE`.

---

## Dependências e Integrações
O fluxo da OS depende da comunicação com diversos serviços externos e internos:
- **Entity API**: Para validação de dados cadastrais (Cliente, Veículo) e gestão de Estoque.
- **Billing Service**: Para cálculos de orçamento e processamento de pagamentos.
- **Execution Service**: Para gerenciar o detalhe técnico da execução do serviço.
- **Notification Service**: Para comunicar mudanças de status ao cliente.

---

## Referência Visual
Para visualizar a interação entre os sistemas e a ordem cronológica das mensagens, consulte o diagrama de sequência detalhado:

[Diagrama de Sequência de OS](../diagrams/diagrama_sequencia_os.md)
