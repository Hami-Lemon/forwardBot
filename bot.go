package forwardBot

import (
	"context"
	"fmt"
	"forwardBot/push"
)

const (
	BiliLiveMsg = 1 << iota
	BiliDynMsg
	TikTokLiveMsg
)

type Bot struct {
	sources []Source
	sinks   []Sink
	ch      chan *push.Msg
}

func NewBot(buf int) *Bot {
	return &Bot{
		ch: make(chan *push.Msg, buf),
	}
}

func (b *Bot) AppendSource(s ...Source) {
	for _, source := range s {
		if source != nil {
			b.sources = append(b.sources, source)
		}
	}
}

func (b *Bot) AppendSink(s ...Sink) {
	for _, sink := range s {
		if sink != nil {
			b.sinks = append(b.sinks, sink)
		}
	}
}

func (b *Bot) Run(ctx context.Context) {
	for _, s := range b.sources {
		go s.Send(ctx, b.ch)
	}
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-b.ch:
			for _, s := range b.sinks {
				go func(s Sink) {
					err := s.Receive(msg)
					if err != nil {
						//TODO
						fmt.Println(err)
					}
				}(s)
			}
		}
	}
}
