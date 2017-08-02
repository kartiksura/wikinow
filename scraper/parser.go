package scraper

import (
	"strings"
)

func getWikiTitles(body string) (ans [][]string) {
	var op []string

	ip := strings.Split(body, "[[")
	for i := range ip {
		if i == 0 {
			continue
		}
		ip1 := strings.Split(ip[i], "]]")
		if len(ip1) > 0 {
			op = append(op, ip1[0])
		}
	}
	for i := range op {
		ans = append(ans, strings.Split(op[i], "|"))
	}
	return
}
