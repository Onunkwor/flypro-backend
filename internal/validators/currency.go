package validators

import "github.com/go-playground/validator/v10"

//TODO: Expand this list or fetch from a reliable source
var validCurrencies = map[string]bool{
	"USD": true,
	"EUR": true,
	"GBP": true,
	"NGN": true,
	"JPY": true,
}

func CurrencyValidator(fl validator.FieldLevel) bool {
	code := fl.Field().String()
	return validCurrencies[code]
}
