package api

import "context"

type ContainerClient interface {
	Create(context.Context, string) error
	Start(context.Context, string) error
	Stop(context.Context, string) error
	Remove(context.Context, string) error
}
