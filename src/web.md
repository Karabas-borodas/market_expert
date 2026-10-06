# Руководство по созданию HTTP-сервера и обработке ошибок в Go 1.22+

Данный документ описывает стандартный способ создания production-ready HTTP-сервера с использованием новой маршрутизации `net/http.ServeMux` из Go 1.22+, с внедрением логгера `*slog.Logger`, обработкой ошибок и поддержкой Graceful Shutdown.

---

## 1. Архитектура и поток данных (Data Flow)

Логгер `*slog.Logger` передается на три уровня серверного слоя:
1. **В `http.Server.ErrorLog`** — перенаправляет внутренние системные ошибки сервера `net/http` в `slog`.
2. **В Middleware (Промежуточное ПО)** — автоматически логирует каждый входящий запрос, время выполнения, HTTP-статус и перехватывает паники (Panic Recovery).
3. **В HTTP-хендлеры (Controllers)** — позволяет логировать ошибки валидации DTO и вызова бизнес-логики.

### Схема движения запроса

```text
Client Request
      │
      ▼
Middleware (Logging & Recovery) ──> Логирует HTTP Метод, Путь, Статус и Паники через slog
      │
      ▼
http.ServeMux (Go 1.22 Router)  ──> Направляет запрос на хендлер (например, "POST /api/v1/users")
      │
      ▼
UserHandler                     ──> Логирует ошибки валидации JSON/БД через свой s.log
      │
      ▼
Client Response
```

---

## 2. Пошаговая реализация

### 1. HTTP Handler с логгером (`internal/transport/http/user_handler.go`)

```go
package http

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"Karabas-borodas/market_expert.git/internal/domain"
)

type UserService interface {
	CreateUser(name string, age uint) (*domain.UserDomain, error)
}

type UserHandler struct {
	log         *slog.Logger
	userService UserService
}

func NewUserHandler(log *slog.Logger, userService UserService) *UserHandler {
	return &UserHandler{
		log:         log,
		userService: userService,
	}
}

func (h *UserHandler) RegisterRoutes(mux *http.ServeMux) {
	// Использование синтаксиса роутинга Go 1.22 (Метод + Путь)
	mux.HandleFunc("POST /api/v1/users", h.CreateUser)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
		Age  uint   `json:"age"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// Логируем ошибку декодирования через переданный логгер
		h.log.Warn("invalid json body in request", "error", err, "remote_addr", r.RemoteAddr)
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	user, err := h.userService.CreateUser(req.Name, req.Age)
	if err != nil {
		// Логируем внутреннюю ошибку бизнес-логики
		h.log.Error("failed to create user", "error", err, "name", req.Name)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(user)
}
```

---

### 2. Middleware для логирования и восстановления после паник (`internal/transport/http/middleware.go`)

```go
package http

import (
	"log/slog"
	"net/http"
	"time"
)

// responseWriterWrapper перехватывает HTTP статус-код для логирования
type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriterWrapper) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// LoggingMiddleware логирует каждый входящий HTTP запрос и предотвращает падение сервера
func LoggingMiddleware(log *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapper := &responseWriterWrapper{ResponseWriter: w, statusCode: http.StatusOK}

		// Перехват паник
		defer func() {
			if rec := recover(); rec != nil {
				log.Error("HTTP handler panic recovered",
					"panic", rec,
					"path", r.URL.Path,
					"method", r.Method,
				)
				http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(wrapper, r)

		// Логируем результат запроса
		log.Info("http request handled",
			"method", r.Method,
			"path", r.URL.Path,
			"status", wrapper.statusCode,
			"duration", time.Since(start).String(),
		)
	})
}
```

---

### 3. Сборка и запуск сервера с Graceful Shutdown (`internal/di/di.go`)

```go
package di

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"Karabas-borodas/market_expert.git/internal/config"
	"Karabas-borodas/market_expert.git/internal/logger"
	transportHTTP "Karabas-borodas/market_expert.git/internal/transport/http"
)

func StartProgramm() {
	// 1. Главный контекст приложения для Graceful Shutdown
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// 2. Загрузка конфигурации и инициализация логгера
	cfg := config.MustLoad()
	log := logger.SetupLogger(cfg.Env)

	log.Info("starting application", "env", cfg.Env)

	// 3. Создание роутера Go 1.22
	mux := http.NewServeMux()

	// 4. Инициализация хендлеров и регистрация маршрутов
	// userHandler := transportHTTP.NewUserHandler(log, userService)
	// userHandler.RegisterRoutes(mux)

	// 5. Обертывание роутера в logging middleware
	handlerWithMiddleware := transportHTTP.LoggingMiddleware(log, mux)

	// 6. Настройка HTTP-сервера с передачей логгера в ErrorLog
	srv := &http.Server{
		Addr:         cfg.HTTPserver.Address,
		Handler:      handlerWithMiddleware,
		ReadTimeout:  cfg.HTTPserver.Timeout,
		WriteTimeout: cfg.HTTPserver.Timeout,
		IdleTimeout:  cfg.HTTPserver.Iddle_timeout,
		// Перенаправляем внутренние ошибки сервера net/http в slog
		ErrorLog: slog.NewLogLogger(log.Handler(), slog.LevelError),
	}

	// 7. Запуск сервера в отдельной горутине
	go func() {
		log.Info("http server started", "address", cfg.HTTPserver.Address)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("http server failed to listen and serve", "error", err)
			os.Exit(1)
		}
	}()

	// 8. Ожидание сигнала выключения (Graceful Shutdown)
	<-ctx.Done()
	log.Info("stopping http server gracefully...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("http server forced shutdown failed", "error", err)
	} else {
		log.Info("http server stopped successfully")
	}
}
```
