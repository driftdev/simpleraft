package rafthandlers

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/hashicorp/raft"
)

// removeRequest request payload for removing node from raft cluster
type removeRequest struct {
	NodeID string `json:"node_id"`
}

// RemoveRaftHandler handling removing raft
func (h handler) RemoveRaftHandler(ctx *fiber.Ctx) error {
	var form = removeRequest{}
	if err := ctx.BodyParser(&form); err != nil {
		return ctx.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": fmt.Sprintf("error binding: %s", err.Error()),
		})
	}

	var nodeID = form.NodeID

	if h.raft.State() != raft.Leader {
		return ctx.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": "not the leader",
		})
	}

	configFuture := h.raft.GetConfiguration()
	if err := configFuture.Error(); err != nil {
		return ctx.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": fmt.Sprintf("failed to get raft configuration: %s", err.Error()),
		})
	}

	future := h.raft.RemoveServer(raft.ServerID(nodeID), 0, 0)
	if err := future.Error(); err != nil {
		return ctx.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": fmt.Sprintf("error removing existing node %s: %s", nodeID, err.Error()),
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": fmt.Sprintf("node %s removed successfully", nodeID),
		"data":    h.raft.Stats(),
	})
}
