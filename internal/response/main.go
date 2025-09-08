package response

import (
	"fmt"
	"io"
	"strconv"

	"github.com/jman-berg/httpfromtcp/internal/headers"
)

type Writer struct {
}

type StatusCode int

const (
	StatusOK                  StatusCode = 200
	StatusBadRequest          StatusCode = 400
	StatusInternalServerError StatusCode = 500
)

var statusMessages = map[StatusCode]string{
	StatusOK:                  "HTTP/1.1 200 OK\r\n",
	StatusBadRequest:          "HTTP/1.1 400 Bad Request\r\n",
	StatusInternalServerError: "HTTP/1.1 500 Internal Server Error\r\n",
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	msg, ok := statusMessages[statusCode]
	if !ok {
		return nil
	}
	_, err := w.Write([]byte(msg))
	if err != nil {
		return err
	}
	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	contentLenString := strconv.Itoa(contentLen)
	headers := headers.Headers{}
	headers.Set("Content-Type", "text/plain\r\n")
	headers.Set("Connection", "closed\r\n")
	headers.Set("Content-Length", contentLenString+"\r\n")
	return headers
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	for k, v := range headers {
		_, err := w.Write([]byte(fmt.Sprintf("%s: %s\r\n", k, v)))
		if err != nil {
			return err
		}
		_, err = w.Write([]byte("\r\n"))
		return err
	}
	return nil
}
