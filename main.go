package main

import (
	"embed"
	"io/fs"
	"log"
	"os"
	"strconv"

	"github.com/casell/signing-service-challenge/api"
)

const (
	ListenAddress = ":8080"
	CorsEnvName   = "CORS_ENABLED"
	CorsDefault   = false
)

//go:embed openapi/openapi.yaml
var spec embed.FS

func getCORSFromEnv() (bool, error) {
	cors, set := os.LookupEnv(CorsEnvName)
	if !set {
		return CorsDefault, nil
	}
	return strconv.ParseBool(cors)
}

func main() {
	specFS, err := fs.Sub(spec, "openapi")
	if err != nil {
		log.Fatal("Unable to Sub on embedded openapi FS", err)
	}

	cors, err := getCORSFromEnv()
	if err != nil {
		log.Fatalf("Unable to parse %s variable: %v", CorsEnvName, err)
	}

	server := api.NewServer(ListenAddress, specFS, cors)

	if err := server.Run(); err != nil {
		log.Fatal("Could not start server on ", ListenAddress)
	}
}
