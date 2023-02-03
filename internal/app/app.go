// Package app реализует сборку и запуск сервиса с Graceful Shutdown.
package app

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"

	"pub_sub_websocket_server/internal/handler"
	"pub_sub_websocket_server/internal/repository"
)

// Run - вход в приложение.
func Run() {
	storage := repository.NewRepository(repository.WithParseFlag())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	server := &http.Server{
		Addr:    ":8080",
		Handler: handler.NewRouter(storage),
	}
	// Запускаем горутину Grace-ful Shutdown.
	go func() {
		<-sig
		shutdownCtx, shutdownCtxCancel := context.WithTimeout(context.Background(), time.Second*20)
		defer shutdownCtxCancel()
		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal().Msg("graceful shutdown timed out and forcing exit.")
			}
		}()
		err := server.Shutdown(context.Background())
		if err != nil {
			log.Fatal().Err(err).Msg("server shutdown error")
		}
	}()
	// Запускаем сервер.
	go func() {
		log.Info().Str("starting server at", server.Addr)
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("failed to run server")
		}
	}()

	fmt.Println("Сервис готов к рассылке сообщений.Введите команду:")
	for {
		select {
		case <-sig:
			fmt.Println("Service goes close.")
			time.Sleep(time.Second * 25)
			return
		default:
			command, err := bufio.NewReader(os.Stdin).ReadString('\n')
			if err != nil {
				log.Err(err).Msg("something goes wrong with reader")
				continue
			}
			str := strings.Fields(command)
			if len(str) < 3 {
				fmt.Println("not enough arguments.Must be a three[command №Hub message].")
				continue
			}
			switch str[0] {
			case "send":
				storage.PubHub(command)
			case "sendc":
				storage.PubOne(command)
			default:
				fmt.Println("Unknown command")
			}
		}
	}
}
