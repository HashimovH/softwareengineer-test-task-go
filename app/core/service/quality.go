package service

import "github.com/HashimovH/softwareengineer-test-task-go/app/core/domain"

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
