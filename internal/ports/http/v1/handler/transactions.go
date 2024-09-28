package handler

import "github.com/gofiber/fiber/v2"

// new transaction handler is used for creating a single
// transaction. often is used for testing systems performance.
func (h Handler) newTransaction(c *fiber.Ctx) error {
	return nil
}

// new transaction file handler accepts a .CSV file and
// creates transactions into the system.
func (h Handler) newTransactionFile(c *fiber.Ctx) error {
	return nil
}
