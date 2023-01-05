// Créer une application qui requête l’API Kraken pour récupérer:
// le statut et timing du serveur,
// les paires de trading,
// les informations relatives à chaque paire de trading

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"sync"

	"github.com/lib/pq"

	//"time"
	"net/http"
	"os"
)

const (
	KrakenAPI = "https://api.kraken.com/0/public/"
)

const (
	host     = "localhost"
	port     = "5432"
	user     = "exam"
	password = "mysecretpassword"
	dbname   = "kraken"
	schema   = "public"
)

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
	LeverageBuy       []int  `json:"leverage_buy"`
	LeverageSell      []int  `json:"leverage_sell"`
	Fees              []Fees `json:"fees"`
	FeesMaker         []Fees `json:"fees_maker"`
	FeeVolumeCurrency string `json:"fee_volume_currency"`
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

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func GetStatus() (ServerStatus, error) {
	// Requête HTTP vers l'API de Kraken
	res, err := http.Get(KrakenAPI + "SystemStatus")
	checkErr(err)
	defer res.Body.Close()

	// Décodage de la réponse JSON
	var serverStatus ServerStatus
	err = json.NewDecoder(res.Body).Decode(&serverStatus)
	checkErr(err)
	fmt.Println(serverStatus)
	return serverStatus, nil
}

func GetKrakenAssetPairs() (AssetPair, error) {
	res, err := http.Get(KrakenAPI + "AssetPairs?pair=XETHXXBT")
	checkErr(err)
	defer res.Body.Close()

	// // Décodage de la réponse JSON
	var assetPairs AssetPair
	err = json.NewDecoder(res.Body).Decode(&assetPairs)
	// get this pairs prices to aggregate with the asset pairs

	checkErr(err)
	//fmt.Println(assetPairs)
	return assetPairs, nil
}

func GetKrakenAssetsPrice(assetName string) (AssetPriceInfo, error) {
	res, err := http.Get(KrakenAPI + "Ticker?pair=" + assetName)
	checkErr(err)
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	checkErr(err)

	// // Décodage de la réponse JSON
	var assetPrice AssetPrice
	err = json.Unmarshal(body, &assetPrice)
	checkErr(err)
	// get this pairs prices to aggregate with the asset pairs
	assetPriceInfo := assetPrice.Result[assetName]
	return assetPriceInfo, nil
}

func InsertIntoFile(assets AssetPair) {
	dir := "Archive"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.Mkdir("Archive", 0755)
		checkErr(err)
	}
	file, err := os.Create("Archive/assets.json")
	checkErr(err)
	defer file.Close()
	// Encodage de l'objet en JSON
	encoder := json.NewEncoder(file)
	if err := encoder.Encode(assets); err != nil {
		panic(err)
	}
}

func ConnectionToDB() (dbase *sql.DB, err error) {
	connectionString := "host=" + host + " port=" + port + " user=" + user + " password=" + password + " dbname=" + dbname + " sslmode=disable"
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	checkErr(err)
	fmt.Println("Successfully connected!")
	return db, nil
}

func InserPairIntoDB(dbase *sql.DB, assets AssetPair, timestamp string) {
	sqlStat := "CREATE TABLE IF NOT EXISTS public.pairs( id SERIAL NOT NULL, altname character varying NOT NULL, wsname character varying, aclassbase character varying,base character varying, aclassquote character varying, quote character varying, lot character varying,costdecimals integer, pairdecimals integer, lotdecimals integer,lotmultiplier integer, leveragebuy integer[],leveragesell integer[], fees real[], feesmaker real[], feevolumecurrency character varying,margincall integer, marginstop integer, ordermin character varying, costmin character varying, ticksize character varying, status character varying, timestamp character varying,prix_achat character varying,prix_vente character varying,prix character varying,volume character varying, volume24h character varying, volume_moyenne character varying,low character varying, high character varying, opening character varying, PRIMARY KEY (id) ); ALTER TABLE IF EXISTS public.pairs OWNER to exam;"
	_, errors := dbase.Exec(sqlStat)
	checkErr(errors)
	for i, asset := range assets.Result {
		// Insertion des données dans la base de données
		arr := pq.Array(asset.LeverageBuy)
		fmt.Println(arr)
		fmt.Println("asset", i)
		assetPrice, err := GetKrakenAssetsPrice(i)
		checkErr(err)

		fmt.Println(assetPrice.PRICE[0])
		_, err = dbase.Exec("INSERT INTO pairs (altname, wsname, aclassbase, base, aclassquote, quote, lot, costdecimals, pairdecimals,lotdecimals, lotmultiplier, leveragebuy, leveragesell, fees, feesmaker, feevolumecurrency,margincall, marginstop, ordermin, costmin, ticksize, status,timestamp,prix_achat,prix_vente,prix ,volume,volume24h ,volume_moyenne,low, high, opening)VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22,$23,$24,$25,$26,$27,$28,$29,$30,$31,$32)", asset.Altname, asset.Wsname, asset.AclassBase, asset.Base, asset.AclassQuote, asset.Quote, asset.Lot,
			asset.CostDecimals, asset.PairDecimals, asset.LotDecimals, asset.LotMultiplier,
			pq.Array(asset.LeverageBuy), pq.Array(asset.LeverageSell), pq.Array(asset.Fees), pq.Array(asset.FeesMaker),
			asset.FeeVolumeCurrency,
			asset.MarginCall, asset.MarginStop, asset.Ordermin, asset.Costmin, asset.Ticksize, asset.Status, timestamp,
			assetPrice.PRIX_ACHAT[0], assetPrice.PRIX_VENTE[0], assetPrice.PRICE[0], assetPrice.VOLUME[0], assetPrice.VWAP[1],
			assetPrice.VOLUME_MOYENNE[0], assetPrice.LOW[0], assetPrice.HIGH[0],
			assetPrice.OPENING)
		checkErr(err)
	}
}

func downloadFile(w http.ResponseWriter, r *http.Request) {
	// Ouvrez le fichier en lecture
	file, err := os.Open("Archive/assets.json")
	if err != nil {
		http.Error(w, "Impossible d'ouvrir le fichier", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Définissez l'en-tête du fichier
	w.Header().Set("Content-Disposition", "attachment; filename=monfichier.txt")
	w.Header().Set("Content-Type", "text/plain")

	// Copiez le contenu du fichier dans la réponse
	io.Copy(w, file)
}

// Get status
// Get Kraken datas
// Insert into file
// Connection to DB
// Insert into DB
func main() {
	
	// var wg sync.WaitGroup
	// wg.Add(3)
	
	serverStatus, err := GetStatus()
	checkErr(err)
	assetsPairs, err := GetKrakenAssetPairs()
	checkErr(err)
	//assetsPrices, err := GetKrakenAssetsPrice()
	checkErr(err)
	InsertIntoFile(assetsPairs)
	db, err := ConnectionToDB()

	checkErr(err)
	InserPairIntoDB(db, assetsPairs, serverStatus.Result.Timestamp)

	http.HandleFunc("/download", downloadFile)
    http.ListenAndServe(":8080", nil)
	//InserPriceIntoDB(db, assetsPrices)

	// Ecrire ces datas dans la bdd
	// connexion a la bdd

}
