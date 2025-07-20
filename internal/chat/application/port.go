package application

import (
	"acai_travel/internal/chat/domain"
	"context"
)

type ChatServiceInterface interface {
	GetDestinationAdvice(ctx context.Context, chat *domain.Chat, injections domain.PromptInjectable, model domain.LLMModel) (*domain.Chat, error)
	PlanBudget(ctx context.Context, chat *domain.Chat, injections domain.PromptInjectable, model domain.LLMModel) (*domain.Chat, error)
	StreamTripSummary(ctx context.Context, chat *domain.Chat, injections domain.PromptInjectable, model domain.LLMModel, streamFn func(eventType, data string) error) error
	InformationExtraction(ctx context.Context, chat *domain.Chat, schema map[string]any, model domain.LLMModel) (map[string]string, error)
}

type DestinationExpertUseCase interface {
	Run(ctx context.Context, chat *domain.Chat, injections domain.PromptInjectable, model domain.LLMModel) (*domain.Chat, error)
}

type BudgetPlannerUseCase interface {
	Run(ctx context.Context, chat *domain.Chat, injections domain.PromptInjectable, model domain.LLMModel) (*domain.Chat, error)
}

type InformationExtractorUsecase interface {
	Run(ctx context.Context, chat *domain.Chat, schema map[string]any, model domain.LLMModel) (map[string]string, error)
}

type TripSynthesizerUseCase interface {
	Stream(ctx context.Context, chat *domain.Chat, injections domain.PromptInjectable, model domain.LLMModel, streamFn func(eventType, data string) error) error
}

type ChatService struct {
	destExpert      DestinationExpertUseCase
	budgetPlanner   BudgetPlannerUseCase
	tripSynthesizer TripSynthesizerUseCase
	infoExtractor   InformationExtractorUsecase
}

func NewChatService(
	dest DestinationExpertUseCase,
	budget BudgetPlannerUseCase,
	synth TripSynthesizerUseCase,
	info InformationExtractorUsecase,
) *ChatService {
	return &ChatService{
		destExpert:      dest,
		budgetPlanner:   budget,
		tripSynthesizer: synth,
		infoExtractor:   info,
	}
}

func (s *ChatService) GetDestinationAdvice(ctx context.Context, chat *domain.Chat, injections domain.PromptInjectable, model domain.LLMModel) (*domain.Chat, error) {
	return s.destExpert.Run(ctx, chat, injections, model)
}

func (s *ChatService) PlanBudget(ctx context.Context, chat *domain.Chat, injections domain.PromptInjectable, model domain.LLMModel) (*domain.Chat, error) {
	return s.budgetPlanner.Run(ctx, chat, injections, model)
}

func (s *ChatService) InformationExtraction(ctx context.Context, chat *domain.Chat, schema map[string]any, model domain.LLMModel) (map[string]string, error) {
	return s.infoExtractor.Run(ctx, chat, schema, model)
}

func (s *ChatService) StreamTripSummary(ctx context.Context, chat *domain.Chat, injections domain.PromptInjectable, model domain.LLMModel, streamFn func(eventType, data string) error) error {
	return s.tripSynthesizer.Stream(ctx, chat, injections, model, streamFn)
}
