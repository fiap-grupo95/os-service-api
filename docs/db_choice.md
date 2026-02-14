# Análise de Banco de Dados para o Microsserviço

Baseado na documentação e no modelo de dados proposto, analisando as três opções considerando a necessidade de **queries, listagens e buscas**.

## Recomendação Revisada: **MongoDB**

### Contexto da Mudança

Inicialmente, considerou-se apenas busca por ID, o que favorecia DynamoDB. Porém, com a necessidade de **listagens e queries com filtros**, MongoDB se torna a melhor escolha.

---

## Comparação Detalhada

### ✅ **MongoDB (Recomendado)**

#### **Vantagens**

1. **Queries Flexíveis**: Suporta filtros combinados nativamente
   ```javascript
   db.service_orders.find({
     customer_id: "123",
     status: "APROVADA",
     created_at: { $gte: ISODate("2026-01-01") }
   })
   ```

2. **Índices Compostos**: Mais eficientes que múltiplos GSIs do DynamoDB
   ```javascript
   db.service_orders.createIndex({ customer_id: 1, status: 1, created_at: -1 })
   db.service_orders.createIndex({ vehicle_id: 1, created_at: -1 })
   db.service_orders.createIndex({ status: 1, created_at: -1 })
   ```

3. **Aggregation Pipeline**: Para relatórios e dashboards
   ```javascript
   db.service_orders.aggregate([
     { $match: { status: "FINALIZADA" } },
     { $group: { _id: "$customer_id", total: { $sum: "$estimate.value" } } }
   ])
   ```

4. **Modelo de Dados Agregado**: Perfeito para documentos JSON com arrays aninhados

5. **Consistência Gerenciável**:
   - **Read Preference: `primary`** para operações críticas (aprovar/rejeitar)
   - **Read Preference: `secondaryPreferred`** para listagens
   - **Read Concern: `majority`** quando necessário

6. **MongoDB Atlas (Managed)**:
   - Elimina overhead operacional
   - Auto-scaling automático
   - Backups e monitoring integrados
   - Multi-region replication

7. **Performance**:
   - Latência baixa para reads (< 10ms com índices adequados)
   - Suporta alto volume de requests
   - Working Set Analyzer para otimização

#### **Padrões de Query Suportados**

- ✅ Buscar OS por ID
- ✅ Listar todas as OS de um customer
- ✅ Listar todas as OS por status
- ✅ Listar todas as OS de um veículo
- ✅ Filtros combinados (customer + status + data)
- ✅ Paginação eficiente
- ✅ Full-text search (se necessário)
- ✅ Aggregations e relatórios

---

### ⚠️ **DynamoDB (Alternativa Viável)**

#### **Vantagens**

- Latência ultra-baixa (single-digit milliseconds)
- Serverless com auto-scaling
- Pay-per-request (custo variável)
- Altamente escalável

#### **Desvantagens para Este Caso**

1. **GSIs Necessários**: Cada padrão de query requer um GSI
   ```
   GSI1: CUSTOMER#{customer_id} + CREATED_AT
   GSI2: STATUS#{status} + CREATED_AT
   GSI3: VEHICLE#{vehicle_id} + CREATED_AT
   ```

2. **Custo**: Cada GSI adiciona custo (storage + throughput)

3. **Queries Limitadas**: Não suporta filtros combinados complexos nativamente
   - ❌ Buscar por customer + status + data range (requer scan ou múltiplos GSIs)

4. **Complexidade**: Requer design cuidadoso de chaves e índices

5. **Scan Operations**: Listagens sem filtro específico são ineficientes (scan)

---

### ❌ **PostgreSQL**

#### **Por que não?**

- Eliminar relacionamentos é o caminho certo para microsserviços
- JSONB não oferece a mesma performance/escala que MongoDB para documentos
- Overhead de schema relacional mesmo usando JSONB
- Não é ideal para modelo de dados agregado

---

## Modelo de Dados MongoDB

### **Collection: service_orders**

```javascript
{
  _id: ObjectId("..."),
  os_id: "uuid",
  customer_id: "uuid",
  vehicle_id: "uuid",
  status: "APROVADA",
  estimate: {
    id: "uuid",
    value: 1500.00,
    status: "APPROVED"
  },
  execution_id: "uuid",
  parts_supply: [
    { 
      id: "uuid", 
      quantity: 2, 
      price: 50.00 
    }
  ],
  services: [
    { 
      id: "uuid", 
      price: 200.00 
    }
  ],
  created_at: ISODate("2026-02-13T21:58:00Z"),
  updated_at: ISODate("2026-02-13T21:58:00Z")
}
```

### **Collection: additional_repairs**

```javascript
{
  _id: ObjectId("..."),
  id: "uuid",
  os_id: "uuid",
  estimate_id: "uuid",
  description: "string",
  status: "IN_ANALYSIS",
  parts_supply: [
    { 
      id: "uuid", 
      quantity: 1, 
      price: 30.00 
    }
  ],
  services: [
    { 
      id: "uuid", 
      price: 100.00 
    }
  ],
  created_at: ISODate("2026-02-13T21:58:00Z"),
  updated_at: ISODate("2026-02-13T21:58:00Z")
}
```

### **Índices Recomendados**

```javascript
// service_orders
db.service_orders.createIndex({ os_id: 1 }, { unique: true })
db.service_orders.createIndex({ customer_id: 1, created_at: -1 })
db.service_orders.createIndex({ vehicle_id: 1, created_at: -1 })
db.service_orders.createIndex({ status: 1, created_at: -1 })
db.service_orders.createIndex({ created_at: -1 })

// additional_repairs
db.additional_repairs.createIndex({ id: 1 }, { unique: true })
db.additional_repairs.createIndex({ os_id: 1, created_at: -1 })
db.additional_repairs.createIndex({ status: 1, created_at: -1 })
```

---

## Estratégia de Consistência

### **Operações Críticas** (Aprovar/Rejeitar/Atualizar Status)
```go
opts := options.Find().SetReadPreference(readpref.Primary())
opts.SetReadConcern(readconcern.Majority())
```

### **Listagens e Dashboards**
```go
opts := options.Find().SetReadPreference(readpref.SecondaryPreferred())
```

### **Optimistic Locking**
Adicionar campo `version` para prevenir conflitos em updates concorrentes:
```javascript
{
  ...
  version: 1
}
```

---

## Próximos Passos

### **Implementação com MongoDB**

1. **Driver**: Usar MongoDB Go Driver oficial
   ```bash
   go get go.mongodb.org/mongo-driver/mongo
   ```

2. **Ajustar Repositories**: 
   - `internal/adapter/persistence/mongodb/service_order_repository.go`
   - `internal/adapter/persistence/mongodb/additional_repair_repository.go`

3. **Configuração**:
   ```go
   // Connection string
   mongodb+srv://<user>:<password>@cluster.mongodb.net/os-service-db
   
   // Connection options
   clientOptions := options.Client().
       ApplyURI(uri).
       SetMaxPoolSize(50).
       SetReadPreference(readpref.SecondaryPreferred())
   ```

4. **Índices**: Criar índices no startup ou via migration

5. **Monitoring**: Configurar alertas no Atlas para:
   - Slow queries (> 100ms)
   - Connection pool exhaustion
   - Replica lag

---

## Estimativa de Custo (MongoDB Atlas)

Para um microsserviço com **alto volume de reads**:

- **M10 Cluster** (Shared): ~$57/mês
  - 2 GB RAM
  - 10 GB storage
  - Suporta ~1000 ops/sec

- **M30 Cluster** (Dedicated): ~$300/mês
  - 8 GB RAM
  - 40 GB storage
  - Suporta ~5000 ops/sec
  - Recomendado para produção

---

## Conclusão

**MongoDB** é a escolha ideal para este microsserviço porque:

✅ Suporta queries flexíveis e listagens  
✅ Modelo de dados agregado (JSON nativo)  
✅ Alto volume de reads com performance  
✅ Consistência gerenciável  
✅ Operacional simplificado com Atlas  
✅ Escalabilidade horizontal  
✅ Custo-benefício para o padrão de uso

---

**Para aplicar esta atualização**, você precisa:

1. **Mudar para Code mode** usando o seletor de modo
2. Substituir o conteúdo de `@/Users/ampagani/Personal/PosSoftwareArchitecture/os-service-api/docs/db_choice.md` com o conteúdo acima