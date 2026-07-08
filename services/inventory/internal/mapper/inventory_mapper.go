package mapper

import (
	"github.com/rudransh/distributed-commerce/inventory/internal/dto"
	"github.com/rudransh/distributed-commerce/inventory/internal/model"
)

func ToInventoryResponse(
	inventory *model.Inventory,
) dto.InventoryResponse {

	return dto.InventoryResponse{
		ProductID: inventory.ProductID,

		AvailableQuantity: inventory.AvailableQuantity,

		ReservedQuantity: inventory.ReservedQuantity,

		Version: inventory.Version,
	}

}