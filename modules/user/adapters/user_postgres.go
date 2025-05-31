package adapters

import (
	"context"
	"database/sql"
	"time"

	"tixgo/modules/user/domain"
	"tixgo/shared/syserr"

	"github.com/jmoiron/sqlx"
)

// UserPostgresRepository implements the UserRepository interface using PostgreSQL
type UserPostgresRepository struct {
	db *sqlx.DB
}

// NewUserPostgresRepository creates a new PostgreSQL user repository
func NewUserPostgresRepository(db *sqlx.DB) *UserPostgresRepository {
	return &UserPostgresRepository{db: db}
}

// Create creates a new user in the database
func (r *UserPostgresRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (email, password_hash, first_name, last_name, phone, date_of_birth, user_type, status, email_verified, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id`

	err := r.db.QueryRowContext(
		ctx,
		query,
		user.Email,
		user.PasswordHash,
		user.FirstName,
		user.LastName,
		user.Phone,
		user.DateOfBirth,
		user.UserType,
		user.Status,
		user.EmailVerified,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID)

	if err != nil {
		return syserr.Wrap(err, syserr.InternalCode, "failed to create user")
	}

	return nil
}

// GetByID retrieves a user by ID
func (r *UserPostgresRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name, phone, date_of_birth, 
		       user_type, status, email_verified, created_at, updated_at, last_login
		FROM users 
		WHERE id = $1`

	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.Phone,
		&user.DateOfBirth,
		&user.UserType,
		&user.Status,
		&user.EmailVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.LastLogin,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrUserNotFound
		}
		return nil, syserr.Wrap(err, syserr.InternalCode, "failed to get user by ID")
	}

	return user, nil
}

// GetByEmail retrieves a user by email
func (r *UserPostgresRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name, phone, date_of_birth, 
		       user_type, status, email_verified, created_at, updated_at, last_login
		FROM users 
		WHERE email = $1`

	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.Phone,
		&user.DateOfBirth,
		&user.UserType,
		&user.Status,
		&user.EmailVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.LastLogin,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrUserNotFound
		}
		return nil, syserr.Wrap(err, syserr.InternalCode, "failed to get user by email")
	}

	return user, nil
}

// Update updates an existing user
func (r *UserPostgresRepository) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users 
		SET email = $2, password_hash = $3, first_name = $4, last_name = $5, 
		    phone = $6, date_of_birth = $7, user_type = $8, status = $9, 
		    email_verified = $10, updated_at = $11, last_login = $12
		WHERE id = $1`

	user.UpdatedAt = time.Now()

	result, err := r.db.ExecContext(
		ctx,
		query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.FirstName,
		user.LastName,
		user.Phone,
		user.DateOfBirth,
		user.UserType,
		user.Status,
		user.EmailVerified,
		user.UpdatedAt,
		user.LastLogin,
	)

	if err != nil {
		return syserr.Wrap(err, syserr.InternalCode, "failed to update user")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return syserr.Wrap(err, syserr.InternalCode, "failed to get rows affected")
	}

	if rowsAffected == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

// Delete deletes a user by ID
func (r *UserPostgresRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return syserr.Wrap(err, syserr.InternalCode, "failed to delete user")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return syserr.Wrap(err, syserr.InternalCode, "failed to get rows affected")
	}

	if rowsAffected == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}
