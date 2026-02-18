# New Relic - Deploy via Manifestos Kubernetes

Este diretório contém os manifestos Kubernetes para deploy do New Relic Infrastructure e coleta de logs.

## Arquivos

- `newrelic-namespace.yml` - Namespace para recursos do New Relic
- `newrelic-secret.yml` - Secret com a License Key
- `newrelic-configmap.yml` - Configurações gerais do cluster
- `newrelic-rbac.yml` - Service Account, ClusterRole e ClusterRoleBinding
- `newrelic-daemonset-infra.yml` - Infrastructure Agent (métricas dos nodes)
- `newrelic-fluent-bit-configmap.yml` - Configuração do Fluent Bit
- `newrelic-daemonset-fluent-bit.yml` - Fluent Bit para coleta de logs
- `newrelic-kube-state-metrics.yml` - Kube State Metrics
- `newrelic-events-configmap.yml` - Configuração do Events Collector
- `newrelic-deployment-events.yml` - Events Collector

## Pré-requisitos

1. Cluster Kubernetes rodando
2. `kubectl` configurado e conectado ao cluster
3. License Key do New Relic ([obtenha aqui](https://one.newrelic.com/admin-portal/api-keys/home))

## Instalação

### Passo 1: Configurar License Key

Edite o arquivo `newrelic-secret.yml` e substitua `YOUR_LICENSE_KEY_HERE` pela sua License Key:

```yaml
stringData:
  license: "sua_license_key_aqui"
```

### Passo 2: Aplicar os Manifestos

#### Aplicação Manual

```bash
# 1. Namespace
kubectl apply -f newrelic-namespace.yml

# 2. Secret
kubectl apply -f newrelic-secret.yml

# 3. ConfigMaps
kubectl apply -f newrelic-configmap.yml
kubectl apply -f newrelic-fluent-bit-configmap.yml
kubectl apply -f newrelic-events-configmap.yml

# 4. RBAC
kubectl apply -f newrelic-rbac.yml

# 5. Infrastructure Agent
kubectl apply -f newrelic-daemonset-infra.yml

# 6. Fluent Bit (Logs)
kubectl apply -f newrelic-daemonset-fluent-bit.yml

# 7. Kube State Metrics
kubectl apply -f newrelic-kube-state-metrics.yml

# 8. Events Collector
kubectl apply -f newrelic-deployment-events.yml
```

## Verificação

### Verificar pods

```bash
kubectl get pods -n newrelic
```

Você deve ver:
- `newrelic-infra-xxxxx` (um por node - DaemonSet)
- `fluent-bit-xxxxx` (um por node - DaemonSet)
- `kube-state-metrics-xxxxx`
- `newrelic-kube-events-xxxxx`

### Verificar logs

```bash
# Logs do Fluent Bit
kubectl logs -n newrelic -l app=fluent-bit --tail=50

# Logs do Infrastructure Agent
kubectl logs -n newrelic -l app=newrelic-infra --tail=50

# Logs do Events Collector
kubectl logs -n newrelic -l app=newrelic-kube-events --tail=50
```

### Status detalhado

```bash
# Todos os recursos no namespace
kubectl get all -n newrelic

# Descrever um pod específico
kubectl describe pod -n newrelic <nome-do-pod>
```

## Customizações

### Alterar nome do cluster

Edite `newrelic-configmap.yml`:
```yaml
data:
  cluster-name: "seu-cluster-name"
```

### Adicionar atributos customizados

Edite `newrelic-daemonset-infra.yml`:
```yaml
- name: NRIA_CUSTOM_ATTRIBUTES
  value: '{"environment":"production","region":"us-east-1","team":"seu-time"}'
```

### Ajustar recursos

Modifique as seções `resources` em cada manifesto:
```yaml
resources:
  limits:
    memory: 256Mi
  requests:
    cpu: 100m
    memory: 128Mi
```

### Filtrar logs específicos

Edite `newrelic-fluent-bit-configmap.yml` e adicione filtros:
```conf
[FILTER]
    Name    grep
    Match   kube.*
    Exclude kubernetes.namespace_name kube-system
```

## Enriquecimento de Logs da Aplicação

Para melhor rastreabilidade, adicione labels aos seus deployments:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: mecanica-xpto-deployment
  namespace: mecanica-xpto
  labels:
    app: mecanica-xpto-api
    version: "1.0"
  annotations:
    newrelic.com/integrations: "true"
spec:
  template:
    metadata:
      labels:
        app: mecanica-xpto-api
        version: "1.0"
```

## Troubleshooting

### Pods não iniciam

```bash
# Ver eventos
kubectl describe pods -n newrelic

# Verificar recursos do node
kubectl top nodes
```

### Logs não aparecem no New Relic

1. Verifique se o Fluent Bit está rodando:
```bash
kubectl get pods -n newrelic -l app=fluent-bit
```

2. Verifique os logs do Fluent Bit:
```bash
kubectl logs -n newrelic -l app=fluent-bit --tail=100
```

3. Verifique a License Key:
```bash
kubectl get secret -n newrelic newrelic-license-key -o yaml
```

### Erro de permissões

Verifique o RBAC:
```bash
kubectl get clusterrolebinding newrelic
kubectl describe clusterrolebinding newrelic
```

## Desinstalação

```bash
# Remover todos os recursos
kubectl delete -f newrelic-deployment-events.yml
kubectl delete -f newrelic-kube-state-metrics.yml
kubectl delete -f newrelic-daemonset-fluent-bit.yml
kubectl delete -f newrelic-daemonset-infra.yml
kubectl delete -f newrelic-rbac.yml
kubectl delete -f newrelic-events-configmap.yml
kubectl delete -f newrelic-fluent-bit-configmap.yml
kubectl delete -f newrelic-configmap.yml
kubectl delete -f newrelic-secret.yml
kubectl delete -f newrelic-namespace.yml

# Ou simplesmente:
kubectl delete namespace newrelic
```

## Acessando os Dados

1. **Logs**: https://one.newrelic.com/logs
   - Filtre por: `cluster.name = "mecanica-xpto-cluster"`
   - Ou: `k8s.namespaceName = "mecanica-xpto"`

2. **Kubernetes Cluster Explorer**: https://one.newrelic.com/kubernetes
   - Selecione seu cluster
   - Navegue por namespaces, deployments e pods

3. **Dashboards**: https://one.newrelic.com/dashboards

## Recursos Consumidos

Estimativa de recursos por node:
- **newrelic-infra**: ~100m CPU, ~150Mi RAM
- **fluent-bit**: ~100m CPU, ~128Mi RAM

Total aproximado: **200m CPU e 280Mi RAM por node**

## Suporte

- Documentação: https://docs.newrelic.com/docs/kubernetes-pixie/kubernetes-integration/
- GitHub: https://github.com/newrelic/helm-charts
- Fórum: https://forum.newrelic.com/
