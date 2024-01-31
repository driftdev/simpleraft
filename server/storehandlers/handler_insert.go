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

// insertRequest payload for storing new data in raft cluster
type insertRequest struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

// Insert handling save to raft cluster. Insert will invoke raft.Apply to make this stored in all cluster
// with acknowledge from n quorum. Insert must be done in raft leader, otherwise return error.
func (h handler) Insert(ctx *fiber.Ctx) error {
	var form = insertRequest{}
	if err := ctx.BodyParser(&form); err != nil {
		return ctx.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": fmt.Sprintf("error binding: %s", err.Error()),
		})
	}

	form.Key = strings.TrimSpace(form.Key)
	if form.Key == "" {
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
		Operation: "SET",
		Key:       form.Key,
		Value:     form.Value,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return ctx.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": fmt.Sprintf("error preparing saving data payload: %s", err.Error()),
		})
	}

	applyFuture := h.raft.Apply(data, 500*time.Millisecond)
	if err := applyFuture.Error(); err != nil {
		return ctx.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": fmt.Sprintf("error persisting data in raft cluster: %s", err.Error()),
		})
	}

	_, ok := applyFuture.Response().(*fsm.ApplyResponse)
	if !ok {
		return ctx.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": fmt.Sprintf("error response is not match apply response"),
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "success persisting data",
		"data":    form,
	})
}
