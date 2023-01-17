package handler

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"net/http"
	"pub_sub_websocket_server/internal/repository"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Пропускаем любой запрос
	},
}

func NewRouter(s *repository.Storage) chi.Router {
	router := chi.NewRouter()

	controller := newServerHandler(s)

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Get("/", controller.Subscribe)

	router.NotFound(NotFound())
	router.MethodNotAllowed(NotAllowed())

	return router
}

type ServerHandler struct {
	Storage *repository.Storage
}

func newServerHandler(s *repository.Storage) *ServerHandler {
	return &ServerHandler{Storage: s}
}

func (h ServerHandler) Subscribe(w http.ResponseWriter, r *http.Request) {
	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Err(err).Msg("something goes wrong with new connection")
		return
	}
	defer connection.Close() // Закрываем соединение
	//Наивная логика распределения.
	h.Storage.Data[h.Storage.Cursor] = append(h.Storage.Data[h.Storage.Cursor], connection)
	log.Info().Int("with hub number", h.Storage.Cursor).Int("with id", h.Storage.Counter).Msg("New Client is connected.")
	if h.Storage.Counter+2-h.Storage.HubLimit-(h.Storage.Cursor*h.Storage.HubLimit) > 0 {
		h.Storage.Cursor++
	}
	h.Storage.Counter++
	for {
		mt, _, err := connection.ReadMessage()

		if err != nil || mt == websocket.CloseMessage {
			break // Выходим из цикла, если клиент пытается закрыть соединение или связь с клиентом прервана
		}
	}
}

// NotFound - обработчик неподдерживаемых маршрутов.
func NotFound() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		w.Write([]byte("route does not exist"))
	}
}

// NotAllowed - обработчик неподдерживаемых методов.
func NotAllowed() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		w.Write([]byte("method does not allowed"))
	}
}