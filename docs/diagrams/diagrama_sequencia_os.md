### Fluxo de Ordem de Serviço (OS)

```mermaid
sequenceDiagram
    participant EndUser as End User<br/>(Usuário Mobile)
    participant Client as Client<br/>(Mecânico)
    participant Auth as Auth Lambda
    participant OS as OS Service
    participant Entity as Entity API<br/>(cliente, veiculo, estoque)
    participant Billing as Billing Service
    participant Execution as Execution Service
    participant Notification as Mobile Notification
    participant Payment as Mercado Pago

    Note over OS: Tudo vai bater na OS service<br/>A única que vai conter o modelo de OS,<br/>e vai manipular status
    Note over Execution: Não altera status de OS,<br/>modelo de execution

    %% Login Flow
    Client->>Auth: POST Login (Usuário e senha)
    alt Falha na autenticação
        Auth-->>Client: Resposta: 500 Internal error
    else Sucesso
        Auth-->>Client: Resposta: 200 OK Token JWT
    end

    %% Criar OS
    Client->>OS: POST /os (dados cliente, veiculo)
    OS->>Entity: Valida se cliente e veículo existe, senão cria
    Entity-->>OS: 200 OK
    OS->>OS: Cria OS com Status: RECEBIDA
    OS-->>Client: Resposta: 201 Created
    OS-->>EndUser: Resposta: 201 Created

    %% Fluxo alternativo - Cancelamento
    alt Cancelamento de OS
        Client->>OS: PATCH /os/{id}/diagnosis (Cliente deseja cancelar OS)
        OS->>OS: Altera Status: RECEBIDA -> CANCELADA
        OS-->>Client: 200 OK
    else Continuar com diagnóstico
        %% Iniciar diagnóstico
        Client->>OS: PATCH /os/{id}/diagnosis (Iniciar diagnóstico)
        OS->>OS: Altera Status: RECEBIDA -> EM DIAGNÓSTICO
        OS-->>Client: Resposta: 200 OK (Diagnóstico iniciado)

        %% Incluir peças e serviços
        Client->>OS: PATCH /os/{id}/diagnosis (incluir peças/serviços)
        OS->>OS: Altera Status: RECEBIDA -> EM DIAGNÓSTICO
        OS->>Entity: GET /servico/{id} (valida se serviço existe)
        OS->>Entity: GET /pecas/{id}/estoque (valida quantidade)
        Entity-->>OS: Resposta: Disponível
        OS->>Entity: PATCH /pecas/{id}/reservar
        Entity-->>OS: Resposta: Peças reservadas

        %% Calcular orçamento
        OS->>Billing: Calcula Orçamento
        Billing-->>OS: 200 OK
        OS->>OS: Altera Status: EM DIAGNÓSTICO -> AGUARDANDO APROVAÇÃO

        %% Decisão sobre orçamento
        alt Aprovação do Orçamento
            EndUser->>OS: PATCH /os/{id}/estimate (ação: aprovar)
            OS->>Billing: PATCH /os/{id}/estimate (ação: aprovar)
            Billing-->>OS: 200 OK
            OS->>OS: Altera Status: AGUARDANDO APROVAÇÃO -> APROVADO
            OS->>Entity: PATCH /pecas/dar-baixa (confirma uso das peças)
            Entity-->>OS: Resposta: Baixa no estoque confirmada
            OS-->>EndUser: Resposta: 200 OK

            %% Iniciar execução
            Client->>OS: PATCH /os/{id}/execution (iniciar execução)
            OS->>OS: Altera Status: APROVADO -> EM EXECUÇÃO
            OS-->>Client: Resposta: 200 OK
            OS->>Execution: PATCH /os/{id}/execution (cria execução)
            Execution-->>OS: 200 OK

            %% Executar serviço
            Note over Client,Execution: ...Serviço é executado...

            %% Finalizar serviço
            Client->>OS: PATCH /os/{id}/execution (finalizar serviço)
            OS->>Execution: POST /execution/os_id/finish (finalizar serviço)
            Execution-->>OS: 200 OK
            OS->>OS: Altera Status: EM EXECUÇÃO -> FINALIZADA
            OS-->>Client: Resposta: 200 OK
            OS->>Notification: Dispara Notificação (Serviço finalizado)

            %% Pagamento
            EndUser->>OS: POST /payment
            OS->>Billing: Cria e processa Pagamento
            Billing->>Payment: Cria e processa Pagamento
            Payment-->>Billing: Resposta: 200 OK (pagamento processado OK)
            Billing-->>OS: Resposta: 200 OK (pagamento processado OK)
            OS-->>EndUser: Resposta: 200 OK (pagamento processado OK)

            %% Entrega
            Client->>OS: PATCH /os/{id}/delivery (Entrega do veículo)
            OS->>Billing: Valida se a OS possui pagamento associado
            OS->>OS: Altera Status: FINALIZADA -> ENTREGUE

        else Rejeição do Orçamento
            EndUser->>OS: PATCH /os/{id}/estimate (ação: rejeitar)
            OS->>Billing: PATCH /os/{id}/estimate (ação: rejeitar orçamento)
            Billing-->>OS: 200 OK
            OS->>Entity: PATCH /pecas/liberar-reserva
            Entity-->>OS: Resposta: Reserva liberada
            OS->>OS: Altera Status: AGUARDANDO APROVAÇÃO -> CANCELADO
            OS-->>EndUser: Resposta: 200 OK

        else Solicitação de Mudança
            EndUser->>OS: PATCH /os/{id}/estimate (ação: solicitar_mudança)
            OS->>Billing: PATCH /os/{id}/estimate (ação: cancelar orçamento)
            OS->>Entity: PATCH /pecas/liberar-reserva
            Entity-->>OS: Resposta: Reserva liberada
            OS->>OS: Altera Status: AGUARDANDO APROVAÇÃO -> EM DIAGNÓSTICO
            OS-->>EndUser: Resposta: 200 OK
            Note over OS,Entity: Fluxo retorna para a etapa de diagnóstico e orçamento.
        end
    end
```

### Descrição do Fluxo

#### Participantes
- **End User (Usuário Mobile)**: Cliente final que acompanha a OS pelo aplicativo mobile
- **Client (Mecânico)**: Profissional que executa o serviço
- **Auth Lambda**: Serviço de autenticação
- **OS Service**: Serviço principal que gerencia as ordens de serviço e seus status
- **Entity API**: API de entidades (clientes, veículos, estoque)
- **Billing Service**: Serviço de cobrança e orçamento
- **Execution Service**: Serviço de execução (não altera status da OS)
- **Mobile Notification**: Serviço de notificações push
- **Mercado Pago**: Gateway de pagamento

#### Fluxos Principais

1. **Autenticação**: Mecânico faz login no sistema
2. **Criação da OS**: Registro inicial com dados do cliente e veículo
3. **Diagnóstico**: Identificação dos problemas e inclusão de peças/serviços necessários
4. **Orçamento**: Cálculo e apresentação do orçamento ao cliente
5. **Aprovação**: Cliente pode aprovar, rejeitar ou solicitar mudanças
6. **Execução**: Realização dos serviços aprovados
7. **Pagamento**: Processamento do pagamento via Mercado Pago
8. **Entrega**: Finalização e entrega do veículo ao cliente

#### Status da OS
- **RECEBIDA**: OS criada inicialmente
- **EM DIAGNÓSTICO**: Análise do veículo em andamento
- **AGUARDANDO APROVAÇÃO**: Orçamento enviado ao cliente
- **APROVADO**: Cliente aprovou o orçamento
- **EM EXECUÇÃO**: Serviços sendo executados
- **FINALIZADA**: Serviços concluídos
- **ENTREGUE**: Veículo entregue ao cliente
- **CANCELADA**: OS cancelada pelo cliente
