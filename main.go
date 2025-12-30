package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"
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

func main() {
	jobs := make(chan net.Conn, 10)

	port := flag.Int("port", 8080, "port number")
	root := flag.String("root", "./", "root of the server")

	flag.Parse()

	address := fmt.Sprintf("127.0.0.1:%d", *port)
	server, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("address already taken")
	}
	fmt.Printf("listening on: %v", server.Addr())
	for i := 0; i < 5; i++ {
		go worker(i, jobs, root)
	}

	for {
		conn, err := server.Accept()
		if err != nil {
			fmt.Println("error in Accept")
			// os.Exit(1)
		}

		jobs <- conn
	}

}

func worker(id int, jobs chan net.Conn, root *string) {
	for conn := range jobs {

		fmt.Printf("Worker %d processing connection\n", id)

		startTime := time.Now()

		status, err := handleConnection(conn, root)
		if err != nil {
			fmt.Printf("error in handleConnection: %s", err)
		}
		duration := time.Since(startTime)

		fmt.Printf("[Worker %d] %d | %s | %v\n", id, status, conn.RemoteAddr(), duration)

	}

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

func handleConnection(c net.Conn, root *string) (int, error) {
	defer c.Close()
	buff := make([]byte, 1024)
	n, err := c.Read(buff)
	if err != nil {
		fmt.Println("error in reading")
		return 0, err
	}

	header := string(buff[:n])

	firstLine := strings.Split(header, "\n")[0]
	requestData, err := parseRequestLine(firstLine)
	if err != nil {
		fmt.Println("error in parseRequestLine fucntion")
		return 0, err
	}

	cleanedPath := filepath.Clean(requestData.path)

	fileName := strings.TrimPrefix(cleanedPath, "/")

	if fileName == "" {
		fileName = "index.html"
	}

	filePath := filepath.Join(*root, fileName)

	fmt.Printf("Requested: %s -> Serving: %s\n", requestData.path, filePath)
	htmlData, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("failed to read the %s file", filePath)
		notFoundPath := filepath.Join(*root, "404.html")
		htmlData, err := os.ReadFile(notFoundPath)
		if err != nil {
			fmt.Println("failed to read the 404.html file")
			return 0, err
		}
		response := Response{
			statusCode: 404,
			status:     "404 Not Found",
			headers:    map[string]string{"Content-Type": "text/html"},
			body:       htmlData,
		}
		res, err := writeResponse(c, response)
		if err != nil {
			fmt.Printf("error in writeResponse function: %s", err)
			return 0, err
		}
		return res.statusCode, nil
	}

	fullFileName := strings.Split(filePath, ".")
	fileExtension := fullFileName[len(fullFileName)-1]
	var contentType string
	switch fileExtension {
	case "html":
		contentType = "text/html"
	case "css":
		contentType = "text/css"
	case "png":
		contentType = "image/png"
	default:
		contentType = "text/plain"

	}

	response := Response{
		statusCode: 200,
		status:     "200 OK",
		headers:    map[string]string{"Content-Type": contentType, "X-Some": "i will die for you"},
		body:       htmlData,
	}
	res, err := writeResponse(c, response)
	if err != nil {
		fmt.Printf("error in writeResponse function: %s", err)
		return 0, err
	}
	return res.statusCode, nil
}
