package application

import (
	"acai_travel/internal/chat/domain"
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type MultiAgentOrchestrator struct {
	service ChatServiceInterface
}

func NewMultiAgentOrchestrator(service ChatServiceInterface) *MultiAgentOrchestrator {
	return &MultiAgentOrchestrator{service: service}
}

type OrchestratorInput struct {
	ConversationID uuid.UUID
	UserID         uuid.UUID
	Role           string
	Content        string
}

type AgentResponse struct {
	Result string
	Error  error
}

func (m *MultiAgentOrchestrator) Run(
	ctx context.Context,
	input OrchestratorInput,
	streamFn func(eventType, data string) error,
) error {
	streamFn("status", "Invoking LLM 1 (extraction)")

	info, err := m.extractInformation(ctx, input)
	if err != nil {
		_ = streamFn("error", fmt.Sprintf("LLM 1 failed: %v", err))
		return fmt.Errorf("LLM 1 failed: %w", err)
	}
	streamFn("status", "Got response from LLM 1 (info extracted)")

	destinationChan := make(chan AgentResponse)
	budgetChan := make(chan AgentResponse)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		streamFn("status", "Invoking LLM 2 (destination expert)")
		destinationChan <- m.runDestinationExpert(ctx, input, info)
	}()

	go func() {
		defer wg.Done()
		streamFn("status", "Invoking LLM 3 (budget planner)")
		budgetChan <- m.runBudgetPlanner(ctx, input, info)
	}()

	var destinationRes, budgetRes AgentResponse
	wgDone := make(chan struct{})
	go func() {
		defer close(wgDone)
		wg.Wait()
	}()

	select {
	case destinationRes = <-destinationChan:
	case <-ctx.Done():
		_ = streamFn("error", "Timeout while waiting for destination expert")
		return ctx.Err()
	}

	select {
	case budgetRes = <-budgetChan:
	case <-ctx.Done():
		_ = streamFn("error", "Timeout while waiting for budget planner")
		return ctx.Err()
	}

	if destinationRes.Error != nil {
		_ = streamFn("error", fmt.Sprintf("LLM 2 failed: %v", destinationRes.Error))
	}
	if budgetRes.Error != nil {
		_ = streamFn("error", fmt.Sprintf("LLM 3 failed: %v", budgetRes.Error))
	}

	streamFn("status", "Invoking LLM 4 (trip synthesizer)")
	return m.streamFinalSummary(ctx, input, streamFn, destinationRes.Result, budgetRes.Result)
}

func (m *MultiAgentOrchestrator) extractInformation(ctx context.Context, input OrchestratorInput) (map[string]string, error) {
	chat := domain.NewChat(input.UserID)
	systemMsg := domain.NewSystemMessage(chat.ID, "Por favor, analiza esta solicitud del usuario. Por favor llena los campos; si no existe alguno, coloca 'ninguna'.")
	userMsg := domain.NewUserMessage(chat.ID, input.Content)

	_ = chat.AddMessage(systemMsg)
	_ = chat.AddMessage(userMsg)

	schema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"Destinations": map[string]string{"type": "string"},
			"Preferences":  map[string]string{"type": "string"},
			"Interest":     map[string]string{"type": "string"},
		},
		"required":             []string{"Destinations", "Preferences", "Interest"},
		"additionalProperties": false,
	}

	return m.service.InformationExtraction(ctx, chat, schema, "gpt-4o")
}

func (m *MultiAgentOrchestrator) runDestinationExpert(ctx context.Context, input OrchestratorInput, info map[string]string) AgentResponse {
	chat := domain.NewChat(input.UserID)
	chat.AddMessage(domain.NewUserMessage(chat.ID, input.Content))

	injection := domain.DestinationExpertInjection{
		Interest:    info["Interest"],
		Destination: info["Destinations"],
	}

	resp, err := m.service.GetDestinationAdvice(ctx, chat, injection, "gpt-4")
	if err != nil {
		return AgentResponse{"No destination advice available.", err}
	}
	if len(resp.Messages) == 0 {
		return AgentResponse{"No destination advice available.", fmt.Errorf("empty response")}
	}
	return AgentResponse{resp.Messages[len(resp.Messages)-1].Content, nil}
}

func (m *MultiAgentOrchestrator) runBudgetPlanner(ctx context.Context, input OrchestratorInput, info map[string]string) AgentResponse {
	chat := domain.NewChat(input.UserID)
	chat.AddMessage(domain.NewUserMessage(chat.ID, "Dadas tus instrucciones responde con mis vacaciones perferctas"))

	injection := domain.BudgetPlannerInjection{
		Preferences: info["Preferences"],
		Destination: info["Destinations"],
	}

	resp, err := m.service.PlanBudget(ctx, chat, injection, "gpt-4")
	if err != nil {
		return AgentResponse{"No budget plan available.", err}
	}
	if len(resp.Messages) == 0 {
		return AgentResponse{"No budget plan available.", fmt.Errorf("empty response")}
	}
	return AgentResponse{resp.Messages[len(resp.Messages)-1].Content, nil}
}

func (m *MultiAgentOrchestrator) streamFinalSummary(
	ctx context.Context,
	input OrchestratorInput,
	streamFn func(eventType, data string) error,
	destination, budget string,
) error {
	chat := domain.NewChat(input.UserID)
	chat.AddMessage(domain.NewUserMessage(chat.ID, budget))
	chat.AddMessage(domain.NewUserMessage(chat.ID, destination))
	chat.AddMessage(domain.NewUserMessage(chat.ID, "Given this messages pelase give me my best vacations"))

	injections := domain.TripSynthesizerInjection{
		Suggestions: "Follow very closely toy instructions, Used all information provided by the user",
	}

	err := m.service.StreamTripSummary(ctx, chat, injections, "gpt-4", streamFn)
	if err != nil {
		_ = streamFn("error", fmt.Sprintf("LLM 4 failed: %v", err))
		return fmt.Errorf("LLM 4 failed: %w", err)
	}

	_ = streamFn("status", "completed")
	return nil
}
