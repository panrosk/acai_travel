package server

import (
	chathttpadapter "acai_travel/internal/chat/adapters/chat_http_adapter"
	"acai_travel/internal/chat/adapters/llm"
	"acai_travel/internal/chat/application"
	"bufio"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func (s *FiberServer) RegisterFiberRoutes() {
	s.App.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Accept,Authorization,Content-Type",
		AllowCredentials: false,
		MaxAge:           300,
	}))

	s.App.Get("/", s.HelloWorldHandler)
	openaiApiKey := os.Getenv("OPENAI_API_KEY")
	openaiClient := llm.NewOpenAIClient(openaiApiKey)

	infoExtractor := application.NewInformationExtractor(openaiClient)
	destExper := application.NewDestinationExpert(openaiClient)
	budgetPlanner := application.NewBudgetPlanner(openaiClient)
	tripSynth := application.NewTripSynthesizer(openaiClient)

	chat_service := application.NewChatService(destExper, budgetPlanner, tripSynth, infoExtractor)

	orchestrator := application.NewMultiAgentOrchestrator(chat_service)
	handler := chathttpadapter.NewTravelHandler(orchestrator)
	handler.RegisterRoutes(s.App)

	s.App.Get("/events", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/event-stream")
		c.Set("Cache-Control", "no-cache")
		c.Set("Connection", "keep-alive")

		c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
			writeEvent := func(event, data string) {
				w.WriteString(fmt.Sprintf("event: %s\n", event))
				w.WriteString(fmt.Sprintf("data: %s\n\n", data))
				w.Flush()
			}

			writeEvent("status", "loading")

			var wg sync.WaitGroup
			var mu sync.Mutex

			var result1, result2 string

			wg.Add(2)

			go func() {
				defer wg.Done()
				res := func1()
				mu.Lock()
				result1 = res
				mu.Unlock()
				writeEvent("status", "func1 done")
			}()

			go func() {
				defer wg.Done()
				res := func2()
				mu.Lock()
				result2 = res
				mu.Unlock()
				writeEvent("status", "func2 done")
			}()

			wg.Wait()

			final := func3(result1, result2)
			writeEvent("status", "completed")
			writeEvent("message", final)
		})

		return nil
	})

}

func func1() string {
	time.Sleep(10 * time.Second)
	return "resultado de func1"
}

func func2() string {
	time.Sleep(3 * time.Second)
	return "resultado de func2"
}

func func3(res1, res2 string) string {
	return fmt.Sprintf("func3 recibió:\n- %s\n- %s", res1, res2)
}

func (s *FiberServer) HelloWorldHandler(c *fiber.Ctx) error {
	resp := fiber.Map{
		"message": "Hello World",
	}
	return c.JSON(resp)
}
