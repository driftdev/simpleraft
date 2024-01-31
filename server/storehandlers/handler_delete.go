package storehandlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ArkamFahry/simpleraft/fsm"
	"github.com/gofiber/fiber/v2"
	"github.com/hashicorp/raft"
)

func (h handler) Delete(ctx *fiber.Ctx) error {
	var key = strings.TrimSpace(ctx.Query("key"))
	if key == "" {
		return ctx.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": "key is empty",
		})
	}

	if h.raft.State() != raft.Leader {
		return ctx.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": "not the leader",
		})
	}

	payload := fsm.CommandPayload{
		Operation: "DELETE",
		Key:       key,
		Value:     nil,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return ctx.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": fmt.Sprintf("error preparing remove data payload: %s", err.Error()),
		})
	}

	applyFuture := h.raft.Apply(data, 500*time.Millisecond)
	if err := applyFuture.Error(); err != nil {
		return ctx.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": fmt.Sprintf("error removing data in raft cluster: %s", err.Error()),
		})
	}

	_, ok := applyFuture.Response().(*fsm.ApplyResponse)
	if !ok {
		return ctx.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": fmt.Sprintf("error response is not match apply response"),
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "success removing data",
		"data": map[string]interface{}{
			"key":   key,
			"value": nil,
		},
	})
}
