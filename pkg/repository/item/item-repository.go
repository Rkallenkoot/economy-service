package itemrepository

import (
	"context"
	"database/sql"
	"log"

	v1 "github.com/GameComponent/economy-service/pkg/api/v1"
)

// ItemRepository struct
type ItemRepository struct {
	db *sql.DB
}

// NewItemRepository constructor
func NewItemRepository(db *sql.DB) *ItemRepository {
	return &ItemRepository{
		db: db,
	}
}

// Create a new item
func (r *ItemRepository) Create(ctx context.Context, name string) (*v1.Item, error) {
	lastInsertUUID := ""
	err := r.db.QueryRowContext(
		ctx,
		`INSERT INTO item(name) VALUES ($1) RETURNING id`,
		name,
	).Scan(&lastInsertUUID)

	if err != nil {
		return nil, err
	}

	return &v1.Item{
		Id:   lastInsertUUID,
		Name: name,
	}, nil
}

// Update an item
func (r *ItemRepository) Update(ctx context.Context, id string, name string, metadata string) (*v1.Item, error) {
	_, err := r.db.ExecContext(
		ctx,
		`UPDATE item SET name = $1, metadata = $2 WHERE id = $3`,
		name,
		metadata,
		id,
	)

	if err != nil {
		return nil, err
	}

	item := &v1.Item{
		Id:   id,
		Name: name,
	}

	return item, nil
}

// List all items
func (r *ItemRepository) List(ctx context.Context) ([]*v1.Item, error) {
	// Query items from the database
	rows, err := r.db.QueryContext(ctx, "SELECT id, name FROM economy.item")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Unwrap rows into items
	items := []*v1.Item{}
	for rows.Next() {
		var item v1.Item
		err := rows.Scan(&item.Id, &item.Name)
		if err != nil {
			log.Fatalln(err)
		}

		items = append(items, &item)
	}

	return items, nil
}

// Get an item
func (r *ItemRepository) Get(ctx context.Context, itemID string) (*v1.Item, error) {
	item := &v1.Item{}

	err := r.db.QueryRowContext(
		ctx,
		`SELECT id, name FROM item WHERE id = $1`,
		itemID,
	).Scan(&item.Id, &item.Name)

	if err != nil {
		return nil, err
	}

	return item, nil
}