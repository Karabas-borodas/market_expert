package web

import (
	"Karabas-borodas/market_expert.git/internal/config"
	"Karabas-borodas/market_expert.git/internal/domain"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log/slog"
	"math"
	"net/http"
	"net/mail"
	"regexp"
	"time"
)

// const httpPort = ":3002"

type UserAuthenticator struct {
	userService *domain.UserService
}

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
		if err := ValidateUserWeb(&req); err != nil {
			l.Debug("validation failed", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
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

// INFO: возможно нужно вынести в отденьный файл
func ValidateUserWeb(u *UserWeb) error {
	if u.Id == uuid.Nil {
		return fmt.Errorf("некорректный ID: %w")
	}
	if _, err := mail.ParseAddress(u.Email); err != nil {
		return fmt.Errorf("некорректный Email: %w", err)
	}
	if !regexp.MustCompile(`^\+[1-9]\d{1,14}$`).MatchString(u.Phone) {
		return fmt.Errorf("некорректный телефон")
	}
	if len(u.Name) < 2 || len(u.Sername) < 2 {
		return fmt.Errorf("имя и фамилия слишком короткие")
	}
	if u.Age < 18 || u.Age > 120 {
		return fmt.Errorf("возраст должен быть от 18 до 120 лет")
	}
	return nil
}
