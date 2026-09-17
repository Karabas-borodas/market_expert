# Руководство по архитектуре и лучшим практикам Go (Production-Ready Layout)

Данный документ содержит разбор рекомендаций по построению архитектуры Go-приложения (монолита / микросервисов), разграничению ответственности слоев и работе с базой данных PostgreSQL.

---

## 1. Основные концепции и ошибки начальной архитектуры

### А. Чтение конфигурации в месте подключения (Антипаттерн)
* **Проблема:** Вызов `cleanenv.ReadEnv` или запуск парсинга конфигурации внутри модуля подключения к БД размывает ответственность. Модуль подключения не должен знать, откуда берутся логин и пароль (из ENV, файлов YAML или стороннего vault-сервиса).
* **Решение:** Вся конфигурация приложения считывается **в одном месте** — в `internal/config/config.go`. Далее в модуль базы данных передается готовая структура конфигурации.

### Б. Именование модулей и разделение ответственности (`storage` vs `db` & `repository`)
* **Проблема:** Называя модуль `storage` и помещая туда параметры подключения, парсинг ENV и методы выполнения запросов, вы создаете «монолитный ком». При добавлении других хранилищ (Redis, S3, Kafka) возникает путаница.
* **Решение:**
  1. **`internal/db/postgres`** — отвечает **только** за создание и проверку пула подключений `*pgxpool.Pool`.
  2. **`internal/repository`** — отвечает за написание конкретных SQL-запросов к таблицам (User, Product, Order и т.д.), используя общий `*pgxpool.Pool`.

### В. Работа с `context.Context` и Graceful Shutdown
* **Проблема:** Создание `context.Background()` внутри функции подключения лишает приложение возможности корректного плавного завершения (Graceful Shutdown).
* **Решение:** Главный `context.Context` создается в `main.go` через `signal.NotifyContext` и пробрасывается во все нижележащие слои. При остановке приложения все соединения закрываются через `defer pool.Close()`.

---

## 2. Пошаговая реализация Production-Ready архитектуры

### Шаг 1. Единая конфигурация в `internal/config/config.go`

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

---

### Шаг 2. Модуль подключения к PostgreSQL в `internal/db/postgres/postgres.go`

```go
package postgres

import (
	"context"
	"fmt"
	"time"

	"Karabas-borodas/market_expert.git/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

// New создает и проверяет пул подключений к PostgreSQL
func New(ctx context.Context, cfg config.PostgresConfig) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse pool config failed: %w", err)
	}

	poolConfig.MaxConns = 20
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

	return pool, nil
}
```

---

### Шаг 3. Слои репозиториев (Repository)

Пример доменного репозитория пользователя в `internal/repository/user.go`:

```go
package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

// Пример метода работы с базой
func (r *UserRepository) GetUserByID(ctx context.Context, id int64) error {
	// r.pool.QueryRow(ctx, "SELECT ... WHERE id = $1", id)
	return nil
}
```

---

### Шаг 4. Точка входа в `cmd/marketplace/main.go`

```go
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"Karabas-borodas/market_expert.git/internal/config"
	"Karabas-borodas/market_expert.git/internal/db/postgres"
	"Karabas-borodas/market_expert.git/internal/logger"
	"Karabas-borodas/market_expert.git/internal/repository"
)

func main() {
	// 1. Контекст приложения с отслеживанием сигналов выключения (Graceful Shutdown)
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// 2. Инициализация единой конфигурации и логгера
	cfg := config.MustLoad()
	log := logger.SetupLogger(cfg.Env)

	// 3. Подключение к PostgreSQL
	pgPool, err := postgres.New(ctx, cfg.Postgres)
	if err != nil {
		log.Error("failed to connect to postgres", "error", err)
		os.Exit(1)
	}
	defer pgPool.Close() // Автоматическое закрытие пула при завершении программы

	log.Info("✅ Successfully connected to PostgreSQL pool")

	// 4. Инициализация репозиториев с передачей общего pgPool
	userRepo := repository.NewUserRepository(pgPool)
	_ = userRepo
}
```

---

## 3. Итоги и правила архитектуры

1. **Единый источник истины для конфигурации:** Все параметры (HTTP, DB, Redis) читаются в `config`.
2. **Переиспользование ресурсов:** Пул подключений `*pgxpool.Pool` создается один раз в `main.go` и передается во все нужные репозитории.
3. **Безопасное завершение:** `defer pgPool.Close()` гарантирует завершение всех незавершенных транзакций при остановке приложения.
