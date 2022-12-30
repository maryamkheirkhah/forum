package main

import (
	"fmt"
	"forum/backend/handlers"
	"forum/s"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	Temp     *template.Template
	HTTPData = s.StatusData{
		StatusCode: 200,
	}
	colour = s.Colours{
		Reset:      "\033[0m",
		Red:        "\033[31m",
		LightRed:   "\033[1;31m",
		Green:      "\033[32m",
		LightGreen: "\033[1;32m",
		Blue:       "\033[0;34m",
		LightBlue:  "\033[1;34m",
		Orange:     "\033[0;33m",
		Yellow:     "\033[1;33m",
	}
	//dbUsers    = map[string]s.User{} // User ID, user struct
	dbSessions = map[string]string{} // Session ID, User ID
)

func main() {
	// Start the loading bar
	block := '\u2588'
	fmt.Print(colour.LightGreen + "\nLoading server...\n\n" + colour.Reset)

	// Initialise templates and fileserver
	Temp = template.Must(template.ParseGlob("frontend/static/*.html"))
	fileServer := http.FileServer(http.Dir("./frontend"))
	http.Handle("/frontend/", http.StripPrefix("/frontend/", fileServer))

	// Register handlers, and continue with loading bar
	http.HandleFunc("/", handlers.Index)
	http.HandleFunc("/main", handlers.Main)
	http.HandleFunc("/user", handlers.User)
	http.HandleFunc("/login", handlers.Login)
	http.HandleFunc("/register", handlers.Register)
	http.HandleFunc("/topics", handlers.Topics)
	http.HandleFunc("/comments", handlers.Comments)
	http.HandleFunc("/post", handlers.Post)

	// Initialise channel which receives OS signals
	osChannel := make(chan os.Signal, 1)
	signal.Notify(osChannel, os.Interrupt, syscall.SIGTERM)

	// Incorporate server timeout
	server := &http.Server{
		Addr:              ":8080",
		ReadHeaderTimeout: 10 * time.Second,
	}

	// Initialise server on separate go-routine, listen / serve on specified port
	go func() {
		i := 1
		for i < 51 {
			time.Sleep(20 * time.Millisecond)
			fmt.Printf(colour.LightGreen+"%c"+colour.Reset, block)
			i++
		}
		fmt.Printf(colour.LightGreen+"\n\n... Complete!\n\nServer listening on port %v\n"+colour.Reset, server.Addr)
		errServer := server.ListenAndServe()
		if errServer != nil {
			HTTPData.StatusCode = 500
			HTTPData.StatusMsg = "Bad Server Request"
			fmt.Println(HTTPData.StatusCode, HTTPData.StatusMsg)
			log.Fatalln(errServer)
		}
	}()

	// Standby for a server shutdown signal
	<-osChannel
	fmt.Printf("\n\n" + colour.Red + "Server shutting down...\n" + colour.Reset)
}
