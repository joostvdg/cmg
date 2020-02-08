package mapgen

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/go-errors/errors"
	"github.com/joostvdg/cmg/pkg/analytics"
	"github.com/joostvdg/cmg/pkg/game"
	"github.com/joostvdg/cmg/pkg/webserver/model"
	log "github.com/sirupsen/logrus"
)

func ProcessMapGenerationRequest(rules game.GameRules, requestInfo model.RequestInfo) (model.Map, error) {
	start := time.Now()

	log.WithFields(log.Fields{
		"GameRules":  rules,
		"RequestId":  requestInfo.RequestId,
		"RequestURI": requestInfo.RequestURI,
		"HOST":       requestInfo.Host,
		"RemoteAddr": requestInfo.RemoteAddr,
	}).Info("Attempt to generate a fair map:")

	gameType := game.NormalGame
	if rules.GameType == 0 {
		gameType = game.NormalGame
	} else if rules.GameType == 1 {
		gameType = game.LargeGame
	}

	gameTypeTime := time.Now()
	gameTypeElapsed := gameTypeTime.Sub(start)
	log.WithFields(log.Fields{
		"Duration": gameTypeElapsed,
	}).Debug("Setup Game Type ")

	verbose := false
	totalGenerations := 0

	board := MapGenerationAttempt(gameType, verbose)
	elapsedGen, _ := time.ParseDuration("0ns")
	for !board.IsValid(rules, gameType) {
		totalGenerations++
		if totalGenerations > rules.Generations {
			return model.Map{}, errors.New("Stuck in generation loop")
		}
		startGen := time.Now()
		board = MapGenerationAttempt(gameType, verbose)
		finishGen := time.Now()
		elapsedGen += finishGen.Sub(startGen)
	}

	var content = model.Map{
		GameType: gameType.Name,
		Board:    board.Board,
		GameCode: board.GetGameCode(requestInfo.Delimiter),
	}

	t := time.Now()
	elapsed := t.Sub(start)
	avgDurationNanaseconds := int(elapsedGen.Nanoseconds()) / totalGenerations
	avgDuration, _ := time.ParseDuration(fmt.Sprintf("%dns", avgDurationNanaseconds))

	log.WithFields(log.Fields{
		"RequestId":              requestInfo.RequestId,
		"Total Generations":      totalGenerations,
		"Avg. Creation Duration": avgDuration,
		"Total Duration":         elapsed,
	}).Info("Created a new map")

	go LogRequestEvent(requestInfo, totalGenerations, elapsed, gameType.Name)

	return content, nil
}

func LogRequestEvent(requestInfo model.RequestInfo, generations int, elapsed time.Duration, gameType string) {
	// retrieve API endpoint of CMG Analytics
	sendAnalytics := false
	analyticsAPIEndpoint, apiOke := os.LookupEnv("ANALYTICS_API_ENDPOINT")
	if apiOke {
		sendAnalytics = true
	}

	duration := int(elapsed / time.Microsecond)

	// construct analytics object
	generationRequest := analytics.GenerationRequest{
		Duration:        duration,
		GameType:        "Normal",
		GenerationCount: generations,
		Host:            requestInfo.RemoteAddr,
		MapType:         gameType,
		Parameters:      nil,
		RequestID:       requestInfo.RequestId.String(),
		UserAgent:       "",
	}
	jsonBody, _ := json.Marshal(generationRequest)

	// retrieve bearer token
	bearerToken := getBearerToken()

	// call http endpoint of CMG Analytics
	if sendAnalytics {
		req, err := http.NewRequest("POST", analyticsAPIEndpoint, bytes.NewBuffer(jsonBody))
		if err != nil {
			log.Warn("Could not create request for Analytics")
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+bearerToken)
		client := &http.Client{}
		client.Timeout = time.Second * 180 // allow for some time for the heroku app to wake up
		_, requestErr := client.Do(req)
		if requestErr != nil {
			log.Error("Could not send analytics", requestErr)
		}
	} else {
		log.Warn("We have no Analytics API Endpoint, cannot send analytics!")
	}
}

func getBearerToken() string {
	endpointAvailable := false
	analyticsLoginEndpoint, apiOke := os.LookupEnv("ANALYTICS_LOGIN_ENDPOINT")
	if apiOke {
		endpointAvailable = true
	}

	// retrieve credentials of CMG Analytics
	analyticsAPIUser, apiUserOke := os.LookupEnv("ANALYTICS_API_USER")
	if !apiUserOke {
		analyticsAPIUser = "test"
	}

	analyticsAPIPassword, apiPasswordOke := os.LookupEnv("ANALYTICS_API_PASSWORD")
	if !apiPasswordOke {
		analyticsAPIPassword = "test"
	}
	loginRequest := analytics.LoginRequest{
		Username: analyticsAPIUser,
		Password: analyticsAPIPassword,
	}
	jsonBody, _ := json.Marshal(loginRequest)

	if !endpointAvailable {
		log.Warn("No Analytics Login Endpoint Available!")
	} else {
		req, err := http.NewRequest("POST", analyticsLoginEndpoint, bytes.NewBuffer(jsonBody))
		if err != nil {
			log.Warn("Could not create login request for Analytics")
		}
		req.SetBasicAuth(analyticsAPIUser, analyticsAPIPassword)
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		client.Timeout = time.Second * 180 // allow for some time for the heroku app to wake up
		resp, requestErr := client.Do(req)
		if requestErr != nil {
			log.Error("Could not login to analytics", requestErr)
		}

		defer resp.Body.Close()
		rawResponseData, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			log.Error("Cannot read analytics login response")
		}

		var loginResponse analytics.LoginResponse

		err = json.Unmarshal([]byte(rawResponseData), &loginResponse) // here!

		if err != nil {
			log.Error("Cannot parse analytics login response")
		}
		return loginResponse.AccessToken

	}
	return ""
}
