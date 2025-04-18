package capital

import (
	"bytes"
	"capital/models"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

func (c *client) CreateSession(demo bool, apiKey, identifier, password string) (*models.CreateSessionResponse, *models.SessionTokens, error) {
	if apiKey == "" {
		return nil, nil, errors.New("capital.com API key is required")
	}

	if identifier == "" || password == "" {
		return nil, nil, errors.New("capital.com credentials are required")
	}

	payload := map[string]interface{}{
		"identifier":        identifier,
		"password":          password,
		"encryptedPassword": false,
	}

	data, headers, err := c.request("POST", demo, "/session", payload, "", "", apiKey)

	if err != nil {
		return nil, nil, fmt.Errorf("error creating session: %w", err)
	}

	var session models.CreateSessionResponse
	err = json.Unmarshal(data, &session)
	if err != nil {
		return nil, nil, fmt.Errorf("error unmarshalling session: %w", err)
	}

	cst := headers.Get("CST")
	securityToken := headers.Get("X-SECURITY-TOKEN")

	if cst == "" || securityToken == "" {
		return nil, nil, errors.New("failed to obtain session tokens")
	}

	return &session, &models.SessionTokens{
		CST:           cst,
		SecurityToken: securityToken,
		Timestamp:     time.Now(),
	}, nil
}

func (c *client) OpenPosition(demo bool, accountId, direction, epic string, size float64, stopLevel, profitLevel *float64, guaranteedStop bool, cst, securityToken string) (string, error) {
	currentAccount, err := c.GetCurrentAccount(demo, cst, securityToken)
	if err != nil {
		return "", fmt.Errorf("error getting current account: %w", err)
	}

	if currentAccount.AccountId != accountId {
		_, sessionTokens, err := c.SwitchActiveAccount(demo, accountId, cst, securityToken)
		if err != nil {
			return "", fmt.Errorf("error switching account: %w", err)
		}
		cst = sessionTokens.CST
		securityToken = sessionTokens.SecurityToken
	}

	payload := map[string]interface{}{
		"epic":           epic,
		"direction":      direction, // "BUY" or "SELL"
		"size":           size,
		"guaranteedStop": guaranteedStop,
	}

	if stopLevel != nil && *stopLevel != 0.0 {
		payload["stopLevel"] = *stopLevel
	}

	if profitLevel != nil && *profitLevel != 0.0 {
		payload["profitLevel"] = *profitLevel
	}

	data, _, err := c.request("POST", demo, "/positions", payload, cst, securityToken, "")
	if err != nil {
		return "", fmt.Errorf("error opening position: %w", err)
	}

	var response models.CapitalDealReference
	if err := json.Unmarshal(data, &response); err != nil {
		return "", fmt.Errorf("error parsing position response: %w", err)
	}

	// Verify that the position was actually opened
	confirm, err := c.ConfirmDeal(demo, accountId, response.DealReference, cst, securityToken)
	if err != nil {
		return "", fmt.Errorf("error confirming deal: %w", err)
	}

	if confirm.DealStatus != "ACCEPTED" {
		return "", fmt.Errorf("deal was not accepted: %s", confirm.Status)
	}

	if len(confirm.AffectedDeals) == 0 {
		return "", fmt.Errorf("no affected deals found")
	}

	return confirm.AffectedDeals[0].DealID, nil
}

func (c *client) ClosePosition(demo bool, accountId, dealID string, cst, securityToken string) (*models.CapitalDealConfirmation, error) {
	currentAccount, err := c.GetCurrentAccount(demo, cst, securityToken)
	if err != nil {
		return nil, fmt.Errorf("error getting current account: %w", err)
	}

	if currentAccount.AccountId != accountId {
		_, sessionTokens, err := c.SwitchActiveAccount(demo, accountId, cst, securityToken)
		if err != nil {
			return nil, fmt.Errorf("error switching account: %w", err)
		}
		cst = sessionTokens.CST
		securityToken = sessionTokens.SecurityToken
	}

	data, _, err := c.request("GET", demo, "/positions/"+dealID, nil, cst, securityToken, "")
	if err != nil {
		return nil, fmt.Errorf("error getting position details: %w", err)
	}

	var posResponse struct {
		Position models.CapitalPosition `json:"position"`
	}

	if err := json.Unmarshal(data, &posResponse); err != nil {
		return nil, fmt.Errorf("error parsing position details: %w", err)
	}

	closeData, _, err := c.request("DELETE", demo, fmt.Sprintf("/positions/%s", dealID), nil, cst, securityToken, "")
	if err != nil {
		return nil, fmt.Errorf("error closing position: %w", err)
	}

	var response models.CapitalDealReference
	if err := json.Unmarshal(closeData, &response); err != nil {
		return nil, fmt.Errorf("error parsing close position response: %w", err)
	}

	// Verify that the position was actually closed
	confirm, err := c.ConfirmDeal(demo, accountId, response.DealReference, cst, securityToken)
	if err != nil {
		return nil, fmt.Errorf("error confirming position close: %w", err)
	}

	if confirm.DealStatus != "ACCEPTED" {
		return nil, fmt.Errorf("position close not accepted: %s", confirm.Status)
	}

	return confirm, nil
}

func (c *client) ConfirmDeal(demo bool, accountId, dealReference string, cst, securityToken string) (*models.CapitalDealConfirmation, error) {
	currentAccount, err := c.GetCurrentAccount(demo, cst, securityToken)
	if err != nil {
		return nil, fmt.Errorf("error getting current account: %w", err)
	}

	if currentAccount.AccountId != accountId {
		_, sessionTokens, err := c.SwitchActiveAccount(demo, accountId, cst, securityToken)
		if err != nil {
			return nil, fmt.Errorf("error switching account: %w", err)
		}
		cst = sessionTokens.CST
		securityToken = sessionTokens.SecurityToken
	}

	data, _, err := c.request("GET", demo, "/confirms/"+dealReference, nil, cst, securityToken, "")
	if err != nil {
		return nil, fmt.Errorf("error confirming deal: %w", err)
	}

	var confirm models.CapitalDealConfirmation
	if err := json.Unmarshal(data, &confirm); err != nil {
		return nil, fmt.Errorf("error parsing confirm response: %w", err)
	}

	return &confirm, nil
}

func (c *client) GetPositions(demo bool, accountId, cst, securityToken string) (*models.PositionsResponse, error) {
	currentAccount, err := c.GetCurrentAccount(demo, cst, securityToken)
	if err != nil {
		return nil, fmt.Errorf("error getting current account: %w", err)
	}

	if currentAccount.AccountId != accountId {
		_, sessionTokens, err := c.SwitchActiveAccount(demo, accountId, cst, securityToken)
		if err != nil {
			return nil, fmt.Errorf("error switching account: %w", err)
		}
		cst = sessionTokens.CST
		securityToken = sessionTokens.SecurityToken
	}

	data, _, err := c.request("GET", demo, "/positions", nil, cst, securityToken, "")
	if err != nil {
		return nil, fmt.Errorf("error getting positions: %w", err)
	}

	var response models.PositionsResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("error parsing positions response: %w", err)
	}

	return &response, nil
}

func (c *client) GetMarketDetails(demo bool, accountId, epic string, cst, securityToken string) (*models.CapitalMarketDetailsResponse, error) {
	currentAccount, err := c.GetCurrentAccount(demo, cst, securityToken)
	if err != nil {
		return nil, fmt.Errorf("error getting current account: %w", err)
	}

	if currentAccount.AccountId != accountId {
		_, sessionTokens, err := c.SwitchActiveAccount(demo, accountId, cst, securityToken)
		if err != nil {
			return nil, fmt.Errorf("error switching account: %w", err)
		}
		cst = sessionTokens.CST
		securityToken = sessionTokens.SecurityToken
	}

	data, _, err := c.request("GET", demo, "/markets/"+epic, nil, cst, securityToken, "")
	if err != nil {
		return nil, fmt.Errorf("error getting market details: %w", err)
	}

	var response models.CapitalMarketDetailsResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("error parsing market details response: %w", err)
	}

	return &response, nil
}

func (c *client) GetAccounts(demo bool, cst, securityToken string) ([]models.CapitalAccount, error) {
	data, _, err := c.request("GET", demo, "/accounts", nil, cst, securityToken, "")
	if err != nil {
		return nil, fmt.Errorf("error getting accounts: %w", err)
	}

	var response models.CapitalAccountsResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("error parsing accounts response: %w", err)
	}

	return response.Accounts, nil
}

func (c *client) SwitchActiveAccount(demo bool, accountId string, cst, securityToken string) (*models.SwitchAccountResponse, *models.SessionTokens, error) {
	payload := map[string]interface{}{
		"accountId": accountId,
	}

	data, headers, err := c.request("PUT", demo, "/session", payload, cst, securityToken, "")
	if err != nil {
		return nil, nil, fmt.Errorf("error switching account: %w", err)
	}

	var response models.SwitchAccountResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, nil, fmt.Errorf("error parsing switch account response: %w", err)
	}

	return &response, &models.SessionTokens{
		CST:           headers.Get("CST"),
		SecurityToken: headers.Get("X-SECURITY-TOKEN"),
		Timestamp:     time.Now(),
	}, nil
}

func (c *client) GetCurrentAccount(demo bool, cst, securityToken string) (*models.CurrentAccount, error) {
	data, _, err := c.request("GET", demo, "/session", nil, cst, securityToken, "")
	if err != nil {
		return nil, fmt.Errorf("error getting session info: %w", err)
	}

	var response models.CurrentAccount
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("error parsing session response: %w", err)
	}

	return &response, nil
}

func (c *client) request(method string, demo bool, endpoint string, payload interface{}, cst, securityToken, apiKey string) ([]byte, http.Header, error) {
	baseURL := c.baseURL
	if demo {
		baseURL = c.demoBaseURL
	}
	url := fmt.Sprintf("%s%s", baseURL, endpoint)

	var bodyReader io.Reader
	if payload != nil {
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			return nil, nil, fmt.Errorf("error marshaling request payload: %w", err)
		}
		bodyReader = bytes.NewBuffer(payloadBytes)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating API request: %w", err)
	}

	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Use the provided tokens for authentication
	if securityToken != "" {
		req.Header.Set("X-SECURITY-TOKEN", securityToken)
	}
	if cst != "" {
		req.Header.Set("CST", cst)
	}

	if apiKey != "" {
		req.Header.Set("X-CAP-API-KEY", apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("error making API request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("error reading response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return body, resp.Header, nil
}
