package postgres

import (
	"application_aggregator/internal/model"
	"context"
	"database/sql"
)

// OrganizationRepository — интерфейс репозитория
type OrganizationRepository interface {
	Create(ctx context.Context, organization *model.Organization) (*model.Organization, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}

// organizationRepository — реализация
type organizationRepository struct {
	db *sql.DB
}

// NewOrganizationRepository создаёт репозиторий
func NewOrganizationRepository(db *sql.DB) OrganizationRepository {
	return &organizationRepository{db: db}
}

// Create сохраняет организацию в бд
func (r *organizationRepository) Create(ctx context.Context, organization *model.Organization) (*model.Organization, error) {
	const query = `
	INSERT INTO organizations (name, email) 
	VALUES ($1, $2)
	RETURNING uuid, name, email
`

	var saved model.Organization
	err := r.db.QueryRowContext(ctx, query, organization.Name, organization.Email).
		Scan(&saved.UUID, &saved.Name, &saved.Email)
	if err != nil {
		return nil, err
	}
	return &saved, nil

}

// ExistsByEmail проверяет, существует ли организация с таким email
func (r *organizationRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	const query = `SELECT 1 FROM organizations WHERE email = $1`
	var dummy int
	err := r.db.QueryRowContext(ctx, query, email).Scan(&dummy)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil // не найдено → не существует
		}
		return false, err // ошибка БД
	}
	return true, nil // найдено → существует
}
