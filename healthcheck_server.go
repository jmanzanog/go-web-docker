package src

import (
	"log"
	"net/http"
)

type HealthCheckServer struct {
	//conn *connection.Connection
	//cfg  *connection.Config
}

func (h *HealthCheckServer) HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	HealthCheckerHandler(w, r)
}

func (h *HealthCheckServer) StartServer() {
	//server := NewServer(util.GetEnv(ServerPort, ServerPortDefaultValue))
	server := NewServer(":4500")
	server.Handle(
		HealthCheckPath,
		GetMethod,
		server.AddMiddleware(h.HandleHealthCheck, LoggingMiddleware()))
	err := server.Start()
	if err != nil {
		log.Println(UnableServerMsg, err)
		return
	}
	log.Println(ServerRunningMsg, server.Port)
}
