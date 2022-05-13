package examples

import "strings"

// Source: https://dave.cheney.net/2019/05/07/prefer-table-driven-tests
func Split(s, sep string) []string {
	var result []string
	i := strings.Index(s, sep)
	for i > -1 {
		result = append(result, s[:i])
		s = s[i+len(sep):]
		i = strings.Index(s, sep)
	}
	return append(result, s)
}