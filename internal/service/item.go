package service

import (
	"context"
	"net/http"

	"github.com/robertt3kuk/xiaoma-test-task/internal/model"
)

type ItemService struct {
	t ItemRepository
}

func NewItemService(t ItemRepository) *ItemService {
	return &ItemService{t: t}
}

func (s *ItemService) Create(ctx context.Context, item model.Item) (int, Status) {
	var status Status
	ID, err := s.t.IDByItemName(ctx, item.ItemName)
	if err != nil {
		return 0, status.withError(
			"ItemService - Create - s.t.ItemNameExists:%w",
			err,
			"error with item name",
			http.StatusInternalServerError,
		)
	}
	if ID != 0 {
		return 0, status.withError(
			"ItemService - Create - s.t.ItemNameExists:%w",
			err,
			"item name already exists",
			http.StatusBadRequest,
		)
	}
	id, err := s.t.Create(ctx, item)
	if err != nil {
		return 0, status.withError(
			"ItemService - Create - s.t.Create:%w",
			err,
			"error with item creation",
			http.StatusInternalServerError,
		)
	}
	return id, status.success("item succesfully created", http.StatusCreated)
}

func (s *ItemService) GetByID(ctx context.Context, id int) (model.Item, Status) {
	var status Status
	var item model.Item
	exist, err := s.t.IDExists(ctx, id)
	if err != nil {
		if !exist {
			return item, status.withError(
				"ItemService - Update -  s.t.IDExists:%w",
				err,
				"item does not exist",
				http.StatusNotFound,
			)
		}
	}
	item, err = s.t.GetByID(ctx, id)
	if err != nil {
		return item, status.withError(
			"ItemService - Update- s.t.GetByID:%w",
			err,
			"couldn't get item",
			http.StatusInternalServerError,
		)
	}
	return item, status.success("item retrieved", http.StatusCreated)
}

func (s *ItemService) GetAll(ctx context.Context, limit, offset int) ([]model.Item, Status) {
	var status Status
	var items []model.Item
	items, err := s.t.GetAll(ctx, limit, offset)
	if err != nil {
		return items, status.withError(
			"ItemService - GetAll - s.t.GetAll:%w",
			err,
			"couldn't get all items",
			http.StatusInternalServerError,
		)
	}
	return items, status.success("items retrieved", http.StatusOK)
}

func (s *ItemService) Update(ctx context.Context, item model.Item) (model.Item, Status) {
	var status Status
	exist, err := s.t.IDExists(ctx, item.ID)
	if err != nil {
		if !exist {
			return item, status.withError(
				"ItemService - Update -  s.t.IDExists:%w",
				err,
				"item does not exist",
				http.StatusNotFound,
			)
		}
	}
	ID, err := s.t.IDByItemName(ctx, item.ItemName)
	if err != nil {
		return item, status.withError(
			"ItemService - Update - s.t.IDByItemName:%w",
			err,
			"couldn't get item id",
			http.StatusInternalServerError,
		)
	}
	if ID != item.ID && ID != 0 {
		// name already in use
		return item, status.withError(
			"ItemService - Update - s.t.IDByItemName:%w",
			err,
			"item name already exists",
			http.StatusBadRequest,
		)
	}

	item, err = s.t.Update(ctx, item)
	if err != nil {
		return item, status.withError(
			"ItemService - Update - s.t.Update:%w",
			err,
			"couldn't update item",
			http.StatusInternalServerError,
		)
	}
	return item, status.success("item updated", http.StatusOK)
}

func (s *ItemService) Delete(ctx context.Context, id int) Status {
	var status Status
	err := s.t.Delete(ctx, id)
	if err != nil {
		return status.withError(
			"ItemService - Delete - s.t.Delete:%w",
			err,
			"couldn't delete item",
			http.StatusInternalServerError,
		)
	}
	return status.success("item deleted", http.StatusOK)
}
