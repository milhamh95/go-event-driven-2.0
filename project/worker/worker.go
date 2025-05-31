package worker

import (
	"context"
	"log"
)

type Task int

const (
	TaskIssueReceipt Task = iota
	TaskAppendToTracker
)

type Workerer interface {
	Send(...Message)
	Run()
}

type Message struct {
	Task     Task
	TicketID string
}

type Worker struct {
	queue           chan Message
	spreadsheetsAPI SpreadsheetsAPI
	receiptsService ReceiptsService
}

type SpreadsheetsAPI interface {
	AppendRow(ctx context.Context, sheetName string, row []string) error
}

type ReceiptsService interface {
	IssueReceipt(ctx context.Context, ticketID string) error
}

func NewWorker(
	spreadsheetsAPI SpreadsheetsAPI,
	receiptsService ReceiptsService,
) *Worker {
	return &Worker{
		queue:           make(chan Message, 100),
		spreadsheetsAPI: spreadsheetsAPI,
		receiptsService: receiptsService,
	}
}

func (w *Worker) Send(msg ...Message) {
	for _, m := range msg {
		w.queue <- m
	}
}

func (w *Worker) Run() {
	for msg := range w.queue {
		switch msg.Task {
		case TaskIssueReceipt:
			go func() {
				for {
					if err := w.receiptsService.IssueReceipt(context.Background(), msg.TicketID); err != nil {
						log.Printf("failed to issue receipt: %v", err)
						continue
					}
					break
				}
			}()
		case TaskAppendToTracker:
			go func() {
				if err := w.spreadsheetsAPI.AppendRow(context.Background(), "tickets-to-print", []string{msg.TicketID}); err != nil {
					log.Printf("failed to append to tracker: %v", err)
				}
			}()
		}
	}
}
