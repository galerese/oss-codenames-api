package domain_util

import (
	"galere.se/oss-codenames-api/pkg/logging"
)

type BaseService struct {
	Logger *logging.Logger
}

func NewBaseService(logger *logging.Logger) *BaseService {
	return &BaseService{Logger: logger}
}
