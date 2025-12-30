package main

import (
	"bytes"
	"fmt"
	"net"
	"strings"
)

type Request struct {
	method      string
	path        string
	httpVersion string
	headers     *map[string]string
}
type Response struct {
	statusCode int
	status     string
	headers    map[string]string
	body       []byte
}

func writeResponse(c net.Conn, response Response) (Response, error) {
	responseLine := []byte(fmt.Sprintf("HTTP/1.1 %s \r\n", response.status))
	var buf bytes.Buffer
	buf.Write(responseLine)

	// adding the headers one-by-one
	for k, v := range response.headers {
		header := fmt.Sprintf("%s: %s\r\n", k, v)
		buf.Write([]byte(header))
	}
	buf.Write([]byte("\r\n")) // there shoud be two `\r\n` after the last header and body

	buf.Write(response.body)
	parsedResponse := buf.Bytes()

	_, err := c.Write(parsedResponse)
	if err != nil {
		return Response{}, err
	}
	return response, nil
}

func parseRequestLine(s string) (Request, error) {
	parsed_string := strings.Split(s, " ")

	if len(parsed_string) != 3 {
		fmt.Println("malformed request")
		return Request{}, fmt.Errorf("malformed request: %s ", s)
	}
	request_data := Request{
		method:      parsed_string[0],
		path:        parsed_string[1],
		httpVersion: parsed_string[2],
	}
	return request_data, nil
}
