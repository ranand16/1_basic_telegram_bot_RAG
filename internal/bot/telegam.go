package bot

import (
	"1_basic_RAG/internal/document"
	"1_basic_RAG/internal/rag"
	"1_basic_RAG/internal/store"
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UserState string

const (
	StateNone    UserState = "none"
	StateFeeding UserState = "feeding"
	StateAsking  UserState = "asking"
)

type TelegramBot struct {
	api        *tgbotapi.BotAPI
	ragService *rag.Service
	store      *store.QdrantStore
	userStates map[int64]UserState
}

func NewTelegramBot(token string, r *rag.Service, s *store.QdrantStore) (*TelegramBot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return &TelegramBot{
		api:        api,
		ragService: r,
		store:      s,
		userStates: make(map[int64]UserState),
	}, nil
}

func (b *TelegramBot) Start(ctx context.Context) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		chatId := update.Message.Chat.ID
		text := update.Message.Text

		switch text {
		case "/start":
			b.reply(chatId, "Welcome to the Personal RAG Assistant! Use /feed to input data, /ask to ask questions and /exit to reset.")
			continue
		case "/feed":
			b.userStates[chatId] = StateFeeding
			b.reply(chatId, "What would you like me to know about you? Send me the information and I'll remember it.")
			continue
		case "/ask":
			b.userStates[chatId] = StateAsking
			b.reply(chatId, "What would you like to ask? Send me your question and I'll do my best to answer it based on what I know about you.")
			continue
		case "/exit":
			delete(b.userStates, chatId)
			b.reply(chatId, "State reset. Use /feed to input data or /ask to ask questions.")
			continue
		}

		state := b.userStates[chatId]
		switch state {
		case StateFeeding:
			docs, err := document.ProcessText(text)
			if err != nil {
				b.reply(chatId, "Sorry, I couldn't process that information. Please try again.")
				continue
			}
			_, err = b.store.AddDocuments(ctx, docs)
			if err != nil {
				b.reply(chatId, "Sorry, I couldn't save that information. Please try again.")
				continue
			}
			b.reply(chatId, "Information saved successfully! You can use /ask to ask questions based on this data.")
		case StateAsking:
			answer, err := b.ragService.Ask(ctx, text)
			if err != nil {
				b.reply(chatId, "Sorry, I couldn't answer that question. Please try again.")
				continue
			}
			response := "------ THE QUESTION ------\n" + text + "\n\n------ THE ANSWER ------\n" + answer
			b.reply(chatId, response)
		default:
			b.reply(chatId, "I didn't understand that. Please use /feed to input data or /ask to ask questions.")
		}
	}
}

func (b *TelegramBot) reply(chatId int64, message string) {
	msg := tgbotapi.NewMessage(chatId, message)
	b.api.Send(msg)
}
