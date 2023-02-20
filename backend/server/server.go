package server

import (
	"context"
	"fmt"
	"forum/backend/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"text/template"
	"time"
)

var (
	Temp   *template.Template
	Colour = Colours{
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
)

/*
initiateHandlers registers the handlers for the server, which are defined in the
handlers package.
*/
func initiateHandlers() {
	http.HandleFunc("/", handlers.LandingPage)
	http.HandleFunc("/main", handlers.MainPage)
	http.HandleFunc("/login", handlers.LoginPage)
	http.HandleFunc("/register", handlers.RegisterPage)
	http.HandleFunc("/content", handlers.ContentPage)
	http.HandleFunc("/post", handlers.PostPage)
	http.HandleFunc("/profile", handlers.ProfilePage)
}

/*
DefineServer takes a port number and a timeout duration as inputs, and returns a pointer
to a http.Server struct which it creates.
*/
func defineServer(port string, timeout time.Duration) *http.Server {
	return &http.Server{
		Addr:              port,
		ReadHeaderTimeout: timeout,
	}
}

/*
LoadingBar takes a channel as input, and displays a loading bar until the channel is
written to.
*/
func loadingBar(serverUp, continuePrinting chan bool) {
	block := '\u2588'
	fmt.Print(Colour.LightGreen + "\nLoading server...\n\n" + Colour.Reset)
	// Loop as long as serverUp is not written to
	for {
		select {
		case <-serverUp:
			continuePrinting <- true
			return
		default:
			fmt.Print(Colour.LightGreen + string(block) + Colour.Reset)
			time.Sleep(10 * time.Millisecond)
		}
	}
}

/*
StartServer takes an input port number and a timeout duration, and starts the server
listening on that port. It also calls the LoadingBar function to display a loading bar
until the server is up.
*/
func StartServer(port string, timeout time.Duration) *http.Server {
	// Initialise templates and fileserver
	Temp = template.Must(template.ParseGlob("./frontend/static/*.html"))
	fileServer := http.FileServer(http.Dir("./frontend"))
	http.Handle("/frontend/", http.StripPrefix("/frontend/", fileServer))

	// Register handlers
	initiateHandlers()

	// Call local function to define server
	server := defineServer(port, timeout)

	// Initialise loading bar and the channel to signal when server is up
	serverUp := make(chan bool)
	continuePrinting := make(chan bool)
	go func() {
		loadingBar(serverUp, continuePrinting)
	}()

	// Start 0.5 second timer
	timer := time.NewTimer(750 * time.Millisecond)
	// Wait for timer to finish
	<-timer.C

	// Start the server
	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf(Colour.LightRed+"Error encountered while listening: %v"+
				Colour.Reset, err)
		}
	}()

	// Signal that server is up
	serverUp <- true

	// Wait for loading bar to finish
	<-continuePrinting
	fmt.Printf(Colour.LightGreen+"\n\nServer listening on port %v\n\n"+Colour.Reset, server.Addr)

	return server
}

/*
WaitForShutdownSignal waits for an interrupt signal from the user, and then prints a
message to the terminal. The server is then gracefully shut down.
*/
func WaitForShutdownSignal(server *http.Server) {
	osChannel := make(chan os.Signal, 1)
	signal.Notify(osChannel, os.Interrupt, syscall.SIGTERM)

	<-osChannel
	fmt.Printf("\n\n" + Colour.Orange + "Server shutting down...\n" + Colour.Reset)

	// Graceful server shutdown, waiting for 5 seconds to first cancel all requests
	// and then shutting down the server
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	server.Shutdown(ctx)

	// Delay printing of message to allow time for server to shut down
	timer := time.NewTimer(2 * time.Second)
	<-timer.C
	fmt.Printf(Colour.Orange + "Server shut down gracefully\n" + Colour.Reset)
}
