package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidPrice     = errors.New("invalid price")
	ErrInvalidStartDate = errors.New("invalid start date")
	ErrInvalidEndDate   = errors.New("invalid end date")
	ErrEmptyServiceName = errors.New("service name is empty")
	ErrInvalidDate      = errors.New("invalid date format, expected MM-YYYY")
	ErrCompareDate      = errors.New("start date cannot be after end date")
)

type Subscription struct {
	id          uuid.UUID
	serviceName string
	price       Price
	userId      uuid.UUID
	startDate   SubDate
	endDate     *SubDate
}

type Price int

type SubDate struct {
	time.Time
}

func ParseSubDate(s string) (*SubDate, error) {
	t, err := time.Parse("01-2006", s)
	if err != nil {
		return nil, ErrInvalidDate
	}

	t = time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
	return &SubDate{t}, nil
}

func NewSubDateFromTime(t time.Time) *SubDate {
	return &SubDate{t}
}

func (p Price) Validate() error {
	if p <= 0 {
		return ErrInvalidPrice
	}
	return nil
}

func (d SubDate) Validate() error {
	if d.Time.IsZero() {
		return ErrInvalidStartDate
	}

	now := time.Now()
	currentYear, currentMonth, _ := now.Date()

	year, month, _ := d.Date()

	if month < 1 || month > 12 {
		return ErrInvalidStartDate
	}

	if year < currentYear || (year == currentYear && month < currentMonth) {
		return ErrInvalidStartDate
	}

	return nil
}

func NewSubscription(
	id uuid.UUID,
	serviceName string,
	price Price,
	userId uuid.UUID,
	startDate SubDate,
	endDate *SubDate,
) (*Subscription, error) {

	if serviceName == "" {
		return nil, ErrEmptyServiceName
	}

	if err := price.Validate(); err != nil {
		return nil, err
	}

	if err := startDate.Validate(); err != nil {
		return nil, err
	}

	if endDate != nil {
		if err := endDate.Validate(); err != nil {
			return nil, ErrInvalidEndDate
		}
		if endDate.Before(startDate.Time) {
			return nil, ErrCompareDate
		}
	}

	if id == uuid.Nil {
		id = uuid.New()
	}

	return &Subscription{
		id:          id,
		serviceName: serviceName,
		price:       price,
		userId:      userId,
		startDate:   startDate,
		endDate:     endDate,
	}, nil
}

// ------------------- Getters ------------------

func (s *Subscription) ID() uuid.UUID {
	return s.id
}
func (s *Subscription) ServiceName() string {
	return s.serviceName
}
func (s *Subscription) Price() int {
	return int(s.price)
}
func (s *Subscription) UserID() uuid.UUID {
	return s.userId
}
func (s *Subscription) StartDate() time.Time {
	return s.startDate.Time
}
func (s *Subscription) StartDateStr() string {
	return s.startDate.Format("01-2006")
}
func (s *Subscription) EndDate() *time.Time {
	if s.endDate == nil {
		return nil
	}
	endDateTime := s.endDate.Time
	return &endDateTime
}
func (s *Subscription) EndDateStr() *string {
	if s.endDate == nil {
		return nil
	}

	endDateStr := s.endDate.Format("01-2006")
	return &endDateStr
}

// ----------------------------- Filters --------------------------
type SubscriptionFilter struct {
	UserID      *uuid.UUID
	ServiceName *string
	StartDate   *SubDate
	EndDate     *SubDate
	Limit       int
	Offset      int
}

func NewSubscriptionFilter(
	userID *uuid.UUID,
	serviceName *string,
	start *string,
	end *string,
	limit int,
	offset int,
) (*SubscriptionFilter, error) {
	var startDate *SubDate
	var endDate *SubDate
	var err error

	if start != nil {
		startDate, err = ParseSubDate(*start)
		if err != nil {
			return nil, err
		}
	}

	if end != nil {
		endDate, err = ParseSubDate(*end)
		if err != nil {
			return nil, err
		}
	}

	if startDate != nil && endDate != nil {
		if startDate.After(endDate.Time) {
			return nil, ErrCompareDate
		}
	}

	return &SubscriptionFilter{
		UserID:      userID,
		ServiceName: serviceName,
		StartDate:   startDate,
		EndDate:     endDate,
		Limit:       limit,
		Offset:      offset,
	}, nil
}

type SumResult struct {
	Rows       []*Subscription
	TotalSum   int
	TotalCount int
}

func NewSumResult(rows []*Subscription, totalSum, totalCount int) *SumResult {
	return &SumResult{
		Rows:       rows,
		TotalSum:   totalSum,
		TotalCount: totalCount,
	}
}

type PagingBase struct {
	Limit  int
	Offset int
}

func NewPagingBase(limit, offset int) *PagingBase {
	return &PagingBase{
		Limit:  limit,
		Offset: offset,
	}
}
