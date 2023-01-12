package repository

import (
	"database/sql"

	"github.com/HashimovH/softwareengineer-test-task-go/app/core/domain"
)

func NewRepository(db *sql.DB) *repository {
	return &repository{DB: db}
}

func (r *repository) GetAggregatedCategoryScores(from string, to string) ([]domain.Score, error) {
	query := `
		WITH dates AS (
			SELECT $1 as dt
			UNION ALL
			SELECT 
			CASE
				WHEN julianday($2) - julianday($1) > 31 THEN 
				date(dt, '+7 days')
				ELSE 
				date(dt, '+1 day')
			END AS dt
			FROM dates
			WHERE 
				CASE 
					WHEN julianday($2) - julianday($1) > 31 THEN 
						dt < date($2, '-7 days')
					ELSE 
						dt < $2
				END
		)
		SELECT dt,cats.name, count(r.rating), ROUND(AVG(r.rating)*20)
		FROM dates
		CROSS JOIN (SELECT id, name from rating_categories) as cats
		LEFT JOIN ratings r ON 
			CASE 
				WHEN julianday($2) - julianday($1) > 31 THEN 
					JULIANDAY(dt) - JULIANDAY(STRFTIME("%Y-%m-%d", r.created_at)) < 7
					AND JULIANDAY(dt) - JULIANDAY(STRFTIME("%Y-%m-%d", r.created_at)) >= 0
				ELSE 
					dt = STRFTIME("%Y-%m-%d", r.created_at) 
			END
			AND 
				cats.id = r.rating_category_id
		GROUP BY cats.name, dt;
	`

	rows, err := r.DB.Query(string(query), from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := []domain.Score{}
	for rows.Next() {
		i := domain.Score{}
		err = rows.Scan(&i.Date, &i.Category, &i.Rating, &i.Score)
		if err != nil {
			return nil, err
		}
		data = append(data, i)
	}

	return data, nil
}

func (r *repository) GetScoresByTicket(from string, to string) ([]*domain.ScoreByTicket, error) {
	query := `
		SELECT r.ticket_id, rc.name, ROUND(AVG(r.rating) * 20) as score
		FROM ratings r
		LEFT JOIN rating_categories rc ON rc.id = r.rating_category_id
		WHERE r.created_at BETWEEN :startDate AND :endDate
		GROUP BY r.ticket_id, r.rating_category_id;
	`

	rows, err := r.DB.Query(string(query), from, to)
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
