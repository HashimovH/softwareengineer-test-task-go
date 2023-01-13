package service

import (
	"github.com/HashimovH/softwareengineer-test-task-go/app/core/domain"
)

type QualityRepository interface {
	GetOveralQualityScore(from string, to string) (*domain.OveralQuality, error)
}

type QualityService struct {
	qualityRepository QualityRepository
}

func NewQualityService(qR QualityRepository) *QualityService {
	return &QualityService{
		qualityRepository: qR,
	}
}

func (service QualityService) GetOveralQualityService(from string, to string) (*domain.OveralQuality, error) {
	return service.qualityRepository.GetOveralQualityScore(from, to)
}

func (service QualityService) GetScoreChangePeriodOverPeriod(current_from string, current_to string, previous_from string, previous_to string) (*domain.PeriodScoreChange, error) {
	previous, err := service.qualityRepository.GetOveralQualityScore(previous_from, previous_to)
	if err != nil {
		return nil, err
	}
	current, err := service.qualityRepository.GetOveralQualityScore(current_from, current_to)
	if err != nil {
		return nil, err
	}

	difference := current.OveralScore - previous.OveralScore
	percentage := float64(difference) / float64(previous.OveralScore) * 100

	return &domain.PeriodScoreChange{ScoreChange: int32(percentage)}, nil
}
