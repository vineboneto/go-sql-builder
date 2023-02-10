package lib

func checkHasValue(v any) bool {
	switch v := v.(type) {
	case int:
		if v != 0 {
			return true
		}
	case string:
		if v != "" {
			return true
		}
	case float64:
		if v != float64(0) {
			return true
		}
	case bool:
		return v
	case nil:
		return false
	default:
		return false
	}
	return false
}
