package chathttpadapter

import (
	"acai_travel/internal/chat/application"
	"bufio"
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type TravelHandler struct {
	orchestrator *application.MultiAgentOrchestrator
}

func NewTravelHandler(orchestrator *application.MultiAgentOrchestrator) *TravelHandler {
	return &TravelHandler{orchestrator: orchestrator}
}

func (h *TravelHandler) RegisterRoutes(app *fiber.App) {
	travelGroup := app.Group("/travel")
	travelGroup.Post("/recommendation", h.multiAgentRecomendation)
}

func (h *TravelHandler) multiAgentRecomendation(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")

	req, err := parseRequest(c)
	if err != nil {
		return err
	}

	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
		defer cancel()

		writeEvent := func(eventType, data string) error {
			if _, err := fmt.Fprintf(w, "event: %s\n", eventType); err != nil {
				return err
			}
			if _, err := fmt.Fprintf(w, "data: %s\n\n", data); err != nil {
				return err
			}
			return w.Flush()
		}

		convoID, err := uuid.Parse(req.ConversationID)
		if err != nil {
			_ = writeEvent("error", "Invalid conversation ID")
			return
		}

		userID, err := uuid.Parse(req.UserID)
		if err != nil {
			_ = writeEvent("error", "Invalid user ID")
			return
		}

		orchInput := application.OrchestratorInput{
			ConversationID: convoID,
			UserID:         userID,
			Role:           req.Message.Role,
			Content:        req.Message.Content,
		}

		err = h.orchestrator.Run(ctx, orchInput, writeEvent)
		if err != nil {
			_ = writeEvent("error", fmt.Sprintf("Error: %v", err))
		}
	})

	return nil
}

func parseRequest(c *fiber.Ctx) (ChatRequestDTO, error) {
	var req ChatRequestDTO
	if err := c.BodyParser(&req); err != nil {
		return req, FormatErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}

	if validationErrors := req.Validate(); len(validationErrors) > 0 {
		return req, FormatErrorResponse(c, fiber.StatusBadRequest, "Validation failed", validationErrors)
	}

	return req, nil
}
