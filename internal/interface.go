package internal

import "context"

type ProviderInf interface {
	Open(ctx context.Context) error
	Close(ctx context.Context) error
}
