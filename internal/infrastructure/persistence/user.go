package persistence

import (
	"context"
	"database/sql"
	"time"

	"finance/internal/domain/user/entity"
	"finance/internal/domain/user/repository"
)

// postgresUserRepository implements repository.UserRepository
type postgresUserRepository struct {
	db *sql.DB
}

// NewPostgresUserRepository creates a new PostgreSQL user repository
func NewPostgresUserRepository(db *sql.DB) repository.UserRepository {
	return &postgresUserRepository{db: db}
}

func (r *postgresUserRepository) Create(ctx context.Context, user *entity.User) error {
	query := `
		INSERT INTO users (username, email, phone, balance_amount, balance_currency, 
			status_code, status_description, status_updated_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`
	return r.db.QueryRowContext(ctx, query,
		user.Username,
		user.Email,
		user.Phone,
		user.Balance.Amount,
		user.Balance.Currency,
		user.Status.Code,
		user.Status.Description,
		user.Status.UpdatedAt,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID)
}

func (r *postgresUserRepository) GetByID(ctx context.Context, id int64) (*entity.User, error) {
	query := `
		SELECT id, username, email, phone, balance_amount, balance_currency,
			status_code, status_description, status_updated_at, created_at, updated_at
		FROM users WHERE id = $1
	`
	user := &entity.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Phone,
		&user.Balance.Amount,
		&user.Balance.Currency,
		&user.Status.Code,
		&user.Status.Description,
		&user.Status.UpdatedAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *postgresUserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `
		SELECT id, username, email, phone, balance_amount, balance_currency,
			status_code, status_description, status_updated_at, created_at, updated_at
		FROM users WHERE email = $1
	`
	user := &entity.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Phone,
		&user.Balance.Amount,
		&user.Balance.Currency,
		&user.Status.Code,
		&user.Status.Description,
		&user.Status.UpdatedAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *postgresUserRepository) Update(ctx context.Context, user *entity.User) error {
	query := `
		UPDATE users SET
			username = $1, email = $2, phone = $3,
			balance_amount = $4, balance_currency = $5,
			status_code = $6, status_description = $7, status_updated_at = $8,
			updated_at = $9
		WHERE id = $10
	`
	_, err := r.db.ExecContext(ctx, query,
		user.Username,
		user.Email,
		user.Phone,
		user.Balance.Amount,
		user.Balance.Currency,
		user.Status.Code,
		user.Status.Description,
		user.Status.UpdatedAt,
		time.Now(),
		user.ID,
	)
	return err
}

func (r *postgresUserRepository) List(ctx context.Context, limit, offset int) ([]*entity.User, error) {
	query := `
		SELECT id, username, email, phone, balance_amount, balance_currency,
			status_code, status_description, status_updated_at, created_at, updated_at
		FROM users ORDER BY id DESC LIMIT $1 OFFSET $2
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*entity.User
	for rows.Next() {
		user := &entity.User{}
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.Phone,
			&user.Balance.Amount,
			&user.Balance.Currency,
			&user.Status.Code,
			&user.Status.Description,
			&user.Status.UpdatedAt,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}
