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
			"message": "No credentials provided.",
		})
	}

	amtFL, okamt := jsonMap["amount"].(float64)
	cur, okcur := jsonMap["currency"].(string)

	if !okamt || !okcur {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid input",
		})
	}

	amt := strconv.FormatFloat(amtFL, 'f', 2, 64)

	amount := paypalsdk.Amount{
		Currency: cur,
		Total:    amt,
	}

	payment := paypalservice.MakePaypalPayment(amount)

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
