package repository

import (
	"context"
	"database/sql"
	"time"

	"erp-system/internal/common/valueobject"
	"erp-system/internal/domain/user/entity"
	"erp-system/internal/domain/user/repository"
)

// userRepositoryImpl implements repository.UserRepository
type userRepositoryImpl struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB) repository.UserRepository {
	return &userRepositoryImpl{db: db}
}

// Create creates a new user
func (r *userRepositoryImpl) Create(ctx context.Context, user *entity.User) error {
	query := `INSERT INTO users (username, email, phone, balance_amount, balance_currency, status_code, status_description, status_updated_at, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`
	
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

// GetByID gets a user by ID
func (r *userRepositoryImpl) GetByID(ctx context.Context, id int64) (*entity.User, error) {
	query := `SELECT id, username, email, phone, balance_amount, balance_currency, 
		status_code, status_description, status_updated_at, created_at, updated_at 
		FROM users WHERE id = $1`
	
	row := r.db.QueryRowContext(ctx, query, id)
	
	var user entity.User
	var updatedAt time.Time
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Phone,
		&user.Balance.Amount,
		&user.Balance.Currency,
		&user.Status.Code,
		&user.Status.Description,
		&updatedAt,
		&user.CreatedAt,
		&updatedAt,
	)
	if err != nil {
		return nil, err
	}
	user.Status.UpdatedAt = updatedAt
	user.UpdatedAt = updatedAt
	
	return &user, nil
}

// GetByEmail gets a user by email
func (r *userRepositoryImpl) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `SELECT id, username, email, phone, balance_amount, balance_currency, 
		status_code, status_description, status_updated_at, created_at, updated_at 
		FROM users WHERE email = $1`
	
	row := r.db.QueryRowContext(ctx, query, email)
	
	var user entity.User
	var updatedAt time.Time
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Phone,
		&user.Balance.Amount,
		&user.Balance.Currency,
		&user.Status.Code,
		&user.Status.Description,
		&updatedAt,
		&user.CreatedAt,
		&updatedAt,
	)
	if err != nil {
		return nil, err
	}
	user.Status.UpdatedAt = updatedAt
	user.UpdatedAt = updatedAt
	
	return &user, nil
}

// Update updates a user
func (r *userRepositoryImpl) Update(ctx context.Context, user *entity.User) error {
	query := `UPDATE users SET username=$1, email=$2, phone=$3, balance_amount=$4, 
		balance_currency=$5, status_code=$6, status_description=$7, status_updated_at=$8, 
		updated_at=$9 WHERE id=$10`
	
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

// List lists users with pagination
func (r *userRepositoryImpl) List(ctx context.Context, limit, offset int) ([]*entity.User, error) {
	query := `SELECT id, username, email, phone, balance_amount, balance_currency, 
		status_code, status_description, status_updated_at, created_at, updated_at 
		FROM users ORDER BY id LIMIT $1 OFFSET $2`
	
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var users []*entity.User
	for rows.Next() {
		var user entity.User
		var updatedAt time.Time
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.Phone,
			&user.Balance.Amount,
			&user.Balance.Currency,
			&user.Status.Code,
			&user.Status.Description,
			&updatedAt,
			&user.CreatedAt,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}
		user.Status.UpdatedAt = updatedAt
		user.UpdatedAt = updatedAt
		users = append(users, &user)
	}
	
	return users, rows.Err()
}

// Ensure interface compliance
var _ repository.UserRepository = (*userRepositoryImpl)(nil)
