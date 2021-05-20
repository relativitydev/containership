package prisma

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

// GetImage retrieves the scan reports for a registry. If no report exists, the response is empty
// Note: This is VERY similar to GET /registry. Each endpoint returns the same data object, but there is more or less detail for som properties.
func (c *Client) GetImage(ctx context.Context, registry string) (*[]ScanReport, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/images?name=%s&compact=true", c.BaseURL, registry), nil)
	if err != nil {
		return nil, errors.Wrap(err, "Error creating http request for prisma")
	}

	req = req.WithContext(ctx)

	res := []ScanReport{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}
