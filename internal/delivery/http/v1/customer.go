package v1

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v3"
	"github.com/robertt3kuk/xiaoma-test-task/init/logger"
	"github.com/robertt3kuk/xiaoma-test-task/internal/model"
	"github.com/robertt3kuk/xiaoma-test-task/internal/service"
)

type CustomerRoutes struct {
	l logger.Interface
	s service.Customer
}

func NewCustomerRoutes(l logger.Interface, s service.Customer) *CustomerRoutes {
	return &CustomerRoutes{l: l, s: s}
}

type CustomerRequest struct {
	Name    string  `json:"customer_name"`
	Balance float64 `json:"balance"`
}

func (c *CustomerRequest) toModel() model.Customer {
	return model.Customer{
		Name:    c.Name,
		Balance: c.Balance,
	}
}

func (c *CustomerRequest) validate() error {
	var err string
	if c.Name == "" || len(c.Name) < 5 {
		err += " name is invalid or shorter than 5,"
	}
	if c.Balance <= 0 {
		err += " balance is invalid or equal or less than zero"
	}
	if len(err) > 0 {
		return errors.New(err)
	} else {
		return nil
	}

}

func (r *CustomerRoutes) Create(c fiber.Ctx) error {
	var requestBody CustomerRequest
	err := c.Bind().JSON(&requestBody)
	if err != nil {
		r.l.Error("CustomerRoutes - Create - c.Bind.JSON:%w", err)
		return c.Status(400).JSON(gin.H{"error": "invalid request"})
	}
	err = requestBody.validate()
	if err != nil {
		r.l.Error("CustomerRoutes - Create - requestBody.validate:%w", err)
		return c.Status(400).JSON(gin.H{"error": err.Error()})
	}
	customer := requestBody.toModel()
	result, status := r.s.Create(c.Context(), customer)
	if !status.Ok() {
		r.l.Error("CustomerRoutes - Create - r.s.Create:%w", status.Err)
		return c.Status(status.Code).JSON(gin.H{"error": status.Msg})
	}
	return c.Status(status.Code).JSON(result)

}

func (r *CustomerRoutes) Update(c fiber.Ctx) error {
	var customer CustomerRequest
	err := c.Bind().JSON(&customer)
	if err != nil {
		r.l.Error("CustomerRoutes - Update - c.Bind.JSON:%w", err)
		return c.Status(400).JSON(gin.H{"error": "invalid request"})
	}
	idParam := c.Params("id")
	if idParam == "" {
		r.l.Error("CustomerRoutes - Update - c.Params.Get:%w", errors.New("missing the id parameter"))
		return c.Status(400).JSON(gin.H{"error": "invalid request missing id in query parameters"})
	}
	idParamInt, err := strconv.Atoi(idParam)
	if err != nil {
		r.l.Error("CustomerRoutes - Update - parseInt:%w", err)
		return c.Status(400).JSON(gin.H{"error": "id is invalid integer"})
	}
	customerBody := customer.toModel()
	customerBody.ID = idParamInt
	result, status := r.s.Update(c.Context(), customerBody)
	if !status.Ok() {
		r.l.Error("CustomerRoutes - Update - r.s.Update:%w", status.Err)
		return c.Status(status.Code).JSON(gin.H{"error": status.Msg})
	}
	return c.Status(status.Code).JSON(result)
}

func (r *CustomerRoutes) GetByID(c fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		r.l.Error("CustomerRoutes - GetByID - c.Params.Get:%w", errors.New("missing the id parameter"))
		return c.Status(400).JSON(gin.H{"error": "invalid request missing id in query parameters"})
	}
	idParamInt, err := strconv.Atoi(idParam)
	if err != nil {
		r.l.Error("CustomerRoutes - GetByID - parseInt:%w", err)
		return c.Status(400).JSON(gin.H{"error": "id is invalid integer"})
	}
	result, status := r.s.GetByID(c.Context(), idParamInt)
	if !status.Ok() {
		r.l.Error("CustomerRoutes - GetByID - r.s.GetByID:%w", status.Err)
		return c.Status(status.Code).JSON(gin.H{"error": status.Msg})
	}
	return c.Status(status.Code).JSON(result)
}

func (r *CustomerRoutes) GetAll(c fiber.Ctx) error {
	limitF := c.FormValue("limit")
	offsetF := c.FormValue("offest")
	limit, err := strconv.Atoi(limitF)
	if err != nil {
		r.l.Error("CustomerRoutes - GetAll - strconv.Atoi:%w", err)
		return c.Status(400).JSON(gin.H{"error": "limit is invalid integer"})
	}
	offset, err := strconv.Atoi(offsetF)
	if err != nil {
		r.l.Error("CustomerRoutes - GetAll - strconv.Atoi:%w", err)
		return c.Status(400).JSON(gin.H{"error": "offset is invalid integer"})
	}
	result, status := r.s.GetAll(c.Context(), limit, offset)
	if !status.Ok() {
		r.l.Error("CustomerRoutes - GetAll - r.s.GetAll:%w", status.Err)
		return c.Status(status.Code).JSON(gin.H{"error": status.Msg})
	}
	return c.Status(status.Code).JSON(result)
}

func (r *CustomerRoutes) Delete(c fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		r.l.Error("CustomerRoutes - Delete - c.Params.Get:%w", errors.New("missing the id parameter"))
		return c.Status(400).JSON(gin.H{"error": "invalid request missing id in query parameters"})
	}
	idParamInt, err := strconv.Atoi(idParam)
	if err != nil {
		r.l.Error("CustomerRoutes - Delete - parseInt:%w", err)
		return c.Status(400).JSON(gin.H{"error": "id is invalid integer"})
	}
	status := r.s.Delete(c.Context(), idParamInt)
	if !status.Ok() {
		r.l.Error("CustomerRoutes - Delete - r.s.Delete:%w", status.Err)
		return c.Status(status.Code).JSON(gin.H{"error": status.Msg})
	}
	return c.Status(status.Code).Send([]byte(status.Msg))
}
