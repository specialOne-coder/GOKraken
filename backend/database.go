package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"

	//"strconv"
	"time"

	"github.com/lib/pq"
)

// Connection to the database
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

// Insert the asset pairs into the database
func InserPairIntoDBAndFile(dbase *sql.DB, assets AssetPair, timestamp string) {
	// _, err := dbase.Exec("DROP TABLE IF EXISTS pairs;")
	// checkErr(err)
	sqlStat := "CREATE TABLE IF NOT EXISTS public.pairs( id SERIAL NOT NULL, altname character varying NOT NULL, wsname character varying, aclassbase character varying,base character varying, aclassquote character varying, quote character varying, lot character varying,costdecimals integer, pairdecimals integer, lotdecimals integer,lotmultiplier integer, leveragebuy integer[],leveragesell integer[],margincall integer, marginstop integer, ordermin character varying, costmin character varying, ticksize character varying, status character varying, timestamp character varying,prix_achat character varying,prix_vente character varying,prix character varying,volume character varying, volume24h character varying, volume_moyenne character varying,low character varying, high character varying, opening character varying, PRIMARY KEY (id) ); ALTER TABLE IF EXISTS public.pairs OWNER to exam;"
	_, errors := dbase.Exec(sqlStat)

	checkErr(errors)
	var datasInFile []Aggregation
	for i, asset := range assets.Result {
		// Insertion des données dans la base de données
		assetPrice, err := GetKrakenAssetsPrice(i)
		checkErr(err)
		_, err = dbase.Exec("INSERT INTO pairs (altname, wsname, aclassbase, base, aclassquote, quote, lot, costdecimals, pairdecimals,lotdecimals, lotmultiplier, leveragebuy, leveragesell,margincall, marginstop, ordermin, costmin, ticksize, status,timestamp,prix_achat,prix_vente,prix ,volume,volume24h ,volume_moyenne,low, high, opening)VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22,$23,$24,$25,$26,$27,$28,$29)", asset.Altname, asset.Wsname, asset.AclassBase, asset.Base, asset.AclassQuote, asset.Quote, asset.Lot,
			asset.CostDecimals, asset.PairDecimals, asset.LotDecimals, asset.LotMultiplier,
			pq.Array(asset.LeverageBuy), pq.Array(asset.LeverageSell),
			asset.MarginCall, asset.MarginStop, asset.Ordermin, asset.Costmin, asset.Ticksize, asset.Status, timestamp,
			assetPrice.PRIX_ACHAT[0], assetPrice.PRIX_VENTE[0], assetPrice.PRICE[0], assetPrice.VOLUME[0], assetPrice.VWAP[1],
			assetPrice.VOLUME_MOYENNE[0], assetPrice.LOW[0], assetPrice.HIGH[0],
			assetPrice.OPENING)
		checkErr(err)
		//var dataInFile Aggregation
		aggregate := Aggregation{asset, assetPrice.PRIX_ACHAT[0], assetPrice.PRIX_VENTE[0], assetPrice.PRICE[0], assetPrice.VOLUME[0], assetPrice.VWAP[1],
			assetPrice.VOLUME_MOYENNE[0], assetPrice.LOW[0], assetPrice.HIGH[0],
			assetPrice.OPENING, timestamp}
		datasInFile = append(datasInFile, aggregate)
		checkErr(err)
	}
	t, err := time.Parse("2006-01-02T15:04:05Z", timestamp)
	checkErr(err)
	Paris := time.FixedZone("Paris Time", int((1 * time.Hour).Seconds()))
	trueTime := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, time.UTC).In(Paris)
	format := trueTime.Format("2006-01-02") + "_" + strconv.Itoa(trueTime.Hour()) + "H"
	InsertIntoFile(datasInFile, format)
	fmt.Println("Successfully inserted!")
}

// Get the asset pairs from the database
func GetDbPairs() (datas []Aggregation, bytes []byte, err error) {
	dbase, err := ConnectionToDB()
	checkErr(err)
	rows, err := dbase.Query("SELECT * FROM pairs")
	checkErr(err)
	defer rows.Close()
	var dbDatas []Aggregation

	for rows.Next() {

		var id int
		var altname string
		var wsname string
		var aclassbase string

		var base string
		var aclassquote string
		var quote string

		var lot string
		var costdecimals int
		var pairdecimals int

		var lotdecimals int
		var lotmultiplier int
		var leveragebuy pq.Int64Array

		var leveragesell pq.Int64Array
		var margincall int
		var marginstop int

		var ordermin string
		var costmin string
		var ticksize string

		var status string
		var timestamp string
		var prix_achat string

		var prix_vente string
		var prix string
		var volume string

		var volume24h string
		var volume_moyenne int
		var low string

		var high string
		var opening string

		var dbData Aggregation

		err = rows.Scan(&id, &altname, &wsname, &aclassbase,
			&base, &aclassquote, &quote,
			&lot, &costdecimals, &pairdecimals,
			&lotdecimals, &lotmultiplier, &leveragebuy,
			&leveragesell, &margincall, &marginstop,
			&ordermin, &costmin, &ticksize,
			&status, &timestamp, &prix_achat,
			&prix_vente, &prix, &volume,
			&volume24h, &volume_moyenne, &low,
			&high, &opening)
		checkErr(err)

		assets := AssetPairInfo{altname, wsname, aclassbase, base, aclassquote, quote, lot, costdecimals, pairdecimals, lotdecimals, lotmultiplier, leveragebuy, leveragesell, margincall, marginstop, ordermin, costmin, ticksize, status}
		dbData = Aggregation{assets, prix_achat, prix_vente, prix, volume, volume24h, volume_moyenne, low, high, opening, timestamp}
		dbDatas = append(dbDatas, dbData)

	}
	datasByte, _ := json.MarshalIndent(dbDatas, "", "\t")

	return dbDatas, datasByte, nil
}
