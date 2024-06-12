package steam

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const baseURL = "https://api.steampowered.com"
const userService = "ISteamUser"
const playerService = "IPlayerService"

// Player contains details about the Steam User
type Player struct {
	SteamID      string `json:"steamid"`
	PersonaName  string `json:"personaname"`
	AvatarSmall  string `json:"avatar"`
	AvatarMedium string `json:"avatarmedium"`
	AvatarFull   string `json:"avatarfull"`
}

// PlayersList contains a slice of Player objects.
type PlayersList struct {
	Players []Player `json:"players"`
}

// PlayersResult contains a "response" object with relevant data
type PlayersResult struct {
	Response PlayersList `json:"response"`
}

// Game contains details about a Steam game
type Game struct {
	AppID           int    `json:"appid"`
	Name            string `json:"name"`
	PlaytimeForever int    `json:"playtime_forever"`
	ImgIconURL      string `json:"img_icon_url"`
	ImgLogoURL      string `json:"img_logo_url"`
}

// GamesList contains a slice of Game objects
type GamesList struct {
	Games []Game `json:"games"`
}

// GamesResult contains a "response" object with relevant data
type GamesResult struct {
	Response GamesList `json:"response"`
}

// A Friend is a reference to a Player who is friends with a particular user
type Friend struct {
	SteamID     string `json:"steamid"`
	FriendSince int    `json:"friend_since"`
}

// FriendsList contains an array of Friend objects
type FriendsList struct {
	Friends []Friend `json:"friends"`
}

// FriendsResult contains a "friendslist"" object with relevant data
type FriendsResult struct {
	FriendsList FriendsList `json:"friendslist"`
}

// Service is the type that owns methods for fetching steam data
type Service struct {
	client *http.Client
	apiKey string
}

// NewService returns a service object configured with the provided Steam web API Key
func NewService(apiKey string) *Service {
	return &Service{
		client: &http.Client{},
		apiKey: apiKey,
	}
}

// Players accepts one or more steamIDs and returns a PlayersResult
func (s *Service) Players(steamIDs []string) ([]Player, error) {

	req, err := http.NewRequest(
		http.MethodGet,
		baseURL+"/"+userService+"/GetPlayerSummaries/v0002?"+url.Values{
			"steamids": {strings.Join(steamIDs, ",")},
			"key":      {s.apiKey},
		}.Encode(),
		nil,
	)
	if err != nil {
		fmt.Println("got an error forming the request")
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := s.client.Do(req)
	if res != nil {
		fmt.Println("got a response, so defering the body close")
		defer res.Body.Close()
	}

	if err != nil {
		fmt.Println("got an error on .Do of http request")
		return nil, err
	}

	parsedPlayers := PlayersResult{}
	// data, _ := ioutil.ReadAll(res.Body)
	// fmt.Println(string(data))
	if err := json.NewDecoder(res.Body).Decode(&parsedPlayers); err != nil {
		fmt.Println("got an error decoding the http response body")
		fmt.Printf("the err is: %s\n", err)
		return nil, err
	}

	return parsedPlayers.Response.Players, nil
}

// Player accepts one steamID and returns that player's object
func (s *Service) Player(steamID string) (Player, error) {

	req, err := http.NewRequest(
		http.MethodGet,
		baseURL+"/"+userService+"/GetPlayerSummaries/v0002?"+url.Values{
			"steamids": {steamID},
			"key":      {s.apiKey},
		}.Encode(),
		nil,
	)
	if err != nil {
		fmt.Println("got an error forming the request")
		return Player{}, err
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := s.client.Do(req)
	if res != nil {
		fmt.Println("got a response, so defering the body close")
		defer res.Body.Close()
	}

	if err != nil {
		fmt.Println("got an error on .Do of http request")
		return Player{}, err
	}

	parsedPlayers := PlayersResult{}
	// data, _ := ioutil.ReadAll(res.Body)
	// fmt.Println(string(data))
	if err := json.NewDecoder(res.Body).Decode(&parsedPlayers); err != nil {
		fmt.Println("got an error decoding the http response body")
		fmt.Printf("the err is: %s\n", err)
		return Player{}, err
	}
	playerResult := parsedPlayers.Response.Players[0]

	return playerResult, nil
}

// Games accepts one or more steamIDs and returns a GamesResult
func (s *Service) Games(steamID string) ([]Game, error) {

	req, err := http.NewRequest(
		http.MethodGet,
		baseURL+"/"+playerService+"/GetOwnedGames/v0001/?"+url.Values{
			"steamid":                   {steamID},
			"key":                       {s.apiKey},
			"include_appinfo":           {"1"},
			"include_played_free_games": {"1"},
		}.Encode(),
		nil,
	)
	if err != nil {
		fmt.Println("got an error forming the request")
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := s.client.Do(req)
	if res != nil {
		fmt.Println("got a response, so defering the body close")
		defer res.Body.Close()
	}

	if err != nil {
		fmt.Println("got an error on .Do of http request")
		return nil, err
	}

	parsedGames := GamesResult{}

	if err := json.NewDecoder(res.Body).Decode(&parsedGames); err != nil {
		fmt.Println("got an error decoding the http response body")
		return nil, err
	}
	gamesResult := parsedGames.Response.Games

	return gamesResult, nil
}

// Friends accepts a steamID and returns all friends for that ID
func (s *Service) Friends(steamID string) ([]Player, error) {

	req, err := http.NewRequest(
		http.MethodGet,
		baseURL+"/"+userService+"/GetFriendList/v0001/?"+url.Values{
			"steamid":      {steamID},
			"key":          {s.apiKey},
			"relationship": {"friend"},
		}.Encode(),
		nil,
	)
	if err != nil {
		fmt.Println("got an error forming the request")
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := s.client.Do(req)
	if res != nil {
		fmt.Println("got a response, so defering the body close")
		defer res.Body.Close()
	}

	if err != nil {
		fmt.Println("got an error on .Do of http request")
		return nil, err
	}

	parsedFriends := FriendsResult{}

	if err := json.NewDecoder(res.Body).Decode(&parsedFriends); err != nil {
		fmt.Println("got an error decoding the http response body")
		return nil, err
	}

	friendsResult := parsedFriends.FriendsList.Friends

	var idList []string
	for _, friend := range friendsResult {
		idList = append(idList, friend.SteamID)
	}
	playerFriendsList, err := s.Players(idList)
	if err != nil {
		return nil, err
	}
	return playerFriendsList, nil
}
