package service

import (
	"context"
	"errors"
	domain "testingtask/internal/domain/subscription"
	myerrors "testingtask/internal/errors"
	"testingtask/internal/repository"
	logger "testingtask/pkg"

	"github.com/google/uuid"
)

type SubService interface {
	Create(ctx context.Context, sub *domain.Subscription) (uuid.UUID, error)
	Get(ctx context.Context, id uuid.UUID) (*domain.Subscription, error)
	List(ctx context.Context, paging *domain.PagingBase) ([]*domain.Subscription, int64, error)
	Sum(ctx context.Context, filters *domain.SubscriptionFilter) (*domain.SumResult, error)
	Update(ctx context.Context, id uuid.UUID, sub *domain.Subscription) (uuid.UUID, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type subService struct {
	repo repository.SubRepository
}

func NewSubService(r repository.SubRepository) SubService {
	return &subService{repo: r}
}

func (s *subService) Create(ctx context.Context, sub *domain.Subscription) (uuid.UUID, error) {
	logger.Info(ctx, "service: creating subscription", map[string]interface{}{
		"user_id":      sub.UserID(),
		"service_name": sub.ServiceName(),
		"price":        sub.Price(),
		"start_date":   sub.StartDate(),
		"end_date":     sub.EndDate(),
	})
	if err := s.repo.Create(ctx, sub); err != nil {
		logger.Error(ctx, "service: create failed", err, map[string]interface{}{"user_id": sub.UserID()})
		return uuid.Nil, err
	}

	logger.Info(ctx, "service: subscription created", map[string]interface{}{
		"id": sub.ID(),
	})

	return sub.ID(), nil
}

func (s *subService) Get(ctx context.Context, id uuid.UUID) (*domain.Subscription, error) {
	logger.Debug(ctx, "service: getting subscription by ID", map[string]interface{}{
		"id": id,
	})

	sub, err := s.repo.Get(ctx, id)
	if err != nil {
		logger.Error(ctx, "service: get failed", err, map[string]interface{}{"id": id})
		return nil, err
	}

	return sub, nil
}

func (s *subService) List(ctx context.Context, paging *domain.PagingBase) ([]*domain.Subscription, int64, error) {
	logger.Debug(ctx, "service: getting list subscirptions", map[string]interface{}{
		"paging": paging,
	})
	subs, err := s.repo.List(ctx, paging)

	if err != nil {
		logger.Error(ctx, "service: list failed", err, map[string]interface{}{
			"limit":  paging.Limit,
			"offset": paging.Offset,
		})
		return nil, 0, err
	}

	totalCount, err := s.repo.Count(ctx)
	if err != nil {
		logger.Error(ctx, "service: total count failed", err, nil)
		return nil, 0, err
	}

	return subs, totalCount, nil
}

func (s *subService) Sum(ctx context.Context, filters *domain.SubscriptionFilter) (*domain.SumResult, error) {
	logger.Debug(ctx, "service: getting sum subscriptions prices", map[string]interface{}{
		"filters": filters,
	})

	rows, totalSum, err := s.repo.Sum(ctx, filters)
	if err != nil {
		logger.Error(ctx, "service: sum failed", err, map[string]interface{}{
			"user_id":      filters.UserID,
			"service_name": filters.ServiceName,
		})
		return nil, err
	}

	totalCount, err := s.repo.Count(ctx)
	if err != nil {
		logger.Error(ctx, "service: total count failed", err, nil)
		return nil, err
	}

	return &domain.SumResult{
		Rows:       rows,
		TotalSum:   totalSum,
		TotalCount: int(totalCount),
	}, nil
}

func (s *subService) Update(ctx context.Context, id uuid.UUID, sub *domain.Subscription) (uuid.UUID, error) {
	logger.Debug(ctx, "service: updating subscription", map[string]interface{}{
		"id":   id,
		"data": sub,
	})

	existingSub, err := s.repo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, myerrors.ErrNotFound) {
			logger.Warn(ctx, "service: subscription not found for update", map[string]interface{}{
				"id": id,
			})
			return uuid.Nil, myerrors.ErrNotFound
		}
		logger.Error(ctx, "service: failed to check subscription existence", err, map[string]interface{}{
			"id": id,
		})

		return uuid.Nil, err
	}

	if existingSub.ID() != sub.ID() {
		logger.Warn(ctx, "service: id mismatch", map[string]interface{}{
			"path_id": id,
			"body_id": sub.ID(),
		})
		return uuid.Nil, myerrors.ErrInvalidData
	}

	if err := s.repo.Update(ctx, sub); err != nil {
		logger.Error(ctx, "service: update failed", err, map[string]interface{}{
			"id": id,
		})
		return uuid.Nil, err
	}

	logger.Info(ctx, "service: subscription updated", map[string]interface{}{
		"id": id,
	})

	return id, nil
}

func (s *subService) Delete(ctx context.Context, id uuid.UUID) error {
	logger.Info(ctx, "service: deleting subscription", map[string]interface{}{
		"id": id,
	})

	if err := s.repo.Delete(ctx, id); err != nil {
		logger.Error(ctx, "service: delete failed", err, map[string]interface{}{
			"id": id,
		})
		return err
	}

	logger.Info(ctx, "service: subscription deleted", map[string]interface{}{
		"id": id,
	})

	return nil
}
