package utils

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

//SplitServer is
func SplitServer(server string) (string, uint16) {
	s := strings.Split(server, ":")
	if len(s) != 2 {
		fmt.Println("Error", s)
		return "", 0
	}
	fmt.Printf("IP: %s, Port: %s", s[0], s[1])
	p, err := strconv.Atoi(s[1])
	if err != nil {
		fmt.Println("Error", err)
		return "", 0
	}
	return s[0], uint16(p)
}

//IsHTTPAPIHealth is check http error
func IsHTTPAPIHealth(ip, port, path, schem string) bool {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	url := fmt.Sprintf("%s://%s:%s%s", schem, ip, port, path)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer resp.Body.Close()

	_ = resp
	return true
}
