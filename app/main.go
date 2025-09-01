package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
)

var CRLF = "\r\n"

type HttpRequest struct {
	Method  string
	URL     string
	Version string
	Headers map[string]string
}

type HttpResponse struct {
	Status  int
	Version string
	Headers map[string]string
	Body    []byte
}

func main() {
	fmt.Println("Logs from HTTP Server")

	l, err := net.Listen("tcp", "0.0.0.0:4221")

	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()

		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go do(conn)
	}
}

func do(conn net.Conn) {
	defer conn.Close()

	buff := make([]byte, 1024)

	n, err := conn.Read(buff)

	if err != nil {
		fmt.Errorf("Error reading from connection: " + err.Error())
	}

	httpRequest := parseRequest(buff[0:n])

	respone := route(httpRequest)

	encodedResp := respone.Encode()

	sendResponse(encodedResp, conn)
}

func sendResponse(resp []byte, conn net.Conn) {
	conn.Write(resp)
}

func parseRequest(buff []byte) *HttpRequest {
	lines := splitter(buff, CRLF)
	reqLine := splitter(lines[0], " ")
	method := reqLine[0]
	url := reqLine[1]
	ver := reqLine[2]
	hdrs := createHeadersMap(lines[1 : len(lines)-1])
	return &HttpRequest{
		Method:  string(method),
		URL:     string(url),
		Version: string(ver),
		Headers: hdrs,
	}
}

func route(req *HttpRequest) *HttpResponse {
	if req.URL == "/" {
		return buildResponse(200, req.Version)
	} else if hasPrefix(req.URL, "/echo") {
		path := extractEcho(req.URL[5:])
		return buildResponseWithBody(200, req.Version, []byte(path), "text/plain")
	} else {
		return buildResponse(404, req.Version)
	}
}

func (r *HttpResponse) Encode() []byte {
	res := make([]byte, 0)
	res = append(res, r.Version...)
	res = append(res, ' ')
	res = append(res, buildStatus(r.Status)...)
	res = append(res, CRLF...)

	for k, v := range r.Headers {
		res = append(res, k...)
		res = append(res, ": "...)
		res = append(res, v...)
		res = append(res, CRLF...)
	}

	res = append(res, CRLF...)
	res = append(res, r.Body...)
	return res
}

func createHeadersMap(hrs [][]byte) map[string]string {
	m := make(map[string]string)

	for i := range hrs {
		hr := splitter(hrs[i], ": ")
		k, v := string(hr[0]), string(hr[1])
		m[k] = v
	}

	return m
}

func splitter(buff []byte, sep string) [][]byte {
	s := make([][]byte, 0)

	a := 0

	for i := 0; i < len(buff)-len(sep); i++ {
		if string(buff[i:i+len(sep)]) == sep {
			s = append(s, buff[a:i])
			a = i + len(sep)
		}
	}

	if a != len(buff) {
		s = append(s, buff[a:])
	}

	return s
}

func buildStatus(status int) string {
	switch status {
	case 200:
		return "200 OK"
	case 404:
		return "404 Not Found"
	}

	return ""
}

func buildResponse(status int, version string) *HttpResponse {
	return &HttpResponse{
		Status:  status,
		Version: version,
		Headers: make(map[string]string),
		Body:    []byte(""),
	}
}

func buildResponseWithBody(status int, version string, body []byte, ctype string) *HttpResponse {
	hrs := make(map[string]string)
	hrs["Content-Type"] = ctype
	hrs["Content-Length"] = strconv.Itoa(len(body))
	return &HttpResponse{
		Status:  status,
		Version: version,
		Headers: hrs,
		Body:    body,
	}
}

func extractEcho(url string) string {
	if len(url) <= 1 {
		return ""
	}

	return url[1:]
}

func hasPrefix(s, prefix string) bool {
	if len(prefix) > len(s) {
		return false
	}
	return s[0:len(prefix)] == prefix
}
