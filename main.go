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
)

type Request struct {
	method   string
	path     string
	protocol string
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

		// No 'go' keyword here!
		// We want THIS worker to be busy doing this task.
		handleConnection(conn, root)

		fmt.Printf("Worker %d finished. Ready for next.\n", id)
	}

}

func writeResponse(c net.Conn, body []byte, status_code string, contentType string) {
	responseLine := []byte(fmt.Sprintf("HTTP/1.1 %s\r\n", status_code))
	contentTypeHeader := []byte(fmt.Sprintf("Content-Type: %s\r\n\r\n", contentType))
	var buf bytes.Buffer
	buf.Write(responseLine)
	buf.Write(contentTypeHeader)
	buf.Write(body)
	response := buf.Bytes()

	c.Write(response)
}

func parseRequestLine(s string) (Request, error) {
	parsed_string := strings.Split(s, " ")

	if len(parsed_string) != 3 {
		fmt.Println("malformed request")
		return Request{}, fmt.Errorf("malformed request: %s ", s)
	}
	request_data := Request{
		method:   parsed_string[0],
		path:     parsed_string[1],
		protocol: parsed_string[2],
	}
	return request_data, nil
}

func handleConnection(c net.Conn, root *string) {
	defer c.Close()
	buff := make([]byte, 1024)
	n, err := c.Read(buff)
	if err != nil {
		fmt.Println("error in reading")
		return
	}

	header := string(buff[:n])

	firstLine := strings.Split(header, "\n")[0]
	requestData, err := parseRequestLine(firstLine)
	if err != nil {
		fmt.Println("error in parseRequestLine fucntion")
		return
	}

	cleanedPath := filepath.Clean(requestData.path)

	// 2. Remove the leading "/" to make it relative to your project folder
	// "/secret.pem" becomes "secret.pem"
	fileName := strings.TrimPrefix(cleanedPath, "/")

	if fileName == "" {
		fileName = "index.html"
	}

	filePath := filepath.Join(*root, fileName)

	fmt.Printf("Requested: %s -> Serving: %s\n", requestData.path, filePath)
	htmlData, err := os.ReadFile(filePath)
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

	if err != nil {
		fmt.Printf("failed to read the %s file", filePath)
		htmlData, err := os.ReadFile("./404.html")
		if err != nil {
			fmt.Println("failed to read the 404.html file")
			return
		}
		writeResponse(c, htmlData, "404 Not Found", "text/html")
		return
	}
	writeResponse(c, htmlData, "200 OK", contentType)
}
