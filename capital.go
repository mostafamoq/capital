package capital

import (
	"capital/models"
	"net/http"
)

type Client interface {
	CreateSession(demo bool, apiKey, identifier, password string) (*models.CreateSessionResponse, *models.SessionTokens, error)
	OpenPosition(demo bool, accountId, direction, epic string, size float64, stopLevel, profitLevel *float64, guaranteedStop bool, cst, securityToken string) (string, error)
	ClosePosition(demo bool, accountId, dealID string, cst, securityToken string) (*models.CapitalDealConfirmation, error)
	ConfirmDeal(demo bool, accountId, dealReference string, cst, securityToken string) (*models.CapitalDealConfirmation, error)
	GetPositions(demo bool, accountId, cst, securityToken string) (*models.PositionsResponse, error)
	GetMarketDetails(demo bool, accountId, epic string, cst, securityToken string) (*models.CapitalMarketDetailsResponse, error)
	GetAccounts(demo bool, cst, securityToken string) ([]models.CapitalAccount, error)
	SwitchActiveAccount(demo bool, accountId string, cst, securityToken string) (*models.SwitchAccountResponse, *models.SessionTokens, error)
	GetCurrentAccount(demo bool, cst, securityToken string) (*models.CurrentAccount, error)
}

type client struct {
	httpClient  *http.Client
	baseURL     string
	demoBaseURL string
}

func New(baseUrl, demoBaseUrl string) Client {
	return &client{
		httpClient:  &http.Client{},
		baseURL:     baseUrl,
		demoBaseURL: demoBaseUrl,
	}
}
