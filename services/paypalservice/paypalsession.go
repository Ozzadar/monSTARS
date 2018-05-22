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

// MakePaypalDonation will create a paypal payment and return the CreatePaymentResponse object with relevant links
// It will also store the payment object within the database for retrieval in the future.
func MakePaypalDonation(amount paypalsdk.Amount) *paypalsdk.CreatePaymentResp {
	redirectURI := "http://localhost:4000/donate"
	cancelURI := "http://localhost:4000/donate"

	p := paypalsdk.Payment{
		Intent: "sale",
		Payer: &paypalsdk.Payer{
			PaymentMethod: "paypal",
		},
		RedirectURLs: &paypalsdk.RedirectURLs{
			CancelURL: cancelURI,
			ReturnURL: redirectURI,
		},
		Transactions: []paypalsdk.Transaction{
			paypalsdk.Transaction{
				Amount:      &amount,
				Description: "Donation to monSTARS",
				ItemList: &paypalsdk.ItemList{
					Items: []paypalsdk.Item{
						paypalsdk.Item{
							Quantity: 1,
							Name:     "Donation",
							Price:    amount.Total,
							Currency: amount.Currency,
							SKU:      "donate",
						},
					},
				},
			}},
	}

	paymentResponse, err := PaypalSession.CreatePayment(p)

	if err != nil {
		logging.Log("payment result failed")
		return nil
	}

	return paymentResponse
}

//ExecutePayment will finalize and execute an approved PayPal payment
func ExecutePayment(PayerID string, PaymentID string) *paypalsdk.ExecuteResponse {
	executeResult, err := PaypalSession.ExecuteApprovedPayment(PaymentID, PayerID)

	if err != nil {
		return nil
	}

	return executeResult
}
