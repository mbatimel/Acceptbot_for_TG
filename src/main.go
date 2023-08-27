package main

import (
	tgclient "example/main/src/clients/telegram"
	"example/main/src/consumer/event-consumer"
	tgevent "example/main/src/events/telegram"
	"example/main/src/storage/files"
	"flag"
	"log"
)
const (
	tgBotHost = "api.telegram.org"
	storagePath = "storage"
	batchSize = 1000
)
func main() {
	eventsProcessor := tgevent.New(tgclient.New(tgBotHost,mustToken()),files.NewStorage(storagePath))

	 log.Print("server started")
	 
	 err:= event_consumer.New(eventsProcessor, eventsProcessor, batchSize).Start()
	 if err!= nil{
		log.Print("server is soported", err)
	 }
	 

}
func mustToken() string {
	token := flag.String("tg-bot-token", "", "bot token sdf")
	flag.Parse()
	if *token == "" {
		log.Fatal("empty token")
	}
	return *token
}