package main

import (
	"flag"
	"fmt"
	"github.com/SlyMarbo/gmail"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"time"
)

type Form struct {
	Email   string
	Name    string
	Message string
}

//Password contains the password that will be used to log into the SMTP server
var Password string

//Password contains the password that will be used to log into the SMTP server
var Username string

// DefaultHandler is the http handler which will handle the form post request
func DefaultHandler(w http.ResponseWriter, h *http.Request) {
	h.ParseForm()
	form := Form{
		Email:   h.PostFormValue("email"),
		Name:    h.PostFormValue("name"),
		Message: h.PostFormValue("message"),
	}

	log.Println(form)
	email := gmail.Compose("chaisefarrar.xyz Contact-me", fmt.Sprintln(form))
	email.From = Username
	email.Password = Password

	email.AddRecipient(Username)

	err := email.Send()
	if err != nil {
		log.Println(err)
		fmt.Fprint(w,"failure")
	} else {
		fmt.Fprint(w,"Success")
	}

}

func main() {
	//Flags
	flag.StringVar(&Password, "p", "", "password to log into account")
	flag.StringVar(&Username, "u", "", "username to log into account")
	flag.Parse()

	log.Println("starting server...")

	r := mux.NewRouter()
	r.HandleFunc("/", DefaultHandler)
	r.Methods("POST")
	loggedRouter := handlers.LoggingHandler(os.Stdout, r)

	srv := &http.Server{
		Handler:      loggedRouter,
		Addr:         ":8181",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
