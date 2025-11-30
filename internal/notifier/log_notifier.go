package notifier

import (
	"context"
	"log"
)

type LogNotifier struct{}

func (l LogNotifier) Notify(ctx context.Context, msg string) error {
	log.Println("{ALERTA ISS}", msg)
	return nil
}