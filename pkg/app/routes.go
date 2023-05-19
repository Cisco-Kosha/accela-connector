package app

import (
	"encoding/json"
	"github.com/kosha/accela-connector/pkg/httpclient"
	"github.com/kosha/accela-connector/pkg/logger"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (a *App) commonMiddleware(log logger.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//Allow CORS here By * or specific origin
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")

		if (*r).Method == "OPTIONS" {
			w.WriteHeader(200)
			return
		}

		var err error

		serverUrl := a.Cfg.GetServerURL()
		requestUri := r.RequestURI
		method := r.Method
		queryParams := r.URL.Query().Encode()

		var contentTypeHeaderFound bool

		url := serverUrl + requestUri
		if queryParams != "" && !strings.Contains(requestUri, "?") {
			url += "?" + queryParams
		}

		var c interface{}
		decoder := json.NewDecoder(r.Body)
		_ = decoder.Decode(&c)
		defer r.Body.Close()

		headers := make(map[string]string)
		// Loop over header names
		if len(r.Header) > 0 {
			for name, values := range r.Header {
				// Loop over all values for the name.
				if strings.ToLower(name) == "content-type" {
					contentTypeHeaderFound = true
				}
				for _, value := range values {
					if name != "" && value != "" {
						headers[name] = value
					}
				}
			}
		}
		// use application/json as default content type
		if !contentTypeHeaderFound {
			headers["Content-Type"] = "application/json; charset=utf-8"
		}
		clientId, clientSecret := a.Cfg.GetClientIDAndSecret()
		username, password := a.Cfg.GetUsernameAndPassword()
		scope := a.Cfg.GetScopes()
		env := a.Cfg.GetEnv()

		var token string
		var statusCode int
		accessTokenMap, ok := a.TokenMap[clientId]
		if !ok || accessTokenMap.AccessToken == "" {
			token = a.getToken(clientId, clientSecret, username, password, scope, env, a.Cfg.GetServerURL(), log)
		} else {
			expiryIn := accessTokenMap.ExpiresIn
			if expiryIn.Before(time.Now()) {
				token = a.getToken(clientId, clientSecret, username, password, scope, env, a.Cfg.GetServerURL(), log)
			} else {
				token = accessTokenMap.AccessToken
			}
		}

		res, statusCode, err := httpclient.MakeHttpCall(headers, clientId, clientSecret, method, a.Cfg.GetServerURL(), url, c, token, log)
		if err != nil {
			a.Log.Errorf("Encountered an error while making a call: %v\n", err)
			respondWithError(w, statusCode, err.Error())
			return
		}
		if res == nil {
			respondWithJSON(w, statusCode, res)
		}
		respondWithJSON(w, statusCode, res)
		return

	})
}

func (a *App) getToken(clientId, clientSecret, username, password, scope, env, serverUrl string, log logger.Logger) string {
	token, expiresIn, _ := httpclient.GenerateToken(clientId, clientSecret, username, password, scope, env, serverUrl, log)

	s := strconv.Itoa(expiresIn)
	duration, _ := time.ParseDuration(s + "s")

	a.TokenMap[clientId] = &TokenExpires{
		AccessToken: token,
		AppId:       clientId,
		ExpiresIn:   time.Now().Add(duration),
	}
	return token
}

func (a *App) InitializeRoutes(log logger.Logger) {
	a.Router.PathPrefix("/").Handler(a.commonMiddleware(log)).Methods("GET", "POST", "PUT", "DELETE", "OPTIONS")
}
