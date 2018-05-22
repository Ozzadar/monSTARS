package models

import (
	"github.com/logpacker/PayPal-Go-SDK"
)

// Transaction struct for holding transactions with database
type Transaction struct {
	ID        string                       `json:"id" gorethink:"id"`
	Owner     string                       `json:"owner" gorethink:"owner"`
	State     string                       `json:"state" gorethink:"state"`
	Payment   *paypalsdk.CreatePaymentResp `json:"payment" gorethink:"payment"`
	Execution *paypalsdk.ExecuteResponse   `json:"execution" gorethink:"execution"`
}

//NewTransaction will create a new transaction in the database
func NewTransaction(owner User, payment *paypalsdk.CreatePaymentResp) *Transaction {
	newTransaction := &Transaction{}

	newTransaction.ID = payment.ID
	newTransaction.Owner = owner.Username
	newTransaction.State = "Created"
	newTransaction.Payment = payment

	return newTransaction
}

//Executed will save a successful execution response from paypal into the transaction object
func (t *Transaction) Executed(response *paypalsdk.ExecuteResponse) {
	t.Execution = response
	t.State = "Complete"
}
