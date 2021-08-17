package data

import "context"

type (
	Store interface {
		Get(context.Context, string) (string, error)
		Save(context.Context, string, string) error
	}
)
