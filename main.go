package main

import (
	"github.com/gofiber/fiber/v2"
	socketio "github.com/googollee/go-socket.io"
	"log"
)

func main() {
	// * ============ init fiber ============
	app := fiber.New()

	//* ============ socket io ============
	server := socketio.NewServer(nil)

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		log.Println("connected:", s.ID())

		s.Emit("สวัสดีนี่คือข้อความจาก server")
		return nil
	})

	server.OnEvent("/", "notice", func(s socketio.Conn, msg string) {
		log.Println("notice:", msg)
		s.Emit("reply", "have "+msg)
	})

	server.OnEvent("/chat", "msg", func(s socketio.Conn, msg string) string {
		s.SetContext(msg)
		return "recv " + msg
	})

	server.OnEvent("/", "bye", func(s socketio.Conn) string {
		last := s.Context().(string)
		s.Emit("bye", last)
		s.Close()
		return last
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		log.Println("meet error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		log.Println("closed", reason)
	})

	go func() {
		if err := server.Serve(); err != nil {
			log.Fatalf("socketio listen error: %s\n", err)
		}
	}()
	defer server.Close()

	// * ============ route ============
	app.Get("/socket.io", server)

	// * Serve static files
	app.Static("/", "./index.html")

	log.Fatal(app.Listen("127.0.0.1:3000"))
}
