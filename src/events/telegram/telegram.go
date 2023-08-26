package telegram

import (
	"errors"
	"example/main/src/clients/telegram"
	"example/main/src/events"
	"example/main/src/lib/e"
	"example/main/src/storage"
)
type Processor struct {
	tg *telegram.Client
	offset int
	storage storage.Storage
}
type Meta struct{
	ChatID int
	Username string
}
var ErrUnknownEventType = errors.New("unknown event type")
func New(client *telegram.Client,storage storage.Storage) *Processor {
	return &Processor{
		tg: client,
		storage: storage,
	}
}
func ( p *Processor) Fetch(limit int)([]events.Event,error){
	udates,err := p.tg.Updates(p.offset,limit)
	if err != nil{
		return nil,e.Wrap("can't get events", err)
	}
	if len(udates) == 0{
		return nil,nil
	}
	res:=make([]events.Event,0,len(udates))
	for _,u:= range udates{
		res=append(res, event(u))
	}
	p.offset=udates[len(udates)-1].ID+1
	return res,nil
}
func (p *Processor) Process( event events.Event) error {
	// сюда надо будет добавить что у нас будут кнопки чтобы код смог работать с кнопками 
	// не забудь добавить еще сюда работу с кнопками которые должны будут приниматься в методе fetch
	switch event.Type {
	case events.Message:
		p.processMessage(event)
	default:
		return e.Wrap("can't process message", ErrUnknownEventType)
	}
	
}
func (p *Processor) processMessage(event events.Event){
	meta,err:=meta(event)
}
func meta(event events.Event)(Meta, error)  {
	res, ok:=event.Meta.(Meta)
	if !ok{
		
	}
}
func event(upd telegram.Update) events.Event{
	updType:= fetchType(upd)
	res:=events.Event{
		Type: updType,
		Text: fetchText(upd),
	}
	if updType == events.Message{
		res.Meta =Meta{
			ChatID: upd.Message.Chat.ID,
			Username: upd.Message.From.Username,
		}
	}
	return res

}
func fetchType(upd telegram.Update) events.Type{
	if upd.Message == nil{
		return events.Unknown
	}
	return events.Message
}
func fetchText(upd telegram.Update) string{
	if upd.Message == nil{
			return ""
		}
	return upd.Message.Text
}