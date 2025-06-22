package kernel

// ReverseBytes reverses bytes for display
func ReverseBytes(data []byte) []byte {
	result := make([]byte, len(data))
	for i, b := range data {
		result[len(data)-1-i] = b
	}
	return result
}
