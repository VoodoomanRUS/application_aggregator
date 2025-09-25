package backup

import (
	"application_aggregator/internal/models"

	"gorm.io/gorm"
)

func CreateOrganizationTable() Migration {
	return Migration{
		Version:     "001",
		Description: "Create organization table",
		Up: func(db *gorm.DB) error {
			if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
				return err
			}
			return db.AutoMigrate(&models.Organization{})
		},
		Down: func(db *gorm.DB) error {
			return db.Migrator().DropTable("organizations")
		},
	}

}
