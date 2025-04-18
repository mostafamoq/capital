package models

import "time"

type (
	SessionTokens struct {
		CST           string    `json:"cst"`
		SecurityToken string    `json:"securityToken"`
		Timestamp     time.Time `json:"timestamp"`
	}

	CapitalMarket struct {
		Epic                     string  `json:"epic"`
		MarketStatus             string  `json:"marketStatus"`
		NetChange                float64 `json:"netChange"`
		PercentageChange         float64 `json:"percentageChange"`
		Bid                      float64 `json:"bid"`
		Offer                    float64 `json:"offer"`
		UpdateTime               string  `json:"updateTime"`
		DelayTime                int     `json:"delayTime"`
		StreamingPricesAvailable bool    `json:"streamingPricesAvailable"`
	}

	CapitalPositionsResponse struct {
		Positions []CapitalPosition `json:"positions"`
	}

	CapitalPosition struct {
		DealID         string  `json:"dealId"`
		Direction      string  `json:"direction"`
		Size           float64 `json:"size"`
		Epic           string  `json:"epic"`
		OpenLevel      float64 `json:"openLevel"`
		CurrentLevel   float64 `json:"currentLevel"`
		ProfitLoss     float64 `json:"profitLoss"`
		Currency       string  `json:"currency"`
		CreatedDate    string  `json:"createdDate"`
		StopLevel      float64 `json:"stopLevel,omitempty"`
		ProfitLevel    float64 `json:"profitLevel,omitempty"`
		GuaranteedStop bool    `json:"guaranteedStop"`
	}

	CapitalDealReference struct {
		DealReference string `json:"dealReference"`
	}

	CapitalDealConfirmation struct {
		Status        string         `json:"status"`
		DealStatus    string         `json:"dealStatus"`
		DealReference string         `json:"dealReference"`
		AffectedDeals []AffectedDeal `json:"affectedDeals"`
		Reason        string         `json:"reason,omitempty"`
	}

	AffectedDeal struct {
		DealID string `json:"dealId"`
		Status string `json:"status"`
	}

	CapitalMarketDetailsResponse struct {
		DealingRules DealingRules `json:"dealingRules"`
		Instrument   Instrument   `json:"instrument"`
	}

	DealingRules struct {
		MinDealSize             DealSize `json:"minDealSize"`
		MaxDealSize             DealSize `json:"maxDealSize"`
		MinStopOrProfitDistance DealSize `json:"minStopOrProfitDistance"`
	}

	DealSize struct {
		Value float64 `json:"value"`
		Unit  string  `json:"unit"`
	}

	Instrument struct {
		Name                     string  `json:"name"`
		Type                     string  `json:"type"`
		MarketID                 string  `json:"marketId"`
		SpotBid                  float64 `json:"spotBid"`
		SpotAsk                  float64 `json:"spotAsk"`
		MinDealSize              float64 `json:"minDealSize"`
		MaxDealSize              float64 `json:"maxDealSize"`
		OtcTradeable             bool    `json:"otcTradeable"`
		MarketStatus             string  `json:"marketStatus"`
		StreamingPricesAvailable bool    `json:"streamingPricesAvailable"`
	}

	CapitalSessionRequest struct {
		Identifier        string `json:"identifier"`
		Password          string `json:"password"`
		EncryptedPassword bool   `json:"encryptedPassword"`
	}

	CapitalAccountsResponse struct {
		Accounts []CapitalAccount `json:"accounts"`
	}

	CapitalAccount struct {
		AccountID   string  `json:"accountId"`
		AccountName string  `json:"accountName"`
		AccountType string  `json:"accountType"`
		Preferred   bool    `json:"preferred"`
		Balance     Balance `json:"balance"`
		Currency    string  `json:"currency"`
		Symbol      string  `json:"symbol"`
		Status      string  `json:"status"`
	}

	CapitalSessionResponse struct {
		ClientID       string           `json:"clientId"`
		AccountType    string           `json:"accountType"`
		Currency       string           `json:"currency"`
		CurrentAccount CapitalAccount   `json:"currentAccount"`
		Accounts       []CapitalAccount `json:"accounts"`
	}

	Balance struct {
		Balance    float64 `json:"balance"`
		Deposit    float64 `json:"deposit"`
		ProfitLoss float64 `json:"profitLoss"`
		Available  float64 `json:"available"`
	}

	SwitchAccountResponse struct {
		TrailingStopsEnabled  bool `json:"trailingStopsEnabled"`
		DealingEnabled        bool `json:"dealingEnabled"`
		HasActiveDemoAccounts bool `json:"hasActiveDemoAccounts"`
		HasActiveLiveAccounts bool `json:"hasActiveLiveAccounts"`
	}

	CurrentAccount struct {
		ClientId       string `json:"clientId"`
		AccountId      string `json:"accountId"`
		TimezoneOffset int    `json:"timezoneOffset"`
		Locale         string `json:"locale"`
		Currency       string `json:"currency"`
		Symbol         string `json:"symbol"`
		StreamEndpoint string `json:"streamEndpoint"`
	}

	PositionsResponse struct {
		Positions []PositionObj `json:"positions"`
	}

	PositionObj struct {
		Position Position `json:"position"`
		Market   Market   `json:"market"`
	}

	Position struct {
		ContractSize   int     `json:"contractSize"`
		CreatedDate    string  `json:"createdDate"`
		CreatedDateUTC string  `json:"createdDateUTC"`
		DealId         string  `json:"dealId"`
		DealReference  string  `json:"dealReference"`
		WorkingOrderId string  `json:"workingOrderId"`
		Size           float64 `json:"size"`
		Leverage       int     `json:"leverage"`
		Upl            float64 `json:"upl"`
		Direction      string  `json:"direction"`
		Level          float64 `json:"level"`
		Currency       string  `json:"currency"`
		GuaranteedStop bool    `json:"guaranteedStop"`
	}

	Market struct {
		InstrumentName           string  `json:"instrumentName"`
		Expiry                   string  `json:"expiry"`
		MarketStatus             string  `json:"marketStatus"`
		Epic                     string  `json:"epic"`
		InstrumentType           string  `json:"instrumentType"`
		LotSize                  int     `json:"lotSize"`
		High                     float64 `json:"high"`
		Low                      float64 `json:"low"`
		PercentageChange         float64 `json:"percentageChange"`
		NetChange                float64 `json:"netChange"`
		Bid                      float64 `json:"bid"`
		Offer                    float64 `json:"offer"`
		UpdateTime               string  `json:"updateTime"`
		UpdateTimeUTC            string  `json:"updateTimeUTC"`
		DelayTime                int     `json:"delayTime"`
		StreamingPricesAvailable bool    `json:"streamingPricesAvailable"`
		ScalingFactor            int     `json:"scalingFactor"`
	}

	CapitalErrorResponse struct {
		ErrorCode string `json:"errorCode"`
	}

	CreateSessionResponse struct {
		AccountType string `json:"accountType"`
		AccountInfo struct {
			Balance    float64 `json:"balance"`
			Deposit    float64 `json:"deposit"`
			ProfitLoss float64 `json:"profitLoss"`
			Available  float64 `json:"available"`
		} `json:"accountInfo"`
		CurrencyIsoCode  string `json:"currencyIsoCode"`
		CurrencySymbol   string `json:"currencySymbol"`
		CurrentAccountId string `json:"currentAccountId"`
		StreamingHost    string `json:"streamingHost"`
		Accounts         []struct {
			AccountId   string `json:"accountId"`
			AccountName string `json:"accountName"`
			Preferred   bool   `json:"preferred"`
			AccountType string `json:"accountType"`
		} `json:"accounts"`
		ClientId              string `json:"clientId"`
		TimezoneOffset        int    `json:"timezoneOffset"`
		HasActiveDemoAccounts bool   `json:"hasActiveDemoAccounts"`
		HasActiveLiveAccounts bool   `json:"hasActiveLiveAccounts"`
		TrailingStopsEnabled  bool   `json:"trailingStopsEnabled"`
	}
)
