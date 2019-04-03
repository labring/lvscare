package utils

import (
	"fmt"
	"strings"
)

//SplitServer is
func SplitServer(server string) (string, string) {
	s := strings.Split(server, ":")
	if len(s) != 2 {
		fmt.Println("Error", s)
		return "", ""
	}
	fmt.Printf("IP: %s, Port: %s", s[0], s[1])
	return s[0], s[1]
}
