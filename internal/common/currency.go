package common

import "math/rand"

// Constants for all supported currencies
const (
	IDR = "IDR"
	USD = "USD"
	EUR = "EUR"
	CAD = "CAD"
)

// IsSupportedCurrency returns true if the currency is supported
func IsSupportedCurrency(currency string) bool {
	switch currency {
	case IDR, USD, EUR, CAD:
		return true
	}
	return false
}

// RandomCurrency generates a random currency code
func RandomCurrency() string {
	currencies := []string{IDR, EUR, USD, CAD}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}
