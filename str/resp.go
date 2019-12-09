package str

import "fmt"

func BuildRespStr(words []string) string{
	var count = getKeyWordsLen(words)
	var str = ""
	str += fmt.Sprintf("*%d\r\n", count)

	for _, value := range words {
		if value != "" {
			str += fmt.Sprintf("$%d\r\n%s\r\n", len(value), value)
		}
	}
	return str
}

func getKeyWordsLen(words []string) int {
	count := 0
	for _, value := range words {
		if value != "" {
			count++
		}
	}
	return count
}
