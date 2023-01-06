package main

import (
	//"net/http"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Get datas from Kraken every hour
// send datas into channel and get it
// Insert datas into file and DB
func getDataEveryHourAndInsertIntoFileAndDB(wg sync.WaitGroup) {
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()
	serverStatus := make(chan ServerStatus)
	assets := make(chan AssetPair)

	for range ticker.C {
		wg.Add(2)
		go func() {
			serverStatus <- GetStatus(&wg)
			assets <- GetKrakenAssetPairs(&wg)
		}()
		result := <-serverStatus
		assetsPairs := <-assets
		fmt.Println("Get datas from Kraken and insert into file and DB")
		db, err := ConnectionToDB()
		checkErr(err)
		InserPairIntoDBAndFile(db, assetsPairs, result.Result.Timestamp)
	}

}

func main() {

	var wg sync.WaitGroup
	go func() {
		fmt.Println("Server is running on port 8080")
		http.HandleFunc("/download", downloadFile)
		http.HandleFunc("/getDatas", getDatas)
		http.HandleFunc("/getFiles", getFileNames)
		http.ListenAndServe(":8080", nil)
	}()

	getDataEveryHourAndInsertIntoFileAndDB(wg)

	//checkErr(err)
	wg.Wait()
}

// l'user download le fichier qu'il souhaite en fonction
// loadenv pour les variable d'environnement ?
// logique de trading : get, set (condition du trade), eval,
