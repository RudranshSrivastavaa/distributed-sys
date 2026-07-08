package service

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/rudransh/distributed-commerce/inventory/internal/dto"
	inventory_error "github.com/rudransh/distributed-commerce/inventory/internal/errors"
	"github.com/rudransh/distributed-commerce/inventory/internal/mapper"
	"github.com/rudransh/distributed-commerce/inventory/internal/model"
	"github.com/rudransh/distributed-commerce/inventory/internal/repository"
	"github.com/rudransh/distributed-commerce/inventory/internal/state"

	"gorm.io/gorm"
)

type inventoryService struct {
	productRepository     repository.ProductRepository
	inventoryRepository   repository.InventoryRepository
	reservationRepository repository.ReservationRepository
	txManager             *repository.TransactionManager
}

func NewInventoryService(
	productRepo repository.ProductRepository,
	inventoryRepo repository.InventoryRepository,
	reservationRepo repository.ReservationRepository,
	txManager *repository.TransactionManager,
) InventoryService {

	return &inventoryService{
		productRepository:     productRepo,
		inventoryRepository:   inventoryRepo,
		reservationRepository: reservationRepo,
		txManager:             txManager,
	}
}

func (s *inventoryService) CreateProduct(
	request dto.CreateProductRequest,
) (dto.ProductResponse, error) {

	if err := validateCreateProduct(request); err != nil {
		return dto.ProductResponse{}, err
	}

	if _, err := s.productRepository.FindBySKU(request.SKU); err == nil {
		return dto.ProductResponse{}, errors.New("sku already exists")
	}

	product := mapper.ToProduct(request)

	err := s.txManager.Execute(func(tx *gorm.DB) error {

		productRepo := repository.NewProductRepository(tx)
		inventoryRepo := repository.NewInventoryRepository(tx)

		if err := productRepo.Create(product); err != nil {
			return err
		}

		inventory := &model.Inventory{
			ProductID:         product.ID,
			AvailableQuantity: 0,
			ReservedQuantity:  0,
			Version:           1,
		}

		if err := inventoryRepo.Create(inventory); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return dto.ProductResponse{}, err
	}

	return mapper.ToProductResponse(product), nil
}

func (s *inventoryService) AddStock(productID uuid.UUID, request dto.AddStockRequest) (dto.InventoryResponse, error) {

	if err := validateStockQuantity(
		request.Quantity,
	); err != nil {

		return dto.InventoryResponse{}, err
	}

	inventory, err := s.inventoryRepository.
		FindByProductID(productID)

	if err != nil {
		return dto.InventoryResponse{}, err
	}

	inventory.AvailableQuantity += request.Quantity

	if err := s.inventoryRepository.
		Update(inventory); err != nil {

		return dto.InventoryResponse{}, err
	}

	return mapper.ToInventoryResponse(
		inventory,
	), nil

}

func (s *inventoryService) Reserve(
	request dto.CreateReservationRequest,
) (dto.ReservationResponse, error) {

	if request.Quantity <= 0 {
		return dto.ReservationResponse{},
			errors.New("quantity must be greater than zero")
	}

	var reservation *model.Reservation

	const maxRetries = 3

	var err error

	for attempt := 0; attempt < maxRetries; attempt++ {
		
		err = s.txManager.Execute(func(tx *gorm.DB) error {

			inventoryRepo := repository.NewInventoryRepository(tx)
			reservationRepo := repository.NewReservationRepository(tx)

			inventory, err := inventoryRepo.FindByProductID(
				request.ProductID,
			)
			if err != nil {
				return err
			}

			if err := inventory.Reserve(request.Quantity); err != nil {
				return err
			}

			if err := inventoryRepo.UpdateWithVersion(inventory); err != nil {
				return err
			}

			reservation = &model.Reservation{
				OrderID:   request.OrderID,
				ProductID: request.ProductID,
				Quantity:  request.Quantity,
				Status:    model.StatusReserved,
				ExpiresAt: time.Now().Add(15 * time.Minute),
			}

			if err := reservationRepo.Create(reservation); err != nil {
				return err
			}

			return nil
		})

		if err == nil {
			break
		}

		if !errors.Is(err, inventory_error.ErrOptimisticLockConflict) {
			break
		}
	}


	if err != nil {
		return dto.ReservationResponse{}, err
	}

	return mapper.ToReservationResponse(
		reservation,
	), nil

}

func (s *inventoryService) Release(
	reservationID uuid.UUID,
) (dto.ReservationResponse, error) {

	var reservation *model.Reservation

	err := s.txManager.Execute(func(tx *gorm.DB) error {

		reservationRepo := repository.NewReservationRepository(tx)

		inventoryRepo := repository.NewInventoryRepository(tx)

		var err error

		reservation, err = reservationRepo.FindByID(reservationID)

		if err != nil {
			return err
		}

		inventory, err := inventoryRepo.FindByProductID(
			reservation.ProductID,
		)

		if err != nil {
			return err
		}

		if err := inventory.Release(reservation.Quantity); err != nil {
			return err
		}

		if err := state.Transition(

			reservation,

			model.StatusReleased,
		); err != nil {

			return err

		}

		if err := inventoryRepo.UpdateWithVersion(inventory); err != nil {
			return err
		}

		if err := reservationRepo.Update(reservation); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return dto.ReservationResponse{}, err
	}

	return mapper.ToReservationResponse(reservation), nil
}

func (s *inventoryService) Confirm(
	reservationID uuid.UUID,
) (dto.ReservationResponse, error) {

	var reservation *model.Reservation

	err := s.txManager.Execute(func(tx *gorm.DB) error {

		reservationRepo := repository.NewReservationRepository(tx)

		inventoryRepo := repository.NewInventoryRepository(tx)

		var err error

		reservation, err = reservationRepo.FindByID(reservationID)

		if err != nil {
			return err
		}

		inventory, err := inventoryRepo.FindByProductID(
			reservation.ProductID,
		)

		if err != nil {
			return err
		}

		if err := inventory.Confirm(
			reservation.Quantity,
		); err != nil {
			return err
		}

		if err := state.Transition(

			reservation,

			model.StatusConfirmed,
		); err != nil {

			return err

		}

		if err := inventoryRepo.UpdateWithVersion(inventory); err != nil {
			return err
		}

		if err := reservationRepo.Update(reservation); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return dto.ReservationResponse{}, err
	}

	return mapper.ToReservationResponse(reservation), nil
}


func (s *inventoryService) GetExpiredReservations() (
	[]model.Reservation,
	error,
) {

	return s.reservationRepository.
		FindExpiredReservations(
			time.Now(),
		)

}

func (s *inventoryService) ExpireReservation(
	reservationID uuid.UUID,
) error {

	return s.txManager.Execute(func(tx *gorm.DB) error {

		reservationRepo := repository.NewReservationRepository(tx)

		inventoryRepo := repository.NewInventoryRepository(tx)

		reservation, err := reservationRepo.FindByID(reservationID)
		if err != nil {
			return err
		}

		if reservation.Status != model.StatusReserved {
			return nil
		}

		inventory, err := inventoryRepo.FindByProductID(
			reservation.ProductID,
		)
		if err != nil {
			return err
		}

		if err := inventory.Release(
			reservation.Quantity,
		); err != nil {
			return err
		}

		if err := state.Transition(
			reservation,
			model.StatusExpired,
		); err != nil {
			return err
		}

		if err := inventoryRepo.UpdateWithVersion(
			inventory,
		); err != nil {
			return err
		}

		if err := reservationRepo.Update(
			reservation,
		); err != nil {
			return err
		}

		return nil
	})
}