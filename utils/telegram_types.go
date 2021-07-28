package utils

type TelegramWebhookUpdate struct {
	UpdateId int64                        `json:"update_id"`
	Message  TelegramWebhookUpdateMessage `json:"message"`
}

type TelegramWebhookUpdateMessage struct {
	MessageId int64                            `json:"message_id"`
	From      TelegramWebhookUpdateMessageUser `json:"from"`
	Chat      TelegramWebhookUpdateMessageChat `json:"chat"`
	Text      string                           `json:"text"`
}

type TelegramWebhookUpdateMessageUser struct {
	Id        int64  `json:"id"`
	IsBot     bool   `json:"is_bot"`
	FirstName string `json:"first_name"`
	Username  string `json:"username"`
}

type TelegramWebhookUpdateMessageChat struct {
	Id   int64  `json:"id"`
	Type string `json:"type"`
}

type TelegramRequestSendMessage struct {
	ChatId                int64  `json:"chat_id"`
	Text                  string `json:"text"`
	ParseMode             string `json:"parse_mode"`
	DisableWebPagePreview bool   `json:"disable_web_page_preview"`
	ReplyToMessageId      int64  `json:"reply_to_message_id"`
	DisableNotification   bool   `json:"disable_notification"`
}
