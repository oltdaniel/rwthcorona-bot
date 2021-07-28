package workers

import (
	"fmt"

	"github.com/oltdaniel/rwthcorona-bot/utils"
)

var TelegramResponderQueue chan (*utils.TelegramWebhookUpdate) = make(chan *utils.TelegramWebhookUpdate)

func TelegramResponder() {
	for {
		update := <-TelegramResponderQueue
		msg := update.Message.Text
		switch {
		case msg == "/gestern":
			dataset := utils.DATASET.Get()
			day, err := dataset.Yesterday()
			if err != nil {
				continue
			}
			total := day.Total()
			response := fmt.Sprintf("Anzahl: `%.2f`\nRate: `%.2f`", total.AnzahlWoche, total.RateWoche)
			utils.TelegramSendMessage(&utils.TelegramRequestSendMessage{
				ChatId:           update.Message.Chat.Id,
				Text:             response,
				ReplyToMessageId: update.Message.MessageId,
				ParseMode:        "MarkdownV2",
			})
		default:
			utils.TelegramSendMessage(&utils.TelegramRequestSendMessage{
				ChatId:           update.Message.Chat.Id,
				Text:             "not implemented yet",
				ReplyToMessageId: update.Message.MessageId,
			})
		}
	}
}
