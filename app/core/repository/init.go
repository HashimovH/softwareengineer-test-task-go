package repository

import "database/sql"

type TicketRepository struct {
	DB *sql.DB
}
