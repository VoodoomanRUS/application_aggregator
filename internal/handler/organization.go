package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// OrganizationRequest represents the input for creating an organization.
// @Description Input data for creating an organization
type OrganizationRequest struct {
	// @Description Name of the organization
	Name string `json:"name" binding:"required"`
	// @Description Email of the organization
	Email string `json:"email" binding:"required,email"`
}

// CreateOrganization creates a new organization.
// @Summary Create organization
// @Description Create a new organization with name and email
// @Tags organizations
// @Accept json
// @Produce json
// @Param organization body OrganizationRequest true "Organization data"
// @Success 201 {object} model.Organization
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /organizations [post]
func (h *Handler) CreateOrganization(c *gin.Context) {
	// 1. Прочитать JSON из запроса
	var input OrganizationRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Validation error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. Валидация (опционально)
	if input.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is required"})
		return
	}
	// 3. Вызов сервиса (бизнес-логики)
	org, err := h.organizationService.Create(c.Request.Context(), input.Name, input.Email)
	if err != nil {
		log.Printf("ERROR creating organization: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create organization"})
		return
	}
	// 4. Отправить ответ
	c.JSON(http.StatusCreated, org)
}
