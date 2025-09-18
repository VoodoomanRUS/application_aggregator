package migrations

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Migration struct {
	Version     string
	Description string
	Up          func(db *gorm.DB) error
	Down        func(db *gorm.DB) error
}
type MigrationRecord struct {
	ID          uint64    `gorm:"primary_key"`
	Version     string    `gorm:"uniqueIndex;not null"`
	Description string    `gorm:"not null"`
	AppliedAt   time.Time `gorm:"not null"`
}

type Migrator struct {
	db         *gorm.DB
	migrations []Migration
	tableName  string
}

func NewMigrator(db *gorm.DB) *Migrator {
	return &Migrator{
		db:         db,
		migrations: make([]Migration, 0),
		tableName:  "schema_migrations",
	}
}

func (m *Migrator) AddMigration(migration Migration) {
	m.migrations = append(m.migrations, migration)
}

func (m *Migrator) addMigrationsTables() error {
	return m.db.AutoMigrate(&MigrationRecord{})
}

func (m *Migrator) getAppliedMigrations() (map[string]bool, error) {
	var records []MigrationRecord
	if err := m.db.Find(&records).Error; err != nil {
		return nil, fmt.Errorf("failed to get applied migrations: %w", err)
	}
	appliedMigrations := make(map[string]bool)
	for _, record := range records {
		appliedMigrations[record.Version] = true
	}
	return appliedMigrations, nil
}

func (m *Migrator) Up() error {

	if err := m.addMigrationsTables(); err != nil {
		return fmt.Errorf("failed to add migrations: %w", err)
	}

	appliedMigrations, err := m.getAppliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	for _, migration := range m.migrations {
		if appliedMigrations[migration.Version] {
			continue
		}

		fmt.Printf("Applying migration %s %s\n", migration.Version, migration.Description)

		tx := m.db.Begin()
		if tx.Error != nil {
			return fmt.Errorf("failed to start transaction: %w", tx.Error)
		}

		if err := migration.Up(tx); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to apply migration %s: %w", migration.Version, err)
		}

		migrationRecord := MigrationRecord{
			Version:     migration.Version,
			Description: migration.Description,
			AppliedAt:   time.Now(),
		}

		if err := tx.Create(&migrationRecord).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration %s: %w", migration.Version, err)
		}

		if err := tx.Commit().Error; err != nil {
			return fmt.Errorf("failed to commit migration %s: %w", migration.Version, err)
		}

		fmt.Printf("Successfully applied migration %s %s\n", migration.Version, migration.Description)
	}
	return nil
}

func (m *Migrator) Down() error {
	if err := m.addMigrationsTables(); err != nil {
		return fmt.Errorf("failed to add migrations: %w", err)
	}

	var lastRecord MigrationRecord
	if err := m.db.Order("applied_at DESC").First(&lastRecord).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("no migration to rollback %q", lastRecord.Version)
		}
		return fmt.Errorf("failed to get last migration record: %w", err)
	}

	var migration *Migration
	for i := range m.migrations {
		if m.migrations[i].Version == lastRecord.Version {
			migration = &m.migrations[i]
			break
		}
	}

	if migration == nil {
		return fmt.Errorf("migration %s not found in migrations list", lastRecord.Version)
	}

	if migration.Down == nil {
		return fmt.Errorf("migration %s does not suppoer rollback", lastRecord.Version)
	}
	fmt.Printf("Rollbacking migration %s %s\n", migration.Version, migration.Description)

	tx := m.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to start transaction: %w", tx.Error)
	}

	if err := migration.Down(tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to rollback migration %s: %w", migration.Version, err)
	}

	if err := tx.Commit().Error; err != nil {

		return fmt.Errorf("failed to commit rollback migration %s: %w", migration.Version, err)
	}
	fmt.Printf("Successfully rollbacked migration %s %s\n", migration.Version, migration.Description)
	return nil
}

func (m *Migrator) Status(ctx context.Context) error {
	if err := m.addMigrationsTables(); err != nil {
		return fmt.Errorf("failed to add migrations: %w", err)
	}
	appliedMigrations, err := m.getAppliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	fmt.Println("Migration status:")
	fmt.Println("-----------------")

	for _, migration := range m.migrations {
		status := "Pending"
		if appliedMigrations[migration.Version] {
			status = "Applied"
		}
		fmt.Printf("%s - %s: %s\n\n", migration.Version, status, migration.Description)
	}
	return nil
}
