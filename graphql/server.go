package graphql

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

func Handler(schema *Schema) http.Handler {
	upgrader := &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		socket, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("upgrader.Upgrade: %v", err)
			return
		}
		defer socket.Close()

		//makeCtx := func(ctx context.Context) context.Context {
		//	return ctx
		//}

		//ServeJSONSocket(r.Context(), socket, schema, makeCtx, &simpleLogger{})
	})
}
