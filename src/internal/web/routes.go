package web

import (
	// "Karabas-borodas/market_expert.git/internal/domain"
	"net/http"
)

const httpPort = ":3002"

func StartMarkerWeb() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world!"))
	})

	_ = http.ListenAndServe("localhost:5000", mux)

	// lc.Append(fx.Hook{
	// 	OnStart: func(ctx context.Context) error {
	// 		log.Printf("Запуск сервера на %s\n", srv.Addr)
	// 		go func() {
	// 			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
	// 				log.Fatalf("Ошибка сервера: %s\n", err)
	// 			}
	// 		}()
	// 		return nil
	// 	},
	// 	OnStop: func(ctx context.Context) error {
	// 		log.Println("Остановка сервера...")
	// 		return srv.Shutdown(ctx)
	// 	},
	// })
}
