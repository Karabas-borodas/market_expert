package web

import (
	"Karabas-borodas/market_expert.git/internal/config"
	"log/slog"
	"net/http"
	"time"
)

// const httpPort = ":3002"

type HTTPserverWEB struct {
	Address       string        `yaml:"address" env-default:"localhost:20005"`
	Timeout       time.Duration `yaml:"timeout" env-default:"10s"`
	Iddle_timeout time.Duration `yaml:"iddle_timeout" env-default:"80s"`
}

func HTTPserverDomainToWeb(serv config.HTTPserver) *HTTPserverWEB {
	return &HTTPserverWEB{
		Address:       serv.Address,
		Timeout:       serv.Timeout,
		Iddle_timeout: serv.Iddle_timeout,
	}
}
func StartMarkerWeb(log *slog.Logger, cfg config.HTTPserver) {

	configWeb := HTTPserverDomainToWeb(cfg)
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world!"))
	})

	// mux.HandleFunc("POST /game/create", auth.Middleware(func(w http.ResponseWriter, r *http.Request) {
	// 	userUUID := r.Context().Value(UserContextKey).(uuid.UUID)
	// 	u, _ := user.GetUserInfo(userUUID)
	//
	// 	var req struct {
	// 		Type string `json:"type"`
	// 	}
	// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
	// 		req.Type = domain.TypePVE
	// 	}
	//
	// 	game, err := service.CreateGame(*u, req.Type)
	// 	if err != nil {
	// 		http.Error(w, err.Error(), http.StatusInternalServerError)
	// 		return
	// 	}
	// 	w.Header().Set("Content-Type", "application/json")
	// 	json.NewEncoder(w).Encode(DomainToWeb(*game))
	// }))
	_ = http.ListenAndServe(configWeb.Address, mux)
}
