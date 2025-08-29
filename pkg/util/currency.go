package util

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
