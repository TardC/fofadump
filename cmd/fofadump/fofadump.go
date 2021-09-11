package main

import (
	"encoding/json"
	"flag"
	"github.com/tardc/fofadump"
	"github.com/tardc/fofadump/config"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	cfg := config.NewFofaConfig()
	flag.StringVar(&cfg.Email, "email", "", "Email of fofa account")
	flag.StringVar(&cfg.Key, "key", "", "Key of fofa account")

	var size int
	flag.IntVar(&size, "size", 100, "The number size to query")
	var page int
	flag.IntVar(&page, "page", 1, "The number page to query")
	var fofaQuery string
	flag.StringVar(&fofaQuery, "q", "", "Fofa query rule")
	var fields string
	flag.StringVar(&fields, "fields", "host,ip,port", "Fields to query")

	flag.Parse()

	if fofaQuery == "" {
		log.Println("Please specify the fofa query rule")
		os.Exit(1)
	}
	// If not setting account, use default config in config file.
	if cfg.Email == "" || cfg.Key == "" {
		if defaultCfg, err := ioutil.ReadFile("config.json"); err == nil {
			if err := json.Unmarshal(defaultCfg, cfg); err == nil {
				if cfg.Email == "" || cfg.Key == "" {
					log.Println("Please set email and key of fofa account")
					os.Exit(1)
				}
			} else {
				log.Println("Set fofa account failed:", err)
				os.Exit(1)
			}
		} else {
			log.Println("Read fofa config file failed:", err)
			os.Exit(1)
		}
	}

	fc := fofadump.NewFofaClient(cfg)
	fc.DoWork(fofaQuery, page, size, fields)
}
