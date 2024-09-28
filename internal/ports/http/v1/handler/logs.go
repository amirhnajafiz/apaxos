package handler

import "github.com/gofiber/fiber/v2"

// get logs returns the list of current node logs.
func (h Handler) getLogs(c *fiber.Ctx) error {
	return nil
}

// get logs history pulls all previous commited logs
// that are sotred in database.
func (h Handler) getLogsHistory(c *fiber.Ctx) error {
	return nil
}
