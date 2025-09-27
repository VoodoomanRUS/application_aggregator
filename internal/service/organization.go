package service

import (
	"application_aggregator/internal/model"
	"application_aggregator/internal/repository/postgres"
	"context"
	"errors"
	"regexp"
)

// OrganizationService - интерфейс сервиса
type OrganizationService interface {
	Create(ctx context.Context, name, email string) (*model.Organization, error)
}

// organizationService — реализация сервиса
type organizationService struct {
	repo postgres.OrganizationRepository
}

// NewOrganizationService - создаёт новый сервис
func NewOrganizationService(repo postgres.OrganizationRepository) OrganizationService {
	return &organizationService{repo: repo}
}

// Create — бизнес-логика создания организации
func (s *organizationService) Create(ctx context.Context, name string, email string) (*model.Organization, error) {

	//1. Валидация
	if name == "" {
		return nil, errors.New("name is required")
	}
	if !isValidEmail(email) {
		return nil, errors.New("invalid email")
	}

	//2. Проверка уникальности Email (через репозиторий)
	exists, err := s.repo.ExistsByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, errors.New("organization with this email already exists")
	}

	// 4. Сохранение в БД (через репозиторий)

	org := &model.Organization{Name: name, Email: email}
	return s.repo.Create(ctx, org)
}
func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}
