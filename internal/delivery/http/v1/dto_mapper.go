package v1

import (
	domain "testingtask/internal/domain/subscription"
	"testingtask/internal/web/subscriptions"

	"github.com/google/uuid"
)

func DTOToFilter(dto ListSubscriptionsRequestDTO) (*domain.SubscriptionFilter, error) {
	return domain.NewSubscriptionFilter(
		dto.UserID,
		dto.ServiceName,
		dto.Start,
		dto.End,
		dto.Limit,
		dto.Offset,
	)
}

func SubscriptionToDTO(d *domain.Subscription) *SubscriptionResponseDTO {
	return NewSubscriptionResponseDTO(
		d.ID(),
		d.ServiceName(),
		d.Price(),
		d.UserID(),
		d.StartDate(),
		d.EndDate(),
	)
}

// --------------------
// DTO -> Domain
// --------------------

func DTOToDomain(uid *uuid.UUID, dto SubscriptionDTO) (*domain.Subscription, error) {
	var id uuid.UUID
	if uid != nil {
		id = *uid
	} else {
		id = uuid.New()
	}
	start, err := domain.ParseSubDate(dto.StartDate)
	if err != nil {
		return nil, err
	}

	var end *domain.SubDate
	if dto.EndDate != nil {
		end, err = domain.ParseSubDate(*dto.EndDate)
		if err != nil {
			return nil, err
		}
	}

	return domain.NewSubscription(
		id,
		dto.ServiceName,
		domain.Price(dto.Price),
		dto.UserID,
		*start,
		end,
	)
}

func SumDTOToDomain(dto ListSubscriptionsRequestDTO) (*domain.SubscriptionFilter, error) {
	return domain.NewSubscriptionFilter(
		dto.UserID,
		dto.ServiceName,
		dto.Start,
		dto.End,
		dto.Limit,
		dto.Offset,
	)
}

func NewPagingBase(dto PagingBase) *domain.PagingBase {
	return domain.NewPagingBase(dto.Limit, dto.Offset)
}

func NewRows(l []SubscriptionResponseDTO) []subscriptions.Subscription {
	rows := make([]subscriptions.Subscription, 0, len(l))
	for _, r := range l {
		rows = append(rows, subscriptions.Subscription{
			Id:          r.ID,
			ServiceName: r.ServiceName,
			Price:       r.Price,
			UserId:      r.UserID,
			StartDate:   r.StartDate,
			EndDate:     r.EndDate,
		})
	}

	return rows
}

func SumPaging(p Paging) subscriptions.Paging {
	return subscriptions.Paging{
		Limit:  &p.Limit,
		Offset: &p.Offset,
		Total:  &p.Total,
	}
}

func ListRequestToDTO(req subscriptions.ListRequestObject) PagingBase {
	return PagingBase{
		Limit:  req.Params.Limit,
		Offset: req.Params.Offset,
	}
}

func CreateRequestToDTO(req subscriptions.CreateJSONRequestBody) *SubscriptionDTO {
	return NewSubscriptionDTO(
		req.ServiceName,
		req.Price,
		req.UserId,
		req.StartDate,
		req.EndDate,
	)
}

func UpdateRequestToDTO(req subscriptions.UpdateJSONRequestBody) *SubscriptionDTO {
	return NewSubscriptionDTO(
		req.ServiceName,
		req.Price,
		req.UserId,
		req.StartDate,
		req.EndDate,
	)
}

func SumRequestToDTO(req subscriptions.SumRequestObject) ListSubscriptionsRequestDTO {
	return NewListSubscriptionsRequestDTO(
		req.Params.UserId,
		req.Params.ServiceName,
		req.Params.Start,
		req.Params.End,
		*req.Params.Limit,
		*req.Params.Offset,
	)
}

// --------------------
// Domain -> DTO
// --------------------

func DomainToDTO(d *domain.Subscription) *SubscriptionResponseDTO {
	return NewSubscriptionResponseDTO(
		d.ID(),
		d.ServiceName(),
		d.Price(),
		d.UserID(),
		d.StartDate(),
		d.EndDate(),
	)
}

func DomainListToDTO(list []*domain.Subscription) []SubscriptionResponseDTO {
	res := make([]SubscriptionResponseDTO, 0, len(list))
	for _, s := range list {
		res = append(res, *DomainToDTO(s))
	}
	return res
}

// --------------------
// DTO -> subscriptions (response)
// --------------------

func SumDTOToResponse(l ListSubscriptionsResponseDto, paging Paging) subscriptions.Sum200JSONResponse {
	rows := NewRows(l.Rows)
	p := SumPaging(paging)
	return subscriptions.Sum200JSONResponse{
		Paging:   p,
		TotalSum: l.TotalSum,
		Rows:     rows,
	}
}

func ListDTOToResponse(s []SubscriptionResponseDTO, dto PagingBase, totalCount int64) subscriptions.List200JSONResponse {
	rows := NewRows(s)
	total := int(totalCount)
	return subscriptions.List200JSONResponse{
		Paging: &subscriptions.Paging{
			Limit:  &dto.Limit,
			Offset: &dto.Offset,
			Total:  &total,
		},
		Rows: &rows,
	}
}

func GetDTOToResponse(s *SubscriptionResponseDTO) subscriptions.Get200JSONResponse {
	return subscriptions.Get200JSONResponse{
		Id:          s.ID,
		ServiceName: s.ServiceName,
		UserId:      s.UserID,
		Price:       s.Price,
		StartDate:   s.StartDate,
		EndDate:     s.EndDate,
	}
}

func UpdateToResponse(id uuid.UUID) subscriptions.Update200JSONResponse {
	return subscriptions.Update200JSONResponse{Id: &id}
}

func CreateToResponse(id uuid.UUID) subscriptions.Create201JSONResponse {
	return subscriptions.Create201JSONResponse{Id: &id}
}
