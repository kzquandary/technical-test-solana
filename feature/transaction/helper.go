package transaction

import "strconv"

func ConvertBalances(balances []interface{}) []float64 {
	var convertedBalances []float64
	for _, balance := range balances {
		switch v := balance.(type) {
		case float64:
			convertedBalances = append(convertedBalances, v)
		case string:
			// if balance is a string that represents a number, convert it
			if floatValue, err := strconv.ParseFloat(v, 64); err == nil {
				convertedBalances = append(convertedBalances, floatValue)
			}
		}
	}
	return convertedBalances
}
