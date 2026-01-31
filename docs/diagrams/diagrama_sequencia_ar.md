### Fluxo de Reparos Adicionais (AR)

```mermaid
sequenceDiagram
    participant EndUser as End User<br/>(Usuário Mobile)
    participant Client as Client<br/>(Mecânico)
    participant OS as OS Service
    participant Entity as Entity API<br/>(cliente, veiculo, estoque)
    participant Billing as Billing Service
    participant Notification as Mobile Notification

    Note over OS: Uma OS principal já existe e está APROVADA ou EM EXECUÇÃO.

    %% Solicitação de Reparo Adicional
    Client->>OS: POST /os/{id}/reparos (incluir peças/serviços adicionais)
    OS->>OS: Cria entidade "Reparo Adicional"<br/>Status: RECEBIDO
    
    OS->>Entity: GET /pecas/{id}/estoque (valida quantidade)
    Entity-->>OS: Resposta: Disponível
    OS->>Entity: PATCH /pecas/{id}/reservar
    Entity-->>OS: Resposta: Peças reservadas

    OS->>Billing: Recalcula Orçamento Total (OS original + Reparo Adicional)
    Billing-->>OS: 200 OK

    OS->>OS: Altera Status do Reparo: RECEBIDO -> AGUARDANDO APROVAÇÃO
    OS->>Notification: Dispara Notificação (Novo orçamento para aprovação)

    alt Aprovação do Novo Orçamento
        EndUser->>OS: PATCH /reparos/{id_reparo} (ação: aprovar)
        OS->>OS: Altera Status do Reparo: AGUARDANDO APROVAÇÃO -> APROVADO
        OS->>OS: Atualiza novo orçamento na OS
        OS->>Entity: PATCH /pecas/dar-baixa (confirma uso das novas peças)
        Entity-->>OS: Resposta: Baixa no estoque confirmada
        OS-->>EndUser: Resposta: 200 OK
        Note over OS: A OS principal continua sua execução com o novo escopo.
    else Rejeição do Novo Orçamento
        EndUser->>OS: PATCH /reparos/{id_reparo} (ação: rejeitar)
        OS->>OS: Altera Status do Reparo: AGUARDANDO APROVAÇÃO -> CANCELADO
        OS->>Entity: PATCH /pecas/liberar-reserva (libera peças do reparo)
        Entity-->>OS: Resposta: Reserva liberada
        OS-->>EndUser: Resposta: 200 OK
        Note over OS: A OS principal continua com seu escopo original,<br/>ignorando o reparo adicional.
    end
```

### Descrição do Fluxo

#### Participantes
- **End User (Usuário Mobile)**: Cliente final que deve aprovar ou rejeitar o novo orçamento.
- **Client (Mecânico)**: Profissional que identifica a necessidade de reparo extra.
- **OS Service**: Gerencia o ciclo de vida do reparo adicional independente da OS principal.
- **Entity API**: Valida e reserva peças do estoque.
- **Billing Service**: Recalcula o orçamento total da OS considerando os novos itens.
- **Mobile Notification**: Notifica o cliente sobre a nova solicitação pendente.

#### Fluxos Principais

1. **Solicitação**: Mecânico adiciona novos serviços/peças a uma OS já em andamento.
2. **Reserva e Cálculo**: Sistema reserva itens e calcula impacto financeiro.
3. **Notificação**: Cliente é avisado de uma alteração no orçamento.
4. **Decisão**:
    - **Aprovação**: Baixa no estoque e atualização do valor da OS.
    - **Rejeição**: Liberação do estoque e cancelamento do reparo adicional.

#### Status do Reparo Adicional
- **RECEBIDO** (IN_ANALYSIS): Reparo criado, aguardando validações internas.
- **AGUARDANDO APROVAÇÃO**: Orçamento calculado e enviado ao cliente.
- **APROVADO**: Cliente aceitou o custo adicional.
- **CANCELADO**: Cliente rejeitou o reparo adicional.