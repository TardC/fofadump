package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"

	"github.com/tardc/fofadump"
)

func main() {
	cfg := fofadump.NewFofaConfig()
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
	var full bool
	flag.BoolVar(&full, "full", false, "Whether to search all data")

	flag.Parse()

	if fofaQuery == "" {
		log.Fatalln("Please specify the fofa query rule")
	}
	// If not setting account, use default config in config file.
	if cfg.Email == "" || cfg.Key == "" {
		if defaultCfg, err := ioutil.ReadFile("config.json"); err == nil {
			if err := json.Unmarshal(defaultCfg, cfg); err == nil {
				if cfg.Email == "" || cfg.Key == "" {
					log.Fatalln("Please set email and key of fofa account")
				}
			} else {
				log.Fatalln("Set fofa account failed:", err)
			}
		} else {
			log.Fatalln("Read fofa config file failed:", err)
		}
	}

	fc := fofadump.NewFofaClient(cfg)
	fc.DoWork(fofaQuery, page, size, fields, full)
}
