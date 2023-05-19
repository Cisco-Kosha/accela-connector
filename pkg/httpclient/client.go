package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/kosha/accela-connector/pkg/logger"
	"github.com/kosha/accela-connector/pkg/models"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	url2 "net/url"
)

func makePasswordCredentialsReq(method, url string, body interface{}, log logger.Logger) ([]byte, int) {

	var req *http.Request
	if body != nil {
		bytes1, _ := body.(string)
		req, _ = http.NewRequest(method, url, bytes.NewBuffer([]byte(bytes1)))
	} else {
		req, _ = http.NewRequest(method, url, nil)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	req.Header.Set("Accept-Encoding", "identity")

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		log.Error(err)
		return nil, 500
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
	}
	return bodyBytes, resp.StatusCode
}

func setOauth2Header(newReq *http.Request, tokenMap map[string]string) {
	fmt.Println(tokenMap["access_token"])
	newReq.Header.Set("Authorization", tokenMap["access_token"])
	newReq.Header.Set("Content-Type", "application/json")
	newReq.Header.Set("Accept", "application/json")
	newReq.Header.Set("x-accela-appid", tokenMap["app_id"])

	newReq.Header.Set("Accept-Encoding", "identity")

	return
}

func Oauth2ApiRequest(headers map[string]string, method, url string, data interface{}, tokenMap map[string]string, log logger.Logger) ([]byte, int) {
	var client = &http.Client{
		Timeout: time.Second * 10,
	}
	var body io.Reader
	if data == nil {
		body = nil
	} else {
		var requestBody []byte
		requestBody, err := json.Marshal(data)
		if err != nil {
			log.Error(err)
			return nil, 500
		}
		body = bytes.NewBuffer(requestBody)
	}

	request, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Error(err)
		return nil, 500
	}
	for k, v := range headers {
		request.Header.Add(k, v)
	}
	setOauth2Header(request, tokenMap)
	response, err := client.Do(request)

	if err != nil {
		log.Error(err)
		return nil, 500
	}
	defer response.Body.Close()
	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Error(err)
		return nil, 500
	}
	return respBody, response.StatusCode
}

func MakeHttpCall(headers map[string]string, consumerId, consumerSecret, method, serverUrl, url string, body interface{}, token string, log logger.Logger) (interface{}, int, error) {

	var response interface{}
	var payloadRes []byte

	var statusCode int
	tokenMap := make(map[string]string)

	if token != "" {
		tokenMap["access_token"] = token
		payloadRes, statusCode = Oauth2ApiRequest(headers, method, url, body, tokenMap, log)
		if string(payloadRes) == "" {
			return nil, statusCode, fmt.Errorf("nil")
		}
		// Convert response body to target struct
		err := json.Unmarshal(payloadRes, &response)
		if err != nil {
			log.Error("Unable to parse response as json")
			log.Error(err)
			return nil, http.StatusInternalServerError, err
		}
		if statusCode == 200 && response != nil {
			return response, statusCode, nil
		}
	}
	return nil, http.StatusInternalServerError, fmt.Errorf("token invalid")
}

func GenerateToken(clientId, clientSecret, username, password, scopes, env, serverUrl string, log logger.Logger) (string, int, error) {
	// token is not generated, or is invalid so get new token
	grantType := "password"
	token, expiresIn, _ := getToken(clientId, clientSecret, username, password, scopes, env, serverUrl, log, grantType, nil)
	if token == "" {
		return "", 0, fmt.Errorf("error generating token")
	}
	return token, expiresIn, nil
}

func getToken(clientId, clientSecret, username, password, scopes, env, serverUrl string, log logger.Logger,
	grantType string, body interface{}) (string, int, int) {

	var tokenResponse models.AccessToken

	url := serverUrl + "/oauth2/token"

	data := url2.Values{}
	data.Set("client_id", clientId)
	data.Set("client_secret", clientSecret)
	data.Add("username", username)
	data.Add("password", password)
	data.Add("scope", scopes)
	data.Add("environment", env)
	data.Add("grant_type", grantType)
	data.Add("agency_name", "NULLISLAND")

	encodedData := data.Encode()

	res, _ := makePasswordCredentialsReq("POST", url, encodedData, log)
	if string(res) == "" {
		return "", 0, 500
	}
	// Convert response body to target struct
	err := json.Unmarshal(res, &tokenResponse)
	if err != nil {
		log.Error("Unable to parse auth token response as json")
		log.Error(err)
		return "", 0, 500
	}
	return tokenResponse.AccessToken, tokenResponse.ExpiresIn, 200
}
