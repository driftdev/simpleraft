package rafthandlers

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/hashicorp/raft"
)

// joinRequest request payload for joining raft cluster
type joinRequest struct {
	NodeID      string `json:"node_id"`
	RaftAddress string `json:"raft_address"`
}

func (h handler) JoinRaftHandler(ctx *fiber.Ctx) error {
	var form = joinRequest{}
	if err := ctx.BodyParser(&form); err != nil {
		return ctx.Status(http.StatusUnprocessableEntity).JSON(map[string]interface{}{
			"error": fmt.Sprintf("error binding: %s", err.Error()),
		})
	}

	var (
		nodeID   = form.NodeID
		raftAddr = form.RaftAddress
	)

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

	// This must be run on the leader or it will fail.
	f := h.raft.AddVoter(raft.ServerID(nodeID), raft.ServerAddress(raftAddr), 0, 0)
	if f.Error() != nil {
		return ctx.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": fmt.Sprintf("error adding voter: %s", f.Error().Error()),
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": fmt.Sprintf("node %s at %s joined successfully", nodeID, raftAddr),
		"data":    h.raft.Stats(),
	})
}
