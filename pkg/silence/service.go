package silence

import (
	"context"
	"fmt"
	"time"

	"github.com/dacamposol/monitoring-operator/api/v1alpha1"
)

type Dao struct {
	ID        string             `json:"id"`
	Matchers  []v1alpha1.Matcher `json:"matchers"`
	StartsAt  time.Time          `json:"startsAt"`
	EndsAt    time.Time          `json:"endsAt"`
	CreatedBy string             `json:"createdBy,omitempty"`
	Comment   string             `json:"comment,omitempty"`
}

type service struct {
	HttpClient HttpClient
	URL        string
}

func NewService(httpClient HttpClient, url string) Writer {
	return &service{
		HttpClient: httpClient,
		URL:        url,
	}
}

func (r *service) UpdateSilence(ctx context.Context, src *Dao) error {
	if src == nil {
		return fmt.Errorf("cannot update Silence with nil pointer")
	}

	req, err := post[Dao](ctx, r.HttpClient, fmt.Sprintf("%s/silences", r.URL), *src)
	if err != nil {
		return fmt.Errorf("unable to update Silence in AlertManager: %w", err)
	}

	src.ID = req.ID
	src.Matchers = req.Matchers
	src.StartsAt = req.StartsAt
	src.EndsAt = req.EndsAt
	src.CreatedBy = req.CreatedBy
	src.Comment = req.Comment

	return nil
}

func (r *service) DeleteSilence(ctx context.Context, uuid string) error {
	if uuid == "" {
		return fmt.Errorf("cannot delete Silence with empty uuid")
	}

	err := remove(ctx, r.HttpClient, fmt.Sprintf("%s/silence/%s", r.URL, uuid))
	if err != nil {
		return fmt.Errorf("unable to delete Silence in AlertManager: %w", err)
	}

	return nil
}
