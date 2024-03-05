package v1

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v3"
	"github.com/robertt3kuk/xiaoma-test-task/init/logger"
	"github.com/robertt3kuk/xiaoma-test-task/internal/model"
	"github.com/robertt3kuk/xiaoma-test-task/internal/service"
)

type TransactionRoutes struct {
	l logger.Interface
	s service.Transaction
}

func NewTransactionRoutes(l logger.Interface, s service.Transaction) *TransactionRoutes {
	return &TransactionRoutes{l: l, s: s}
}

type TransactionRequest struct {
	CustomerID int     `json:"customer_id"`
	ItemID     int     `json:"item_id"`
	Qty        int     `json:"qty"`
	Price      float64 `json:"price"`
}

func (t *TransactionRequest) toModel() model.Transaction {
	return model.Transaction{
		CustomerID: t.CustomerID,
		ItemID:     t.ItemID,
		Qty:        t.Qty,
		Price:      t.Price,
		Amount:     t.Price * float64(t.Qty),
		DeletedAt:  nil,
	}
}

func (t *TransactionRequest) validate() error {
	// i need to validate the struct
	var err string
	if t.CustomerID < 1 {
		err += " customer id is invalid,"
	}
	if t.ItemID < 1 {
		err += " item id is invalid,"
	}
	if t.Qty < 1 {
		err += " qty is invalid,"
	}
	if t.Price <= 0 {
		err += " price is invalid or under zero,"
	}
	if len(err) != 0 {
		return errors.New(err)
	} else {
		return nil
	}
}

func (r *TransactionRoutes) Create(c fiber.Ctx) error {
	var transaction TransactionRequest
	if err := c.Bind().JSON(&transaction); err != nil {
		r.l.Error("TransactionRoutes - Create - c.Bind.JSON:%w", err)
		return c.Status(http.StatusBadRequest).JSON(gin.H{"error": "invalid request"})
	}
	id, status := r.s.Create(c.Context(), transaction.toModel())
	if !status.Ok() {
		r.l.Error("TransactionRoutes - Create - r.s.Create:%w", status.Err)
		return c.Status(status.Code).JSON(gin.H{"error": status.Msg})
	}
	return c.Status(status.Code).JSON(gin.H{"id": id})
}

func (r *TransactionRoutes) Update(c fiber.Ctx) error {
	var transaction TransactionRequest
	if err := c.Bind().JSON(&transaction); err != nil {
		r.l.Error("TransactionRoutes - Update - ctx.ShouldBindJSON:%w", err)
		return c.Status(http.StatusBadRequest).JSON(gin.H{"error": "invalid request"})
	}
	idParam := c.Params("id")
	idParamInt, err := strconv.Atoi(idParam)
	if err != nil {
		r.l.Error("TransactionRoutes - Update - parseInt:%w", err)
		return c.Status(http.StatusBadRequest).JSON(gin.H{"error": "id is invalid integer"})
	}
	transactionModel := transaction.toModel()
	transactionModel.ID = idParamInt
	result, status := r.s.Update(c.Context(), transactionModel)
	if !status.Ok() {
		r.l.Error("TransactionRoutes - Update - r.s.Update:%w", status.Err)
		return c.Status(status.Code).JSON(gin.H{"error": status.Msg})
	}
	return c.Status(status.Code).JSON(result)
}

func (r *TransactionRoutes) GetByID(c fiber.Ctx) error {
	idParam := c.Params("id")
	idParamInt, err := strconv.Atoi(idParam)
	if err != nil {
		r.l.Error("TransactionRoutes - GetByID - parseInt:%w", err)
		return c.Status(http.StatusBadRequest).JSON(gin.H{"error": "id is invalid integer"})
	}
	result, status := r.s.GetByID(c.Context(), idParamInt)
	if !status.Ok() {
		r.l.Error("TransactionRoutes - GetByID - r.s.GetByID:%w", status.Err)
		return c.Status(status.Code).JSON(gin.H{"error": status.Msg})
	}
	return c.Status(status.Code).JSON(result)
}

func (r *TransactionRoutes) GetAll(c fiber.Ctx) error {
	limitF := c.FormValue("limit")
	offsetF := c.FormValue("offset")
	limit, err := strconv.Atoi(limitF)
	if limitF != "" {
		if err != nil {
			r.l.Error("TransactionRoutes - GetAll - strconv.Atoi:%w", err)
			return c.Status(http.StatusBadRequest).JSON(gin.H{"error": "limit is invalid integer"})
		}
	}
	offset, err := strconv.Atoi(offsetF)
	if offsetF != "" {
		if err != nil {
			r.l.Error("TransactionRoutes - GetAll - strconv.Atoi:%w", err)
			return c.Status(http.StatusBadRequest).JSON(gin.H{"error": "offset is invalid integer"})
		}
	}
	result, status := r.s.GetAll(c.Context(), limit, offset)
	if !status.Ok() {
		r.l.Error("TransactionRoutes - GetAll - r.s.GetAll:%w", status.Err)
		return c.Status(status.Code).JSON(gin.H{"error": status.Msg})
	}
	return c.Status(status.Code).JSON(result)
}

func (r *TransactionRoutes) Delete(c fiber.Ctx) error {
	idParam := c.Params("id")
	idParamInt, err := strconv.Atoi(idParam)
	if err != nil {
		r.l.Error("TransactionRoutes - Delete - parseInt:%w", err)
		return c.Status(http.StatusBadRequest).JSON(gin.H{"error": "id is invalid integer"})
	}
	status := r.s.Delete(c.Context(), idParamInt)
	if !status.Ok() {
		r.l.Error("TransactionRoutes - Delete - r.s.Delete:%w", status.Err)
		return c.Status(status.Code).JSON(gin.H{"error": status.Msg})
	}
	return c.Status(status.Code).Send([]byte(status.Msg))
}

func (r *TransactionRoutes) GetTransactionViewByID(c fiber.Ctx) error {
	idParam := c.Params("id")
	idParamInt, err := strconv.Atoi(idParam)
	if err != nil {
		r.l.Error("TransactionRoutes - GetTransactionViewByID - parseInt:%w", err)
		return c.Status(http.StatusBadRequest).JSON(gin.H{"error": "id is invalid integer"})
	}
	result, status := r.s.GetByTransactionID(c.Context(), idParamInt)
	if !status.Ok() {
		r.l.Error(
			"TransactionRoutes - GetTransactionViewByID - r.s.GetByTransactionID:%w",
			status.Err,
		)
		return c.Status(status.Code).JSON(gin.H{"error": status.Msg})
	}
	return c.Status(status.Code).JSON(result)
}

func (r *TransactionRoutes) GetAllTransactionView(c fiber.Ctx) error {
	limitF := c.FormValue("limit")
	offsetF := c.FormValue("offset")
	limit := 0
	offset := 0
	var err error
	if limitF != "" {
		limit, err = strconv.Atoi(limitF)
		if err != nil {
			r.l.Error("TransactionRoutes - GetAllTransactionView - strconv.Atoi:%w", err)
			return c.Status(http.StatusBadRequest).JSON(gin.H{"error": "limit is invalid integer"})
		}
	}
	if offsetF != "" {
		offset, err = strconv.Atoi(offsetF)
		if err != nil {
			r.l.Error("TransactionRoutes - GetAllTransactionView - strconv.Atoi:%w", err)
			return c.Status(http.StatusBadRequest).JSON(gin.H{"error": "offset is invalid integer"})
		}
	}
	result, status := r.s.GetAllTransactionViews(c.Context(), limit, offset)
	if !status.Ok() {
		r.l.Error(
			"TransactionRoutes - GetAllTransactionView - r.s.GetAllTransactionViews:%w",
			status.Err,
		)
		return c.Status(status.Code).JSON(gin.H{"error": status.Msg})
	}
	return c.Status(status.Code).JSON(result)
}

func (r *TransactionRoutes) GetAllTransactionViewByFilters(c fiber.Ctx) error {
	filter := model.TransactionFilter{}
	err := c.Bind().JSON(&filter)
	if err != nil {
		r.l.Error("TransactionRoutes - GetAllTransactionViewByFilters - c.Bind:%w", err)
		return c.Status(http.StatusBadRequest).JSON(gin.H{"error": "invalid request"})
	}
	result, status := r.s.GetAllTransactionViewsByFilters(c.Context(), &filter)
	if !status.Ok() {
		r.l.Error(
			"TransactionRoutes - GetAllTransactionViewByFilters - r.s.GetAllTransactionViewsByFilters:%w",
			status.Err,
		)
		return c.Status(status.Code).JSON(gin.H{"error": status.Msg})
	}
	return c.Status(status.Code).JSON(result)
}
