package api

import (
	"encoding/json"
	"io"
	"messange_handler/src/entities"
	"messange_handler/src/errors/api_errors"
	"messange_handler/src/services"
	"messange_handler/src/utils"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

func InitMessageRoutes(router *mux.Router, pool *pgxpool.Pool, log *logrus.Logger) {
	router.HandleFunc("/message", createMessage(pool, log)).Methods("POST")
	router.HandleFunc("/message/stats", getMessageStats(pool, log)).Methods("GET")
}

func createMessage(pool *pgxpool.Pool, log *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		data, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			resp := entities.Response{
				Error: api_errors.InternalServerError{}.Error(),
			}
			json.NewEncoder(w).Encode(resp)
			return
		}

		if data == nil {
			w.WriteHeader(http.StatusBadRequest)
			resp := entities.Response{
				Error: api_errors.BadRequestError{Detail: "Content is required"}.Error(),
			}
			json.NewEncoder(w).Encode(resp)
			return
		}

		if !utils.IsValidJSON(data) {
			w.WriteHeader(http.StatusBadRequest)
			resp := entities.Response{
				Error: api_errors.BadRequestError{Detail: "Content must be a valid JSON"}.Error(),
			}
			json.NewEncoder(w).Encode(resp)
			return
		}

		message_id, err := services.CreateMessage(
			r.Context(),
			pool,
			log,
			&entities.MessageFromRequest{Content: string(data)},
		)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			resp := entities.Response{
				Error: api_errors.InternalServerError{}.Error(),
			}
			json.NewEncoder(w).Encode(resp)
			return
		}

		w.WriteHeader(http.StatusCreated)
		resp := entities.Response{Data: entities.CreateMessageResponse{
			MessageId: message_id,
		}}
		json.NewEncoder(w).Encode(resp)
	}
}

func getMessageStats(pool *pgxpool.Pool, log *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		stats, err := services.GetMessageStat(
			r.Context(),
			pool,
			log,
		)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			resp := entities.Response{
				Error: api_errors.InternalServerError{}.Error(),
			}
			json.NewEncoder(w).Encode(resp)
			return
		}

		w.WriteHeader(http.StatusOK)
		resp := entities.Response{Data: stats}
		json.NewEncoder(w).Encode(resp)
	}
}
