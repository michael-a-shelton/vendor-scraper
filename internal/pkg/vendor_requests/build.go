package vendor_requests

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

//go:embed config/vendor-data.json
var vendorDataFile []byte

//go:embed config/cookies
var cookies embed.FS

type VendorData []struct {
	AddToCartURL   string            `json:"addToCartURL"`
	CheckoutURL    string            `json:"checkoutURL"`
	Company        string            `json:"company"`
	Verb           string            `json:"verb"`
	Origin         string            `json:"origin"`
	CookieFilePath string            `json:"cookieFilePath"`
	Delay          string            `json:"delay"`
	Payload        string            `json:"payload,omitempty"`
	Headers        map[string]string `json:"headers"`
}

func BuildRequestFromConfig() (*[]http.Request, error) {
	var data VendorData

	err := readDataFromString(&data)
	if err != nil {
		return nil, fmt.Errorf("failed to read data from file: %w", err)
	}

	var requests []http.Request
	for i := 0; i < len(data); i++ {
		parsedURL, _ := url.ParseRequestURI(data[i].AddToCartURL)

		req, err := http.NewRequest(data[i].Verb, parsedURL.String(), bytes.NewBuffer([]byte(data[i].Payload)))
		if err != nil {
			return nil, err
		}

		req.ContentLength = int64(len(data[i].Payload))
		req.Header = make(http.Header)
		req.Host = data[i].Origin

		cookieBytes, err := cookies.ReadFile(data[i].CookieFilePath)
		if err != nil {
			return nil, fmt.Errorf("failed to get cookie: %w", err)
		}

		req.Header.Set("cookie", string(cookieBytes))

		for header, value := range data[i].Headers {
			req.Header.Set(header, value)
		}

		requests = append(requests, *req)

	}

	return &requests, nil
}

func readDataFromString(data *VendorData) error {
	err := json.Unmarshal(vendorDataFile, &data)
	if err != nil {
		return fmt.Errorf("error unmarshaling data: %w", err)
	}

	return nil
}
