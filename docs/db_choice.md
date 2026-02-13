# Análise de Banco de Dados para o Microsserviço

Baseado na documentação e no modelo de dados proposto, vou analisar as três opções:

## Recomendação: **DynamoDB**

### Justificativa

Para este microsserviço específico, o **DynamoDB** é a melhor escolha pelos seguintes motivos:

#### ✅ **Vantagens do DynamoDB**

1. **Padrão de acesso por ID**: Você mencionou que não há necessidade de consultas complexas além de buscar por ID - isso é o cenário ideal para DynamoDB (key-value store)

2. **Alto volume de GETs**: DynamoDB oferece:
   - Latência consistente em single-digit milliseconds
   - Auto-scaling horizontal nativo
   - DAX (DynamoDB Accelerator) para cache in-memory se necessário

3. **Modelo de dados agregado**: Seu modelo já está pensado em agregados (Service Order com arrays de parts_supply e services) - perfeito para documento JSON no DynamoDB

4. **Consistência eventual gerenciável**: 
   - DynamoDB oferece **strongly consistent reads** quando necessário
   - Para GETs críticos (ex: aprovar/rejeitar OS), use `ConsistentRead=true`
   - Para listagens menos críticas, eventual consistency é aceitável

5. **Escalabilidade**: 
   - Serverless (on-demand pricing) ou provisionado
   - Sem limites práticos de escala
   - Ideal para microsserviços com tráfego variável

6. **Custo-benefício**: Pay-per-request em modo on-demand

#### ⚠️ **Por que não MongoDB?**

- **Consistência eventual mais complexa**: Réplicas do MongoDB podem ter lag maior que DynamoDB
- **Operacional**: Requer gerenciamento de cluster (mesmo com Atlas)
- **Custo**: Mais caro que DynamoDB para alto volume de reads simples
- **Over-engineering**: Oferece recursos (queries complexas, aggregations) que você não precisa

#### ❌ **Por que não PostgreSQL?**

Você está correto:
- Eliminar relacionamentos é o caminho certo para microsserviços
- JSONB no Postgres não oferece a mesma performance/escala que DynamoDB para este padrão
- Overhead de joins/relacionamentos mesmo usando JSONB

### Modelo de Dados Sugerido para DynamoDB

```
Table: ServiceOrders
PK: OS#{os_id}
SK: METADATA

Attributes:
{
  "os_id": "uuid",
  "vehicle_id": "uuid",
  "customer_id": "uuid", 
  "status": "RECEBIDA",
  "estimate": {
    "id": "uuid",
    "value": 1500.00
  },
  "execution_id": "uuid",
  "parts_supply": [
    {"id": "uuid", "quantity": 2, "price": 50.00}
  ],
  "services": [
    {"id": "uuid", "price": 200.00}
  ],
  "created_at": "2026-02-13T21:58:00Z",
  "updated_at": "2026-02-13T21:58:00Z"
}

Table: AdditionalRepairs
PK: AR#{id}
SK: METADATA
GSI1PK: OS#{os_id}  // Para buscar todos ARs de uma OS

Attributes:
{
  "id": "uuid",
  "os_id": "uuid",
  "estimate_id": "uuid",
  "description": "string",
  "status": "IN_ANALYSIS",
  "created_at": "2026-02-13T21:58:00Z",
  "updated_at": "2026-02-13T21:58:00Z"
}
```

### Estratégia para Consistência

Para mitigar inconsistência eventual nos GETs críticos:

1. **Operações críticas** (aprovar/rejeitar/atualizar status): Use `ConsistentRead=true`
2. **Listagens/dashboards**: Eventual consistency é aceitável
3. **Considere**: Adicionar versioning (`version` field) para optimistic locking em updates concorrentes

### Próximos Passos

Se optar por DynamoDB, você precisará:
- Ajustar os repositories em `@/Users/ampagani/Personal/PosSoftwareArchitecture/os-service-api/internal/adapter/persistence`
- Usar AWS SDK for Go v2
- Configurar índices GSI para queries por `os_id` em AdditionalRepairs
- Implementar strongly consistent reads para operações críticas