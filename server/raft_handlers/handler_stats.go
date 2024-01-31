package rafthandlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func (h handler) StatsRaftHandler(ctx *fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Here is the raft status",
		"data":    h.raft.Stats(),
	})
}
