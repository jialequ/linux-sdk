package mathx

// AtLeast returns the greater of x or lower.
func AtLeast(x, lower float64) float64 {
	if x < lower {
		return lower
	}
	return x
}

func Between(x, lower, upper float64) float64 {
	if x < lower {
		return lower
	}
	if x > upper {
		return upper
	}
	return x
}
