package api

import (
	"encoding/json"
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
}

func createMessage(pool *pgxpool.Pool, log *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var message entities.CreateMessageRequest
		err := json.NewDecoder(r.Body).Decode(&message)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			resp := entities.Response{
				Error: api_errors.BadRequestError{Detail: err.Error()}.Error(),
			}
			json.NewEncoder(w).Encode(resp)
			return
		}
		if message.Content == "" {
			w.WriteHeader(http.StatusBadRequest)
			resp := entities.Response{
				Error: api_errors.BadRequestError{Detail: "Content is required"}.Error(),
			}
			json.NewEncoder(w).Encode(resp)
			return
		}

		if !utils.IsValidJSON(message.Content) {
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
			&message,
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
		resp := entities.CreateMessageResponse{
			MessageId: message_id,
		}
		json.NewEncoder(w).Encode(resp)
	}
}
