package handler

import (
	"github.com/gorilla/websocket"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"pub_sub_websocket_server/internal/repository"
)

func TestNewRouter(t *testing.T) {
	t.Run("NewRouter", func(t *testing.T) {
		storage := repository.NewRepository()
		require.NotNil(t, storage)
		got := NewRouter(storage)
		require.NotNil(t, got)
	})
}

func TestServerHandler_Subscribe(t *testing.T) {
	t.Run("Subscribe", func(t *testing.T) {
		storage := repository.NewRepository()
		require.NotNil(t, storage)
		router := NewRouter(storage)
		require.NotNil(t, router)
		ts := httptest.NewServer(router)
		defer ts.Close()
		u := "ws" + strings.TrimPrefix(ts.URL, "http")
		ws, _, err := websocket.DefaultDialer.Dial(u, nil)
		if err != nil {
			t.Fatalf("%v", err)
		}
		defer ws.Close()
		require.Equal(t, 0, storage.Counter)
	})
}
