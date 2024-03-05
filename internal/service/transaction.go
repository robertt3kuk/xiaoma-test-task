package service

import (
	"context"
	"net/http"

	"github.com/robertt3kuk/xiaoma-test-task/internal/model"
)

type TransactionService struct {
	t TransactionRepository
	c CustomerRepository
	i ItemRepository
}

func NewTransactionService(
	t TransactionRepository,
	c CustomerRepository,
	i ItemRepository,
) *TransactionService {
	return &TransactionService{
		t: t,
		c: c,
		i: i,
	}
}

func (s *TransactionService) Create(
	ctx context.Context,
	transaction model.Transaction,
) (int, Status) {
	var status Status
	ItemIDExistss, err := s.i.IDExists(ctx, transaction.ItemID)
	if err != nil {
		return 0, status.withError(
			"TransactionService - Create - s.i.IDExists:%w",
			err,
			"error with item id",
			http.StatusInternalServerError,
		)
	}
	if !ItemIDExistss {
		return 0, status.withError(
			"TransactionService - Create - s.i.IDExists:%w",
			err,
			"item id does not exist",
			http.StatusNotFound,
		)
	}
	CustomerIDExistss, err := s.c.IDExists(ctx, transaction.CustomerID)
	if err != nil {
		return 0, status.withError(
			"TransactionService - Create - s.c.IDExists:%w",
			err,
			"error with customer id",
			http.StatusInternalServerError,
		)
	}
	if !CustomerIDExistss {
		return 0, status.withError(
			"TransactionService - Create - s.c.IDExists:%w",
			err,
			"customer id does not exist",
			http.StatusNotFound,
		)
	}

	balance, err := s.c.GetBalance(ctx, transaction.CustomerID)
	if err != nil {
		return 0, status.withError(
			"TransactionService - Create - s.c.GetBalance:%w",
			err,
			"error with customer balance",
			http.StatusInternalServerError,
		)
	}
	if balance < transaction.Amount {
		return 0, status.withError(
			"TransactionService - Create - s.c.GetBalance:%w",
			err,
			"customer balance is not enough",
			http.StatusBadRequest,
		)
	}

	id, err := s.t.Create(ctx, transaction)
	if err != nil {
		return 0, status.withError(
			"TransactionService - Create - s.t.Create:%w",
			err,
			"error with transaction creation",
			http.StatusInternalServerError,
		)
	}
	return id, status.success("transaction succesfully created", http.StatusCreated)
}

func (s *TransactionService) GetByID(ctx context.Context, id int) (model.Transaction, Status) {
	var status Status
	var transaction model.Transaction
	exist, err := s.t.IDExists(ctx, id)
	if err != nil {
		if !exist {
			return transaction, status.withError(
				"TransactionService - Update -  s.t.IDExists:%w",
				err,
				"transaction does not exist",
				http.StatusNotFound,
			)
		}
	}
	transaction, err = s.t.GetByID(ctx, id)
	if err != nil {
		return transaction, status.withError(
			"TransactionService - Update- s.t.GetByID:%w",
			err,
			"couldn't get transaction",
			http.StatusInternalServerError,
		)
	}
	return transaction, status.success("transaction retrieved", http.StatusCreated)
}

func (s *TransactionService) GetAll(
	ctx context.Context,
	limit, offset int,
) ([]model.Transaction, Status) {
	var status Status
	var transactions []model.Transaction
	transactions, err := s.t.GetAll(ctx, limit, offset)
	if err != nil {
		return transactions, status.withError(
			"TransactionService - GetAll - s.t.GetAll:%w",
			err,
			"couldn't get all transactions",
			http.StatusInternalServerError,
		)
	}
	return transactions, status.success("transactions retrieved", http.StatusOK)
}

func (s *TransactionService) Update(
	ctx context.Context,
	transaction model.Transaction,
) (model.Transaction, Status) {
	var status Status
	exist, err := s.t.IDExists(ctx, transaction.ID)
	if err != nil {
		if !exist {
			return transaction, status.withError(
				"TransactionService - Update -  s.t.IDExists:%w",
				err,
				"transaction does not exist",
				http.StatusNotFound,
			)
		}
	}
	balance, err := s.c.GetBalance(ctx, transaction.CustomerID)
	if err != nil {
		return transaction, status.withError(
			"TransactionService - Update - s.c.GetBalance:%w",
			err,
			"error with customer balance",
			http.StatusInternalServerError,
		)
	}
	if balance < transaction.Amount {
		return transaction, status.withError(
			"TransactionService - Update - s.c.GetBalance:%w",
			err,
			"customer balance is not enough",
			http.StatusBadRequest,
		)
	}
	transaction, err = s.t.Update(ctx, transaction)
	if err != nil {
		return transaction, status.withError(
			"TransactionService - Update - s.t.Update:%w",
			err,
			"couldn't update transaction",
			http.StatusInternalServerError,
		)
	}
	return transaction, status.success("transaction updated", http.StatusOK)
}

func (s *TransactionService) Delete(ctx context.Context, id int) Status {
	var status Status
	err := s.t.Delete(ctx, id)
	if err != nil {
		return status.withError(
			"TransactionService - Delete - s.t.Delete:%w",
			err,
			"couldn't delete transaction",
			http.StatusInternalServerError,
		)
	}
	return status.success("transaction deleted", http.StatusOK)
}

func (s *TransactionService) GetAllTransactionViews(
	ctx context.Context,
	limit, offset int,
) ([]model.TransactionView, Status) {
	var status Status
	var transactions []model.TransactionView
	transactions, err := s.t.GetAllTransactionViews(ctx, limit, offset)
	if err != nil {
		return transactions, status.withError(
			"TransactionService - GetAllTransactionViews - s.t.GetAllTransactionViews:%w",
			err,
			"couldn't get all transactions",
			http.StatusInternalServerError,
		)
	}
	return transactions, status.success("transactions retrieved", http.StatusOK)
}

func (s *TransactionService) GetByTransactionID(
	ctx context.Context,
	id int,
) (model.TransactionView, Status) {
	var status Status
	var transaction model.TransactionView
	exist, err := s.t.IDExists(ctx, id)
	if err != nil {
		if !exist {
			return transaction, status.withError(
				"TransactionService - GetByTransactionID -  s.t.IDExists:%w",
				err,
				"transaction does not exist",
				http.StatusNotFound,
			)
		}
	}
	transaction, err = s.t.GetByTransactionID(ctx, id)
	if err != nil {
		return transaction, status.withError(
			"TransactionService - GetByTransactionID - s.t.GetByTransactionID:%w",
			err,
			"couldn't get transaction",
			http.StatusInternalServerError,
		)
	}
	return transaction, status.success("transaction retrieved", http.StatusOK)
}

func (s *TransactionService) GetAllTransactionViewsByFilters(
	ctx context.Context,
	filter *model.TransactionFilter,
) ([]model.TransactionView, Status) {
	var status Status
	// check all fields of transaction filter if all are empty return error
	if filter.CustomerName == "" && filter.ItemName == "" && filter.ID == 0 {
		return nil, status.withError(
			"TransactionService - GetAllTransactionViewsByFilters - s.t.GetAllTransactionViewsByFilters:%w",
			nil,
			"filter is empty",
			http.StatusBadRequest,
		)
	}

	var transactions []model.TransactionView
	transactions, err := s.t.GetAllTransactionViewsByFilters(ctx, filter)
	if err != nil {
		return transactions, status.withError(
			"TransactionService - GetAllTransactionViewsByFilters - s.t.GetAllTransactionViewsByFilters:%w",
			err,
			"couldn't get all transactions",
			http.StatusInternalServerError,
		)
	}
	return transactions, status.success("transactions retrieved", http.StatusOK)
}
