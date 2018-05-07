package models

import (
	"github.com/logpacker/PayPal-Go-SDK"
)

// Transaction struct for holding transactions with database
type Transaction struct {
	ID      string                     `json:"id" gorethink:"id"`
	Owner   User                       `json:"owner" gorethink:"owner"`
	State   string                     `json:"state" gorethink:"state"`
	Payment *paypalsdk.PaymentResponse `json:"payment" gorethink:"payment"`
}

//
func NewTransaction(owner User, payment *paypalsdk.PaymentResponse) *Transaction {
	newTransaction := &Transaction{}

	newTransaction.ID = payment.ID
	newTransaction.Owner = owner
	newTransaction.State = "waiting_for_completion"
	newTransaction.Payment = payment

	return newTransaction
}
