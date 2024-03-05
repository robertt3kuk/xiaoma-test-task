package postgresSQL

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/robertt3kuk/xiaoma-test-task/init/postgres"
	"github.com/robertt3kuk/xiaoma-test-task/internal/model"
)

type ItemPostgres struct {
	pg *postgres.Postgres
}

func NewItemPostgres(pg *postgres.Postgres) *ItemPostgres {
	return &ItemPostgres{pg: pg}
}

const ItemTable = "item"

func (p *ItemPostgres) Create(ctx context.Context, item model.Item) (int, error) {
	// insert
	query := `INSERT INTO ` + ItemTable + ` (item_name, cost, price, sort) 
	VALUES ($1, $2, $3, $4 ) 
	RETURNING id`

	var id int
	err := p.pg.Pool.QueryRow(
		ctx, query, item.ItemName, item.Cost, item.Price, item.Sort,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("postgres - ItemPostgres.Create - p.pg.Pool.QueryRow: %w", err)
	}

	return id, nil
}

func (p *ItemPostgres) IDExists(ctx context.Context, id int) (bool, error) {
	// check if exists
	query := `SELECT EXISTS(SELECT 1 FROM ` + ItemTable + ` WHERE id = $1)`

	var exists bool
	err := p.pg.Pool.QueryRow(ctx, query, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("postgres - ItemPostgres.IDExist - p.pg.Pool.QueryRow: %w", err)
	}

	return exists, nil
}

func (p *ItemPostgres) IDByItemName(ctx context.Context, ItemName string) (int, error) {
	// so check if itemname exists and return it's id if not return 0
	var id int
	err := p.pg.Pool.QueryRow(
		ctx, fmt.Sprintf(
			"SELECT id FROM %s WHERE item_name = $1",
			ItemTable,
		), ItemName,
	).Scan(&id)
	if err != nil {
		// if err is now row return 0 else return error
		if err == pgx.ErrNoRows {
			return 0, nil
		}
		return 0, fmt.Errorf("postgres - ItemPostgres.IDByItemName - p.pg.Pool.QueryRow: %w", err)
	}

	return id, nil
}

func (p *ItemPostgres) GetByID(ctx context.Context, id int) (model.Item, error) {
	query := `SELECT id, item_name, cost, price, sort, created_at, updated_at, deleted_at 
	FROM ` + ItemTable + ` WHERE id = $1 AND WHERE deleted_at IS NULL`

	var item model.Item
	err := p.pg.Pool.QueryRow(ctx, query, id).Scan(
		&item.ID, &item.ItemName, &item.Cost, &item.Price, &item.Sort, &item.CreatedAt, &item.UpdatedAt,
		&item.DeletedAt,
	)
	if err != nil {
		return model.Item{}, fmt.Errorf("postgres - ItemPostgres.GetByID - p.pg.Pool.QueryRow: %w", err)
	}

	return item, nil
}

func (p *ItemPostgres) GetAll(ctx context.Context, limit, offset int) ([]model.Item, error) {

	query := `SELECT id, item_name, cost, price, sort, created_at, updated_at, deleted_at 
FROM ` + ItemTable + getLimitAndOffset(limit, offset) + " WHERE deleted_at IS NULL"

	rows, err := p.pg.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("postgres - ItemPostgres.GetAll - p.pg.Pool.Query: %w", err)
	}
	defer rows.Close()

	var items []model.Item
	for rows.Next() {
		var item model.Item
		err := rows.Scan(
			&item.ID, &item.ItemName, &item.Cost, &item.Price, &item.Sort, &item.CreatedAt, &item.UpdatedAt,
			&item.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("postgres - ItemPostgres.GetAll - rows.Scan: %w", err)
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("postgres - ItemPostgres.GetAll - rows.Err: %w", err)
	}

	return items, nil

}

func (p *ItemPostgres) Update(ctx context.Context, item model.Item) (model.Item, error) {
	query := `UPDATE ` + ItemTable + ` SET  item_name=$1, cost=$2, price=$3, sort=$4, updated_at= now()  WHERE id=$5`

	_, err := p.pg.Pool.Exec(
		ctx, query, item.ItemName, item.Cost, item.Price, item.Sort, item.ID,
	)
	if err != nil {
		return model.Item{}, fmt.Errorf("postgres - ItemPostgres.Update - p.pg.Pool.Exec: %w", err)
	}

	return item, nil
}
func (p *ItemPostgres) Delete(ctx context.Context, id int) error {
	// delete by setting deleted_at to time.Now
	query := `UPDATE ` + ItemTable + ` SET deleted_at= now() WHERE id=$1`

	_, err := p.pg.Pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("postgres - ItemPostgres.Delete - p.pg.Pool.Exec: %w", err)
	}

	return nil
}
