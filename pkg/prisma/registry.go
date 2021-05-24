package prisma

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

// ScanImage forces Prisma to re-scan a specific image
func (c *Client) ScanImage(ctx context.Context, registry string, imageName string, tag string) (bool, error) {
	type tagRequest struct {
		Registry   string `json:"registry"`
		Repository string `json:"repo"`
		Tag        string `json:"tag"`
		Digest     string `json:"digest"`
	}

	type imageScanRequest struct {
		Tag tagRequest `json:"tag"`
	}

	data := imageScanRequest{
		Tag: tagRequest{
			Registry:   registry,
			Tag:        tag,
			Repository: imageName,
		},
	}

	body, err := json.Marshal(data)
	if err != nil {
		return false, errors.Wrap(err, "Failed to marshall object into json")
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/registry/scan", c.BaseURL), bytes.NewBuffer(body))
	if err != nil {
		return false, errors.Wrap(err, "Failed to create http request")
	}

	req = req.WithContext(ctx)

	if err := c.sendRequest(req, nil); err != nil {
		return false, errors.Wrap(err, "Failed to send http request")
	}

	return true, nil
}

// Get retrieves the scan reports for a registry. If no report exists, the response is empty
func (c *Client) Get(ctx context.Context, registry string) (*[]ScanReport, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/registry?name=%s&compact=true", c.BaseURL, registry), nil)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create http request")
	}

	req = req.WithContext(ctx)

	res := []ScanReport{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, errors.Wrap(err, "Failed to send http request")
	}

	return &res, nil
}
