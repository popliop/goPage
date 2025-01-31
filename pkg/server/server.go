package server

import (
	"fmt"
	"net/http"
	"os"
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
		fmt.Println("Error server crashed:", err)
	}
}

func (s *APIServer) registerRoutes() {
	s.router.HandleFunc("/", s.serveHTML)
	//s.router.Handle("/static/", http.FileServer(http.Dir("./static")))
	s.router.HandleFunc("/api/gpt", s.serveGPT)
	s.router.HandleFunc("/test", s.test)
	fmt.Println("Routes OK")
}

func (s *APIServer) serveHTML(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./view/index.html")
	fmt.Println("servering index")

	workingDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting working directory:", err)
		return
	}
	fmt.Println("Current working directory:", workingDir)

	// Print the requested URL path for additional context
	fmt.Println("Request URL Path:", r.URL.Path)

}

func (s *APIServer) test(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("test"))
}

func (s *APIServer) serveGPT(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Extract the product name from the form data
	product := r.FormValue("product")
	if product == "" {
		http.Error(w, "Product field is required", http.StatusBadRequest)
		return
	}

	test := sendtoGPT(product)

	// Generate a test response (you can replace this with GPT logic later)
	responseText := fmt.Sprintf("HS Code: %s", test)

	// Send the response back as plain text
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(responseText))
}
