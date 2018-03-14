package secrets

import (
	"encoding/json"
	"net/http"

	"github.com/gojekfarm/proctor-engine/logger"
	"github.com/gojekfarm/proctor-engine/utility"
)

type handler struct {
	secretsStore Store
}

type Handler interface {
	HandleSubmission() http.HandlerFunc
}

func NewHandler(secretsStore Store) Handler {
	return &handler{
		secretsStore: secretsStore,
	}
}

func (handler *handler) HandleSubmission() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var secret Secret
		err := json.NewDecoder(req.Body).Decode(&secret)
		defer req.Body.Close()
		if err != nil {
			logger.Error("Error parsing request body", err.Error())

			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(utility.ClientError))
			return
		}

		err = handler.secretsStore.CreateOrUpdateJobSecret(secret)
		if err != nil {
			logger.Error("Error updating secrets", err.Error())

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(utility.ServerError))
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}
