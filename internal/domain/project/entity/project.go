package entity

import (
	"time"
)

// Project represents a software project aggregate root
type Project struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Budget      float64   `json:"budget"`
	Cost        float64   `json:"cost"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// NewProject creates a new project
func NewProject(name, code, description string, budget float64, startDate, endDate time.Time) *Project {
	now := time.Now()
	return &Project{
		Name:        name,
		Code:        code,
		Description: description,
		Budget:      budget,
		Status:      "planning",
		StartDate:   startDate,
		EndDate:     endDate,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// UpdateStatus updates the project status
func (p *Project) UpdateStatus(status string) {
	p.Status = status
	p.UpdatedAt = time.Now()
}

// UpdateCost updates the project cost
func (p *Project) UpdateCost(cost float64) {
	p.Cost = cost
	p.UpdatedAt = time.Now()
}

// IsActive checks if project is active
func (p *Project) IsActive() bool {
	return p.Status == "active" || p.Status == "planning"
}

// IsOverBudget checks if project cost exceeds budget
func (p *Project) IsOverBudget() bool {
	return p.Cost > p.Budget
}
