package cqhttpbotd

type Bot struct{}

func (b *Bot) Name() string {
	return "test"
}

func New() (d *Bot) {
	return &Bot{}
}
