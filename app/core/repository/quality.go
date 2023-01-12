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
		select sum(weight*score)/sum(weight) as overal
		from
		(
			SELECT r.ticket_id as ticket_id, rc.name as name, ROUND(AVG(r.rating) * 20) as score,rc.weight as weight
			FROM ratings r
			LEFT JOIN rating_categories rc ON rc.id = r.rating_category_id
			WHERE r.created_at BETWEEN ? AND ?
			GROUP BY r.ticket_id, r.rating_category_id
		)
	`

	rows, err := r.DB.Query(query, from, to)
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
