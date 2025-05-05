package floatutil

import "strconv"

// FormatFloat format output of float
func FormatFloat(inp float64) string {
	if inp == float64(int64(inp)) {
		return strconv.FormatInt(int64(inp), 10)
	}
	return strconv.FormatFloat(inp, 'f', 2, 64)
}
