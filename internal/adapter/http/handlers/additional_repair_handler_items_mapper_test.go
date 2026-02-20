package handlers

import (
	"testing"

	"github.com/fiap-grupo95/os-service-api/internal/adapter/http/dto/request"
	"github.com/stretchr/testify/assert"
)

func TestToAdditionalRepairEntityFromItems(t *testing.T) {
	payload := request.AdditionalRepairItemsRequest{
		ServiceOrderID: "so-1",
		Description:    "desc",
		Services:       []request.AdditionalRepairServiceItem{{ID: "s1"}},
		PartsSupplies:  []request.AdditionalRepairPartsSupplyItem{{ID: "p1", Quantity: 2}},
	}

	entity := toAdditionalRepairEntityFromItems(payload)
	assert.Equal(t, "so-1", entity.ServiceOrderID)
	assert.Equal(t, "desc", entity.Description)
	if assert.Len(t, entity.Services, 1) {
		assert.Equal(t, "s1", entity.Services[0].ID)
	}
	if assert.Len(t, entity.PartsSupplies, 1) {
		assert.Equal(t, "p1", entity.PartsSupplies[0].ID)
		assert.Equal(t, 2, entity.PartsSupplies[0].Quantity)
	}
}
