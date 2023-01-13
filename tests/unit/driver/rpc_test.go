package tests

import (
	"context"
	"database/sql"
	"github.com/HashimovH/softwareengineer-test-task-go/app/core/repository"
	"github.com/HashimovH/softwareengineer-test-task-go/app/core/service"
	"testing"

	_ "github.com/mattn/go-sqlite3"

	driver "github.com/HashimovH/softwareengineer-test-task-go/app/driver/rpc"
	protos "github.com/HashimovH/softwareengineer-test-task-go/app/driver/rpc/protos/tickets_service"
)

var db, _ = sql.Open("sqlite3", "./../../test.db")
var scoreRepo = repository.NewRepository(db)
var qualityRepo = repository.NewQualityRepository(db)

var scoreService = service.NewService(scoreRepo)
var qualityService = service.NewQualityService(qualityRepo)

var rpc = driver.NewRPCAdapter(scoreService, qualityService)

func Test_GetAggregatedCategoryScoresRPC(t *testing.T) {
	input := &protos.DateRange{
		RangeFrom: "2019-07-17",
		RangeTo:   "2019-07-19",
	}
	response, err := rpc.GetAggregatedCategoryScores(context.Background(), input)

	if err != nil {
		t.Errorf("error was not expected while getting scores: %s", err)
	}
	if len(response.Scores) != 4 {
		t.Errorf("expected 4 categories in the score object, but got %d", len(response.Scores))
	}

	for index := range response.Scores {
		if len(response.Scores[index].DateScores) != 3 {
			t.Errorf("expected 3 Date score for each %s, but got %d", response.Scores[index].CategoryName, len(response.Scores[index].DateScores))
		}
	}

	if response.Scores[0].DateScores[0].Date != "2019-07-17" {
		t.Errorf("expected first date 2019-07-17, but got %s", response.Scores[0].DateScores[0].Date)
	}
	if response.Scores[0].DateScores[2].Date != "2019-07-19" {
		t.Errorf("expected last date 2019-07-19, but got %s", response.Scores[0].DateScores[2].Date)
	}
}

func Test_GetAggregatedCategoryScoresWeeklyRPC(t *testing.T) {
	input := &protos.DateRange{
		RangeFrom: "2019-07-17",
		RangeTo:   "2019-08-19",
	}
	response, err := rpc.GetAggregatedCategoryScores(context.Background(), input)

	if err != nil {
		t.Errorf("error was not expected while getting scores: %s", err)
	}
	if len(response.Scores) != 4 {
		t.Errorf("expected 4 categories in the score object, but got %d", len(response.Scores))
	}

	for index := range response.Scores {
		if len(response.Scores[index].DateScores) != 5 {
			t.Errorf("expected 27 Date score for each %s, but got %d", response.Scores[index].CategoryName, len(response.Scores[index].DateScores))
		}
	}

	if response.Scores[0].DateScores[0].Date != "2019-07-17" {
		t.Errorf("expected first date 2019-07-17, but got %s", response.Scores[0].DateScores[0].Date)
	}
	if response.Scores[0].DateScores[4].Date != "2019-08-14" {
		t.Errorf("expected last date 2019-07-19, but got %s", response.Scores[0].DateScores[4].Date)
	}
}

func Test_GetScoresByTicketRPC(t *testing.T) {
	input := &protos.DateRange{
		RangeFrom: "2019-07-17",
		RangeTo:   "2019-07-19",
	}
	response, err := rpc.GetScoresByTicket(context.Background(), input)

	if err != nil {
		t.Errorf("error was not expected while getting scores: %s", err)
	}
	if len(response.Scores) != 58 {
		t.Errorf("expected 58 Tickets in the score object, but got %d", len(response.Scores))
	}

	for index := range response.Scores {
		if len(response.Scores[index].CategoryScores) != 4 {
			t.Errorf("expected 4 Category for each ticket, but got %d", len(response.Scores[index].CategoryScores))
		}
	}
}

func Test_GetScoreOveralForQualityRPC(t *testing.T) {
	input := &protos.DateRange{
		RangeFrom: "2019-07-17",
		RangeTo:   "2019-07-19",
	}
	response, err := rpc.GetScoreOveralForQuality(context.Background(), input)

	if err != nil {
		t.Errorf("error was not expected while getting scores: %s", err)
	}
	if response.Score != int32(50) {
		t.Errorf("expected 50%% overal quality, but got %d", response.Score)
	}
}

func Test_GetScoreChangePeriodOverPeriodRPC(t *testing.T) {
	previous_date := &protos.DateRange{
		RangeFrom: "2019-07-01",
		RangeTo:   "2019-07-30",
	}
	end_date := &protos.DateRange{
		RangeFrom: "2019-08-01",
		RangeTo:   "2019-08-30",
	}

	input := &protos.PeriodRange{
		EndPeriod:      end_date,
		PreviousPeriod: previous_date,
	}

	response, err := rpc.GetScoreChangePeriodOverPeriod(context.Background(), input)

	if err != nil {
		t.Errorf("error was not expected while getting scores: %s", err)
	}
	if response.ChangeScore != int32(-4) {
		t.Errorf("expected -4%% overal quality difference, but got %d", response.ChangeScore)
	}
}
