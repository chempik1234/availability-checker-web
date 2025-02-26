package logsadapters

import (
	"context"
	"fmt"
	"github.com/chempik1234/availability-checker-web/internal/models"
	"github.com/chempik1234/availability-checker-web/pkg/storage/postgres"
	"github.com/jackc/pgx/v5"
	"time"
)

type LogRecordRepositoryDB struct {
	DBInstance *postgres.DBInstance
}

func NewLogRecordRepositoryDB(DBInstance *postgres.DBInstance) *LogRecordRepositoryDB {
	return &LogRecordRepositoryDB{DBInstance: DBInstance}
}

func (l LogRecordRepositoryDB) ListByName(ctx context.Context, nameFilter string) ([]models.LogRecord, error) {
	var query string
	var args pgx.NamedArgs
	if len(nameFilter) > 0 {
		query = `SELECT "name", "result", "datetime" from log_records WHERE "name"=@nameFilter ORDER BY "datetime" DESC`
		args = pgx.NamedArgs{"nameFilter": nameFilter}
	} else {
		query = `SELECT "name", "result", "datetime" from log_records ORDER BY "datetime" DESC`
	}

	var result []models.LogRecord
	rows, err := l.DBInstance.Db.Query(ctx, query, args)
	if err != nil {
		return result, fmt.Errorf("unable to select rows: %w", err)
	}

	defer rows.Close()

	result = make([]models.LogRecord, 0)
	for rows.Next() {
		record := models.LogRecord{}
		err = rows.Scan(&record.Name, &record.Result, &record.Datetime)
		if err != nil {
			return result, fmt.Errorf("unable to scan row: %w", err)
		}
		result = append(result, record)
	}
	return result, nil
}

func (l LogRecordRepositoryDB) Create(ctx context.Context, record models.LogRecord) error {
	query := `INSERT INTO log_records ("name", "result", "datetime") VALUES (@name, @result, @datetime)`
	args := pgx.NamedArgs{
		"name":     record.Name,
		"result":   record.Result,
		"datetime": record.Datetime,
	}
	_, err := l.DBInstance.Db.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("unable to insert row: %w", err)
	}
	return nil
}

func (l LogRecordRepositoryDB) ClearAllBeforeDatetime(ctx context.Context, datetime time.Time) error {
	var query string
	var args pgx.NamedArgs
	if !datetime.IsZero() {
		query = `DELETE FROM log_records WHERE "datetime" < @datetime`
		args = pgx.NamedArgs{"datetime": datetime}
	} else {
		query = `DELETE FROM log_records`
	}
	_, err := l.DBInstance.Db.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("unable to clear before datetime: %w", err)
	}
	return nil
}
