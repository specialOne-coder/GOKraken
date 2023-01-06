package main

import "github.com/lib/pq"

type ServerStatus struct {
	Result struct {
		Status    string `json:"status"`
		Timestamp string `json:"timestamp"`
	} `json:"result"`
}

type AssetPair struct {
	Result map[string]AssetPairInfo `json:"result"`
}

type Fees []float64

type AssetPairInfo struct {
	Altname           string `json:"altname"`
	Wsname            string `json:"wsname"`
	AclassBase        string `json:"aclass_base"`
	Base              string `json:"base"`
	AclassQuote       string `json:"aclass_quote"`
	Quote             string `json:"quote"`
	Lot               string `json:"lot"`
	CostDecimals      int    `json:"cost_decimals"`
	PairDecimals      int    `json:"pair_decimals"`
	LotDecimals       int    `json:"lot_decimals"`
	LotMultiplier     int    `json:"lot_multiplier"`
	LeverageBuy       pq.Int64Array  `json:"leverage_buy"`
	LeverageSell      pq.Int64Array  `json:"leverage_sell"`
	MarginCall        int    `json:"margin_call"`
	MarginStop        int    `json:"margin_stop"`
	Ordermin          string `json:"ordermin"`
	Costmin           string `json:"costmin"`
	Ticksize          string `json:"tick_size"`
	Status            string `json:"status"`
}

type AssetPrice struct {
	Result map[string]AssetPriceInfo `json:"result"`
}

type AssetPriceInfo struct {
	PRIX_ACHAT     []string `json:"a"`
	PRIX_VENTE     []string `json:"b"`
	PRICE          []string `json:"c"`
	VOLUME         []string `json:"v"`
	VWAP           []string `json:"p"`
	VOLUME_MOYENNE []int    `json:"t"`
	LOW            []string `json:"l"`
	HIGH           []string `json:"h"`
	OPENING        string   `json:"o"`
}

type Config struct {
	APIURL string `env:"API_URL"`
}

type Aggregation struct {
	ASSETS AssetPairInfo
	PRIX_ACHAT     string
	PRIX_VENTE     string
	PRICE          string
	VOLUME         string
	VWAP           string
	VOLUME_MOYENNE int
	LOW            string
	HIGH           string
	OPENING        string
	TIMESTAMP    string
}
