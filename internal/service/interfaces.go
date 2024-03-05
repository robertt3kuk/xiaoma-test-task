package service

import (
	"context"

	"github.com/robertt3kuk/xiaoma-test-task/init/postgres"
	"github.com/robertt3kuk/xiaoma-test-task/internal/model"
	"github.com/robertt3kuk/xiaoma-test-task/internal/service/repo/postgresSQL"
)

type Service struct {
	Item
	Customer
	Transaction
}

type Repo struct {
	ItemRepository
	CustomerRepository
	TransactionRepository
}

func New(repo *Repo) *Service {
	return &Service{
		Item:     NewItemService(repo.ItemRepository),
		Customer: NewCustomerService(repo.CustomerRepository),
		Transaction: NewTransactionService(
			repo.TransactionRepository,
			repo.CustomerRepository,
			repo.ItemRepository,
		),
	}
}

func NewRepo(pg *postgres.Postgres) *Repo {
	return &Repo{
		ItemRepository:        postgresSQL.NewItemPostgres(pg),
		CustomerRepository:    postgresSQL.NewCustomerPostgres(pg),
		TransactionRepository: postgresSQL.NewTransactionPostgres(pg),
	}
}

type Item interface {
	Create(ctx context.Context, item model.Item) (int, Status)
	GetByID(ctx context.Context, id int) (model.Item, Status)
	GetAll(ctx context.Context, limit, offset int) ([]model.Item, Status)
	Update(ctx context.Context, item model.Item) (model.Item, Status)
	Delete(ctx context.Context, id int) Status
}

type Customer interface {
	Create(ctx context.Context, customer model.Customer) (int, Status)
	GetByID(ctx context.Context, id int) (model.Customer, Status)
	GetAll(ctx context.Context, limit, offset int) ([]model.Customer, Status)
	Update(ctx context.Context, customer model.Customer) (model.Customer, Status)
	Delete(ctx context.Context, id int) Status
}

type Transaction interface {
	Create(ctx context.Context, transaction model.Transaction) (int, Status)
	GetByID(ctx context.Context, id int) (model.Transaction, Status)
	GetAll(ctx context.Context, limit, offset int) ([]model.Transaction, Status)
	Update(ctx context.Context, transaction model.Transaction) (model.Transaction, Status)
	Delete(ctx context.Context, id int) Status
	GetAllTransactionViews(ctx context.Context, limit, offset int) ([]model.TransactionView, Status)
	GetByTransactionID(ctx context.Context, id int) (model.TransactionView, Status)
	GetAllTransactionViewsByFilters(ctx context.Context, filter *model.TransactionFilter) ([]model.TransactionView, Status)
}

type ItemRepository interface {
	Create(ctx context.Context, item model.Item) (int, error)
	IDExists(ctx context.Context, id int) (bool, error)
	IDByItemName(ctx context.Context, ItemName string) (int, error)
	GetByID(ctx context.Context, id int) (model.Item, error)
	GetAll(ctx context.Context, limit, offset int) ([]model.Item, error)
	Update(ctx context.Context, item model.Item) (model.Item, error)
	Delete(ctx context.Context, id int) error
}

type CustomerRepository interface {
	Create(ctx context.Context, customer model.Customer) (int, error)
	IDExists(ctx context.Context, id int) (bool, error)
	IDByName(ctx context.Context, name string) (int, error)
	GetBalance(ctx context.Context, id int) (float64, error)
	GetByID(ctx context.Context, id int) (model.Customer, error)
	GetAll(ctx context.Context, limit, offset int) ([]model.Customer, error)
	Update(ctx context.Context, customer model.Customer) (model.Customer, error)
	Delete(ctx context.Context, id int) error
}

type TransactionRepository interface {
	Create(ctx context.Context, transaction model.Transaction) (int, error)
	IDExists(ctx context.Context, id int) (bool, error)
	GetByID(ctx context.Context, id int) (model.Transaction, error)
	GetAll(ctx context.Context, limit, offset int) ([]model.Transaction, error)
	Update(ctx context.Context, transaction model.Transaction) (model.Transaction, error)
	Delete(ctx context.Context, id int) error
	GetAllTransactionViews(ctx context.Context, limit, offset int) ([]model.TransactionView, error)
	GetByTransactionID(ctx context.Context, id int) (model.TransactionView, error)
	GetAllTransactionViewsByFilters(ctx context.Context, filter *model.TransactionFilter) ([]model.TransactionView, error)
}
