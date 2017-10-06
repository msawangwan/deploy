package strutil

// IsNullOrEmpty ...
func IsNullOrEmpty(s string) bool { return s == "" || len(s) < 1 }
