// services/apiclient/apiclient.go
package apiclient

import (
    "context"
    "fmt"
    "io"
    "net/http"
    "time"
    "backend_rental/services/ratelimiter"
)

type APIClient struct {
    client      *http.Client
    rateLimiter *ratelimiter.APIRateLimiter
    rapidAPIKey string
}

func NewAPIClient(rapidAPIKey string) *APIClient {
    return &APIClient{
        client: &http.Client{
            Timeout: 10 * time.Second,
        },
        rateLimiter: ratelimiter.GetInstance(),
        rapidAPIKey: rapidAPIKey,
    }
}

func (c *APIClient) MakeRequest(ctx context.Context, url string) ([]byte, error) {
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return nil, fmt.Errorf("error creating request: %v", err)
    }

    req.Header.Add("x-rapidapi-host", "booking-com18.p.rapidapi.com")
    req.Header.Add("x-rapidapi-key", c.rapidAPIKey)

    // Wait for rate limiter
    if err := c.rateLimiter.Wait(ctx); err != nil {
        return nil, fmt.Errorf("rate limiter error: %v", err)
    }

    resp, err := c.client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("error sending request: %v", err)
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("error reading response body: %v", err)
    }

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("API request failed with status code: %d, body: %s",
            resp.StatusCode, string(body))
    }

    return body, nil
}