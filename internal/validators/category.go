package validators

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

var validCategories = map[string]bool{
	"travel":   true,
	"meals":    true,
	"office":   true,
	"supplies": true,
}

func CategoryValidator(fl validator.FieldLevel) bool {
	val := strings.ToLower(strings.TrimSpace(fl.Field().String()))
	return validCategories[val]
}
