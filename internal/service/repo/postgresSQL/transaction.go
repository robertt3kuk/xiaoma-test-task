package postgresSQL

import (
	"context"
	"fmt"

	"github.com/robertt3kuk/xiaoma-test-task/init/postgres"
	"github.com/robertt3kuk/xiaoma-test-task/internal/model"
)

type TransactionPostgres struct {
	pg *postgres.Postgres
}

func NewTransactionPostgres(pg *postgres.Postgres) *TransactionPostgres {
	return &TransactionPostgres{pg: pg}
}

const TransactionTable = "transaction"

func (p *TransactionPostgres) Create(ctx context.Context, transaction model.Transaction) (int, error) {
	// return id
	var id int
	tx, err := p.pg.Pool.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("TransactionPostgres - Create - p.pg.Pool.Begin: %w", err)
	}
	// need to minutes the amount from the customer balance in customer table by customer_id  and then insert the transaction into transaction table
	_, err = tx.Exec(
		ctx, `
	UPDATE customer
	SET balance = balance - $1
	WHERE id = $2
`, transaction.Amount, transaction.CustomerID,
	)
	if err != nil {
		tx.Rollback(ctx)
		return 0, fmt.Errorf("TransactionPostgres - Create - tx.Pool.Exec: %w", err)
	}
	_, err = tx.Exec(
		ctx, `
	INSERT INTO `+TransactionTable+`
	(customer_id, item_id, qty, price, amount, created_at, updated_at, deleted_at)
	VALUES ($1, $2, $3, $4, $5, now(), now(), null)
`, transaction.CustomerID, transaction.ItemID, transaction.Qty, transaction.Price, transaction.Amount,
	)
	if err != nil {
		return 0, fmt.Errorf("TransactionPostgres - Create - tx.Pool.Exec: %w", err)
	}
	err = tx.Commit(ctx)
	if err != nil {
		tx.Rollback(ctx)
		return 0, fmt.Errorf("TransactionPostgres - Create - tx.Pool.Commit: %w", err)
	}

	return id, nil
}

func (p *TransactionPostgres) IDExists(ctx context.Context, id int) (bool, error) {
	// if this id exist
	var exists bool
	err := p.pg.Pool.QueryRow(
		ctx, `
	SELECT EXISTS (
		SELECT 1
		FROM `+TransactionTable+`
		WHERE id = $1
	)
`, id,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("TransactionPostgres - IDExist - p.pg.Pool.QueryRow: %w", err)
	}
	return exists, nil
}

func (p *TransactionPostgres) GetByID(ctx context.Context, id int) (model.Transaction, error) {
	var transaction model.Transaction
	err := p.pg.Pool.QueryRow(
		ctx, `
	SELECT id, customer_id, item_id, qty, price, amount, created_at, updated_at, deleted_at
	FROM `+TransactionTable+`
	WHERE id = $1 AND WHERE deleted_at IS NULL
`, id,
	).Scan(
		&transaction.ID,
		&transaction.CustomerID,
		&transaction.ItemID,
		&transaction.Qty,
		&transaction.Price,
		&transaction.Amount,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
		&transaction.DeletedAt,
	)
	if err != nil {
		return model.Transaction{}, fmt.Errorf("TransactionPostgres - GetByID - p.pg.Pool.QueryRow: %w", err)
	}
	return transaction, nil
}

func (p *TransactionPostgres) GetAll(ctx context.Context, limit, offset int) ([]model.Transaction, error) {
	if limit == 0 {
		limit = -1
	}
	rows, err := p.pg.Pool.Query(
		ctx, `
	SELECT id, customer_id, item_id, qty, price, amount, created_at, updated_at, deleted_at
	FROM `+TransactionTable+getLimitAndOffset(limit, offset)+" WHERE deleted_at IS NULL",
	)
	if err != nil {
		return nil, fmt.Errorf("TransactionPostgres - GetAll - p.pg.Pool.Query: %w", err)
	}
	defer rows.Close()

	var transactions []model.Transaction
	for rows.Next() {
		var transaction model.Transaction
		err := rows.Scan(
			&transaction.ID,
			&transaction.CustomerID,
			&transaction.ItemID,
			&transaction.Qty,
			&transaction.Price,
			&transaction.Amount,
			&transaction.CreatedAt,
			&transaction.UpdatedAt,
			&transaction.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("TransactionPostgres - GetAll - rows.Scan: %w", err)
		}
		transactions = append(transactions, transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("TransactionPostgres - GetAll - rows.Err: %w", err)
	}

	return transactions, nil

}

func (p *TransactionPostgres) Update(ctx context.Context, transaction model.Transaction) (model.Transaction, error) {
	// at first i need to get the transaction from db and add it's amount to the customer and then minus it with struct's amount and update transaction it self
	tx, err := p.pg.Pool.Begin(ctx)
	if err != nil {
		return model.Transaction{}, fmt.Errorf("TransactionPostgres - Update - p.pg.Pool.Begin: %w", err)
	}
	_, err = tx.Exec(
		ctx, `
        WITH retrieved_data AS (
            SELECT amount, customer_id 
            FROM TransactionTable 
            WHERE id = $1 
        )
        UPDATE customer
        SET balance = balance + retrieved_data.amount
        FROM retrieved_data
        WHERE id = retrieved_data.customer_id
    `, transaction.ID,
	)

	if err != nil {
		tx.Rollback(ctx)
		return model.Transaction{}, fmt.Errorf("TransactionPostgres - Update - tx.Pool.Exec: %w", err)
	}

	_, err = tx.Exec(ctx,
		`
	UPDATE customer 
	set balance = balance - $1
	WHERE id = $2
	`, transaction.Amount, transaction.CustomerID,
	)

	if err != nil {
		tx.Rollback(ctx)
		return model.Transaction{}, fmt.Errorf("TransactionPostgres - Update - tx.Pool.Exec: %w", err)
	}

	_, err = tx.Exec(
		ctx, `
	UPDATE `+TransactionTable+`
	SET customer_id = $1, item_id = $2, qty = $3, price = $4, amount = $5, updated_at = now()
	WHERE id = $6
`, transaction.CustomerID, transaction.ItemID, transaction.Qty, transaction.Price, transaction.Amount, transaction.ID,
	)
	if err != nil {
		tx.Rollback(ctx)
		return model.Transaction{}, fmt.Errorf("TransactionPostgres - Update - tx.Pool.Exec: %w", err)
	}
	err = tx.Commit(ctx)
	if err != nil {
		tx.Rollback(ctx)
		return model.Transaction{}, fmt.Errorf("TransactionPostgres - Update - tx.Pool.Commit: %w", err)
	}

	return transaction, nil
}

func (p *TransactionPostgres) Delete(ctx context.Context, id int) error {
	// set deleted time to time now
	_, err := p.pg.Pool.Exec(
		ctx, `
	UPDATE `+TransactionTable+
			` SET deleted_at = now()
	WHERE id = $1
`, id,
	)
	if err != nil {
		return fmt.Errorf("TransactionPostgres - Delete - p.pg.Pool.Exec: %w", err)
	}
	return nil
}

func (p *TransactionPostgres) GetAllTransactionViews(ctx context.Context, limit, offset int) (
	[]model.TransactionView, error,
) {
	rows, err := p.pg.Pool.Query(
		ctx, `
	SELECT t.id, t.customer_id, c.name, t.item_id, i.name, t.qty, t.price, t.amount, t.created_at, t.updated_at, t.deleted_at
	FROM `+TransactionTable+` AS t
	INNER JOIN customer AS c ON t.customer_id = c.id
	INNER JOIN item AS i ON t.item_id = i.id`+
			getLimitAndOffset(limit, offset)+` WHERE t.deleted_at IS NULL`,
	)
	if err != nil {
		return nil, fmt.Errorf("TransactionPostgres - GetAllTransactionViews - p.pg.Pool.Query: %w", err)
	}
	defer rows.Close()

	var transactionViews []model.TransactionView
	for rows.Next() {
		var transactionView model.TransactionView
		err := rows.Scan(
			&transactionView.ID,
			&transactionView.CustomerID,
			&transactionView.CustomerName,
			&transactionView.ItemID,
			&transactionView.ItemName,
			&transactionView.Qty,
			&transactionView.Price,
			&transactionView.Amount,
			&transactionView.CreatedAt,
			&transactionView.UpdatedAt,
			&transactionView.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("TransactionPostgres - GetAllTransactionViews - rows.Scan: %w", err)
		}
		transactionViews = append(transactionViews, transactionView)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("TransactionPostgres - GetAllTransactionViews - rows.Err: %w", err)
	}

	return transactionViews, nil
}

func (p *TransactionPostgres) GetByTransactionID(ctx context.Context, id int) (model.TransactionView, error) {
	var transactionView model.TransactionView
	err := p.pg.Pool.QueryRow(
		ctx, `
	SELECT t.id, t.customer_id, c.name, t.item_id, i.name, t.qty, t.price, t.amount, t.created_at, t.updated_at, t.deleted_at
	FROM `+TransactionTable+` AS t
	INNER JOIN customer AS c ON t.customer_id = c.id
	INNER JOIN item AS i ON t.item_id = i.id
	WHERE t.id = $1 AND WHERE t.deleted_at IS NULL
`, id,
	).Scan(
		&transactionView.ID,
		&transactionView.CustomerID,
		&transactionView.CustomerName,
		&transactionView.ItemID,
		&transactionView.ItemName,
		&transactionView.Qty,
		&transactionView.Price,
		&transactionView.Amount,
		&transactionView.CreatedAt,
		&transactionView.UpdatedAt,
		&transactionView.DeletedAt,
	)
	if err != nil {
		return model.TransactionView{}, fmt.Errorf(
			"TransactionPostgres - GetByTransactionID - p.pg.Pool.QueryRow: %w", err,
		)
	}
	return transactionView, nil
}

func (p *TransactionPostgres) GetByCustomerName(ctx context.Context, name string) (model.TransactionView, error) {
	var transactionView model.TransactionView
	err := p.pg.Pool.QueryRow(
		ctx, `
	SELECT t.id, t.customer_id, c.name, t.item_id, i.name, t.qty, t.price, t.amount, t.created_at, t.updated_at, t.deleted_at
	FROM `+TransactionTable+` AS t
	INNER JOIN customer AS c ON t.customer_id = c.id
	INNER JOIN item AS i ON t.item_id = i.id
	WHERE c.name = $1
`, name,
	).Scan(
		&transactionView.ID,
		&transactionView.CustomerID,
		&transactionView.CustomerName,
		&transactionView.ItemID,
		&transactionView.ItemName,
		&transactionView.Qty,
		&transactionView.Price,
		&transactionView.Amount,
		&transactionView.CreatedAt,
		&transactionView.UpdatedAt,
		&transactionView.DeletedAt,
	)
	if err != nil {
		return model.TransactionView{}, fmt.Errorf(
			"TransactionPostgres - GetByCustomerName - p.pg.Pool.QueryRow: %w", err,
		)
	}
	return transactionView, nil
}

func (p *TransactionPostgres) GetByItemName(ctx context.Context, name string) (model.TransactionView, error) {
	var transactionView model.TransactionView
	err := p.pg.Pool.QueryRow(
		ctx, `
	SELECT t.id, t.customer_id, c.name, t.item_id, i.name, t.qty, t.price, t.amount, t.created_at, t.updated_at, t.deleted_at
	FROM `+TransactionTable+` AS t
	INNER JOIN customer AS c ON t.customer_id = c.id
	INNER JOIN item AS i ON t.item_id = i.id
	WHERE i.name = $1
`, name,
	).Scan(
		&transactionView.ID,
		&transactionView.CustomerID,
		&transactionView.CustomerName,
		&transactionView.ItemID,
		&transactionView.ItemName,
		&transactionView.Qty,
		&transactionView.Price,
		&transactionView.Amount,
		&transactionView.CreatedAt,
		&transactionView.UpdatedAt,
		&transactionView.DeletedAt,
	)
	if err != nil {
		return model.TransactionView{}, fmt.Errorf(
			"TransactionPostgres - GetByItemName - p.pg.Pool.QueryRow: %w", err,
		)
	}
	return transactionView, nil
}
