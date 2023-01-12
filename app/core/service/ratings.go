package service

import "github.com/HashimovH/softwareengineer-test-task-go/app/core/domain"

type RatingRepository interface {
	GetAggregatedCategoryScores(from string, to string) ([]domain.Score, error)
	GetScoresByTicket(from string, to string) ([]*domain.ScoreByTicket, error)
}

type RatingService struct {
	userRepository RatingRepository
}

func NewService(ur RatingRepository) *RatingService {
	return &RatingService{
		userRepository: ur,
	}
}

func (service RatingService) GetAggregatedCategoryScoresService(from string, to string) ([]domain.Score, error) {
	return service.userRepository.GetAggregatedCategoryScores(from, to)
}

func (service RatingService) GetScoresByTicketInRangeService(from string, to string) ([]*domain.ScoreByTicket, error) {
	return service.userRepository.GetScoresByTicket(from, to)
}
