package repository

import (
	"context"
	"database/sql"
	"errors"
	domain "testingtask/internal/domain/subscription"
	myerrors "testingtask/internal/errors"
	"testingtask/internal/repository/models"
	logger "testingtask/pkg"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type SubRepository interface {
	Create(ctx context.Context, sub *domain.Subscription) error
	Get(ctx context.Context, id uuid.UUID) (*domain.Subscription, error)
	List(ctx context.Context, paging *domain.PagingBase) ([]*domain.Subscription, error)
	Sum(ctx context.Context, filter *domain.SubscriptionFilter) ([]*domain.Subscription, int, error)
	Count(ctx context.Context) (int64, error)
	Update(ctx context.Context, sub *domain.Subscription) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type subRepository struct {
	DB *gorm.DB
}

func NewSubRepository(db *gorm.DB) SubRepository {
	return &subRepository{DB: db}
}

func (s *subRepository) Create(ctx context.Context, sub *domain.Subscription) error {
	m := models.FromDomain(sub)

	err := s.DB.WithContext(ctx).Create(&m).Error
	if err != nil {
		logger.Error(ctx, "repo: subscription create failed", err, map[string]interface{}{
			"data": sub,
		})
		if errors.Is(err, gorm.ErrInvalidData) {
			return myerrors.ErrInvalidData
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				return myerrors.ErrInvalidData
			case "23503":
				return myerrors.ErrInvalidData
			}
		}

		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return myerrors.ErrDatabase
		}

		return myerrors.ErrCreateFailed
	}

	return nil
}

func (s *subRepository) Get(ctx context.Context, id uuid.UUID) (*domain.Subscription, error) {
	var m models.Subscription

	err := s.DB.WithContext(ctx).First(&m, "id = ?", id).Error
	if err != nil {
		logger.Error(ctx, "repo: subscription get failed", err, map[string]interface{}{
			"id": id,
		})
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, myerrors.ErrNotFound
		}

		if errors.Is(err, gorm.ErrInvalidData) {
			return nil, myerrors.ErrInvalidData
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return nil, myerrors.ErrDatabase
		}

		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return nil, myerrors.ErrDatabase
		}

		return nil, myerrors.ErrDatabase
	}

	return models.ToDomain(&m)
}

func (s *subRepository) List(ctx context.Context, paging *domain.PagingBase) ([]*domain.Subscription, error) {
	var m []*models.Subscription

	err := s.DB.WithContext(ctx).
		Limit(paging.Limit).
		Offset(paging.Offset).
		Find(&m).Error

	if err != nil {
		logger.Error(ctx, "repo: subscription list failed", err, map[string]interface{}{
			"paging": paging,
		})
		if errors.Is(err, gorm.ErrInvalidData) {
			return nil, myerrors.ErrInvalidData
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return nil, myerrors.ErrDatabase
		}

		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return nil, myerrors.ErrDatabase
		}

		return nil, myerrors.ErrListFailed
	}

	return models.ToDomains(m)
}

func (r *subRepository) Sum(ctx context.Context, filter *domain.SubscriptionFilter) ([]*domain.Subscription, int, error) {
	var totalSum sql.NullInt64

	sumQuery := r.DB.WithContext(ctx).Model(&models.Subscription{})

	if filter.UserID != nil {
		sumQuery = sumQuery.Where("user_id = ?", *filter.UserID)
	}
	if filter.ServiceName != nil {
		sumQuery = sumQuery.Where("service_name = ?", *filter.ServiceName)
	}
	if filter.StartDate != nil {
		sumQuery = sumQuery.Where("start_date >= ?", filter.StartDate.Time)
	}
	if filter.EndDate != nil {
		sumQuery = sumQuery.Where("end_date <= ?", filter.EndDate.Time)
	}

	if err := sumQuery.Select("SUM(price)").Scan(&totalSum).Error; err != nil {
		logger.Error(ctx, "repo: subscription sum failed", err, map[string]interface{}{
			"filter": filter,
		})
		if errors.Is(err, gorm.ErrInvalidData) {
			return nil, 0, myerrors.ErrInvalidData
		}
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			return nil, 0, myerrors.ErrDatabase
		}

		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return nil, 0, myerrors.ErrDatabase
		}

		return nil, 0, myerrors.ErrDatabase
	}

	var m []*models.Subscription

	rowsQuery := r.DB.WithContext(ctx).Model(&models.Subscription{})

	if filter.UserID != nil {
		rowsQuery = rowsQuery.Where("user_id = ?", *filter.UserID)
	}
	if filter.ServiceName != nil {
		rowsQuery = rowsQuery.Where("service_name = ?", *filter.ServiceName)
	}
	if filter.StartDate != nil {
		rowsQuery = rowsQuery.Where("start_date >= ?", filter.StartDate.Time)
	}
	if filter.EndDate != nil {
		rowsQuery = rowsQuery.Where("end_date <= ?", filter.EndDate.Time)
	}

	if err := rowsQuery.
		Limit(filter.Limit).
		Offset(filter.Offset).
		Find(&m).Error; err != nil {

		logger.Error(ctx, "repo: subscription rows fetch failed", err, map[string]interface{}{
			"filter": filter,
		})

		if errors.Is(err, gorm.ErrInvalidData) {
			return nil, 0, myerrors.ErrInvalidData
		}

		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			return nil, 0, myerrors.ErrDatabase
		}

		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return nil, 0, myerrors.ErrDatabase
		}

		return nil, 0, myerrors.ErrDatabase
	}

	rows, err := models.ToDomains(m)
	if err != nil {
		logger.Error(ctx, "subscription convert failed", err, map[string]interface{}{
			"models": m,
		})
		return nil, 0, myerrors.ErrDatabase
	}

	sum := 0
	if totalSum.Valid {
		sum = int(totalSum.Int64)
	}

	return rows, sum, nil
}

func (s *subRepository) Count(ctx context.Context) (int64, error) {
	var count int64

	query := s.DB.WithContext(ctx).Model(&models.Subscription{})

	err := query.Count(&count).Error
	if err != nil {
		logger.Error(ctx, "repo: subscription count failed", err, nil)
		return 0, myerrors.ErrDatabase
	}

	return count, nil
}

func (s *subRepository) Update(ctx context.Context, sub *domain.Subscription) error {
	m := models.FromDomain(sub)

	err := s.DB.WithContext(ctx).
		Model(&models.Subscription{}).
		Where("id = ?", sub.ID()).
		Select("*").
		Updates(m).
		Error

	if err != nil {

		logger.Error(ctx, "repo: subscription update failed", err, map[string]interface{}{
			"data": sub,
		})

		if errors.Is(err, gorm.ErrInvalidData) {
			return myerrors.ErrInvalidData
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				return myerrors.ErrInvalidData
			case "23503":
				return myerrors.ErrInvalidData
			default:
				return myerrors.ErrDatabase
			}
		}

		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return myerrors.ErrDatabase
		}

		return myerrors.ErrUpdateFailed
	}

	return nil
}

func (s *subRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.DB.WithContext(ctx).Delete(&models.Subscription{}, id).Error
	if err != nil {

		logger.Error(ctx, "repo: subscription delete failed", err, map[string]interface{}{
			"id": id,
		})

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return myerrors.ErrNotFound
		}

		if errors.Is(err, gorm.ErrInvalidData) {
			return myerrors.ErrInvalidData
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return myerrors.ErrDatabase
		}

		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return myerrors.ErrDatabase
		}

		return myerrors.ErrDeleteFailed
	}

	return nil
}
