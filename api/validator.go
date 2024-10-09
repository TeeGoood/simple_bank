package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/teegoood/simplebank/util"
)

var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	// check if it is string
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		return util.IsSupportedCurrency(currency)
	}
	
	return false 
}