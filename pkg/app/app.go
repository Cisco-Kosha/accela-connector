package app

import (
	"github.com/gorilla/mux"
	"github.com/kosha/accela-connector/pkg/config"
	"github.com/kosha/accela-connector/pkg/logger"
	"log"
	"net/http"
	"time"
)

type App struct {
	Router   *mux.Router
	Log      logger.Logger
	Cfg      *config.Config
	TokenMap map[string]*TokenExpires
}

type TokenExpires struct {
	AccessToken string    `json:"access_token,omitempty"`
	AppId       string    `json:"app_id,omitempty"`
	ExpiresIn   time.Time `json:"expires_in,omitempty"`
}

func router() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	return router
}

//func commonMiddleware(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		// Do stuff here
//		fmt.Println(r.RequestURI)
//		// Call the next handler, which can be another middleware in the chain, or the final handler.
//		respondWithJSON(w, http.StatusOK, "ok")
//		//next.ServeHTTP(w, r)
//	})
//}

// Initialize creates the necessary scaffolding of the app
func (a *App) Initialize(log logger.Logger) {

	cfg := config.Get()

	a.Cfg = cfg
	a.Log = log

	a.TokenMap = make(map[string]*TokenExpires)

	a.Router = router()
}

// Run starts the app and serves on the specified addr
func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}
