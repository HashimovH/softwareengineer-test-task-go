package repository

import (
	"database/sql"

	"github.com/HashimovH/softwareengineer-test-task-go/app/core/domain"
)

func NewRepository(db *sql.DB) *repository {
	return &repository{DB: db}
}

func (r *repository) GetAggregatedCategoryScoresDaily(from string, to string) ([]domain.Score, error) {
	query :=
		`
		SELECT rc.name,strftime ('%d', r.created_at), 
		(CASE
				WHEN rc.weight != 0 THEN ROUND(
							(((r.rating * rc.weight) / 5) * 100) / rc.weight)
				ELSE 0 END)       as score, r.rating
		FROM ratings r 
		JOIN rating_categories rc ON rc.id = r.rating_category_id 
		WHERE created_at >= "2019-07-17T15:49:14"
  			and created_at <= "2019-07-20T15:49:14"
		GROUP BY r.rating_category_id, strftime ('%d', r.created_at)
	`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := []domain.Score{}
	for rows.Next() {
		i := domain.Score{}
		err = rows.Scan(&i.Category, &i.Date, &i.Rating, &i.Score)
		if err != nil {
			return nil, err
		}
		data = append(data, i)
	}

	return data, nil
}

func (r *repository) GetScoresByTicket(from string, to string) ([]*domain.ScoreByTicket, error) {
	query := `
		SELECT r.ticket_id, rc.name,
			(CASE
				WHEN rc.weight != 0 THEN ROUND((((r.rating * rc.weight) / 5) * 100) / rc.weight)
				ELSE 0
			END) AS score
		FROM ratings r
		LEFT JOIN rating_categories rc ON rc.id = r.rating_category_id
		WHERE r.created_at BETWEEN ? AND ?
		GROUP BY r.ticket_id, r.rating_category_id;
	`

	rows, err := r.DB.Query(query, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := []*domain.ScoreByTicket{}
	for rows.Next() {
		s := &domain.ScoreByTicket{}
		err := rows.Scan(&s.TicketId, &s.Category, &s.Score)
		if err != nil {
			return nil, err
		}
		data = append(data, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return data, nil
}
