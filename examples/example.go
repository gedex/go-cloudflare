package main

import (
	"github.com/gedex/go-cloudflare/api"
	"log"
)

func main() {
	conf := &api.Config{
		Email: "YOUR_EMAIL",
		Token: "YOUR_TOKEN",
	}
	client := api.NewClient(conf)
	settings, err := client.ClientAPI.ZoneSettings("example.com")
	if err != nil {
		log.Fatalln(err)
	}
	for _, z := range settings {
		log.Printf("UserSecuritySetting %v, CacheLevel %v\n", z.UserSecuritySetting, z.CacheLevel)
	}
}
