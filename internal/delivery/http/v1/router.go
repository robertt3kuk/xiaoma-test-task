package v1

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"

	log "github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/robertt3kuk/xiaoma-test-task/init/logger"
	"github.com/robertt3kuk/xiaoma-test-task/internal/service"
)

func NewRouter(handler *fiber.App, l logger.Interface, t *service.Service) {
	conf := cors.Config{
		AllowOrigins:     "*", // Equivalent to AllowAllOrigins: true
		AllowMethods:     "POST, PUT, GET, DELETE, FETCH",
		AllowHeaders:     "Origin, Content-type, X-API-Key",
		AllowCredentials: false,
		ExposeHeaders:    "Content-Length",
		MaxAge:           3600,
	}
	handler.Use(cors.New(conf))
	// handler.Use(recover.New())
	handler.Use(log.New(log.Config{
		// For more options, see the Config section
		Format: "${pid} ${locals:requestid} ${status} - ${method} ${path}​\n",
	}))

	h := handler.Group("/v1")
	itemRoutes := NewItemRoutes(l, t.Item)
	customerRoutes := NewCustomerRoutes(l, t.Customer)
	transactionRoutes := NewTransactionRoutes(l, t.Transaction)
	h.Get(
		"/healthz",
		func(c fiber.Ctx) error { return c.Status(http.StatusOK).SendString("up and running") },
	)
	items := h.Group("/item")
	items.Post("", itemRoutes.Create)
	items.Put("/:id", itemRoutes.Update)
	items.Get("/:id", itemRoutes.GetByID)
	items.Get("", itemRoutes.GetAll)
	items.Delete("/:id", itemRoutes.Delete)

	customers := h.Group("/customer")
	customers.Post("", customerRoutes.Create)
	customers.Put("/:id", customerRoutes.Update)
	customers.Get("/:id", customerRoutes.GetByID)
	customers.Get("", customerRoutes.GetAll)
	customers.Delete("/:id", customerRoutes.Delete)
	// catch erros

	transactions := h.Group("/transaction")
	transactions.Post("", transactionRoutes.Create)
	transactions.Put("/:id", transactionRoutes.Update)
	transactions.Get("/:id", transactionRoutes.GetByID)
	transactions.Get("", transactionRoutes.GetAll)
	transactions.Delete("/:id", transactionRoutes.Delete)

	transactionsView := h.Group("/transaction-view")
	h.Get("/transaction-view-filter", transactionRoutes.GetAllTransactionViewByFilters)

	transactionsView.Get("/:id", transactionRoutes.GetTransactionViewByID)
	transactionsView.Get("", transactionRoutes.GetAllTransactionView)
}
