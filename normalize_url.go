package main

import (
	"errors"
	"net/url"
	"strings"
)

var (
	ErrInvalidURL = errors.New("invalid url")
)

func NormalizeURL(u string) (string, error) {
	uu, err := url.Parse(u)
	if err != nil {
		return "", errors.Join(ErrInvalidURL, err)
	}

	uuu := strings.ToLower(uu.Host + uu.Path)
	return strings.TrimSuffix(uuu, "/"), nil
}
