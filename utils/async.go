package utils

import (
	"context"
	"log"
	"time"
)

type AsyncProcessor struct {
	auditCh  chan string
	notifyCh chan string
}

func NewAsyncProcessor(auditBuf, notifyBuf int) *AsyncProcessor {
	return &AsyncProcessor{
		auditCh:  make(chan string, auditBuf),
		notifyCh: make(chan string, notifyBuf),
	}
}

func (p *AsyncProcessor) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-p.auditCh:
				log.Printf("AUDIT %s", msg)
			}
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-p.notifyCh:
				time.Sleep(5 * time.Millisecond)
				log.Printf("NOTIFY %s", msg)
			}
		}
	}()
}

func (p *AsyncProcessor) Audit(msg string) {
	select {
	case p.auditCh <- msg:
	default:
		log.Printf("AUDIT_DROP %s", msg)
	}
}

func (p *AsyncProcessor) Notify(msg string) {
	select {
	case p.notifyCh <- msg:
	default:
		log.Printf("NOTIFY_DROP %s", msg)
	}
}
