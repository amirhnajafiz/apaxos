package handler

import "github.com/gofiber/fiber/v2"

// Handler is our HTTP handler.
// It has all endpoints that are being used
// as an API.
type Handler struct{}

// Register method is used to create HTTP endpoints
// for a given fiber group router.
func (h Handler) Register(app fiber.Router) {
	// users
	app.Get("/user/balance", h.getUserBalance)

	// logs
	app.Get("/logs", h.getLogs)
	app.Get("/logs/history", h.getLogsHistory)

	// transactions
	app.Put("/transaction", h.newTransaction)
	app.Post("/transaction", h.newTransactionFile)
}
