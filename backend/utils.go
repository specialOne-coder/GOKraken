package main

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

// func loadenv() (config Config, err error) {
// 	cfg := Config{}
// 	if err := env.Parse(&cfg); err != nil {
// 		log.Fatalf("Error loading .env file: %v", err)
// 	}
// 	fmt.Println("url",cfg.APIURL)
// 	return cfg, nil
// }
