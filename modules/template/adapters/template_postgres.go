package adapters

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"tixgo/modules/template/domain"

	"github.com/duongptryu/gox/pagination"
	"github.com/duongptryu/gox/syserr"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// TemplatePostgresRepository implements the TemplateRepository interface using PostgreSQL
type TemplatePostgresRepository struct {
	db *sqlx.DB
}

// NewTemplatePostgresRepository creates a new PostgreSQL template repository
func NewTemplatePostgresRepository(db *sqlx.DB) *TemplatePostgresRepository {
	return &TemplatePostgresRepository{db: db}
}

// Create creates a new template in the database
func (r *TemplatePostgresRepository) Create(ctx context.Context, template *domain.Template) error {
	query := `
		INSERT INTO templates (name, slug, subject, content, type, status, variables, description, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id`

	err := r.db.QueryRowContext(
		ctx,
		query,
		template.Name,
		template.Slug,
		template.Subject,
		template.Content,
		template.Type,
		template.Status,
		pq.Array(template.Variables),
		template.Description,
		template.CreatedBy,
		template.CreatedAt,
		template.UpdatedAt,
	).Scan(&template.ID)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "unique constraint") {
			return domain.ErrTemplateAlreadyExists
		}
		return syserr.Wrap(err, syserr.InternalCode, "failed to create template")
	}

	return nil
}

// GetByID retrieves a template by ID
func (r *TemplatePostgresRepository) GetByID(ctx context.Context, id int64) (*domain.Template, error) {
	query := `
		SELECT id, name, slug, subject, content, type, status, variables, description, 
		       created_by, created_at, updated_at
		FROM templates 
		WHERE id = $1`

	template := &domain.Template{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&template.ID,
		&template.Name,
		&template.Slug,
		&template.Subject,
		&template.Content,
		&template.Type,
		&template.Status,
		pq.Array(&template.Variables),
		&template.Description,
		&template.CreatedBy,
		&template.CreatedAt,
		&template.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrTemplateNotFound
		}
		return nil, syserr.Wrap(err, syserr.InternalCode, "failed to get template by ID")
	}

	return template, nil
}

// GetBySlug retrieves a template by slug
func (r *TemplatePostgresRepository) GetBySlug(ctx context.Context, slug string) (*domain.Template, error) {
	query := `
		SELECT id, name, slug, subject, content, type, status, variables, description, 
		       created_by, created_at, updated_at
		FROM templates 
		WHERE slug = $1`

	template := &domain.Template{}
	err := r.db.QueryRowContext(ctx, query, slug).Scan(
		&template.ID,
		&template.Name,
		&template.Slug,
		&template.Subject,
		&template.Content,
		&template.Type,
		&template.Status,
		pq.Array(&template.Variables),
		&template.Description,
		&template.CreatedBy,
		&template.CreatedAt,
		&template.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrTemplateNotFound
		}
		return nil, syserr.Wrap(err, syserr.InternalCode, "failed to get template by slug")
	}

	return template, nil
}

// List retrieves templates with pagination and filters
func (r *TemplatePostgresRepository) List(ctx context.Context, filters domain.ListTemplateFilters, paging *pagination.Paging) ([]*domain.Template, error) {
	// Build WHERE clause
	var conditions []string
	var args []interface{}
	argCount := 0

	if filters.Type != nil {
		argCount++
		conditions = append(conditions, fmt.Sprintf("type = $%d", argCount))
		args = append(args, *filters.Type)
	}

	if filters.Status != nil {
		argCount++
		conditions = append(conditions, fmt.Sprintf("status = $%d", argCount))
		args = append(args, *filters.Status)
	}

	if filters.CreatedBy != nil {
		argCount++
		conditions = append(conditions, fmt.Sprintf("created_by = $%d", argCount))
		args = append(args, *filters.CreatedBy)
	}

	if filters.Search != "" {
		argCount++
		conditions = append(conditions, fmt.Sprintf("(name ILIKE $%d OR description ILIKE $%d OR slug ILIKE $%d)", argCount, argCount, argCount))
		args = append(args, "%"+filters.Search+"%")
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count query
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM templates %s", whereClause)
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, syserr.Wrap(err, syserr.InternalCode, "failed to count templates")
	}

	// Set total in paging
	paging.Total = total

	// Main query
	argCount++
	limitArg := argCount
	argCount++
	offsetArg := argCount

	query := fmt.Sprintf(`
		SELECT id, name, slug, subject, content, type, status, variables, description, 
		       created_by, created_at, updated_at
		FROM templates 
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d`, whereClause, limitArg, offsetArg)

	args = append(args, paging.Limit, paging.GetOffset())

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, syserr.Wrap(err, syserr.InternalCode, "failed to list templates")
	}
	defer rows.Close()

	var templates []*domain.Template
	for rows.Next() {
		template := &domain.Template{}
		err := rows.Scan(
			&template.ID,
			&template.Name,
			&template.Slug,
			&template.Subject,
			&template.Content,
			&template.Type,
			&template.Status,
			pq.Array(&template.Variables),
			&template.Description,
			&template.CreatedBy,
			&template.CreatedAt,
			&template.UpdatedAt,
		)
		if err != nil {
			return nil, syserr.Wrap(err, syserr.InternalCode, "failed to scan template")
		}
		templates = append(templates, template)
	}

	if err = rows.Err(); err != nil {
		return nil, syserr.Wrap(err, syserr.InternalCode, "error iterating template rows")
	}

	return templates, nil
}

// Update updates an existing template
func (r *TemplatePostgresRepository) Update(ctx context.Context, template *domain.Template) error {
	query := `
		UPDATE templates 
		SET name = $2, subject = $3, content = $4, status = $5, variables = $6, 
		    description = $7, updated_at = $8
		WHERE id = $1`

	template.UpdatedAt = time.Now()

	result, err := r.db.ExecContext(
		ctx,
		query,
		template.ID,
		template.Name,
		template.Subject,
		template.Content,
		template.Status,
		pq.Array(template.Variables),
		template.Description,
		template.UpdatedAt,
	)

	if err != nil {
		return syserr.Wrap(err, syserr.InternalCode, "failed to update template")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return syserr.Wrap(err, syserr.InternalCode, "failed to get rows affected")
	}

	if rowsAffected == 0 {
		return domain.ErrTemplateNotFound
	}

	return nil
}

// Delete deletes a template by ID
func (r *TemplatePostgresRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM templates WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return syserr.Wrap(err, syserr.InternalCode, "failed to delete template")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return syserr.Wrap(err, syserr.InternalCode, "failed to get rows affected")
	}

	if rowsAffected == 0 {
		return domain.ErrTemplateNotFound
	}

	return nil
}
