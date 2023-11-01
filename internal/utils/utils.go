package utils

import "strconv"

func MapFloatToString(m map[float64]float64) map[string]string {
	sMap := make(map[string]string)

	for k, v := range m {
		sKey := strconv.FormatFloat(k, 'f', -1, 64)
		sVal := strconv.FormatFloat(v, 'f', -1, 64)
		sMap[sKey] = sVal
	}
	return sMap
}

func MapStringToFloat(m map[string]string) (map[float64]float64, error) {
	fMap := make(map[float64]float64)

	for k, v := range m {
		fKey, err := strconv.ParseFloat(k, 64)
		if err != nil {
			return nil, err
		}
		fVal, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return nil, err
		}

		fMap[fKey] = fVal
	}

	return fMap, nil
}