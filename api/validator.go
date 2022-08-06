package api

import (
	"bank-mvp/util"

	"github.com/go-playground/validator/v10"
)

var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if currency, ok := fieldLevel.Field().Interface().(string); ok { // type assertion
		return util.IsSupportedCurrency(currency)
	}

	return false
}

var validPassword validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if password, ok := fieldLevel.Field().Interface().(string); ok {
		if model, ok := fieldLevel.Parent().Interface().(createUserRequest); ok {
			return password == model.PasswordRepeated
		}
	}

	return false
}