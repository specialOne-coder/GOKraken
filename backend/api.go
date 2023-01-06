package main

import (
	"encoding/json"
	"fmt"

	//"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	//"time"
)

// Get the status of the Kraken server
func GetStatus(wg *sync.WaitGroup) (servers ServerStatus) {

	// Requête HTTP vers l'API de Kraken
	res, err := http.Get(KrakenAPI + "SystemStatus")
	checkErr(err)
	defer res.Body.Close()

	// Décodage de la réponse JSON
	var serverStatus ServerStatus
	err = json.NewDecoder(res.Body).Decode(&serverStatus)
	checkErr(err)
	//fmt.Println(serverStatus)
	wg.Done()
	return serverStatus
}

// Get the asset pairs of the Kraken server
func GetKrakenAssetPairs(wg *sync.WaitGroup) (assets AssetPair) {

	res, err := http.Get(KrakenAPI + "AssetPairs?pair=XXBTZUSD,XETHXXBT")
	checkErr(err)
	defer res.Body.Close()

	// // Décodage de la réponse JSON
	var assetPairs AssetPair
	err = json.NewDecoder(res.Body).Decode(&assetPairs)
	// get this pairs prices to aggregate with the asset pairs

	checkErr(err)
	//fmt.Println(assetPairs)
	wg.Done()
	return assetPairs
}

// Get the price of the asset pairs of the Kraken server
func GetKrakenAssetsPrice(assetName string) (AssetPriceInfo, error) {
	fmt.Println("Name", assetName)
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
