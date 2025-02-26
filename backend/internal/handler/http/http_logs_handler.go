package http

import (
	"encoding/json"
	"fmt"
	"github.com/chempik1234/availability-checker-web/internal/models"
	"github.com/chempik1234/availability-checker-web/internal/ports/logs"
	"log"
	"net/http"
	"time"
)

type LogsHttpHandler struct {
	logsRepository logs.LogRecordRepository
}

func NewLogsHttpHandler(logsRepository logs.LogRecordRepository) LogsHttpHandler {
	return LogsHttpHandler{logsRepository}
}

func (h LogsHttpHandler) NewReceiveLogsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			decoder := json.NewDecoder(r.Body)
			var newLogRecord models.LogRecord
			err := decoder.Decode(&newLogRecord)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte(fmt.Errorf("Failed to serialize response: %v", err.Error()).Error()))
				break
			}
			err = h.logsRepository.Create(r.Context(), newLogRecord)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(fmt.Errorf("Failed to save response: %v", err.Error()).Error()))
				break
			}
			w.WriteHeader(http.StatusNoContent)
			log.Printf("Received new log record %v\n", newLogRecord)
		case http.MethodGet, http.MethodPut, http.MethodPatch, http.MethodDelete:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func (h LogsHttpHandler) NewLogsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			nameFilter := r.URL.Query().Get("name_filter")
			logsList, err := h.logsRepository.ListByName(r.Context(), nameFilter)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(fmt.Errorf("Failed to retrieve data from DB: %v", err.Error()).Error()))
				break
			}
			jsonResponse, err := json.Marshal(logsList)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(fmt.Errorf("Failed to serialize response: %v", err.Error()).Error()))
				break
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(jsonResponse)
		case http.MethodDelete:
			clearBeforeTimestampString := r.URL.Query().Get("clear_before")
			var clearBeforeTimestamp time.Time
			var err error
			if clearBeforeTimestampString != "" {
				clearBeforeTimestamp, err = time.Parse(time.RFC3339, clearBeforeTimestampString)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					_, _ = w.Write([]byte("timestamp must be in RFC3339 format or empty string"))
					break
				}
			} else {
				clearBeforeTimestamp = time.Now()
			}
			err = h.logsRepository.ClearAllBeforeDatetime(r.Context(), clearBeforeTimestamp)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(fmt.Errorf("failed to clear before datetime: %v", err.Error()).Error()))
				break
			}
			w.WriteHeader(http.StatusNoContent)
		case http.MethodPost, http.MethodPut, http.MethodPatch:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}
