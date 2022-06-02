package buffer

func intMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func mod(a, b int) int {
	return (a%b + b) % b
}
