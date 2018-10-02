package socketio

var AndroidServer *Server
var ClawServer *Server

func InitAndroidWebSocket() *Server {
	AndroidServer = NewServer("/android")
	go AndroidServer.AndroidServerListen()
	return AndroidServer
}

func InitClawWebSocket() *Server {
	ClawServer = NewServer("/claw")
	go ClawServer.ClawServerListen()
	return ClawServer
}
