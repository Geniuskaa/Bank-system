package dto

import (
	"time"
)

type Status string

const (
	Completed Status = "COMPLETED"
	Refunded  Status = "REFUNDED"
	Blocked   Status = "BLOCKED"
)

type TransactionDto struct {
	Amount int64     `json:"amount"`
	MCC    string    `json:"mcc"`
	Status Status    `json:"status"`
	Date   time.Time `json:"date"`
}
