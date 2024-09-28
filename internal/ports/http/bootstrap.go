package http

import (
	"fmt"
	"log"

	"github.com/f24-cse535/apaxos/internal/config/http"
	"github.com/f24-cse535/apaxos/internal/ports/http/v1/handler"

	"github.com/gofiber/fiber/v2"
)

func Bootstrap(cfg http.Config) {
	// create a new fiber app
	app := fiber.New()

	// init api handlers
	v1 := handler.Handler{}

	// register app endpoints
	v1.Register(app.Group("/v1"))

	// start the HTTP server
	if err := app.Listen(fmt.Sprintf(":%d", cfg.Port)); err != nil {
		log.Fatalf("failed to start node HTTP server on port %d: %v\n", cfg.Port, err)
	}
}
