package web

import (
	"Karabas-borodas/market_expert.git/internal/config"
	"Karabas-borodas/market_expert.git/internal/domain"
	"context"
	"encoding/json"
	// "github.com/google/uuid"
	"encoding/base64"
	"log/slog"
	"math"
	"net/http"
	"time"
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
func StartMarkerWeb(ctx context.Context, l *slog.Logger, cfg config.HTTPserver, user *domain.UserService) {

	configWeb := HTTPserverDomainToWeb(cfg)
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world!"))
	})

	mux.HandleFunc("POST /createInfo", func(w http.ResponseWriter, r *http.Request) {
		var req UserWeb
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			l.Debug("user put incorrect JSON")
			http.Error(w, "incorrect JSON", http.StatusBadRequest)
			return
		}
		if req.Phone == "" || req.Email == "" || req.Age < 0 || req.Age > math.MaxInt {
			l.Debug("A user entered an ampty rassword or username, or age")
			http.Error(w, "Empty Phon number, Email, or age <18", http.StatusBadRequest)
			return
		}

		webUser := &UserWeb{
			Phone:   req.Phone,
			Email:   req.Email,
			Id:      req.Id,
			Name:    req.Name,
			Sername: req.Sername,
			Age:     req.Age,
		}
		//FIX: вставить функци заполнения информации
		err := user.SaveUserInfo(l, ctx, *UserWebToDomain(*webUser))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		l.Info("user info add")
		// w.Write([]byte("Uer created"))
	})

	mux.HandleFunc("POST /create", func(w http.ResponseWriter, r *http.Request) {
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

		webUser := &SignUpRequestWeb{
			Login: req.Login,
			Pass:  req.Pass,
		}
		webUser.Pass = base64.StdEncoding.EncodeToString([]byte(webUser.Pass))
		err := user.RegistUser(l, ctx, *SignUpWebToDomain(webUser))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			w.Write([]byte("User exist"))
			return
		}
		w.WriteHeader(http.StatusCreated)
		l.Info("user created")
		w.Write([]byte("Uer created"))
	})

	server := &http.Server{
		Addr:    configWeb.Address,
		Handler: mux,
	}

	// Запускаем сервер в горутине, чтобы не блокировать основной поток
	go func() {
		l.Info("✅ Successfully server started", "addr", configWeb.Address)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			l.Error("server listen error", "error", err)
		}
	}()

	// Ждём сигнала завершения (Ctrl+C / SIGTERM)
	<-ctx.Done()

	l.Info("shutting down server...")

	// Даём 5 секунд на завершение активных запросов
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		l.Error("server shutdown error", "error", err)
	} else {
		l.Info("server stopped gracefully")
	}
}
