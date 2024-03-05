package v1

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v3"
	"github.com/robertt3kuk/xiaoma-test-task/init/logger"
	"github.com/robertt3kuk/xiaoma-test-task/internal/model"
	"github.com/robertt3kuk/xiaoma-test-task/internal/service"
)

type ItemRoutes struct {
	l logger.Interface
	s service.Item
}

func NewItemRoutes(l logger.Interface, s service.Item) *ItemRoutes {
	return &ItemRoutes{l: l, s: s}
}

//swagger:model
type ItemRequest struct {
	ItemName string  `json:"item_name" required:"true"`
	Cost     float64 `json:"cost"      required:"true"`
	Price    float64 `json:"price"     required:"true"`
	Sort     int     `json:"sort"      required:"true"`
}

func (i *ItemRequest) toModel() model.Item {
	return model.Item{
		ItemName:  i.ItemName,
		Cost:      i.Cost,
		Price:     i.Price,
		Sort:      i.Sort,
		DeletedAt: nil,
	}
}

func (i *ItemRequest) validate() error {
	var err string
	if i.ItemName == "" || len(i.ItemName) < 3 {
		err += " item name is invalid or shorter than 3,"
	}
	if i.Cost <= 0 {
		err += " cost is invalid or under zero,"
	}
	if i.Price <= 0 {
		err += " price is invalid or under zero,"
	}
	if i.Sort <= 0 {
		err += " sort is invalid or under zero,"
	}
	if len(err) != 0 {
		return errors.New(err)
	} else {
		return nil
	}
}

func (r *ItemRoutes) Create(c fiber.Ctx) error {
	var item ItemRequest
	if err := c.Bind().JSON(&item); err != nil {
		r.l.Error("ItemRoutes - Create - c.Bind.JSON:%w", err)
		return c.Status(400).JSON(gin.H{"error": "invalid request"})
	}
	err := item.validate()
	if err != nil {
		r.l.Error("ItemRoutes - Create - item.validate:%w", err)
		return c.Status(400).JSON(gin.H{"error": err.Error()})
	}
	itemb := item.toModel()
	result, status := r.s.Create(c.Context(), itemb)
	if !status.Ok() {
		r.l.Error("ItemRoutes - Create - r.s.Create:%w", status.Err)
		return c.Status(status.Code).JSON(gin.H{"error": status.Msg})
	}
	return c.Status(status.Code).JSON(result)
}

func (r *ItemRoutes) Update(c fiber.Ctx) error {
	var item ItemRequest
	if err := c.Bind().JSON(&item); err != nil {
		r.l.Error("ItemRoutes - Update - c.Bind.JSON:%w", err)
		return c.Status(400).JSON(gin.H{"error": "invalid request"})
	}
	idParam := c.Params("id")
	if idParam == "" {
		r.l.Error("ItemRoutes - Update - c.Params.Get:%w", errors.New("missing the id parameter"))
		return c.Status(400).JSON(gin.H{"error": "invalid request missing id in query parameters"})
	}
	idParamInt, err := strconv.Atoi(idParam)
	if err != nil {
		r.l.Error("ItemRoutes - Update - parseInt:%w", err)
		return c.Status(400).JSON(gin.H{"error": "id is invalid integer"})
	}
	itemb := item.toModel()
	itemb.ID = idParamInt
	result, status := r.s.Update(c.Context(), itemb)
	if !status.Ok() {
		r.l.Error("ItemRoutes - Update - r.s.Update:%w", status.Err)
		return c.Status(status.Code).JSON(gin.H{"error": status.Msg})
	}
	return c.Status(status.Code).JSON(result)
}

func (r *ItemRoutes) GetByID(c fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		r.l.Error("ItemRoutes - GetByID - c.Params.Get:%w", errors.New("missing the id parameter"))
		return c.Status(400).JSON(gin.H{"error": "invalid request missing id in query parameters"})
	}
	idParamInt, err := strconv.Atoi(idParam)
	if err != nil {
		r.l.Error("ItemRoutes - GetByID - parseInt:%w", err)
		return c.Status(400).JSON(gin.H{"error": "id is invalid integer"})
	}
	result, status := r.s.GetByID(c.Context(), idParamInt)
	if !status.Ok() {
		r.l.Error("ItemRoutes - GetByID - r.s.GetByID:%w", status.Err)
		return c.Status(status.Code).JSON(gin.H{"error": status.Msg})
	}
	return c.Status(status.Code).JSON(result)
}

func (r *ItemRoutes) GetAll(c fiber.Ctx) error {
	limitF := c.FormValue("limit")
	offsetF := c.FormValue("offest")
	fmt.Println(limitF, offsetF)
	limit, err := strconv.Atoi(limitF)
	if limitF != "" {
		if err != nil {
			r.l.Error("ItemRoutes - GetAll - strconv.Atoi:%w", err)
			return c.Status(400).JSON(gin.H{"error": "limit is invalid integer"})
		}
	}
	offset, err := strconv.Atoi(offsetF)
	if offsetF != "" {
		if err != nil {
			r.l.Error("ItemRoutes - GetAll - strconv.Atoi:%w", err)
			return c.Status(400).JSON(gin.H{"error": "offset is invalid integer"})
		}
	}
	result, status := r.s.GetAll(c.Context(), limit, offset)
	if !status.Ok() {
		r.l.Error("ItemRoutes - GetAll - r.s.GetAll:%w", status.Err)
		return c.Status(status.Code).JSON(gin.H{"error": status.Msg})
	}
	return c.Status(status.Code).JSON(result)
}

func (r *ItemRoutes) Delete(c fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		r.l.Error("ItemRoutes - Delete - c.Params.Get:%w", errors.New("missing the id parameter"))
		return c.Status(400).JSON(gin.H{"error": "invalid request missing id in query parameters"})
	}
	idParamInt, err := strconv.Atoi(idParam)
	if err != nil {
		r.l.Error("ItemRoutes - Delete - parseInt:%w", err)
		return c.Status(400).JSON(gin.H{"error": "id is invalid integer"})
	}
	status := r.s.Delete(c.Context(), idParamInt)
	if !status.Ok() {
		r.l.Error("ItemRoutes - Delete - r.s.Delete:%w", status.Err)
		return c.Status(status.Code).JSON(gin.H{"error": status.Msg})
	}
	return c.Status(status.Code).Send([]byte(status.Msg))
}
