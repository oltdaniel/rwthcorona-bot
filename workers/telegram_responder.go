package workers

import (
	"github.com/oltdaniel/rwth-coronabot/utils"
)

var TelegramResponderQueue chan (*utils.TelegramWebhookUpdate) = make(chan *utils.TelegramWebhookUpdate)

func TelegramResponder() {
	for {
		update := <-TelegramResponderQueue

		utils.TelegramSendMessage(&utils.TelegramRequestSendMessage{
			ChatId:           update.Message.Chat.Id,
			Text:             "pong",
			ReplyToMessageId: update.Message.MessageId,
		})
	}
}
