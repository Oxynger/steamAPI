package steamapi

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

// profile contain data for working with steam api
type profile struct {
	apiKey, steamID string
}

// NewProfile return new steam profile
func NewProfile(apiKey, steamID string) profile {
	return profile{
		apiKey:  apiKey,
		steamID: steamID,
	}
}

func (p *profile) requsetGetOwnedGemes() string {
	builder := strings.Builder{}

	builder.WriteString("http://api.steampowered.com/IPlayerService/GetOwnedGames/v0001/")
	builder.WriteString("?key=" + p.apiKey)
	builder.WriteString("&steamid=" + p.steamID)
	builder.WriteString("&format=json")

	return builder.String()
}

func getRequest(request string) ([]byte, error) {
	client := http.Client{Timeout: time.Second * 2}

	req, err := http.NewRequest(http.MethodGet, request, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New("Failed send request")
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("Could not read response body")
	}

	return body, err
}

func castGames(userData []byte) ([]interface{}, error) {
	response := make(map[string]interface{})
	if err := json.Unmarshal(userData, &response); err != nil {
		return nil, errors.New("Not marshaling json, please check given url \n" + err.Error())
	}

	userInfo, ok := response["response"].(map[string]interface{})
	if !ok {
		return nil, errors.New("Response key in json does not exist or not a map")
	}

	games, ok := userInfo["games"].([]interface{})
	if !ok {
		return nil, errors.New("Games key in json does not exist or not array")
	}

	return games, nil
}

func (p *profile) NoLife() float64 {
	request := p.requsetGetOwnedGemes()

	body, err := getRequest(request)
	if err != nil {
		log.Fatalln("Failed to send request on url: ", request+"\n")
	}

	games, err := castGames(body)
	if err != nil {
		log.Fatalln("Failed to parse json", err)
	}

	var totalTime float64
	for _, val := range games {
		game := val.(map[string]interface{})
		totalTime += game["playtime_forever"].(float64)
	}

	return totalTime / 60

}
