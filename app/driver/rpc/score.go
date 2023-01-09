package driver

import (
	"context"
	// "fmt"
	// "time"

	"github.com/HashimovH/softwareengineer-test-task-go/app/core/domain"
	protos "github.com/HashimovH/softwareengineer-test-task-go/app/driver/rpc/protos/tickets_service"
	// "github.com/golang/protobuf/ptypes/timestamp"
	"github.com/hashicorp/go-hclog"
)

type RatingService interface {
	GetAggregatedCategoryScoresService(from string, to string) ([]domain.Score, error)
	GetScoresByTicketInRangeService(from string, to string) ([]*domain.ScoreByTicket, error)
}

type QualityService interface {
	GetOveralQualityService(from string, to string) (*domain.OveralQuality, error)
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

func (rpc *RPCAdapter) GetScoreOveralForQuality(ctx context.Context, rr *protos.DateRange) (*protos.QualityResponse, error) {
	overal_score, err := rpc.qualityService.GetOveralQualityService(rr.GetRangeFrom(), rr.GetRangeTo())
	if err != nil {
		logger.Error(err.Error())
	}
	return &protos.QualityResponse{
		Score: overal_score.OveralScore,
	}, nil
}

func (rpc *RPCAdapter) GetScoresByTicket(ctx context.Context, rr *protos.DateRange) (*protos.ScoresByTicketResponse, error) {
	scores, err := rpc.ratingService.GetScoresByTicketInRangeService(rr.GetRangeFrom(), rr.GetRangeTo())
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
	scores, err := rpc.ratingService.GetAggregatedCategoryScoresService(rr.GetRangeFrom(), rr.GetRangeTo())
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	response := &protos.ScoresByCategoryResponse{}
	categories := make(map[string][]domain.Score)

	// Group scores by category
	for _, score := range scores {
		categories[score.Category] = append(categories[score.Category], score)
	}

	// Convert each category to a ScoresByCategory message
	for category, categoryScores := range categories {
		scoresByCategory := &protos.ScoresByCategory{
			CategoryName: category,
			RatingsCount: int32(len(categoryScores)),
			TotalScore:   0,
			DateScores:   []*protos.DateScore{},
		}

		// Convert each score to a DateScore message
		for _, score := range categoryScores {
			dateScore := &protos.DateScore{
				Date:  score.Date,
				Score: int32(score.Score),
			}
			scoresByCategory.DateScores = append(scoresByCategory.DateScores, dateScore)
			scoresByCategory.TotalScore += int32(score.Score)
		}

		response.Scores = append(response.Scores, scoresByCategory)
	}

	return response, nil

	// scoresByCategory := make(map[string]map[time.Time][]domain.Score)
	// for _, score := range scores {
	// 	if _, ok := scoresByCategory[score.Category]; !ok {
	// 		scoresByCategory[score.Category] = make(map[time.Time][]domain.Score)
	// 	}
	// 	parsed_date, _ := time.Parse(score.Date, "2022-07-28T09:30:00Z")
	// 	date := parsed_date.Truncate(24 * time.Hour)
	// 	scoresByCategory[score.Category][date] = append(scoresByCategory[score.Category][date], score)
	// }

	// for categoryName, dateScoresMap := range scoresByCategory {
	// 	// Create a set of DateScore values
	// 	dateScores := make([]*protos.DateScore, 0, len(dateScoresMap))
	// 	var categoryScores []domain.Score
	// 	for date, scores := range dateScoresMap {
	// 		// Convert the date to a Timestamp
	// 		timestamp := &timestamp.Timestamp{
	// 			Seconds: date.Unix(),
	// 			Nanos:   int32(date.UnixNano() % 1e9),
	// 		}

	// 		// Calculate the average score for the date
	// 		var totalScore int
	// 		for _, score := range scores {
	// 			totalScore += score.Score
	// 		}
	// 		averageScore := int32(totalScore) / int32(len(scores))

	// 		// Add a DateScore value to the set
	// 		dateScores = append(dateScores, &protos.DateScore{
	// 			Date:  timestamp,
	// 			Score: averageScore,
	// 		})
	// 		categoryScores = append(categoryScores, scores...)
	// 	}

	// 	var totalScore int32

	// 	for _, score := range categoryScores {
	// 		totalScore += int32(score.Score)
	// 	}
	// 	averageScore := totalScore / int32(len(categoryScores))
	// 	fmt.Println("Avg Score", averageScore)

	// 	// Build the CategoryResultResponse
	// 	response = &protos.ScoresByCategoryResponse{
	// 		Scores: &protos.ScoresByCategory{
	// 			CategoryName: categoryName,
	// 			RatingsCount: int32(len(categoryScores)),
	// 			DateScores:   dateScores,
	// 			TotalScore:   averageScore,
	// 		},
	// 	}

	// }
	// return response, nil
}
