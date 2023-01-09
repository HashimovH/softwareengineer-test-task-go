package repository

import (
	"database/sql"
	// "errors"
	"github.com/HashimovH/softwareengineer-test-task-go/app/core/domain"
)

func NewQualityRepository(db *sql.DB) *repository {
	return &repository{DB: db}
}

func (r *repository) GetOveralQualityScore(from string, to string) (*domain.OveralQuality, error) {
	query := `
		SELECT ROUND(SUM(
				CASE
					WHEN rc.weight != 0 THEN ROUND((((r.rating * rc.weight)/5) * 100)/rc.weight)
					ELSE 0
				END
			)/COUNT(*)) AS score
		FROM ratings r
		LEFT JOIN rating_categories rc ON rc.id = r.rating_category_id
		WHERE r.created_at BETWEEN "2019-07-20" AND "2019-07-26";
	`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := domain.OveralQuality{}
	for rows.Next() {

		err = rows.Scan(&data.OveralScore)
		if err != nil {
			return nil, err
		}
	}
	return &data, nil
}
