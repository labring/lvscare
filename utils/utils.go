package utils

import (
	"crypto/tls"
	"fmt"
	"net/http"
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

//IsKubeAPIHealth check Http GET is return 200 OK
/*
func IsKubeAPIHealth(ip, port, path, schem string) bool {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	url := fmt.Sprintf("%s://%s:%s%s", schem, ip, port, path)
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println(err)
		return false
	}
	body, err := ioutils.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return string(body) == "ok"
}
*/

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
