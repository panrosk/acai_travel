# 🌍 Multi-Agent Travel Planner

This project is a streaming backend API built in Go using the Fiber framework that powers an AI-assisted travel planning experience using multiple Large Language Models (LLMs). The system simulates multiple AI agents (destination expert, budget planner, trip synthesizer) working together to provide rich travel recommendations in real-time.

---

## 🧠 Key Features

- **Multi-Agent Architecture:** Implements a multi-step agent reasoning system.
- **Streaming Responses:** Uses Server-Sent Events (SSE) for real-time LLM outputs.
- **LLM Provider Abstraction:** Works with OpenAI (can be extended to other providers).
- **Prompt Injection System:** Custom prompts per agent with safe variable injection.
- **Clean Hexagonal Architecture:** Clear separation between handlers, use cases, domain, and adapters.

---

## ✈️ API Endpoint

### `POST /travel/multiagent`

Starts a new multi-agent travel planning session and streams responses.

#### Request Body

```json
{
  "conversationId": "uuid-string",
  "userId": "uuid-string",
  "message": {
    "role": "user",
    "content": "I want a relaxing beach vacation in indonesia, i dont like fish but i like fresh pasas"
  },
}


## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

## MakeFile

Run build make command with tests
```bash
make all
```

Build the application
```bash
make build
```

Run the application
```bash
make run
```

Live reload the application:
```bash
make watch
```

Run the test suite:
```bash
make test
```

Clean up binary from the last build:
```bash
make clean
```


To test the enpoint:

```bash
curl -N -X POST http://localhost:8080/travel/multiagent \
  -H "Content-Type: application/json" \
  -d '{
    "conversationId": "c8f8b94e-f2c4-4d1e-8e1d-e6f7a5b7c2a2",
    "userId": "1d5cbf80-9f49-44fd-a0d0-1f7bba36a2fa",
    "message": {
      "role": "user",
      "content": "Quiero unas vacaciones relajantes cerca del mar"
    },
  }'
```


# Notes for the acai travel temam. 

As this is an exercise this 
