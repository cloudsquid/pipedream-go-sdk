package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
)

type GetCurrentUserResponse struct {
	Data UserData `json:"data"`
}

type UserData struct {
	ID                    string `json:"id"`
	Username              string `json:"username"`
	Email                 string `json:"email"`
	DailyComputeTimeQuota int64  `json:"daily_compute_time_quota,omitempty"`
	DailyComputeTimeUsed  int64  `json:"daily_compute_time_used,omitempty"`
	DailyInvocationsQuota int64  `json:"daily_invocations_quota,omitempty"`
	DailyInvocationsUsed  int64  `json:"daily_invocations_used,omitempty"`
	Orgs                  []Org  `json:"orgs,omitempty"`

	// For paid users
	BillingPeriodStartTS int64 `json:"billing_period_start_ts,omitempty"`
	BillingPeriodEndTS   int64 `json:"billing_period_end_ts,omitempty"`
	BillingPeriodCredits int64 `json:"billing_period_credits,omitempty"`
}

type Org struct {
	Name                  string `json:"name"`
	ID                    string `json:"id"`
	OrgName               string `json:"orgname"`
	Email                 string `json:"email"`
	DailyCreditsQuota     int64  `json:"daily_credits_quota,omitempty"`
	DailyCreditsUsed      int64  `json:"daily_credits_used,omitempty"`
	DailyComputeTimeQuota int64  `json:"daily_compute_time_quota,omitempty"`
	DailyComputeTimeUsed  int64  `json:"daily_compute_time_used,omitempty"`
	DailyInvocationsQuota int64  `json:"daily_invocations_quota,omitempty"`
	DailyInvocationsUsed  int64  `json:"daily_invocations_used,omitempty"`
}

// GetCurrentUser Retrieves information on the authenticated user
func (c *Client) GetCurrentUser(
	ctx context.Context,
) (*GetCurrentUserResponse, error) {
	baseURL := c.RestURL().ResolveReference(&url.URL{
		Path: path.Join(c.RestURL().Path, "users", "me")})

	endpoint := baseURL.String()

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("creating get current user request: %w", err)
	}

	response, err := c.doRequestViaApiKey(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("executing request to get current user: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(response.Body)
		return nil,
			fmt.Errorf("unexpected status code %d: %s", response.StatusCode, string(bodyBytes))
	}

	var userInfo GetCurrentUserResponse
	err = json.NewDecoder(response.Body).Decode(&userInfo)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling response for request to get current user: %w", err)
	}

	return &userInfo, nil
}
