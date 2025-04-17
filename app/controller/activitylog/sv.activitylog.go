package activitylog

import (
	"context"

	"app/app/model" // Ensure to import the model package
	"app/internal/logger"
)

func (s *Service) Create(ctx context.Context, req model.ActivityLog) (*model.ActivityLog, error) {
	if _, err := s.db.NewInsert().Model(&req).Exec(ctx); err != nil {
		logger.Infof("[error]: %v", err)
		return nil, err
	}
	return &req, nil
}
