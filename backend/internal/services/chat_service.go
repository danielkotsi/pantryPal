package services

import (
	"context"
	"fmt"
	"strings"

	"pantrypal/backend/internal/modules/ai"
	"pantrypal/backend/internal/repositories"
	"pantrypal/backend/internal/transport/http/dto"
)

type ChatService struct {
	chat            *repositories.ChatRepository
	generateService *GenerateService
	aiClient        *ai.Client
	profileService  *ProfileService
	pantryService   *PantryService
}

func NewChatService(chat *repositories.ChatRepository, generateService *GenerateService, aiClient *ai.Client, profileService *ProfileService, pantryService *PantryService) *ChatService {
	return &ChatService{
		chat:            chat,
		generateService: generateService,
		aiClient:        aiClient,
		profileService:  profileService,
		pantryService:   pantryService,
	}
}

func (s *ChatService) SendMessage(ctx context.Context, userID string, req dto.ChatSendRequest) (dto.ChatSendResponse, error) {
	userMsg, err := s.chat.InsertUserMessage(ctx, userID, req.Message, req.Action)
	if err != nil {
		return dto.ChatSendResponse{}, err
	}

	userMsgResp := dto.ChatMessageResponse{
		ID: userMsg.ID, Role: userMsg.Role, Action: userMsg.Action,
		Content: userMsg.Content, CreatedAt: userMsg.CreatedAt,
	}

	var botContent string

	if req.Action != "" && s.generateService != nil {
		result, planErr := s.generateService.GeneratePlan(ctx, userID, req.Action, req.Message)
		if planErr == nil {
			botContent = formatPlanSummary(result)
		} else {
			botContent = "I couldn't generate a plan right now. Please try again."
		}
	} else if s.aiClient != nil {
		promptReq := buildPromptContext(ctx, s.profileService, s.pantryService, userID, "", req.Message)
		prompt, buildErr := ai.BuildPrompt(ai.PromptTemplateChat, promptReq)
		if buildErr == nil {
			geminiResp, geminiErr := s.aiClient.Generate(ctx, ai.GenerateRequest{Prompt: prompt, ResponseMIMEType: "text/plain"})
			if geminiErr == nil {
				botContent = geminiResp.Text
			}
		}
		if botContent == "" {
			botContent = "I'm having trouble responding right now. Please try again."
		}
	} else {
		botContent = "AI assistant is not available. Use the action buttons above to generate meal plans."
	}

	var botMsgResp dto.ChatMessageResponse
	if botContent != "" {
		botStored, storeErr := s.chat.InsertBotMessage(ctx, userID, botContent)
		if storeErr == nil {
			botMsgResp = dto.ChatMessageResponse{
				ID: botStored.ID, Role: botStored.Role, Content: botStored.Content, CreatedAt: botStored.CreatedAt,
			}
		}
	}

	return dto.ChatSendResponse{UserMessage: userMsgResp, BotMessage: botMsgResp}, nil
}

func formatPlanSummary(result GeneratePlanResult) string {
	proposal := result.Proposal
	plan := proposal.Plan

	if result.FallbackActive {
		var b strings.Builder
		b.WriteString("Fallback mode active — using seeded data instead of AI.\n")
		b.WriteString(fmt.Sprintf("%s plan: %s → %s\n", plan.PeriodType, plan.StartDate, plan.EndDate))
		b.WriteString(fmt.Sprintf("%d days, %d meals", len(proposal.Days), countMeals(proposal)))
		return b.String()
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf("%s plan: %s → %s\n", plan.PeriodType, plan.StartDate, plan.EndDate))
	b.WriteString(fmt.Sprintf("Status: %s\n", plan.Status))
	b.WriteString(fmt.Sprintf("%d days, %d meals", len(proposal.Days), countMeals(proposal)))
	if plan.AICostCentsTotal > 0 {
		b.WriteString(fmt.Sprintf("\nCost: $%.2f", float64(plan.AICostCentsTotal)/100))
	}
	return b.String()
}

func countMeals(proposal dto.ProposalResponse) int {
	count := 0
	for _, d := range proposal.Days {
		if d.Sections.Breakfast != nil {
			count++
		}
		if d.Sections.Lunch != nil {
			count++
		}
		if d.Sections.Dinner != nil {
			count++
		}
		if d.Sections.Snacks != nil {
			count++
		}
	}
	return count
}

func (s *ChatService) GetHistory(ctx context.Context, userID string, limit int) (dto.ChatHistoryResponse, error) {
	messages, err := s.chat.ListRecent(ctx, userID, limit)
	if err != nil {
		return dto.ChatHistoryResponse{}, err
	}
	out := make([]dto.ChatMessageResponse, 0, len(messages))
	for _, msg := range messages {
		out = append(out, dto.ChatMessageResponse{
			ID:        msg.ID,
			Role:      msg.Role,
			Action:    msg.Action,
			Content:   msg.Content,
			CreatedAt: msg.CreatedAt,
		})
	}
	return dto.ChatHistoryResponse{Messages: out}, nil
}
