package web

import (
	"Karabas-borodas/market_expert.git/internal/config"
	"Karabas-borodas/market_expert.git/internal/domain"
	"encoding/json"
	// "github.com/google/uuid"
	"log/slog"
	"net/http"
)

// const httpPort = ":3002"

type UserAuthenticator struct {
	userService *domain.UserService
}

//	func (a *UserAuthenticator) Middleware(next http.HandlerFunc) http.HandlerFunc {
//		return func(w http.ResponseWriter, r *http.Request) {
//			// Извлекаем токен из куки или заголовка
//			var userUUID uuid.UUID
//			var err error
//
//			cookie, err := r.Cookie("session_token")
//			if err == nil {
//				userUUID, err = uuid.Parse(cookie.Value)
//			} else {
//				authHeader := r.Header.Get("Authorization")
//				if authHeader != "" {
//					userUUID, err = a.validateBasicAuth(authHeader)
//				} else {
//					err = http.ErrNoCookie
//				}
//			}
//
//			if err != nil {
//				http.Error(w, "неавторизован", http.StatusUnauthorized)
//				return
//			}
//
//			// Проверяем существование пользователя
//			_, err = a.userService.GetUserInfo(userUUID)
//			if err != nil {
//				http.Error(w, "неавторизован", http.StatusUnauthorized)
//				return
//			}
//
//			// Добавляем UUID в контекст
//			ctx := context.WithValue(r.Context(), UserContextKey, userUUID)
//			next.ServeHTTP(w, r.WithContext(ctx))
//		}
//	}
func StartMarkerWeb(l *slog.Logger, cfg config.HTTPserver) {

	configWeb := HTTPserverDomainToWeb(cfg)
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world!"))
	})

	l.Info("✅ Successfully server started")

	mux.HandleFunc("POST /registr", func(w http.ResponseWriter, r *http.Request) {
		var req SignUpRequestWeb
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			l.Debug("user put incorrect JSON")
			http.Error(w, "incorrect JSON", http.StatusBadRequest)
			return
		}
		if req.Login == "" || req.Pass == "" {
			l.Debug("A user entered an ampty rassword or username")
			http.Error(w, "Empty pass or login", http.StatusBadRequest)
			return
		}

		// webUser := &SignUpRequestWeb{
		// 	Login: req.Login,
		// 	Pass:  req.Pass,
		// }
		// err := user.RegistUser(*SignUpWebToDomain(webUser))
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return
		// }
		w.WriteHeader(http.StatusCreated)
		l.Info("user add")
		w.Write([]byte("USERRRR"))
	})

	_ = http.ListenAndServe(configWeb.Address, mux)
}
