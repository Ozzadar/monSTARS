package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ozzadar/monSTARS/db"

	"github.com/logpacker/PayPal-Go-SDK"
	"github.com/ozzadar/monSTARS/models"
	"github.com/ozzadar/monSTARS/services/paypalservice"

	"github.com/labstack/echo"
)

func Donate(c echo.Context) error {

	user, ok := c.Get("token_user").(*models.User)

	if !ok {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Internal error, please try again",
		})
	}

	jsonMap := make(map[string]interface{})

	err := json.NewDecoder(c.Request().Body).Decode(&jsonMap)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "No payment information provided",
		})
	}

	amtFL, amterr := strconv.ParseFloat(jsonMap["amount"].(string), 64)
	cur, okcur := jsonMap["currency"].(string)

	if amterr != nil || !okcur {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid input",
		})
	}

	amt := strconv.FormatFloat(amtFL, 'f', 2, 64)

	amount := paypalsdk.Amount{
		Currency: cur,
		Total:    amt,
	}

	payment := paypalservice.MakePaypalDonation(amount)

	if payment == nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Internal error, please try again",
		})
	}

	transaction := models.NewTransaction(*user, payment)

	success, _ := db.NewTransaction(transaction)

	if success {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"approval_url": payment.Links[1].Href,
		})

	}

	return c.JSON(http.StatusInternalServerError, map[string]string{
		"message": "Internal error, please try again",
	})
}

func CompletePayment(c echo.Context) error {
	jsonMap := make(map[string]interface{})

	err := json.NewDecoder(c.Request().Body).Decode(&jsonMap)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "No payment information provided",
		})
	}

	payerID, okpayer := jsonMap["payerID"].(string)
	paymentID, okpayment := jsonMap["paymentId"].(string)

	if !okpayer || !okpayment {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid payment information",
		})
	}

	executeResponse := paypalservice.ExecutePayment(payerID, paymentID)

	if executeResponse == nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Failed to complete payment",
		})
	}
	transaction := db.GetTransaction(paymentID)

	if transaction == nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Transaction doesnt exist",
		})
	}

	transaction.Executed(executeResponse)
	ok, errorMsg := db.SaveTransaction(transaction)

	if !ok {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": errorMsg,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":  "Payment completed successfully.",
		"amount":   transaction.Execution.Transactions[0].Amount.Total,
		"currency": transaction.Execution.Transactions[0].Amount.Currency,
	})

}
