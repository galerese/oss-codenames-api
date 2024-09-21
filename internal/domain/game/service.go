package game

import (
	"galere.se/oss-codenames-api/pkg/domain_util"
	"galere.se/oss-codenames-api/pkg/logging"
)

type Service struct {
	repository  GameRepository
	baseService *domain_util.BaseService
}

func NewService(repository GameRepository, logger *logging.Logger) *Service {
	return &Service{
		repository:  repository,
		baseService: domain_util.NewBaseService(logger),
	}
}
