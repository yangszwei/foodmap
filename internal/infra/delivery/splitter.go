package delivery

import "strings"

// split split string by comma and remove empty ones
func Split(str string) (result []string) {
	for _, i := range strings.Split(str, ",") {
		if i != "" {
			result = append(result, i)
		}
	}
	return
}
