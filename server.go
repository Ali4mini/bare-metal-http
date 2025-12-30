package main

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"
)

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
