// Command server serves archived hackathon data only. Use course-server for new work.
package main

import (
	"log"
	"os"

	"agent-arena/backend/internal/api"
)

func main() {
	if os.Getenv("ENABLE_LEGACY_SERVER") != "true" {
		log.Fatal("The hackathon server is retired. Run go run ./cmd/course-server. For historical read-only access set ENABLE_LEGACY_SERVER=true and LEGACY_DATABASE_PATH explicitly.")
	}
	if os.Getenv("LEGACY_DATABASE_PATH") == "" {
		log.Fatal("LEGACY_DATABASE_PATH is required for historical read-only access")
	}
	server, err := api.NewServer()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Historical read-only API listening on 127.0.0.1:8081; no signer or reward loop is started")
	log.Fatal(server.Run("127.0.0.1:8081"))
}
