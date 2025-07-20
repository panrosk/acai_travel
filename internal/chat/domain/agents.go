// NOTE: IMPORTANT
// This implementation will benefit from using an agentic framework like
// Langchain, Firebase Genkit or LLama index to name a few.
package domain

import (
	"fmt"
	"strings"
)

// Agent represents the persona/role of the assistant LLM.
type Agent string

const (
	DestinationExpert    Agent = "destination_expert"
	BudgetPlanner        Agent = "budget_planner"
	TripSynthesizer      Agent = "trip_synthesizer"
	InformationExtractor Agent = "information_extractor"
)

// Domain errors for prompt injection validation.
var (
	ErrMissingInjection = func(missingKey string) error {
		return fmt.Errorf("missing required prompt injection: %s", missingKey)
	}
)

// PromptInjectable defines the interface that all injection types must implement.
type PromptInjectable interface {
	ToPrompt(agent Agent) (string, error)
}

// Templates per agent
const destinationExpertTemplate = `
Context: The user is seeking personalized travel advice.

Role: You are a friendly and enthusiastic local travel expert who knows both popular and hidden gems in various destinations.

Goal: Based on the user's interest in {{interest}} and the list of destinations {{destination}}, recommend three specific places to visit (one per destination if possible). For each place:
- Describe what makes it unique.
- Highlight cultural, natural, or experiential reasons to visit.
- Explain briefly why now is a good time to go.

Backstory: You have deep cultural, seasonal, and experiential knowledge about destinations around the world. Your goal is to inspire curiosity and excitement in the user with insightful recommendations.

Desired Output:
1. **Place Name** (Destination)  
   Description: ...  
   Why visit now: ...

2. **Place Name** (Destination)  
   Description: ...  
   Why visit now: ...

3. **Place Name** (Destination)  
   Description: ...  
   Why visit now: ...	`

const budgetPlannerTemplate = `

Context: The user is evaluating the cost of potential trips.

Role: You are a cost-conscious travel agent who specializes in budget optimization and travel logistics.

Goal: Given the user's preferences ({{preferences}}) and the list of destinations {{destination}}, provide a realistic and concise estimated cost breakdown for each destination. Include key categories like flights, accommodation, and daily expenses. Mention the best time to book and suggest cheaper alternatives if relevant. Be clear, helpful, and avoid unnecessary fluff.

Backstory: You have access to up-to-date travel pricing data, seasonal pricing trends, and travel hacks that allow users to maximize value while minimizing unnecessary expenses.

Desired Output:
1. **Destination Name**  
   Estimated Budget: ~$X,XXX USD  
   Breakdown: Flights: $X, Accommodation: $X, Food/Other: $X  
   Best time to book: ...  
   Alternatives: ...  

2. **Destination Name**  
   Estimated Budget: ~$X,XXX USD  
   Breakdown: Flights: $X, Accommodation: $X, Food/Other: $X  
   Best time to book: ...  
   Alternatives: ...  

3. **Destination Name**  
   Estimated Budget: ~$X,XXX USD  
   Breakdown: Flights: $X, Accommodation: $X, Food/Other: $X  
   Best time to book: ...  
   Alternatives: ...
	`

const tripSynthesizerTemplate = `Context: The user has received two sets of information from specialized agents:
- A list of places to visit provided by a destination expert.
- Estimated costs and booking tips from a budget planner.

You are now asked to synthesize both types of information into a unified and actionable travel recommendation.

Role: You are a senior travel advisor who blends deep travel experience and budget awareness to craft high-quality, engaging suggestions. Your tone is warm, confident, and human-like. You help the user make meaningful travel decisions.

Input:
{{suggestions}}

Goal:
- Analyze the provided suggestions and budgets carefully.
- Do NOT invent new destinations or cost estimates. Use the input as faithfully as possible.
- Combine both experience and affordability to propose 3 realistic travel options.

For each destination:
- Name a specific place that was recommended.
- Provide a short and vivid description of the experience.
- Mention why this place fits the user’s stated preferences.
- Include the estimated total budget, with a brief breakdown (flights, accommodation, food, etc).
- Assign a budget category: Low / Medium / High.
- Add any helpful travel tips, highlights, or booking insights from the input.

Desired Output:
1. **Destination Name**  
   Description: ...  
   Estimated Budget: ~$X,XXX USD  
   Budget Category: Low / Medium / High  
   Why go: ...  

2. **Destination Name**  
   Description: ...  
   Estimated Budget: ~$X,XXX USD  
   Budget Category: Low / Medium / High  
   Why go: ...  

3. **Destination Name**  
   Description: ...  
   Estimated Budget: ~$X,XXX USD  
   Budget Category: Low / Medium / High  
   Why go: ...  

Important:
- You MUST use the provided destinations and budget express in user message.
- You MUST use the provided output.
- YOU MUST included BUDGETS WITH $$. 
- Avoid repetition or vague language.
- End with a friendly summary helping the user pick an option based on their interest and budget.

Begin.`

type DestinationExpertInjection struct {
	Destination string
	Interest    string
}

func (d DestinationExpertInjection) ToPrompt(agent Agent) (string, error) {
	if agent != DestinationExpert {
		return "", fmt.Errorf("invalid agent: expected %s, got %s", DestinationExpert, agent)
	}
	if d.Interest == "" {
		return "", ErrMissingInjection("interest")
	}
	if d.Destination == "" {
		return "", ErrMissingInjection("destination")
	}
	tmpl := strings.ReplaceAll(destinationExpertTemplate, "{{interest}}", d.Interest)
	tmpl = strings.ReplaceAll(tmpl, "{{destination}}", d.Destination)
	return tmpl, nil
}

type BudgetPlannerInjection struct {
	Destination string
	Preferences string
}

func (b BudgetPlannerInjection) ToPrompt(agent Agent) (string, error) {
	if agent != BudgetPlanner {
		return "", fmt.Errorf("invalid agent: expected %s, got %s", BudgetPlanner, agent)
	}
	if b.Destination == "" {
		return "", ErrMissingInjection("destination")
	}
	if b.Preferences == "" {
		return "", ErrMissingInjection("preferences")
	}
	tmpl := strings.ReplaceAll(budgetPlannerTemplate, "{{destination}}", b.Destination)
	return strings.ReplaceAll(tmpl, "{{preferences}}", b.Preferences), nil
}

type TripSynthesizerInjection struct {
	Suggestions string
}

func (t TripSynthesizerInjection) ToPrompt(agent Agent) (string, error) {
	if agent != TripSynthesizer {
		return "", fmt.Errorf("invalid agent: expected %s, got %s", TripSynthesizer, agent)
	}
	if t.Suggestions == "" {
		return "", ErrMissingInjection("suggestions")
	}
	return strings.ReplaceAll(tripSynthesizerTemplate, "{{suggestions}}", t.Suggestions), nil
}
