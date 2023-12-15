package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func main() {
	fileIn, err := os.Open("data.json")
	if err != nil {
		log.Fatalln(err)
	}
	defer fileIn.Close()

	out := make([]string, 0)

	err = json.NewDecoder(fileIn).Decode(&out)
	if err != nil {
		log.Fatalln(err)
	}

	handler := func(ctx context.Context, b *bot.Bot, update *models.Update) {
		if update.Message == nil {
			return
		}

		if strings.Index(update.Message.Text, "/comment") == 0 {
			if update.Message.ReplyToMessage == nil {
				return
			}

			_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:           update.Message.Chat.ID,
				Text:             randString(out),
				ReplyToMessageID: update.Message.ReplyToMessage.ID,
			})

			_, _ = b.DeleteMessage(ctx, &bot.DeleteMessageParams{
				ChatID:    update.Message.Chat.ID,
				MessageID: update.Message.ID,
			})
		}

		if strings.Index(update.Message.Text, "/quote") == 0 {
			replyToMessageID := update.Message.ID

			ok, _ := b.DeleteMessage(ctx, &bot.DeleteMessageParams{
				ChatID:    update.Message.Chat.ID,
				MessageID: update.Message.ID,
			})

			if ok {
				replyToMessageID = 0
			}

			_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:           update.Message.Chat.ID,
				Text:             randString(out),
				ReplyToMessageID: replyToMessageID,
			})
		}
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
	}

	b, err := bot.New(os.Getenv("TG_TOKEN"), opts...)
	if err != nil {
		panic(err)
	}

	b.Start(ctx)
}

func randString(arr []string) string {
	return arr[rand.Intn(len(arr))]
}
