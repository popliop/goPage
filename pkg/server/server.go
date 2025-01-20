package server

import (
	"fmt"
	"net/http"
)

type APIServer struct {
	serverPort string
}

func NewAPIServer(listenAddr string) *APIServer {
	return &APIServer{
		serverPort: listenAddr,
	}
}

func (s *APIServer) Run() {

	s.registerRoutes()

	fmt.Println("Server is running on port", s.serverPort)
	if err := http.ListenAndServe(s.serverPort, nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func (s *APIServer) registerRoutes() {
	http.HandleFunc("/", s.serveHTML)
}

func (s *APIServer) serveHTML(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./view/index.html")
}
