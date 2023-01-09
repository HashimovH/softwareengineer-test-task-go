package repository

import (
	"database/sql"
)

func QualityRepository(db *sql.DB) *repository {
	return &repository{DB: db}
}

func (r *repository) GetOveralQualityScore(from string, to string) {

}
