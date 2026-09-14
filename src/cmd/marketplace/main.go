package main

import (
	"Karabas-borodas/market_expert.git/internal/config"
	"fmt"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg)
}
