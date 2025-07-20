# 🌍 Multi-Agent Travel Planner

This project is a streaming backend API built in Go using the Fiber framework. It powers an AI-assisted travel planning experience by simulating multiple AI agents (LLMs) collaborating in real time.

---

## 🧠 Architecture Overview

The system follows a **structured multi-agent reasoning flow**:

1. **Structured Output Extraction**:  
   A first LLM extracts key information (destinations, preferences, interests) from the user's message.
2. **Parallel Agents**:
   - **Destination Expert**: recommends destinations based on user interests.
   - **Budget Planner**: estimates the cost of the trip given preferences and destination.
3. **Trip Synthesizer**:  
   A final model synthesizes previous agent outputs into a unified travel recommendation.

All agent interactions stream responses incrementally using **Server-Sent Events (SSE)**.

> ✅ Built using a **hybrid architecture** combining **Domain-Driven Design (DDD)** with **Ports and Adapters (Hexagonal)** to enforce clear boundaries between domain logic, use cases, adapters, and HTTP handlers.

---

## ✨ Key Features

- 🧠 **Structured Multi-Agent Reasoning** using LLMs.
- ⚡ **Parallel Agent Execution** for faster response times.
- 📡 **Streaming with SSE** for real-time feedback.
- 🔌 **Pluggable LLM Provider Layer** (currently OpenAI).
- 🧼 **Clean Hexagonal Structure** with DDD principles.

---

## ✈️ API Endpoints

### `POST /travel/multiagent`

Launches a full multi-agent reasoning session: extraction → parallel agents → trip synthesis.

#### Example Payload

```json
{
  "conversationId": "c8f8b94e-f2c4-4d1e-8e1d-e6f7a5b7c2a2",
  "userId": "1d5cbf80-9f49-44fd-a0d0-1f7bba36a2fa",
  "message": {
    "role": "user",
    "content": "I want to go on vacations to either Panama, Costa Rica or Guatemala."
  }
}
```

#### Curl to stream full response (raw):

```bash
curl -N -X POST http://localhost:8080/travel/multiagent \
  -H "Content-Type: application/json" \
  -d '{
    "conversationId": "c8f8b94e-f2c4-4d1e-8e1d-e6f7a5b7c2a2",
    "userId": "1d5cbf80-9f49-44fd-a0d0-1f7bba36a2fa",
    "message": {
      "role": "user",
      "content": "I want to go on vacations to either Panama, Costa Rica or Guatemala."
    }
  }'
```

---

### `POST /travel/recommendation`

Executes **only the structured information extraction** phase.

#### Example:

```bash
curl -N -X POST http://localhost:8080/travel/recommendation \
  -H "Content-Type: application/json" \
  -d '{
    "conversationId": "c8f8b94e-f2c4-4d1e-8e1d-e6f7a5b7c2a2",
    "userId": "1d5cbf80-9f49-44fd-a0d0-1f7bba36a2fa",
    "message": {
      "role": "user",
      "content": "I want to go on vacations to either Panama, Costa Rica or Guatemala."
    }
  }' | awk -F'data: ' '/^data:/ { printf "%s", $2 }'
```

> 🛠 **Note**: The `awk` filter strips the SSE `data:` prefix to show raw content. Remove it to view the full event stream, including `status` messages.

---

## 🛠 Getting Started

NOTE: DONT FORGET TO CONFIG THE ENV.

### Build and Run

```bash
make build     # Compile the project
make run       # Run the server
make watch     # Run with live reload (development)
make clean     # Remove build artifacts
```

### Test Suite

```bash
make test      # Run all tests
make all       # Build and test
```

---

## 📝 Notes for the Acai Travel Team

This exercise demonstrates how structured LLM outputs can orchestrate multi-agent reasoning using clean Go architecture.

Although only three stages were required, we intentionally separated the **information extractor** into its own usecase to help debugging and future extensibility.

A **minimal test suite** is included to validate integration with LLMs and ensure correct behavior of agent workflows. However, full test coverage was deemed **out of scope** for this prototype and is not exhaustive.
