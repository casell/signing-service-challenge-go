package api

import (
	"encoding/json"
	"io/fs"
	"log"
	"net/http"

	"github.com/rs/cors"

	"github.com/casell/signing-service-challenge/domain"
	"github.com/casell/signing-service-challenge/generated/signingapi"
	"github.com/casell/signing-service-challenge/persistence"
)

// Response is the generic API response container.
type Response struct {
	Data interface{} `json:"data"`
}

// ErrorResponse is the generic error API response container.
type ErrorResponse struct {
	Errors []string `json:"errors"`
}

// Server manages HTTP requests and dispatches them to the appropriate services.
type Server struct {
	listenAddress string
	spec          fs.FS
	cors          bool
	store         persistence.Storage
	deviceFactory domain.SigningDeviceFactory
}

// NewServer is a factory to instantiate a new Server.
func NewServer(listenAddress string, spec fs.FS, cors bool) *Server {
	return &Server{
		listenAddress: listenAddress,
		spec:          spec,
		cors:          cors,
		store:         persistence.NewMemoryStore(),
		deviceFactory: domain.NewDefaultDeviceFactory(),
	}
}

// Run registers all HandlerFuncs for the existing HTTP routes and starts the Server.
func (s *Server) Run() error {
	srv, err := signingapi.NewServer(NewDeviceHandler(s.store, s.deviceFactory))
	if err != nil {
		return err
	}
	mux := http.NewServeMux()

	mux.Handle("/api/v0/health", http.HandlerFunc(s.Health))

	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", srv))

	mux.Handle("/api/v1/openapi.yaml", http.StripPrefix("/api/v1", http.FileServer(http.FS(s.spec))))

	var h http.Handler

	if s.cors {
		h = cors.Default().Handler(mux)
	} else {
		h = mux
	}

	log.Printf("Server listening at %s, CORS enabled: %v...\n", s.listenAddress, s.cors)

	return http.ListenAndServe(s.listenAddress, h)
}

// WriteInternalError writes a default internal error message as an HTTP response.
func WriteInternalError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
}

// WriteErrorResponse takes an HTTP status code and a slice of errors
// and writes those as an HTTP error response in a structured format.
func WriteErrorResponse(w http.ResponseWriter, code int, errors []string) {
	w.WriteHeader(code)

	errorResponse := ErrorResponse{
		Errors: errors,
	}

	bytes, err := json.Marshal(errorResponse)
	if err != nil {
		WriteInternalError(w)
	}

	w.Write(bytes)
}

// WriteAPIResponse takes an HTTP status code and a generic data struct
// and writes those as an HTTP response in a structured format.
func WriteAPIResponse(w http.ResponseWriter, code int, data interface{}) {
	w.WriteHeader(code)

	response := Response{
		Data: data,
	}

	bytes, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		WriteInternalError(w)
	}

	w.Write(bytes)
}
