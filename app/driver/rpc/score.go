package driver

import (
	"context"
	"errors"

	"github.com/HashimovH/softwareengineer-test-task-go/app/core/domain"
	protos "github.com/HashimovH/softwareengineer-test-task-go/app/driver/rpc/protos/tickets_service"

	"github.com/hashicorp/go-hclog"
)

type RatingService interface {
	GetAggregatedCategoryScoresService(from string, to string) ([]domain.Score, error)
	GetScoresByTicketInRangeService(from string, to string) ([]*domain.ScoreByTicket, error)
}

type QualityService interface {
	GetOveralQualityService(from string, to string) (*domain.OveralQuality, error)
	GetScoreChangePeriodOverPeriod(current_from string, current_to string, previous_from string, previous_to string) (*domain.PeriodScoreChange, error)
}

type RPCAdapter struct {
	ratingService  RatingService
	qualityService QualityService
}

func NewRPCAdapter(rS RatingService, qS QualityService) *RPCAdapter {
	return &RPCAdapter{ratingService: rS, qualityService: qS}
}

var logger = hclog.New(&hclog.LoggerOptions{
	Name:  "tickets-service",
	Level: hclog.LevelFromString("DEBUG"),
})

func (rpc *RPCAdapter) GetScoreChangePeriodOverPeriod(ctx context.Context, rr *protos.PeriodRange) (*protos.ChangeOverPeriodResponse, error) {
	from := rr.GetEndPeriod().GetRangeFrom()
	to := rr.GetEndPeriod().GetRangeTo()

	prev_from := rr.GetPreviousPeriod().GetRangeFrom()
	prev_to := rr.GetPreviousPeriod().GetRangeTo()

	if from == "" || to == "" || prev_from == "" || prev_to == "" {
		logger.Error("Input Validation Error")
		return nil, errors.New("validation Error: date range is not correct")
	}

	score_difference, err := rpc.qualityService.GetScoreChangePeriodOverPeriod(from, to, prev_from, prev_to)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	return &protos.ChangeOverPeriodResponse{
		ChangeScore: score_difference.ScoreChange,
	}, nil
}

func (rpc *RPCAdapter) GetScoreOveralForQuality(ctx context.Context, rr *protos.DateRange) (*protos.QualityResponse, error) {
	from := rr.GetRangeFrom()
	to := rr.GetRangeTo()

	if from == "" || to == "" {
		logger.Error("Input Validation Error")
		return nil, errors.New("validation Error: date range is not correct")
	}

	overal_score, err := rpc.qualityService.GetOveralQualityService(from, to)
	if err != nil {
		logger.Error(err.Error())
	}
	return &protos.QualityResponse{
		Score: overal_score.OveralScore,
	}, nil
}

func (rpc *RPCAdapter) GetScoresByTicket(ctx context.Context, rr *protos.DateRange) (*protos.ScoresByTicketResponse, error) {
	from := rr.GetRangeFrom()
	to := rr.GetRangeTo()

	if from == "" || to == "" {
		logger.Error("Input Validation Error")
		return nil, errors.New("validation Error: date range is not correct")
	}

	scores, err := rpc.ratingService.GetScoresByTicketInRangeService(from, to)
	if err != nil {
		logger.Error(err.Error())
	}

	converted := map[int32]*protos.ScoresByTicket{}
	for _, s := range scores {
		if _, ok := converted[s.TicketId]; !ok {
			converted[s.TicketId] = &protos.ScoresByTicket{
				TicketId: s.TicketId,
				CategoryScores: []*protos.CategoryAndScorePairs{
					{
						CategoryName: s.Category,
						Score:        s.Score,
					},
				},
			}
		} else {
			converted[s.TicketId].CategoryScores = append(converted[s.TicketId].CategoryScores, &protos.CategoryAndScorePairs{
				CategoryName: s.Category,
				Score:        s.Score,
			})
		}
	}

	response := &protos.ScoresByTicketResponse{
		Scores: make([]*protos.ScoresByTicket, 0, len(converted)),
	}
	for _, v := range converted {
		response.Scores = append(response.Scores, v)
	}

	return response, nil

}

func (rpc *RPCAdapter) GetAggregatedCategoryScores(ctx context.Context, rr *protos.DateRange) (*protos.ScoresByCategoryResponse, error) {

	from := rr.GetRangeFrom()
	to := rr.GetRangeTo()

	if from == "" || to == "" {
		logger.Error("Input Validation Error")
		return nil, errors.New("validation Error: date range is not correct")
	}
	scores, err := rpc.ratingService.GetAggregatedCategoryScoresService(from, to)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	result := &protos.ScoresByCategoryResponse{}
	// map to store the intermediate data
	categoryData := make(map[string]*protos.ScoresByCategory)
	// iterate through the input scores
	for _, score := range scores {
		_, ok := categoryData[score.Category]
		if !ok {
			// if category does not exist, create a new entry
			categoryData[score.Category] = &protos.ScoresByCategory{
				CategoryName: score.Category,
				RatingsCount: 0,
				DateScores:   []*protos.DateScore{},
				TotalScore:   0,
			}
		}
		// increment the count
		categoryData[score.Category].RatingsCount++
		// append score to date_scores
		categoryData[score.Category].DateScores = append(categoryData[score.Category].DateScores,
			&protos.DateScore{Date: score.Date, Score: score.Score})
	}
	// convert the map to array of ScoresByCategory
	for _, value := range categoryData {
		result.Scores = append(result.Scores, value)
	}
	return result, nil

}
