package tests

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"

	"testing"

	// sqlmock "github.com/DATA-DOG/go-sqlmock"
	// "github.com/HashimovH/softwareengineer-test-task-go/app/config"
	"github.com/HashimovH/softwareengineer-test-task-go/app/core/repository"
)

var db, _ = sql.Open("sqlite3", "./../../test.db")
var scoreRepo = repository.NewRepository(db)
var qualityRepo = repository.NewQualityRepository(db)

func Test_GetAggregatedCategoryScores(t *testing.T) {
	// Create a new repository and call the function
	scores, err := scoreRepo.GetAggregatedCategoryScores("2019-07-17", "2019-07-19")

	// Assert the results
	if err != nil {
		t.Errorf("error was not expected while getting scores: %s", err)
	}
	if len(scores) != 12 {
		t.Errorf("expected 12 scores, but got %d", len(scores))
	}

	if scores[0].Date != "2019-07-17" {
		t.Errorf("expected 2019-07-17 as first score's date, but got %s", scores[0].Date)
	}
	if scores[0].Category != "GDPR" {
		t.Errorf("expected GDPR as first score's category, but got %s", scores[0].Category)
	}
	if scores[0].Rating != 27 {
		t.Errorf("expected 10 as first score's rating, but got %d", scores[0].Rating)
	}
	if *scores[0].Score != int32(52) {
		t.Errorf("expected 52 as first score's score, but got %d", scores[0].Score)
	}
}

func Test_GetAggregatedCategoryScoresWeekly(t *testing.T) {
	// Create a new repository and call the function
	scores, err := scoreRepo.GetAggregatedCategoryScores("2019-07-17", "2019-08-19")

	// Assert the results
	if err != nil {
		t.Errorf("error was not expected while getting scores: %s", err)
	}
	if len(scores) != 20 {
		t.Errorf("expected 20 scores, but got %d", len(scores))
	}

	if scores[0].Date != "2019-07-17" && scores[1].Date != "2019-07-24" {
		t.Errorf("expected 2019-07-17 and 2019-07-24 as first score's date with 1 week difference, but got %s and %s", scores[0].Date, scores[1].Date)
	}
	if scores[0].Category != "GDPR" {
		t.Errorf("expected GDPR as first score's category, but got %s", scores[0].Category)
	}
	if scores[0].Rating != 160 {
		t.Errorf("expected 160 as first score's rating, but got %d", scores[0].Rating)
	}
	if *scores[0].Score != int32(54) {
		t.Errorf("expected 54 as first score's score, but got %d", scores[0].Score)
	}
}

func Test_GetScoresByTicket(t *testing.T) {
	// Create a new repository and call the function
	scores, err := scoreRepo.GetScoresByTicket("2019-07-17", "2019-07-19")

	// Assert the results
	if err != nil {
		t.Errorf("error was not expected while getting scores: %s", err)
	}
	if len(scores) != 232 {
		t.Errorf("expected 232 scores, but got %d", len(scores))
	}

	if scores[0].TicketId != int32(19479) {
		t.Errorf("expected 19479 as first ticket id, but got %d", scores[0].TicketId)
	}
	if scores[0].Category != "Spelling" && scores[1].Category != "Grammar" && scores[2].Category != "GDPR" && scores[3].Category != "Randomness" {
		t.Errorf("expected ordering Spelling, Grammar, GDPR, Randomness first ticket's categories")
	}
	if scores[0].Score != 60 {
		t.Errorf("expected 60 as first ticket's score, but got %d", scores[0].Score)
	}
}

func Test_GetOveralScore(t *testing.T) {
	// Create a new repository and call the function
	qualityScore, err := qualityRepo.GetOveralQualityScore("2019-07-17", "2019-07-19")

	// Assert the results
	if err != nil {
		t.Errorf("error was not expected while getting scores: %s", err)
	}
	if qualityScore.OveralScore != int32(50) {
		t.Errorf("expected 50, but got %d", qualityScore.OveralScore)
	}
}
