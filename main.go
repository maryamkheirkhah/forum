package main

import (
	"forum/backend/db"
	"forum/backend/server"
	"forum/backend/sessions"
	"time"
)

func main() {
	// Perform database check, initialise if not found
	db.Check("./backend/db/forum.db", "./backend/db/createDb.sql")

	// Initialise sessions struct, and start go-routine for periodic sessions cleanup
	sessions.ActiveSessions.Initialise()
	go sessions.CleanUpRoutine()

	// Start server
	theRealDealForum := server.StartServer(":8080", 5*time.Second)

	// Wait for interrupt signal to shut down server
	server.WaitForShutdownSignal(theRealDealForum)
}
