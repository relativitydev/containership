package prisma

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/pkg/errors"
	"github.com/relativitydev/containership/pkg/utils"
)

// Client is functionality for interacting with Prisma
type Client struct {
	BaseURL  string
	Username string
	Password string
	// Note: we are using the default retry values defined in the package.
	*retryablehttp.Client
}

// NewClient creates a new Client
func NewClient(registryName string) (*Client, error) {
	config, err := utils.GetRegistryConfig(utils.Before(registryName, "."))
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve registry config: "+registryName, err)
	}

	client := &Client{
		BaseURL:  "https://" + config.PrismaURL + ":8083/api/v1",
		Password: config.PrismaPassword,
		Username: config.PrismaUsername,
		Client:   retryablehttp.NewClient(),
	}

	client.Client.HTTPClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, //nolint
		},
	}

	return client, nil
}

// SendRequest sends http request and uses basic auth if values are set on the client
func (c *Client) sendRequest(req *http.Request, v interface{}) error {
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json")

	if c.Username != "" && c.Password != "" {
		req.SetBasicAuth(c.Username, c.Password)
	}

	retryReq, err := retryablehttp.FromRequest(req)
	if err != nil {
		return errors.Wrap(err, "Error creating http request")
	}

	res, err := c.Do(retryReq)
	if err != nil {
		return errors.Wrap(err, "Error sending http request")
	}

	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		var errRes errorResponse
		if err = json.NewDecoder(res.Body).Decode(&errRes); err == nil {
			return errors.New(errRes.Message)
		}

		return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}

	if v != nil {
		if err = json.NewDecoder(res.Body).Decode(&v); err != nil {
			return errors.Wrap(err, "Error unmarshalling the response")
		}
	}

	return nil
}

type errorResponse struct {
	Message string `json:"err"`
}
