# AI Chat Integration Subtasks

**Problem:** The chat endpoint (`POST /chat`) currently acts as a dumb message store. It saves whatever the user types and returns it unchanged. No AI is involved. Users cannot ask questions or get AI-generated meal plans through the chat interface.

**Goal:** Make `POST /chat` an AI-powered conversational endpoint. When a user sends a plain message, Gemini should reply conversationally (with context about their profile, preferences, and pantry). When a user sends an `action` field (`meal`/`day`/`week`/`month`), it should generate a structured meal plan and return a proposal preview.

---

## Task 1: Add `InsertBotMessage` to Chat Repository

**File:** `backend/internal/repositories/chat_repository.go`

Add a method to insert a bot (AI) response message, mirroring `InsertUserMessage` but with role `"bot"`:

```go
func (r *ChatRepository) InsertBotMessage(ctx context.Context, userID, content, action string) (StoredChatMessage, error)
```

- Uses `role = 'bot'`
- Stores the AI-generated content and optional action
- Returns the `StoredChatMessage` with its generated ID and timestamp
- The INSERT should use `NULLIF(?, '')` for action (same pattern as `InsertUserMessage`)

---

## Task 2: Create Chat-Specific Prompt Template

**File:** `backend/internal/modules/ai/prompts.go`

Add a new prompt template constant for general chat (non-plan messages):

```go
PromptTemplateChat PromptTemplate = "chat"
```

In the `BuildPrompt` switch, add a case for `PromptTemplateChat`:

```go
case PromptTemplateChat:
    objective = "Provide helpful, concise advice about meal planning, nutrition, recipes, and the user's dietary preferences."
```

Also add a guardrail specific to chat:

```go
"Keep responses friendly, concise, and actionable. Do not output JSON unless the user explicitly asks for structured data."
```

The chat prompt should still include the full user context (body metrics, preferences, budget, pantry snapshot) so the AI can give personalized advice.

---

## Task 3: Update `ChatService.SendMessage` to Invoke AI

**File:** `backend/internal/services/chat_service.go`

**Changes:**

1. **Inject dependencies:**
   - `generateService *GenerateService` — for action-triggered meal plan generation
   - `aiClient *ai.Client` — for general conversational AI replies
   - `profileService *ProfileService` — to build user context for AI prompts
   - `pantryService *PantryService` — to build pantry snapshot for AI prompts

   Update `NewChatService()` constructor accordingly.

2. **Update `SendMessage` logic:**

```
SendMessage(ctx, userID, req):
  1. Save user message to DB (existing InsertUserMessage call)
  2. If req.Action is a known period type ("meal"/"day"/"week"/"month"):
     a. Call generateService.GeneratePlan(ctx, userID, req.Action, req.Message)
     b. Build a summary string from the ProposalResponse (period type, dates, meal count, cost, week totals)
     c. Save AI reply as bot message with the action field preserved
     d. Return both user message and bot message
  3. If req.Action is empty (plain chat message):
     a. Build a PromptRequest with user profile context and pantry snapshot
     b. Call BuildPrompt(PromptTemplateChat, promptReq)
     c. Call aiClient.Generate(ctx, GenerateRequest{Prompt: prompt})
     d. Save AI reply as bot message
     e. Return both user message and bot message
  4. If aiClient is nil (no API key configured):
     a. Save a canned bot reply: "I'm running in offline mode. Try the action buttons to generate a meal plan."
     b. Return both messages
```

3. **Return type:** Change the return type from `dto.ChatMessageResponse` to a new struct that includes both the user message and the bot response, so the frontend can render both immediately:

```go
type ChatSendResult struct {
    UserMessage ChatMessageResponse `json:"userMessage"`
    BotMessage  ChatMessageResponse `json:"botMessage"`
}
```

---

## Task 4: Update `ChatHandler` and DTOs

**File:** `backend/internal/transport/http/handlers/chat_handler.go`

- Update `SendMessage` to call the new `ChatService.SendMessage()` signature
- Return the `ChatSendResult` as JSON (HTTP 201)
- Keep `GetHistory` unchanged

**File:** `backend/internal/transport/http/dto/types.go`

Add the new response DTO:

```go
type ChatSendResponse struct {
    UserMessage ChatMessageResponse `json:"userMessage"`
    BotMessage  ChatMessageResponse `json:"botMessage"`
}
```

---

## Task 5: Wire New Dependencies in `app.go`

**File:** `backend/internal/app/app.go`

Update `ChatService` construction to inject `generateService`, `geminiClient`, `profileService`, and `pantryService`:

```go
chatService := services.NewChatService(chatRepo, generateService, geminiClient, profileService, pantryService)
```

This must happen **after** all those services are initialized (move the `chatService` init line to after `generateService` is created).

---

## Task 6: Update Frontend Chat Page to Handle New Response Shape

**File:** `frontend/src/js/pages/chat.js`

In `handleChatSubmit()` (around line 74-91), the response now contains both `userMessage` and `botMessage` objects:

```javascript
const data = await api.sendChatMessage(message);
// data = { userMessage: { id, role, content, createdAt }, botMessage: { id, role, content, createdAt } }
const replyText = data.botMessage.content || 'No response';
this.appendMessage(container, 'bot', replyText);
```

Remove the existing `appendMessage` call for the user message in `handleChatSubmit` — the bot reply already confirmed what was sent, and the user message was added to the DOM optimistically. Alternatively, let the backend response drive both renders:

- Replace the optimistic user message with the server-confirmed one
- Append the bot message from `data.botMessage`

**Also:** Remove or refactor the direct `api.generatePlan()` calls from chat action buttons (`handleGenerate` method at line 104). Instead, route them through `sendChatMessage` with the action field:

```javascript
async handleGenerate(periodType) {
    const container = document.getElementById('chatMessages');
    // ... generate message with action='periodType'
    const response = await api.sendChatMessage(`Generate ${periodType} plan`, periodType);
    // response.botMessage contains the response (which is handled by the backend as action-based generation)
    // Render the bot response and proposal preview
    this.renderBotResponse(response.botMessage);
}
```

This keeps the full conversation history in chat_messages.

---

## Task 7: Chat History Loading — Preserve Action Context

**File:** `frontend/src/js/pages/chat.js`

In `renderMessage()` (line 61-72), display an action badge when `msg.action` is set:

```javascript
const actionBadge = msg.action ? `<span class="action-badge">${msg.action}</span>` : '';
```

Add the action badge into the message template.

---

## Done Criteria

- [ ] A user can type "give me a weekly meal plan" in chat and get an AI-generated plan proposal back
- [ ] A user can type "what should I eat for dinner with chicken and rice?" and get a conversational AI reply with context from their profile/pantry
- [ ] Action buttons ("Generate Meal/Day/Week/Month") store both user message and bot reply with `action` set in chat_history
- [ ] Chat history correctly loads and displays both user and bot messages with timestamps
- [ ] Without `GEMINI_API_KEY`, the chat responds with a friendly offline-mode message
- [ ] The frontend renders the `ChatSendResponse` shape correctly
- [ ] No regressions: `GET /chat` history still works, `POST /plans/generate` still works standalone
