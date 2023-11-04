package utils
import (
	"strconv"
)

func MapFloatToString(m map[float64]float64) map[string]string {
	stringMap := make(map[string]string)

	for k, v := range m {
		stringKey := strconv.FormatFloat(k, 'f', -1, 64)
		stringValue := strconv.FormatFloat(v, 'f', -1, 64)
		stringMap[stringKey] = stringValue
	}
	return stringMap
}

func MapStringToFloat(m map[string]string) (map[float64]float64, error) {
	floatMap := make(map[float64]float64)

	for k, v := range m {
		floatKey, err := strconv.ParseFloat(k, 64)
		if err != nil {
			return nil, err
		}
		floatVal, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return nil, err
		}

		floatMap[floatKey] = floatVal
	}

	return floatMap, nil
}