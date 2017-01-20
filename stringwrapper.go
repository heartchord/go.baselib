package goblazer

// CheckStringEmpty :
func CheckStringEmpty(s *string) bool {
	if s == nil || len(*s) <= 0 {
		return false
	}

	return true
}
