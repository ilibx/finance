package repository

import (
	"database/sql"
	"finance/internal/domain/project/entity"
	"fmt"
)

// ProjectRepository defines the interface for project data access
type ProjectRepository interface {
	Create(project *entity.Project) (*entity.Project, error)
	GetByID(id int64) (*entity.Project, error)
	GetByCode(code string) (*entity.Project, error)
	Update(project *entity.Project) (*entity.Project, error)
	Delete(id int64) error
	List(status string) ([]*entity.Project, error)
}

// postgresProjectRepository implements ProjectRepository using PostgreSQL
type postgresProjectRepository struct {
	db *sql.DB
}

// NewProjectRepository creates a new project repository
func NewProjectRepository(db *sql.DB) ProjectRepository {
	return &postgresProjectRepository{db: db}
}

// Create creates a new project in the database
func (r *postgresProjectRepository) Create(project *entity.Project) (*entity.Project, error) {
	query := `
		INSERT INTO projects (name, code, description, status, budget, cost, start_date, end_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id`
	
	err := r.db.QueryRow(query,
		project.Name,
		project.Code,
		project.Description,
		project.Status,
		project.Budget,
		project.Cost,
		project.StartDate,
		project.EndDate,
		project.CreatedAt,
		project.UpdatedAt,
	).Scan(&project.ID)
	
	if err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}
	
	return project, nil
}

// GetByID retrieves a project by ID
func (r *postgresProjectRepository) GetByID(id int64) (*entity.Project, error) {
	query := `
		SELECT id, name, code, description, status, budget, cost, start_date, end_date, created_at, updated_at
		FROM projects
		WHERE id = $1`
	
	project := &entity.Project{}
	err := r.db.QueryRow(query, id).Scan(
		&project.ID,
		&project.Name,
		&project.Code,
		&project.Description,
		&project.Status,
		&project.Budget,
		&project.Cost,
		&project.StartDate,
		&project.EndDate,
		&project.CreatedAt,
		&project.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("project not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}
	
	return project, nil
}

// GetByCode retrieves a project by code
func (r *postgresProjectRepository) GetByCode(code string) (*entity.Project, error) {
	query := `
		SELECT id, name, code, description, status, budget, cost, start_date, end_date, created_at, updated_at
		FROM projects
		WHERE code = $1`
	
	project := &entity.Project{}
	err := r.db.QueryRow(query, code).Scan(
		&project.ID,
		&project.Name,
		&project.Code,
		&project.Description,
		&project.Status,
		&project.Budget,
		&project.Cost,
		&project.StartDate,
		&project.EndDate,
		&project.CreatedAt,
		&project.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("project not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}
	
	return project, nil
}

// Update updates an existing project
func (r *postgresProjectRepository) Update(project *entity.Project) (*entity.Project, error) {
	query := `
		UPDATE projects
		SET name = $1, code = $2, description = $3, status = $4, budget = $5, cost = $6, 
		    start_date = $7, end_date = $8, updated_at = $9
		WHERE id = $10`
	
	result, err := r.db.Exec(query,
		project.Name,
		project.Code,
		project.Description,
		project.Status,
		project.Budget,
		project.Cost,
		project.StartDate,
		project.EndDate,
		project.UpdatedAt,
		project.ID,
	)
	
	if err != nil {
		return nil, fmt.Errorf("failed to update project: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return nil, fmt.Errorf("project not found")
	}
	
	return project, nil
}

// Delete deletes a project by ID
func (r *postgresProjectRepository) Delete(id int64) error {
	query := `DELETE FROM projects WHERE id = $1`
	
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("project not found")
	}
	
	return nil
}

// List retrieves projects with optional status filter
func (r *postgresProjectRepository) List(status string) ([]*entity.Project, error) {
	var query string
	var args []interface{}
	
	if status == "" {
		query = `
			SELECT id, name, code, description, status, budget, cost, start_date, end_date, created_at, updated_at
			FROM projects
			ORDER BY created_at DESC`
	} else {
		query = `
			SELECT id, name, code, description, status, budget, cost, start_date, end_date, created_at, updated_at
			FROM projects
			WHERE status = $1
			ORDER BY created_at DESC`
		args = append(args, status)
	}
	
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list projects: %w", err)
	}
	defer rows.Close()
	
	var projects []*entity.Project
	for rows.Next() {
		project := &entity.Project{}
		err := rows.Scan(
			&project.ID,
			&project.Name,
			&project.Code,
			&project.Description,
			&project.Status,
			&project.Budget,
			&project.Cost,
			&project.StartDate,
			&project.EndDate,
			&project.CreatedAt,
			&project.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan project: %w", err)
		}
		projects = append(projects, project)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating projects: %w", err)
	}
	
	return projects, nil
}
