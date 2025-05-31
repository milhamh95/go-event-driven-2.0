package http

import (
	"net/http"
	"tickets/worker"

	"github.com/labstack/echo/v4"
)

type ticketsConfirmationRequest struct {
	Tickets []string `json:"tickets"`
}

func (h Handler) PostTicketsConfirmation(c echo.Context) error {
	var request ticketsConfirmationRequest
	err := c.Bind(&request)
	if err != nil {
		return err
	}

	for _, ticket := range request.Tickets {
		taskIssueReceiptMessage := worker.Message{
			Task:     worker.TaskIssueReceipt,
			TicketID: ticket,
		}
		h.worker.Send(taskIssueReceiptMessage)

		taskAppendToTrackerMessage := worker.Message{
			Task:     worker.TaskAppendToTracker,
			TicketID: ticket,
		}

		h.worker.Send(taskAppendToTrackerMessage)
	}

	return c.NoContent(http.StatusOK)
}
