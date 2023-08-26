package main

import (
	"example/main/src/clients/telegram"
	"flag"
	"log"
)
const (
	tgBotHost = "api.telegram.org"
)
func main() {
	tgClient := telegram.New(tgBotHost,mustToken())
	
}
func mustToken() string {
	token := flag.String("token-bot-token", " ", "bot token")
	flag.Parse()
	if *token == "" {
		log.Fatal("empty token")
	}
	return *token
}