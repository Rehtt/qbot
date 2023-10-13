package main

import (
	"fmt"
	"log/slog"

	"github.com/Rehtt/qbot"
	cqhttpbotd "github.com/Rehtt/qbot/cqhttp_bot_d"
)

func main() {
	b, _ := qbot.New(cqhttpbotd.New(), &qbot.Config{
		Logger: slog.Default(),
	})
	b.Logger.Info("test")
	fmt.Println(b.Name())
}
