package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type APIServer struct {
	serverPort string
	router     http.ServeMux
}

func NewAPIServer(listenAddr string) *APIServer {

	return &APIServer{
		serverPort: listenAddr,
		router:     *http.NewServeMux(),
	}
}

func (s *APIServer) Run() {

	s.registerRoutes()

	fmt.Println("Server is running on port", s.serverPort)
	if err := http.ListenAndServe(s.serverPort, &s.router); err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func (s *APIServer) registerRoutes() {
	s.router.HandleFunc("/2", s.serveHTML)
	s.router.HandleFunc("/", s.serveHTML)
	s.router.HandleFunc("/translate", s.translate)
	s.router.HandleFunc("/test/{id}", s.test)

	s.router.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

}

func (s *APIServer) serveHTML(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./view/index.html")
}

func (s *APIServer) translate(w http.ResponseWriter, r *http.Request) {

	writeJSON(w, 200, "hello World")
}

func (s *APIServer) test(w http.ResponseWriter, r *http.Request) {
	data := r.PathValue("id")
	fmt.Printf("Value from URL: %s\n", data)
	writeJSON(w, 200, data)
}

// Helper
func writeJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	// Set the Content-Type header
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	// Encode the data to JSON and send it
	if err := json.NewEncoder(w).Encode(data); err != nil {
		// Handle encoding errors
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}
