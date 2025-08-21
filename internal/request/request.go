package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"
)

type Request struct {
	RequestLine RequestLine
	State       requestState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type requestState int

const (
	requestStateInitialized requestState = iota
	requestStateDone
)

const bufferSize = 8

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize)

	readToIndex := 0

	request := &Request{
		State: requestStateInitialized,
	}

	for request.State != requestStateDone {
		if len(buf) <= readToIndex {
			oldBuf := buf
			buf = make([]byte, len(buf)*2)
			copy(buf, oldBuf)
		}

		numBytesRead, err := reader.Read(buf[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				request.State = requestStateDone
				break
			}
			return nil, err
		}
		readToIndex += numBytesRead
		numBytesParsed, err := request.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}
		copy(buf, buf[numBytesParsed:])
		readToIndex -= numBytesParsed

	}
	return request, nil
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	idx := bytes.Index(data, []byte("\r\n"))
	if idx == -1 {
		return nil, 0, nil
	}
	requestLineText := string(data[:idx])

	requestLine, err := requestLineFromString(requestLineText)
	if err != nil {
		return nil, 0, err
	}

	return requestLine, idx + 2, nil

}

func requestLineFromString(str string) (*RequestLine, error) {

	parts := strings.Split(str, " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("Poorly formatted request-line: %s", str)
	}

	method := parts[0]
	for _, r := range method {
		if !(unicode.IsUpper(r) && unicode.IsLetter(r)) {
			return nil, errors.New("Invalid method, should be capital letters")
		}
	}

	versionParts := strings.Split(parts[2], "/")

	httpPart := versionParts[0]
	if httpPart != "HTTP" {
		return nil, fmt.Errorf("Poorly formatted http part in version: %s", httpPart)
	}

	httpVersion := versionParts[1]
	if !(httpVersion == "1.1") {
		return nil, fmt.Errorf("Http-version should be 1.1, found: %s", httpVersion)
	}

	requestTarget := parts[1]
	if strings.Contains(requestTarget, " ") {
		return nil, fmt.Errorf("Spaces are not allowed in request-target")
	}

	return &RequestLine{
		HttpVersion:   httpVersion,
		RequestTarget: requestTarget,
		Method:        method,
	}, nil

}

func (r *Request) parse(data []byte) (int, error) {
	switch r.State {
	case requestStateDone:
		return 0, errors.New("Error: trying to read data in a done state")
	case requestStateInitialized:
		requestLine, processedBytes, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}

		if processedBytes == 0 {
			return 0, nil
		}

		r.RequestLine = *requestLine
		r.State = 1
		return processedBytes, nil
	default:
		return 0, errors.New("Error: unknown state")
	}
}
