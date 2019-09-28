package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/Oxynger/steamAPI/steamapi"
)

func requirementFlagCheck(apiKey, userID string) error {

	if apiKey == " " {
		return errors.New("key not set. You need get steam key on https://steamcommunity.com/dev/apikey")
	}
	if userID == " " {
		return errors.New("id not set. You need set flag --id and enter steam id")
	}

	return nil

}

func main() {
	apiKey := flag.String("key", " ", "Steam api key")
	userID := flag.String("id", " ", "ID of steam profile")
	flag.Parse()

	if err := requirementFlagCheck(*apiKey, *userID); err != nil {
		panic(err.Error())
	}

	session := steamapi.NewProfile(*apiKey, *userID)

	fmt.Println(session.NoLife())
}
