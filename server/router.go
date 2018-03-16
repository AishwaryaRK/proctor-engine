package server

import (
	"fmt"
	"net/http"

	"github.com/gojekfarm/proctor-engine/audit"
	"github.com/gojekfarm/proctor-engine/jobs/execution"
	"github.com/gojekfarm/proctor-engine/jobs/logs"
	"github.com/gojekfarm/proctor-engine/jobs/metadata"
	"github.com/gojekfarm/proctor-engine/jobs/secrets"
	"github.com/gojekfarm/proctor-engine/kubernetes"
	"github.com/gojekfarm/proctor-engine/redis"
	"github.com/gojekfarm/proctor-engine/storage"
	"github.com/gojekfarm/proctor-engine/storage/postgres"

	"github.com/gorilla/mux"
)

var postgresClient postgres.Client

func NewRouter() *mux.Router {
	router := mux.NewRouter()

	redisClient := redis.NewClient()
	postgresClient = postgres.NewClient()

	store := storage.New(postgresClient)
	metadataStore := metadata.NewStore(redisClient)
	secretsStore := secrets.NewStore(redisClient)

	kubeConfig := kubernetes.KubeConfig()
	kubeClient := kubernetes.NewClient(kubeConfig)

	auditor := audit.New(store)
	jobExecutioner := execution.NewExecutioner(kubeClient, metadataStore, secretsStore, auditor)
	jobLogger := logs.NewLogger(kubeClient)
	jobMetadataHandler := metadata.NewHandler(metadataStore)
	jobSecretsHandler := secrets.NewHandler(secretsStore)

	router.HandleFunc("/ping", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "pong")
	})

	router.HandleFunc("/jobs/execute", jobExecutioner.Handle()).Methods("POST")
	router.HandleFunc("/jobs/logs", jobLogger.Stream()).Methods("GET")
	router.HandleFunc("/jobs/metadata", jobMetadataHandler.HandleSubmission()).Methods("POST")
	router.HandleFunc("/jobs/metadata", jobMetadataHandler.HandleBulkDisplay()).Methods("GET")
	router.HandleFunc("/jobs/secrets", jobSecretsHandler.HandleSubmission()).Methods("POST")

	return router
}
