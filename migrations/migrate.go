package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
	"test_go/models"

)

func RunMigrations(db *gorm.DB) error {
	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "2025050701_rename_female_to_gender",
			Migrate: func(tx *gorm.DB) error {

				return tx.Migrator().RenameColumn(&models.Author{}, "female", "gender")
			},
			Rollback: func(tx *gorm.DB) error {
				
				return tx.Migrator().RenameColumn(&models.Author{}, "gender", "female")
			},
		},
	})

	return m.Migrate()
}