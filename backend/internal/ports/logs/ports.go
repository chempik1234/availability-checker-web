package logs

import (
	"context"
	"github.com/chempik1234/availability-checker-web/internal/models"
	"time"
)

type LogRecordRepository interface {
	ListByName(ctx context.Context, nameFilter string) ([]models.LogRecord, error)
	Create(ctx context.Context, record models.LogRecord) error
	ClearAllBeforeDatetime(ctx context.Context, datetime time.Time) error
}
