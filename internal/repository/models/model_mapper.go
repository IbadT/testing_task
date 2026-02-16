package models

import (
	domain "testingtask/internal/domain/subscription"
)

// --------------------
// Domain -> Model
// --------------------

func FromDomain(d *domain.Subscription) *Subscription {
	return &Subscription{
		ID:          d.ID(),
		ServiceName: d.ServiceName(),
		Price:       d.Price(),
		UserID:      d.UserID(),
		StartDate:   d.StartDate(),
		EndDate:     d.EndDate(),
	}
}

// --------------------
// Model -> Domain
// --------------------

func ToDomain(m *Subscription) (*domain.Subscription, error) {
	var endDate *domain.SubDate
	if m.EndDate != nil {
		endDate = &domain.SubDate{Time: *m.EndDate}
	}

	return domain.NewSubscription(
		m.ID,
		m.ServiceName,
		domain.Price(m.Price),
		m.UserID,
		domain.SubDate{Time: m.StartDate},
		endDate,
	)
}

func ToDomains(models []*Subscription) ([]*domain.Subscription, error) {
	res := make([]*domain.Subscription, 0, len(models))
	for _, m := range models {
		d, err := ToDomain(m)
		if err != nil {
			return nil, err
		}
		res = append(res, d)
	}
	return res, nil
}
