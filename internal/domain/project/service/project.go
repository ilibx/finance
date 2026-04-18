package service

import (
	"finance/internal/domain/project/entity"
	"finance/internal/domain/project/repository"
	"fmt"
	"time"
)

// ProjectService handles project business logic
type ProjectService struct {
	repo repository.ProjectRepository
}

// NewProjectService creates a new project service
func NewProjectService(repo repository.ProjectRepository) *ProjectService {
	return &ProjectService{repo: repo}
}

// CreateProjectInput represents input for creating a project
type CreateProjectInput struct {
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
	Budget      float64   `json:"budget"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
}

// UpdateProjectInput represents input for updating a project
type UpdateProjectInput struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Budget      float64   `json:"budget"`
	Cost        float64   `json:"cost"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
}

// CreateProject creates a new project
func (s *ProjectService) CreateProject(input CreateProjectInput) (*entity.Project, error) {
	// Validate input
	if input.Name == "" {
		return nil, fmt.Errorf("project name is required")
	}
	if input.Code == "" {
		return nil, fmt.Errorf("project code is required")
	}
	if input.Budget < 0 {
		return nil, fmt.Errorf("budget cannot be negative")
	}
	if input.EndDate.Before(input.StartDate) {
		return nil, fmt.Errorf("end date cannot be before start date")
	}

	// Check if code already exists
	existing, err := s.repo.GetByCode(input.Code)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("project code already exists")
	}

	project := entity.NewProject(input.Name, input.Code, input.Description, input.Budget, input.StartDate, input.EndDate)
	return s.repo.Create(project)
}

// GetProject retrieves a project by ID
func (s *ProjectService) GetProject(id int64) (*entity.Project, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid project ID")
	}
	return s.repo.GetByID(id)
}

// UpdateProject updates an existing project
func (s *ProjectService) UpdateProject(input UpdateProjectInput) (*entity.Project, error) {
	if input.ID <= 0 {
		return nil, fmt.Errorf("invalid project ID")
	}

	project, err := s.repo.GetByID(input.ID)
	if err != nil {
		return nil, fmt.Errorf("project not found")
	}

	// Update fields
	if input.Name != "" {
		project.Name = input.Name
	}
	if input.Code != "" {
		project.Code = input.Code
	}
	project.Description = input.Description
	if input.Status != "" {
		project.UpdateStatus(input.Status)
	}
	project.Budget = input.Budget
	if input.Cost > 0 {
		project.UpdateCost(input.Cost)
	}
	project.StartDate = input.StartDate
	project.EndDate = input.EndDate
	project.UpdatedAt = time.Now()

	return s.repo.Update(project)
}

// DeleteProject deletes a project by ID
func (s *ProjectService) DeleteProject(id int64) error {
	if id <= 0 {
		return fmt.Errorf("invalid project ID")
	}
	return s.repo.Delete(id)
}

// ListProjects lists projects with optional status filter
func (s *ProjectService) ListProjects(status string) ([]*entity.Project, error) {
	return s.repo.List(status)
}

// UpdateProjectStatus updates the status of a project
func (s *ProjectService) UpdateProjectStatus(id int64, status string) (*entity.Project, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid project ID")
	}

	project, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("project not found")
	}

	project.UpdateStatus(status)
	return s.repo.Update(project)
}

// TrackProjectProgress updates project cost and checks budget status
func (s *ProjectService) TrackProjectProgress(id int64, cost float64) (*entity.Project, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid project ID")
	}

	project, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("project not found")
	}

	project.UpdateCost(cost)
	updatedProject, err := s.repo.Update(project)
	if err != nil {
		return nil, err
	}

	return updatedProject, nil
}
