package paypalservice

import (
	"os"

	"github.com/klouds/kDaemon/logging"
	"github.com/logpacker/PayPal-Go-SDK"
	"github.com/ozzadar/monSTARS/config"
)

var (
	// PaypalSession The global paypal session
	PaypalSession *paypalsdk.Client
)

// Init the paypal service
func Init() {
	var clientID string
	var secretID string
	var SanboxOrLive string

	clientID, err := config.Config.GetString("default", "paypal_client")
	if err != nil {
		logging.Log("Problem with config file! (paypal_client)")
	}

	secretID, err = config.Config.GetString("default", "paypal_secret")
	if err != nil {
		logging.Log("Problem with config file! (paypal_secret)")
	}

	SanboxOrLive = paypalsdk.APIBaseSandBox

	session, err := paypalsdk.NewClient(clientID, secretID, SanboxOrLive)

	if err != nil {
		logging.Log("Paypal initialize failed")
		PaypalSession = nil
		panic("FAILED")
	}
	session.SetLog(os.Stdout)

	PaypalSession = session
	_, err = PaypalSession.GetAccessToken()

	if err != nil {
		logging.Log("GetAccessToken failed")
		PaypalSession = nil

		panic("FAILED")
	}

	logging.Log("Paypal service initilalized!")
}

// MakePaypalPayment will create a paypal payment and return the PaymentResponse object with relevant links
// It will also store the payment object within the database for retrieval in the future.
func MakePaypalPayment(amount paypalsdk.Amount) *paypalsdk.PaymentResponse {
	redirectURI := "http://localhost:4000/transaction/complete"
	cancelURI := "http://localhost:4000/transaction/cancelled"
	description := "Description for this payment"

	paymentResult, err := PaypalSession.CreateDirectPaypalPayment(amount, redirectURI, cancelURI, description)

	if err != nil {
		logging.Log("payment result failed")
		return nil
	}

	return paymentResult
}
