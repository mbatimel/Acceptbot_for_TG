package telegram

import (
	"errors"
	"example/main/src/lib/e"
	"example/main/src/storage"
	"log"
	"net/url"
	"strings"
)
const (
	RndCmd = "/rnd"
	HelpCmd = "/help"
	StartCmd = "/start"
)
func (p *Processor) doCmd(text string, chatId int, username string)error {
	text = strings.TrimSpace(text)
	log.Printf("got new command '%s' from '%s'", text, username)
	if isAddCmd(text) {
		return p.savePage(chatId,text,username)
	}
	switch text {
		case RndCmd:
			return p.Sendrandom(chatId, username)
		case HelpCmd:
			return p.SendHelp(chatId)
		case StartCmd:
			return p.SendHello(chatId)

		default:
			return p.tg.SendMassage(chatId, msgUnknownCommand)

	}
}
func (p *Processor) savePage(chatID int, pageURL string, username string) (err error) {
	defer func() {err=e.Wrap("can't do command: save page", err) }()
	page := &storage.Page{
		URL: pageURL,
		UserName: username,

	}
	isExists,err :=p.storage.IsExistst(page)
	if err != nil {
		return err
	}
	if isExists{
		return p.tg.SendMassage(chatID, msgAlreadyExists)
	}
	if err :=p.storage.Save(page); err != nil {
		return err
	}
	if err:=p.tg.SendMassage(chatID, msgSaved);err != nil {
		return err
	}
	return nil
	
}
func (p *Processor) Sendrandom(chatId int, username string) (err error) {
	defer func() {err=e.Wrap("can't do command: send random msg", err) }()
	page, err:= p.storage.PickRandom(username)
	if err != nil && !errors.Is(err,storage.ErrNoSavedPages) {
		return err
	}
	if errors.Is(err,storage.ErrNoSavedPages){
		return p.tg.SendMassage(chatId,msgNoSavedPages)
	}
	if err:=p.tg.SendMassage(chatId,page.URL); err != nil{
		return err
	}
	return p.storage.Remove(page)
}
func (p *Processor) SendHello(chatId int) (error) {
	return p.tg.SendMassage(chatId, msgHello)
}
func (p *Processor) SendHelp(chatId int) (error) {
	return p.tg.SendMassage(chatId, msgHelp)

}


func isAddCmd(text string) bool {
	return isUrl(text)
}
func isUrl(text string) bool {
	u, err:=url.Parse(text)
	return err == nil && u.Host != ""
}