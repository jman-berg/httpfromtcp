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
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	rawBytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("Error reading request %s", err.Error())
	}

	requestLine, err := parseRequestLine(rawBytes)
	if err != nil {
		return nil, fmt.Errorf("Error parsing request: %s", err.Error())
	}

	return &Request{
		RequestLine: *requestLine,
	}, nil
}

func parseRequestLine(data []byte) (*RequestLine, error) {
	idx := bytes.Index(data, []byte("\r\n"))
	if idx == -1 {
		return nil, fmt.Errorf("Couldn't find CRLF in request-line")
	}
	requestLineText := string(data[:idx])

	requestLine, err := requestLineFromString(requestLineText)
	if err != nil {
		return nil, err
	}

	return requestLine, nil

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
