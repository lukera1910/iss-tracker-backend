package notifier

import "context"

type Notifier interface {
	Notify(ctx context.Context, msg string) error
}