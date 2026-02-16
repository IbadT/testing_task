package v1

import (
	domain "testingtask/internal/domain/subscription"
	"time"

	"github.com/google/uuid"
)

type SubscriptionDTO struct {
	ServiceName string    `json:"service_name"`
	Price       int       `json:"price"`
	UserID      uuid.UUID `json:"user_id"`
	StartDate   string    `json:"start_date"`
	EndDate     *string   `json:"end_date,omitempty"`
}

func NewSubscriptionDTO(
	serviceName string,
	price int,
	userId uuid.UUID,
	startDate string,
	endDate *string,
) *SubscriptionDTO {
	return &SubscriptionDTO{
		ServiceName: serviceName,
		Price:       price,
		UserID:      userId,
		StartDate:   startDate,
		EndDate:     endDate,
	}
}

type SubscriptionResponseDTO struct {
	ID          uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	ServiceName string    `json:"service_name" example:"Netflix"`
	Price       int       `json:"price" example:"99900"`
	UserID      uuid.UUID `json:"user_id" example:"987f6543-e21b-34d5-c678-426614174999"`
	StartDate   string    `json:"start_date" example:"07-2026"`
	EndDate     *string   `json:"end_date,omitempty" example:"08-2026"`
}

func NewSubscriptionResponseDTO(id uuid.UUID, serviceName string, price int, userId uuid.UUID, start time.Time, end *time.Time) *SubscriptionResponseDTO {

	startStr := start.Format("01-2006")

	var endStr *string
	if end != nil {
		s := end.Format("01-2006")
		endStr = &s
	}

	return &SubscriptionResponseDTO{
		ID:          id,
		ServiceName: serviceName,
		Price:       price,
		UserID:      userId,
		StartDate:   startStr,
		EndDate:     endStr,
	}
}

// -------------- Filters --------------

// Request

type ListSubscriptionsRequestDTO struct {
	UserID      *uuid.UUID `json:"user_id,omitempty"`
	ServiceName *string    `json:"service_name,omitempty"`
	Start       *string    `json:"start,omitempty"`
	End         *string    `json:"end,omitempty"`
	Limit       int        `json:"limit"`
	Offset      int        `json:"offset"`
}

func NewListSubscriptionsRequestDTO(
	userID *uuid.UUID,
	serviceName *string,
	start *string,
	end *string,
	limit int,
	offset int,
) ListSubscriptionsRequestDTO {
	if limit <= 0 {
		limit = 10
	}

	if offset <= 0 {
		offset = 0
	}

	return ListSubscriptionsRequestDTO{
		UserID:      userID,
		ServiceName: serviceName,
		Start:       start,
		End:         end,
		Limit:       limit,
		Offset:      offset,
	}
}

type PagingBase struct {
	Limit  int
	Offset int
}

func NewBasePaging(limit, offset int) *PagingBase {
	if limit <= 0 {
		limit = 10
	}

	if offset <= 0 {
		offset = 0
	}

	return &PagingBase{
		Limit:  limit,
		Offset: offset,
	}
}

// Response

type ListSubscriptionsResponseDto struct {
	Paging   Paging                    `json:"paging"`
	Rows     []SubscriptionResponseDTO `json:"rows"`
	TotalSum int                       `json:"total_sum" example:"15900"`
}

type Paging struct {
	Total  int `json:"total" example:"42"`
	Offset int `json:"offset" example:"5"`
	Limit  int `json:"limit" example:"10"`
}

func NewPagingDTO(limit, offset, total int) Paging {
	return Paging{
		Limit:  limit,
		Offset: offset,
		Total:  total,
	}
}

func DomainToSumDTO(
	res *domain.SumResult,
) ListSubscriptionsResponseDto {
	rows := make([]SubscriptionResponseDTO, 0, len(res.Rows))
	for _, d := range res.Rows {
		rows = append(rows, *SubscriptionToDTO(d))
	}

	return ListSubscriptionsResponseDto{
		Rows:     rows,
		TotalSum: res.TotalSum,
	}
}

type SubscriptionID struct {
	ID uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
}
