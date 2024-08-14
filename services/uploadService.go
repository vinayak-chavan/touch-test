package services

import (
	"fmt"
	"touch-test/models"

	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

const batchSize = 500

func ProcessXlsxFile(db *gorm.DB, filePath string) error {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return fmt.Errorf("could not open XLSX file: %w", err)
	}

	sheetName := f.GetSheetName(0)
	if sheetName == "" {
		return fmt.Errorf("the sheet name cannot be blank")
	}

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return fmt.Errorf("could not read rows from the XLSX file: %w", err)
	}

	expectedHeader := []string{
		"first_name", "last_name", "company_name", "address", "city",
		"county", "postal", "phone", "email", "web",
	}

	if len(rows) == 0 || !validateHeader(rows[0], expectedHeader) {
		return fmt.Errorf("the file is not in the expected format")
	}

	var users []models.User

	for _, row := range rows[1:] {
		if len(row) < len(expectedHeader) {
			return fmt.Errorf("row does not contain enough columns")
		}

		user := models.User{
			FirstName:   row[0],
			LastName:    row[1],
			CompanyName: row[2],
			Address:     row[3],
			City:        row[4],
			County:      row[5],
			Postal:      row[6],
			Phone:       row[7],
			Email:       row[8],
			Web:         row[9],
		}

		users = append(users, user)
		if len(users) >= batchSize {
			if err := insertBatch(db, users); err != nil {
				return err
			}
			users = nil
		}
	}

	if len(users) > 0 {
		if err := insertBatch(db, users); err != nil {
			return err
		}
	}
	return nil
}

func validateHeader(actual, expected []string) bool {
	if len(actual) != len(expected) {
		return false
	}
	for i, v := range actual {
		if v != expected[i] {
			return false
		}
	}
	return true
}

// Helper function to insert a batch of users into the database
func insertBatch(db *gorm.DB, users []models.User) error {
	tx := db.Begin()
	if err := tx.Create(&users).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("could not insert batch into MySQL: %w", err)
	}
	tx.Commit()
	return nil
}
