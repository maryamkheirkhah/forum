package handlers

import (
	"fmt"
	"forum/backend/db"
	"forum/s"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
)

/*
HandleError is a global http error handler function which in addition to the standard http.Responsewriter
and *http.Request variables, also takes an HTTP status data struct as input. If the StatusCode element of
this struct is not 200, error handling is initiated and the associated StatusMsg string is written to the
response body.
*/
func Error(w http.ResponseWriter, r *http.Request, err s.StatusData) {
	var handleErr error

	// Load error template if error is valid
	if err.StatusCode != 200 {
		w.WriteHeader(err.StatusCode)
		handleErr = s.Temp.ExecuteTemplate(w, "err.html", err)
		if handleErr != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}

func Index(w http.ResponseWriter, r *http.Request) {
	// Retrieve cookie / initialise cookie
	c, errCookie := r.Cookie("session")
	if errCookie != nil {
		sessionID, _ := uuid.NewV4()
		c = &http.Cookie{
			Name:  "session",
			Value: sessionID.String(),
		}
		http.SetCookie(w, c)
	}

	if r.URL.Path != "/" {
		//errorHandler(w, http.StatusNotFound, []string{"Page not found"})
		return
	}

	tmpl, err := template.ParseFiles("./frontend/static/landingPage.html")
	if err != nil {
		log.Println(string(err.Error()))
		//errorHandler(w, http.StatusInternalServerError, []string{"Could not parse template"})
		return
	}
	tmpl.Execute(w, nil)
	// If user exists, retrieve user
	/* 	var user s.User
	   	if userName, exists := db.Sessions[c.Value]; exists {
	   		user = db.dbUsers[userName]
	   	} */
}

func Main(w http.ResponseWriter, r *http.Request) {

}

func User(w http.ResponseWriter, r *http.Request) {

}

func Login(w http.ResponseWriter, r *http.Request) {
	var msg string
	var dataArray []string
	tmpl, err := template.ParseFiles("./frontend/static/login.html")
	if err != nil {
		log.Println(string(err.Error()))
		return
	}
	if r.Method == "POST" {
		r.ParseForm()
		username := r.PostFormValue("username")
		password := r.PostFormValue("password")
		fmt.Println(username)
		user, err := db.SelectDataHandler("users", "userName", username)
		if err != nil {
			dataArray = []string{err.Error(), username}
			tmpl.Execute(w, dataArray)
		} else if user.(s.User).Pass != password {
			msg = "password is incorrect"
			dataArray = []string{msg}
			tmpl.Execute(w, dataArray)
		} else {
			msg = "You are successfully logged in"
			dataArray = []string{msg, "", ""}
			tmpl.Execute(w, dataArray)
		}
	} else {
		dataArray = []string{"", "", ""}
		tmpl.Execute(w, dataArray)
	}
}

func Register(w http.ResponseWriter, r *http.Request) {
	var msg string
	var dataArray []string
	tmpl, err := template.ParseFiles("./frontend/static/register.html")
	if err != nil {
		log.Println(string(err.Error()))
		return
	}
	if r.Method == "POST" {
		/* _, userName, email, _, _ := db.SelectUser()
		r.ParseForm()
		username := r.PostFormValue("username")
		inputEmail := r.PostFormValue("email")
		password := r.PostFormValue("user_password")
		if userName == username {
		msg = "An account with that username already exists"
		dataArray = []string{msg, username, inputEmail}
		tmpl.Execute(w, dataArray)
		} else if email == r.PostFormValue("email") {
		msg = "An account with that email already exists"
		dataArray = []string{msg, username, inputEmail}
		tmpl.Execute(w, dataArray)
		} else {
		dt := time.Now()
		db.InsertUser(username, inputEmail, password, dt.Format("01-02-2006 15:04:05"))
		msg = "Your account has been created"
		dataArray = []string{msg, "", ""}
		tmpl.Execute(w, dataArray)
		} */
		r.ParseForm()
		username := r.PostFormValue("username")
		inputEmail := r.PostFormValue("email")
		password := r.PostFormValue("user_password")
		errUserName := db.NotExistData("users", "userName", username)
		errEmail := db.NotExistData("users", "email", inputEmail)
		if errUserName != nil {
			msg = "An account with that username already exists"
			dataArray = []string{msg, username, inputEmail}
			tmpl.Execute(w, dataArray)
		} else if errEmail != nil {
			msg = "An account with that email already exists"
			dataArray = []string{msg, username, inputEmail}
			tmpl.Execute(w, dataArray)
		} else {
			dt := time.Now()
			err := db.InsertData("users", username, inputEmail, password, dt.Format("01-02-2006 15:04:05"))
			if err != nil {
				msg = err.Error()
				fmt.Println(err)
				dataArray = []string{msg, username, inputEmail}
				tmpl.Execute(w, dataArray)
			}
			msg = "Your account has been created"
			dataArray = []string{msg, "", ""}
			tmpl.Execute(w, dataArray)
		}
	} else {
		dataArray = []string{"", "", ""}
		tmpl.Execute(w, dataArray)
	}
}

func Topics(w http.ResponseWriter, r *http.Request) {

}

func Comments(w http.ResponseWriter, r *http.Request) {

}

func Post(w http.ResponseWriter, r *http.Request) {

}
