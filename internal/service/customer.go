package service

import (
	"context"
	"net/http"

	"github.com/robertt3kuk/xiaoma-test-task/internal/model"
)

type CustomerService struct {
	t CustomerRepository
}

func NewCustomerService(t CustomerRepository) *CustomerService {
	return &CustomerService{t: t}
}

func (s *CustomerService) Create(ctx context.Context, customer model.Customer) (int, Status) {
	var status Status
	ID, err := s.t.IDByName(ctx, customer.Name)
	if err != nil {
		return 0, status.withError(
			"CustomerService - Create - s.t.NameExists:%w", err, "error with customer name", http.StatusInternalServerError,
		)
	}
	if ID != 0 {
		return 0, status.withError(
			"CustomerService - Create - s.t.NameExists:%w", err, "customer name already exists", http.StatusBadRequest,
		)
	}

	id, err := s.t.Create(ctx, customer)
	if err != nil {
		return 0, status.withError(
			"CustomerService - Create - s.t.Create:%w", err, "error with customer creation", http.StatusInternalServerError,
		)
	}
	return id, status.success("customer succesfully created", http.StatusCreated)
}

func (s *CustomerService) GetByID(ctx context.Context, id int) (model.Customer, Status) {
	var status Status
	var customer model.Customer
	exist, err := s.t.IDExists(ctx, id)
	if err != nil {
		if !exist {
			return customer, status.withError(
				"CustomerService - Update -  s.t.IDExists:%w", err, "customer does not exist", http.StatusNotFound,
			)
		}
	}
	customer, err = s.t.GetByID(ctx, id)
	if err != nil {
		return customer, status.withError(
			"CustomerService - Update- s.t.GetByID:%w", err, "couldn't get customer", http.StatusInternalServerError,
		)
	}
	return customer, status.success("customer retrieved", http.StatusCreated)
}

func (s *CustomerService) GetAll(ctx context.Context, limit, offset int) ([]model.Customer, Status) {
	var status Status
	var customers []model.Customer
	customers, err := s.t.GetAll(ctx, limit, offset)
	if err != nil {
		return customers, status.withError(
			"CustomerService - GetAll - s.t.GetAll:%w", err, "couldn't get all customers", http.StatusInternalServerError,
		)
	}
	return customers, status.success("customers retrieved", http.StatusOK)
}

func (s *CustomerService) Update(ctx context.Context, customer model.Customer) (model.Customer, Status) {
	var status Status
	exist, err := s.t.IDExists(ctx, customer.ID)
	if err != nil {
		if !exist {
			return customer, status.withError(
				"CustomerService - Update -  s.t.IDExists:%w", err, "customer does not exist", http.StatusNotFound,
			)
		}
	}
	ID, err := s.t.IDByName(ctx, customer.Name)
	if err != nil {
		return customer, status.withError(
			"CustomerService - Update - s.t.IDByName:%w", err, "couldn't get customer id", http.StatusInternalServerError,
		)
	}
	if ID != customer.ID {
		// name already in use
		return customer, status.withError(
			"CustomerService - Update - s.t.IDByName:%w", err, "customer name already exists", http.StatusBadRequest,
		)
	}

	customer, err = s.t.Update(ctx, customer)
	if err != nil {
		return customer, status.withError(
			"CustomerService - Update - s.t.Update:%w", err, "couldn't update customer", http.StatusInternalServerError,
		)
	}
	return customer, status.success("customer updated", http.StatusOK)
}

func (s *CustomerService) Delete(ctx context.Context, id int) Status {
	var status Status
	err := s.t.Delete(ctx, id)
	if err != nil {
		return status.withError(
			"CustomerService - Delete - s.t.Delete:%w", err, "couldn't delete customer", http.StatusInternalServerError,
		)
	}
	return status.success("customer deleted", http.StatusOK)
}
