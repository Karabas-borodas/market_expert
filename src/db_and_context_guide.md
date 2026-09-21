# Руководство: Подключение к PostgreSQL и передача Context, Logger и Pool в Go

Документ содержит подробный разбор работы с БД (PostgreSQL через `pgxpool.Pool`), правила передачи контекста и логгера, а также разбор паттернов и anti-patterns.

---

## 1. Главный архитектурный вопрос: Context, Logger и Pool в одной структуре

### Почему НЕЛЬЗЯ хранить `context.Context` внутри структуры?

Хранение `context.Context` в поле структуры — это **опасный антипаттерн в Go**, нарушающий ключевые принципы языка и архитектурные паттерны:

1. **Разный жизненный цикл (Lifecycle Mismatch):**
   * **Структура сервиса/репозитория (например, `UserService`):** Это **долгоживущий объект** (`Singleton` / Application Scope). Он инициализируется один раз при старте приложения и существует всё время работы сервера.
   * **`context.Context`:** Это **короткоживущий объект** (`Request Scope`). Он привязан к конкретной операции, HTTP-запросу или фоновой задаче. Контекст содержит тайм-ауты (deadlines), сигналы отмены (cancellation) и данные конкретного запроса (trace_id, user_id).

2. **Официальная рекомендация Go Team:**
   В документации к стандартному пакету `context` указано:
   > *"Do not store Contexts inside a struct type; instead, pass a Context explicitly to each function that needs it. The Context should be the first parameter, typically named `ctx`."*

3. **Гонки состояний (Data Races) и непредсказуемое поведение:**
   Если `ctx` записать в структуру, то при параллельной обработке нескольких HTTP-запросов (горутин) они будут либо перезаписывать этот контекст, либо переиспользовать отменённый/просроченный контекст предыдущего запроса.

---

## 2. Как делать правильно: Разделение по жизненному циклу

### А. Долгоживущие зависимости (`*slog.Logger` и `*pgxpool.Pool`)
Их следует передавать в структуру при её создании через конструктор (**Dependency Injection**):

```go
package repository

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	log  *slog.Logger   // Приватное поле (unexported)
	pool *pgxpool.Pool  // Приватное поле (unexported)
}

func NewUserRepository(log *slog.Logger, pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		log:  log,
		pool: pool,
	}
}
```

> **Примечание про заглавные/строчные буквы:**
> Поля `log` и `pool` делаются **со строчной буквы (приватными)**. Логгер и пул подключений — это внутренняя деталь реализации репозитория. Внешнему коду (HTTP-хендлерам) не нужен доступ к `repo.pool` или `repo.log`, он должен вызывать только методы репозитория.

### Б. Короткоживущий контекст (`context.Context`)
Контекст передается **первым параметром каждого метода**:

```go
func (r *UserRepository) GetUserByID(ctx context.Context, id int64) (*User, error) {
	r.log.Info("getting user from db", "id", id)

	var u User
	err := r.pool.QueryRow(ctx, "SELECT id, name FROM users WHERE id = $1", id).Scan(&u.ID, &u.Name)
	if err != nil {
		r.log.Error("failed to get user", "error", err, "id", id)
		return nil, err
	}

	return &u, nil
}
```

---

## 3. Правильная организация подключения к PostgreSQL (`pgxpool.Pool`)

### Конфигурация в `internal/config/config.go`

```go
package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string         `yaml:"env" env-default:"local"`
	HTTPserver HTTPserver     `yaml:"http_server"`
	Postgres   PostgresConfig `yaml:"postgres"`
}

type HTTPserver struct {
	Address       string        `yaml:"address" env-default:"localhost:20005"`
	Timeout       time.Duration `yaml:"timeout" env-default:"10s"`
	Iddle_timeout time.Duration `yaml:"iddle_timeout" env-default:"80s"`
}

type PostgresConfig struct {
	User     string `yaml:"user" env-default:"Mihail"`
	Password string `yaml:"password" env-default:"Chupa-lupa"`
	Host     string `yaml:"host" env-default:"127.0.0.1"`
	Port     string `yaml:"port" env-default:"5432"`
	Database string `yaml:"database" env-default:"marketplace"`
}

func MustLoad() Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config/config.yaml"
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("failed to read config: %s", err)
	}

	return cfg
}
```

### Модуль подключения в `internal/storage/postgres/postgres.go`

```go
package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"Karabas-borodas/market_expert.git/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

// New создает и проверяет пул подключений к PostgreSQL
func New(ctx context.Context, log *slog.Logger, cfg config.PostgresConfig) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse pool config failed: %w", err)
	}

	poolConfig.MaxConns = 20
	poolConfig.MinConns = 2
	poolConfig.MaxConnIdleTime = 5 * time.Minute
	poolConfig.HealthCheckPeriod = time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("create pgxpool failed: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping postgres failed: %w", err)
	}

	log.Info("✅ Successfully connected to PostgreSQL pool")
	return pool, nil
}
```

### Инициализация в `cmd/marketplace/main.go`

```go
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"Karabas-borodas/market_expert.git/internal/config"
	"Karabas-borodas/market_expert.git/internal/logger"
	"Karabas-borodas/market_expert.git/internal/storage/postgres"
)

func main() {
	// 1. Контекст приложения с отслеживанием сигналов отмены (Graceful Shutdown)
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// 2. Инициализация конфигурации и логгера
	cfg := config.MustLoad()
	log := logger.SetupLogger(cfg.Env)

	log.Info("starting application", "env", cfg.Env)

	// 3. Создание пула подключений к PostgreSQL
	pgPool, err := postgres.New(ctx, log, cfg.Postgres)
	if err != nil {
		log.Error("failed to connect to postgres", "error", err)
		os.Exit(1)
	}
	defer pgPool.Close() // Закрытие пула при остановке сервера

	log.Info("application initialized successfully")
}
```

---

## 4. Чек-лист проверки качества кода (Best Practices)

- [x] Контекст (`context.Context`) передается первыми аргументами в функции и методы, а не хранится в структурах.
- [x] Поля зависимостей (`log`, `pool`) в структурах репозиториев/сервисов объявлены со строчной буквы (unexported).
- [x] Пул подключений `pgxpool.Pool` создается единожды при старте в `main.go` и освобождается через `defer pool.Close()`.
- [x] Доступность базы проверяется вызовом `pool.Ping(ctx)`.
- [x] Ошибки инициализации на старте приводят к аварийному завершению через `os.Exit(1)`.
