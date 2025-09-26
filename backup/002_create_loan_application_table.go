package backup

import (
	"application_aggregator/backup/models"

	"gorm.io/gorm"
)

func CreateLoanApplicationTable() Migration {
	return Migration{
		Version:     "002",
		Description: "Create loan application table",
		Up: func(db *gorm.DB) error {
			if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
				return err
			}
			return db.AutoMigrate(&models.LoanApplication{})
		},
		Down: func(db *gorm.DB) error {
			return db.Migrator().DropTable("loan_applications")
		},
	}
}
