package repository

import (
	"flag"
	"strconv"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

type Storage struct {
	Data     map[int][]*websocket.Conn
	HubLimit int
	Counter  int
	Cursor   int
	sync.Mutex
}

func NewRepository() *Storage {
	s := &Storage{
		Data:     make(map[int][]*websocket.Conn),
		Counter:  0,
		Cursor:   0,
		HubLimit: 10,
	}
	s.ParseFlag()
	return s
}

func (s *Storage) ParseFlag() {
	flag.IntVar(&s.HubLimit, "n", s.HubLimit, "HUB_LIMIT")
	flag.Parse()
}

func (s *Storage) PubHub(str string) {
	strF := strings.Fields(str)
	hubNumber, err := strconv.Atoi(strF[1])
	if err != nil {
		log.Err(err).Msg("Something goes wrong with convert hub number.")
		return
	}
	if hubNumber > s.Cursor {
		log.Info().Msg("Out of quantity hubs.")
		return
	}
	message := strings.Join(strF[2:], " ")
	for _, v := range s.Data[hubNumber] {
		err = v.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Err(err).Msg("Something goes wrong with write message to user via hub .")
			continue
		}
	}
	log.Info().Str("hub number", strF[1]).Str("message", message).Msg("Successfully publishing message to hub.")
}

func (s *Storage) PubOne(str string) {
	strF := strings.Fields(str)
	userNumber, err := strconv.Atoi(strF[1])
	if err != nil {
		log.Err(err).Msg("Something goes wrong with convert user number.")
		return
	}
	if userNumber > s.Counter {
		log.Info().Msg("Out of quantity users.")
		return
	}
	curs := userNumber / s.HubLimit
	ind := userNumber - (curs * s.HubLimit)
	conn := s.Data[curs][ind]
	message := strings.Join(strF[2:], " ")
	err = conn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Err(err).Msg("Something goes wrong with write message to user .")
		return
	}
	log.Info().Str("user ID", strF[1]).Str("message", message).Msg("Successfully publishing message to user.")
}
