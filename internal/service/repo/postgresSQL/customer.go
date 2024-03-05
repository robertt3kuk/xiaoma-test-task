package postgresSQL

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/robertt3kuk/xiaoma-test-task/init/postgres"
	"github.com/robertt3kuk/xiaoma-test-task/internal/model"
)

type CustomerPostgres struct {
	pg *postgres.Postgres
}

func NewCustomerPostgres(pg *postgres.Postgres) *CustomerPostgres {
	return &CustomerPostgres{pg: pg}
}

const CustomerTable = "customer"

func (p *CustomerPostgres) Create(ctx context.Context, customer model.Customer) (int, error) {
	// insert and return id
	var id int
	err := p.pg.Pool.QueryRow(
		ctx, fmt.Sprintf(
			"INSERT INTO %s (customer_name, balance) VALUES ($1, $2) RETURNING id",
			CustomerTable,
		), customer.Name, customer.Balance,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("postgres - CustomerPostgres - Create: %w", err)
	}
	return id, nil
}

func (p *CustomerPostgres) IDExists(ctx context.Context, id int) (bool, error) {
	// check if exists
	var exists bool
	err := p.pg.Pool.QueryRow(
		ctx, fmt.Sprintf(
			"SELECT EXISTS(SELECT 1 FROM %s WHERE id = $1 AND deleted_at IS NULL)",
			CustomerTable,
		), id,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("postgres - CustomerPostgres - IDExists: %w", err)
	}
	return exists, nil
}

func (p *CustomerPostgres) IDByName(ctx context.Context, name string) (int, error) {
	// so check if customername exists and return it's id if not return 0
	var id int
	err := p.pg.Pool.QueryRow(
		ctx, fmt.Sprintf(
			"SELECT id FROM %s WHERE customer_name = $1",
			CustomerTable,
		), name,
	).Scan(&id)
	if err != nil {
		// if err is now row return 0 else return error
		if err == pgx.ErrNoRows {
			return 0, nil
		}
		return 0, fmt.Errorf("postgres - CustomerPostgres - IDByName: %w", err)
	}
	return id, nil
}

func (p *CustomerPostgres) GetBalance(ctx context.Context, id int) (float64, error) {
	var balance float64
	err := p.pg.Pool.QueryRow(
		ctx, fmt.Sprintf(
			"SELECT balance FROM %s WHERE id = $1 AND deleted_at IS NULL",
			CustomerTable,
		), id,
	).Scan(&balance)
	if err != nil {
		return 0, fmt.Errorf("postgres - CustomerPostgres - GetBalance: %w", err)
	}
	return balance, nil
}

func (p *CustomerPostgres) GetByID(ctx context.Context, id int) (model.Customer, error) {
	// get by id
	var customer model.Customer
	err := p.pg.Pool.QueryRow(
		ctx, fmt.Sprintf(
			"SELECT id, customer_name, balance, created_at, updated_at, deleted_at FROM %s WHERE id = $1 AND  deleted_at IS NULL",
			CustomerTable,
		), id,
	).Scan(
		&customer.ID,
		&customer.Name,
		&customer.Balance,
		&customer.CreatedAt,
		&customer.UpdatedAt,
		&customer.DeletedAt,
	)
	if err != nil {
		return model.Customer{}, fmt.Errorf("postgres - CustomerPostgres - GetByID: %w", err)
	}
	return customer, nil
}

func (p *CustomerPostgres) GetAll(
	ctx context.Context,
	limit, offset int,
) ([]model.Customer, error) {
	rows, err := p.pg.Pool.Query(
		ctx, fmt.Sprintf(
			"SELECT id, customer_name, balance, created_at, updated_at, deleted_at FROM %s WHERE deleted_at IS NULL"+getLimitAndOffset(
				limit,
				offset,
			),
			CustomerTable,
		),
	)
	if err != nil {
		return nil, fmt.Errorf("postgres - CustomerPostgres - GetAll: %w", err)
	}
	defer rows.Close()

	var customers []model.Customer
	for rows.Next() {
		var customer model.Customer
		err := rows.Scan(
			&customer.ID,
			&customer.Name,
			&customer.Balance,
			&customer.CreatedAt,
			&customer.UpdatedAt,
			&customer.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("postgres - CustomerPostgres - GetAll: %w", err)
		}
		customers = append(customers, customer)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("postgres - CustomerPostgres - GetAll: %w", err)
	}

	return customers, nil
}

func (p *CustomerPostgres) Update(
	ctx context.Context,
	customer model.Customer,
) (model.Customer, error) {
	// update
	_, err := p.pg.Pool.Exec(
		ctx, fmt.Sprintf(
			"UPDATE %s SET customer_name = $1, balance = $2, updated_at = now() WHERE id = $3",
			CustomerTable,
		), customer.Name, customer.Balance, customer.ID,
	)
	if err != nil {
		return model.Customer{}, fmt.Errorf("postgres - CustomerPostgres - Update: %w", err)
	}
	return customer, nil
}

func (p *CustomerPostgres) Delete(ctx context.Context, id int) error {
	// delete by seting deleted_at time.Now

	_, err := p.pg.Pool.Exec(
		ctx, fmt.Sprintf(
			"UPDATE %s SET deleted_at = now() WHERE id = $1",
			CustomerTable,
		), id,
	)
	if err != nil {
		return fmt.Errorf("postgres - CustomerPostgres - Delete: %w", err)
	}
	return nil
}
