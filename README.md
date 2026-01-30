# OS-Service-API

## Operações do OS service

### Criação da OS: /os
	- Validação/Cadastro de cliente e veículo
	- Registro de OS: Recebida

### Cancelamento de OS: /os/{os_id}/cancel
	- Mudança de status: Cancelada

### Processo de diagnóstico: /os/{os_id}/diagnosis
	- Mudança de status: Em Diagnóstico, apenas
	OU
	Patch
	- Mudança de status: Em Diagnóstico
	- Validação de serviços e peças e insumos
	- Operações de reserva de peças do estoque
	Post
	- Solicitação de cálculo de orçamento
	- Mudança de status: Aguardando aprovação

### Processo de Avaliação de orçamento: /os/{os_id}/estimate
    1. Aprovação do orçamento
        - Solicitar aprovação de orçamento
        - Mudança de status: Aprovado
        - Solicitar baixa de estoque (confirma uso de peças)
        OU

    2. Rejeição do orçamento
        - Solicitar rejeição do orçamento
        - Solicitar liberação da reserva de peças e insumos no estoque
        - Mudança de status: Cancelado
        OU

    3. Solicitação de alterações no orçamento
        - Solicitar cancelamento do orçamento
        - Solicitar liberação da reserva de peças e insumos no estoque
        - Mudança de status: Em diagnóstico

### Processo de inicio de execução: /os/{os_id}/execution/create
    - Solicita a criação de uma nova execution
    - Mudança de status: Em Execução

### Proceso de finalização da execução: /os/{os_id}/execution/finish
    - Solicita finalização da execution associada a OS
    - Mudança de status: Finalizada
    - Notifica usuário que a OS foi finalizada

### Processo de pagamento: /os/{os_id}/payment
    - Solicita criação e processamento de pagamento associado a uma OS

### Processo de entrega: /os/{os_id}/delivery
    - Validação da existencia de um pagamento associado a OS
    - Mudança de status: Entregue

### Consulta de status e histórico