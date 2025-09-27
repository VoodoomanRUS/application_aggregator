package handler

import "application_aggregator/internal/service"

// Handler — структура хендлера для организаций
type Handler struct {
	organizationService service.OrganizationService
}

// NewOrganizationHandler создаёт хендлер
func NewOrganizationHandler(organizationService service.OrganizationService) *Handler {
	return &Handler{organizationService: organizationService}
}
