package handlers

import (
	"fmt"
	"forum/backend/db"
	"forum/backend/security"
	"forum/backend/sessions"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

const (
	MESSAGE_COOKIE_NAME = "forum-message"
)

var (
	Status = StatusData{
		Code: 200,
		Msg:  "",
	}
)

/*
RenderTemplate is a global http template rendering function which in addition to the standard http.Responsewriter
and *http.Request variables, also takes an interface{} for the data to be passed to the template, a string for
the file path of the template, and a string for the possible error message to be displayed (if template rendering
fails). The function parses the template file, executes the template with the given data, and writes the resultƒ
to the response writer. If an error occurs, the function invokes the error handling procedure via the SendError()
function.
*/
func renderTemplate(w http.ResponseWriter, r *http.Request, data any, filePath string) {
	// Initialise the template
	tmpl, err := template.ParseFiles(filePath)
	if err != nil {
		// Invoke error hadnling procedure
		fmt.Println("Error in parsing \""+filePath+"\": ", err)
		SendError(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	// Execute the template
	err = tmpl.Execute(w, data)
	if err != nil {
		// Invoke error handling procedure
		SendError(w, r, http.StatusInternalServerError, "Internal Server Error")
		return
	}
}

/*
redirectHandler is a global http redirect handler function which in addition to the standard http.Responsewriter
and *http.Request variables, also takes a string for the target page name and a string for the message to be
displayed on the target page. The message is stored in a cookie with a 10 second lifespan, and the user is
redirected to the target page.
*/
func redirectHandler(w http.ResponseWriter, r *http.Request, pageName string, message string) {
	// Create a new cookie with the message value
	var messageCookie = http.Cookie{
		Name:    MESSAGE_COOKIE_NAME,
		Value:   message,
		MaxAge:  1, // The cookie will last 10 seconds
		Expires: time.Now().Add(1 * time.Second),
		Path:    "/",
	}
	// Set the cookie on the response writer
	http.SetCookie(w, &messageCookie)

	// Redirect to the target page
	http.Redirect(w, r, pageName, http.StatusMovedPermanently)
}

/*
HandleError is a global http error handler function which in addition to the standard http.Responsewriter
and *http.Request variables, also takes an HTTP statusCode integer and http statusMessage string as inputs.
The function then compiles the error template with the given status code and message, and writes the output
to the http.ResponseWriter. If the template compilation fails, the function will return an internal server
error.
*/
func ErrorPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./frontend/static/error.html")
	if err != nil {
		// Extreme case: If the error template execution fails, return an internal server error and log the error
		log.Println("Error in parsing \"./frontend/static/error.html\": ", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, Status)
	if err != nil {
		// Extreme case: If the error template execution fails, return an internal server error and log the error
		log.Printf("Could not execute error template: %s", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

/*
LandingPage is a global http handler function which handles the landing page. The function checks for a logged-in
session cookie, and if found, redirects the user to the main page. If no session cookie is found, the function
renders the landing page template.
*/
func LandingPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// Handle invalid paths
		if r.URL.Path != "/" {
			SendError(w, r, http.StatusNotFound, "Page not found")
			return
		}

		// Check for logged-in session cookie, renew / update if found, return username if found
		userName, err := sessions.Check(w, r)
		if err != nil {
			// Redirect to login page if CheckSessions returns an error
			redirectHandler(w, r, "/login", err.Error()+": please try logging in or registering")
		} else if userName != "" {
			// Redirect to main page if user is logged in
			redirectHandler(w, r, "/main", "You are already logged in")
		}
		// Execute the template
		renderTemplate(w, r, nil, "./frontend/static/landingPage.html")
	} else {
		SendError(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
}

/*
RegisterPage is a global http handler function which handles the register page. The function checks for a logged-in
session cookie, and if found, redirects the user to the main page. If no session cookie is found, the function
renders the register page template. If the user submits the registration form, the function checks if the username
and password fields are valid, and if so, registers the user in the database. If the registration is successful,
the user is redirected to the login page.
*/
func RegisterPage(w http.ResponseWriter, r *http.Request) {
	// Check for logged-in session cookie, renew / update if found, return username if found
	userName, err := sessions.Check(w, r)
	if err != nil {
		// Redirect to login page if CheckSessions returns an error
		redirectHandler(w, r, "/login", err.Error()+": please try logging in or registering")
	} else if userName != "" {
		// Redirect to main page if user is logged in
		redirectHandler(w, r, "/main", "You are already logged in")
	}

	if r.Method == "POST" {
		// Get the form values and perform database checks for validity
		username := r.PostFormValue("username")
		inputEmail := r.PostFormValue("email")
		password, errPwd := security.HashPwd([]byte(r.PostFormValue("user_password")), 8)
		errUserName := db.NotExistData("users", "userName", username)
		errEmail := db.NotExistData("users", "email", inputEmail)

		// Handle invalid inputs
		if errPwd != nil {
			renderTemplate(w, r, []string{"An error was encountered with the password entered",
				username, inputEmail}, "./frontend/static/register.html")
		} else if errUserName != nil {
			renderTemplate(w, r, []string{"An account with that username already exists",
				username, inputEmail}, "./frontend/static/register.html")
		} else if errEmail != nil {
			renderTemplate(w, r, []string{"An account with that email already exists",
				username, inputEmail}, "./frontend/static/register.html")
		} else {

			// Register the user in the database
			dt := time.Now()
			_, err := db.InsertData("users", username, inputEmail, password,
				dt.Format("01-02-2006 15:04:05"))
			if err != nil {
				renderTemplate(w, r, []string{err.Error(), username, inputEmail},
					"./frontend/static/register.html")
			}
			// Redirect to the login page upon successful registration
			redirectHandler(w, r, "/login", "Your account has been created")
		}

	} else if r.Method == "GET" {
		renderTemplate(w, r, []string{"", "", ""}, "./frontend/static/register.html")
	} else {
		SendError(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
}

/*
LoginPage is a global http handler function which handles the login page. The function checks for a logged-in
session cookie, and if found, redirects the user to the main page. If no session cookie is found, the function
renders the login page template. If the user submits the login form, the function checks if the username and
password fields are valid, and if so, logs the user in. If the login is successful, the user is redirected to
the main page.
*/
func LoginPage(w http.ResponseWriter, r *http.Request) {
	// Check for logged-in session cookie, renew / update if found, return username if found
	userName, err := sessions.Check(w, r)
	if err != nil {
		fmt.Println(err.Error())
		// Reload login page if CheckSessions returns an error
		renderTemplate(w, r, err.Error()+": please try logging in or registering",
			"./frontend/static/login.html")
	} else if userName != "" {
		// Redirect to main page if user is logged in
		redirectHandler(w, r, "/main", "You are already logged in")
	}

	// Check if redirected from other pages or not
	msg := ""
	messageCookie, err := r.Cookie(MESSAGE_COOKIE_NAME)
	if err == nil {
		msg = messageCookie.Value
	}

	if r.Method == "POST" {
		// Get the form values and perform database checks for validity
		username := r.PostFormValue("username")
		password := r.PostFormValue("user_password")
		user, err := db.SelectDataHandler("users", "userName", username)
		if err != nil {
			msg = "The user doesn't exist"
			renderTemplate(w, r, msg, "./frontend/static/login.html")
		} else if !security.MatchPasswords([]byte(password), user.(db.User).Pass) {
			renderTemplate(w, r, "The password is incorrect", "./frontend/static/login.html")
		} else {
			err := sessions.Login(w, r, username) // Get userName from Login post method data
			if err != nil {
				msg = "An error was encountered while logging in. Please try again"
				renderTemplate(w, r, msg, "./frontend/static/login.html")
			}

			// Redirect to the main page upon successful login
			redirectHandler(w, r, "/main", "You are successfully logged in")
		}
	} else if r.Method == "GET" {
		// Display the login page with the message if redirected from other pages
		renderTemplate(w, r, msg, "./frontend/static/login.html")
	} else {
		SendError(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
	}
}

/*
MainPage is a global http handler function which handles the main page. The function checks for a logged-in
session cookie, and if found, renders the main page template. If no session cookie is found, the function
redirects the user to the login page. If the user submits the logout form, the function logs the user out
and redirects the user to the landing page.
*/
func MainPage(w http.ResponseWriter, r *http.Request) {
	// Check for logged-in session cookie, renew / update if found, return username if found
	userName, err := sessions.Check(w, r)
	if err != nil {
		// Redirect to login page if CheckSessions returns an error
		redirectHandler(w, r, "/login", err.Error()+": please try logging in or registering")
	} else if userName == "" {
		userName = "guest"
	}
	if r.Method == "POST" && userName == "guest" {
		redirectHandler(w, r, "/", "You must be logged in to perform this request")
		return
	}

	// Handle Logout POST request
	if r.Method == "POST" && r.PostFormValue("Logout") == "Logout" {
		// Perform cookies and sessions logout
		sessions.Logout(w, r)
		// Redirect to landing page once logged out
		redirectHandler(w, r, "/", "You have been logged out")
	}

	// Handle GET request and render main page
	if r.Method == "GET" {
		// Retreive data for main page
		mainPageData, err := GetMainDataStruct(r, userName)
		if err != nil {
			SendError(w, r, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		// get selected topic name from url
		url := r.URL.Query()
		topicName := url.Get("TopicName")

		// if topicName is not empty, get posts of topic
		if topicName != "" {
			postsTopic, err := getPostsOfTopic(topicName, r)
			if err != nil {
				//if there is no posts in topic, redirect to main page
				if err.Error() == "data doesn't exist in postTopics table" {
					redirectHandler(w, r, "/main", "No posts currently exist in "+topicName+" topic")
				} else {
					//for any other errors
					SendError(w, r, http.StatusInternalServerError, "Internal Server Error")
					return
				}

			}
			//fill mainData with posts of topic
			mainPageData.Posts = postsTopic
		}
		renderTemplate(w, r, mainPageData, "./frontend/static/main.html")
	} else {
		SendError(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
}

/*
ProfilePage is a global http handler for the profile page that displays the user's profile.
It first checks for a logged-in session cookie, and if found, it updates the session cookie.
... TBC ...
*/
func ProfilePage(w http.ResponseWriter, r *http.Request) {
	// Check for logged-in session cookie, renew / update if found, return username if found
	activeUsername, err := sessions.Check(w, r)
	if err != nil {
		// Redirect to login page if CheckSessions returns an error
		redirectHandler(w, r, "/login", err.Error()+": please try logging in or registering")
	} else if activeUsername == "" {
		activeUsername = "guest"
	}
	if r.Method == "POST" && activeUsername == "guest" {
		redirectHandler(w, r, "/", "You must be logged in to post")
		return
	}

	// Handle Logout POST request
	if r.Method == "POST" && r.PostFormValue("Logout") == "Logout" {
		// Perform cookies and sessions logout
		sessions.Logout(w, r)
		// Redirect to landing page once logged out
		redirectHandler(w, r, "/", "You have been logged out")
	}

	// Handle GET request and render profile page
	if r.Method == "GET" {
		url := r.URL.Query()
		username := url.Get("Username")
		// Check if username is valid
		_, err := getUserId(username)
		if err != nil {
			redirectHandler(w, r, "/main", "User "+username+" does not exist")
			return
		}
		// Retreive data for profile page
		profilePageData, err := GetProfileDataStruct(r, activeUsername, username)
		if err != nil {
			SendError(w, r, http.StatusInternalServerError, "Internal Server Error:\n"+err.Error())
			return
		}

		renderTemplate(w, r, profilePageData, "./frontend/static/profile.html")
	} else {
		SendError(w, r, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

}

/*
ContentPage is a global http handler for the content page that displays the content of a post.
It handles GET and POST requests. GET requests are used to display the content of a post.
POST requests are used to handle the user's interaction with the post, such as liking, disliking,
and commenting.
*/
func ContentPage(w http.ResponseWriter, r *http.Request) {
	// Check for logged-in session cookie, renew / update if found, return username if found
	userName, err := sessions.Check(w, r)
	if err != nil {
		// Redirect to login page if CheckSessions returns an error
		redirectHandler(w, r, "/login", err.Error()+": please try logging in or registering")
	} else if userName == "" {
		userName = "guest"
	}

	var postId int
	url := r.URL.Query()
	postId, err = strconv.Atoi(url.Get("postId"))
	if err != nil {
		SendError(w, r, 400, "post url error")
	}
	var contentData ContentData
	// Get content data from database and render template
	contentData, err = GetContentDataStruct(r, userName, postId)
	resetContentData := false
	if err != nil {
		//	SendError(w, r, http.StatusInternalServerError, err.Error())
	}
	if contentData.Content == "" {
		redirectHandler(w, r, "main", "post doesn't exist")
	}
	contentData.ActiveUsername = userName
	// Handle Logout POST request
	if r.Method == "POST" && r.PostFormValue("Logout") == "Logout" {
		// Perform cookies and sessions logout
		sessions.Logout(w, r)
		// Redirect to landing page once logged out
		redirectHandler(w, r, "/", "You have been logged out")

	} else if r.Method == "POST" {
		if userName == "guest" {
			SendError(w, r, http.StatusUnauthorized, "Unauthorized")
			return
		}
		resetContentData = true
		r.ParseForm()
		comment := ""
		comment = r.PostFormValue("comment")
		like := r.PostFormValue("likeButton")
		//adding or something
		if like == "like" {
			err = insertReaction(userName, postId, -1, like)
			if err != nil {
				SendError(w, r, 401, err.Error())
				return
			}
		}
		disLike := r.PostFormValue("dislikeButton")
		if disLike == "dislike" {
			err = insertReaction(userName, postId, -1, disLike)
			if err != nil {
				SendError(w, r, 401, err.Error())
				return
			}
		}
		// regex to check if comment is empty or contains only spaces
		reg := regexp.MustCompile(`^\s*$`)

		if comment != "" && !reg.MatchString(comment) {
			err := insertComment(userName, postId, comment)
			if err != nil {
				SendError(w, r, 400, err.Error())
				return
			}
		}
		// get reaction to comments
		for _, comment := range contentData.Comments {
			likeComment := r.PostFormValue("like" + strconv.Itoa(comment.CommentId))
			dislikeComment := r.PostFormValue("dislike" + strconv.Itoa(comment.CommentId))
			if likeComment == "like" {
				err = insertReaction(userName, -1, comment.CommentId, likeComment)
				if err != nil {
					SendError(w, r, 400, err.Error())
					return
				}
			}
			if dislikeComment == "dislike" {
				err = insertReaction(userName, -1, comment.CommentId, dislikeComment)
				if err != nil {
					SendError(w, r, 400, err.Error())
					return
				}
			}

		}

	} else if r.Method != "GET" {
		SendError(w, r, http.StatusBadRequest, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	if resetContentData == true {
		contentData, _ = GetContentDataStruct(r, userName, postId)
		if contentData.Content == "" {
			redirectHandler(w, r, "main", "post doesn't exist")
		}
	}

	renderTemplate(w, r, contentData, "./frontend/static/content.html")
}

/*
PostPage is a global http handler for the post page. It handles GET and POST requests.
It checks for a logged in session cookie, and redirects to the login page if none is found.
If a session cookie is found, it checks for a POST request, and if found, it parses the form and
inserts the post into the database, and then invokes a redirect to the main page. If a GET
request is found, it renders the post page.
*/
func PostPage(w http.ResponseWriter, r *http.Request) {
	// Check for logged-in session cookie, renew / update if found, return username if found
	userName, err := sessions.Check(w, r)
	if err != nil {
		// Redirect to login page if CheckSessions returns an error
		redirectHandler(w, r, "/login", err.Error()+": please try logging in or registering")
	} else if userName == "" {
		redirectHandler(w, r, "/", "You must be logged in to post")
		return
	}
	// Handle Logout POST request
	if r.Method == "POST" && r.PostFormValue("Logout") == "Logout" {
		// Perform cookies and sessions logout
		sessions.Logout(w, r)
		// Redirect to landing page once logged out
		redirectHandler(w, r, "/", "You have been logged out")

	} else if r.Method == "POST" {
		r.ParseForm()
		title, content, topicStr := r.PostFormValue("title"), r.PostFormValue("content"), r.PostFormValue("values")
		if topicStr == "No Topic" {
			redirectHandler(w, r, "/post", "Please select a topic")
			return
		}
		reg := regexp.MustCompile(`^\s*$`)
		if topicStr == "No Topic" {
			redirectHandler(w, r, "/post", "Please select a topic")
		}
		if !reg.MatchString(title) && !reg.MatchString(content) && !reg.MatchString(topicStr) {
			err = insertPostToDB(userName, title, content, topicStr)
			if err != nil {
				// Reload the page with an error message if post insertion fails
				redirectHandler(w, r, "/post", err.Error())
			}
			redirectHandler(w, r, "/main", "Your post has been created")
		}
	} else if r.Method == "GET" {
		postData, err := GetPostDataStruct(r, userName)
		if err != nil {
			SendError(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		// Render the template
		renderTemplate(w, r, postData, "./frontend/static/post.html")
	} else {
		SendError(w, r, http.StatusBadRequest, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
}
