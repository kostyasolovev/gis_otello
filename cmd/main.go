// Ниже реализован сервис бронирования номеров в отеле. В предметной области
// выделены два понятия: Order — заказ, который включает в себя даты бронирования
// и контакты пользователя, и RoomAvailability — количество свободных номеров на
// конкретный день.
//
// Задание:
// - провести рефакторинг кода с выделением слоев и абстракций
// - применить best-practices там где это имеет смысл
// - исправить имеющиеся в реализации логические и технические ошибки и неточности
package main

import (
	ordersHttp "booking/orders/http"
	"booking/pkg/logger"
	roomRepository "booking/rooms/repository"
	roomService "booking/rooms/service"
	"context"
	"os/signal"
	"syscall"

	orderRepository "booking/orders/repository"
	orderService "booking/orders/service"
	"net/http"
	"time"
)

const serverPort = ":8080"

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	log := logger.New()
	roomsDB := roomRepository.New()
	roomsSrv := roomService.NewRoomsService(roomsDB)

	ordersDB := orderRepository.New()
	ordersSrv := orderService.NewService(ordersDB, roomsSrv)
	handler := ordersHttp.New(ordersSrv, log)

	server := &http.Server{
		Addr:    serverPort,
		Handler: handler,
	}

	if err := listenHttp(ctx, server, log); err != nil {
		log.LogErrorf("server stopped with error: %v", err)
	}
}

const serverShutdownTimeout = 5 * time.Second

func listenHttp(ctx context.Context, server *http.Server, logger *logger.Logger) error {
	go func() {
		<-ctx.Done()

		shutctx, cancel := context.WithTimeout(context.Background(), serverShutdownTimeout)
		defer cancel()

		if err := server.Shutdown(shutctx); err != nil {
			logger.LogErrorf("error shutting down server: %s", err)
		}
	}()

	logger.LogInfo("starting server at port %s", serverPort)

	return server.ListenAndServe()
}
