package silence

import "context"

//go:generate go run -mod=mod github.com/vektra/mockery/v2 --all --case=underscore --with-expecter

type Writer interface {
	UpdateSilence(ctx context.Context, src *Dao) error
	DeleteSilence(ctx context.Context, uuid string) error
}
