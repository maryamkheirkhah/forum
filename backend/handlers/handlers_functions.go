package handlers

import (
	"errors"
	"forum/backend/db"
	"forum/backend/sessions"
	"net/http"
	"strings"
	"time"
)

/*
SendError is a function which takes a response writer, request, http error code and error
message as input and writes the error code and message to the global Status struct. It then
invokes the error handler, which executes the error template and writes the error page to the
response writer (it also has its own error handling, thus no need to check for errors here).
*/
func SendError(w http.ResponseWriter, r *http.Request, errorCode int, errorMsg string) {
	// Write error details to global Status struct
	Status.Code = errorCode
	Status.Msg = errorMsg

	// Invoke error handler, which has its own error handling
	ErrorPage(w, r)
}

/*
sortPostSlice is a function that takes a slice of Post structs as input and sorts
them in descending order of creation time using the "bubble-sort algorithm". It returns
a slice of Post structs.
*** NOTE THIS FUNCTION IS A QUICK FIX AS THE READS FROM THE DATABASE ARE NOT SORTED ***
*/
func sortPostSlice(posts []Post) []Post {
	// Initialise flag that is to indicate whether the slice is sorted or not
	sorted := false
	for !sorted {
		// Reset flag
		sorted = true
		for i := 0; i < len(posts)-1; i++ {
			for j := i + 1; j < len(posts); j++ {
				// Parse creation time of posts
				time1, err := time.Parse("02-01-2006 15:04:05", posts[i].CreationTime)
				if err != nil {
					continue
				}
				time2, err := time.Parse("02-01-2006 15:04:05", posts[j].CreationTime)
				if err != nil {
					continue
				}
				// If time of post i is after time of post j, swap them
				if time1.Before(time2) {
					posts[i], posts[j] = posts[j], posts[i]
					sorted = false
				}
			}
		}
	}
	return posts
}

/*
sortPostSummarySlice is a function that takes a slice of PostSummary structs as input and sorts
them in descending order of creation time using the "bubble-sort algorithm". It returns a slice
of PostSummary structs.
*/
// *** NOTE THIS FUNCTION IS A QUICK FIX AS THE READS FROM THE DATABASE ARE NOT SORTED ***
func sortSummarySlice(posts []Summary) []Summary {
	// Initialise flag that is to indicate whether the slice is sorted or not
	sorted := false
	for !sorted {
		// Reset flag
		sorted = true
		for i := 0; i < len(posts)-1; i++ {
			for j := i + 1; j < len(posts); j++ {
				// Parse creation time of posts
				time1, err := time.Parse("02-01-2006 15:04:05", posts[i].CreationTime)
				if err != nil {
					continue
				}
				time2, err := time.Parse("02-01-2006 15:04:05", posts[j].CreationTime)
				if err != nil {
					continue
				}
				// If time of post i is after time of post j, swap them
				if time1.Before(time2) {
					posts[i], posts[j] = posts[j], posts[i]
					sorted = false
				}
			}
		}
	}
	return posts
}

/*
sortCommentSlice is a function that takes a slice of Comment structs as input and sorts
them in descending order of creation time using the "bubble-sort algorithm". It returns
a slice of Comment structs.
*** NOTE THIS FUNCTION IS A QUICK FIX AS THE READS FROM THE DATABASE ARE NOT SORTED ***
*/
func sortCommentSlice(comments []Comment) []Comment {
	// Initialise flag that is to indicate whether the slice is sorted or not
	sorted := false
	for !sorted {
		// Reset flag
		sorted = true
		for i := 0; i < len(comments)-1; i++ {
			for j := i + 1; j < len(comments); j++ {
				// Parse creation time of comments
				time1, err := time.Parse("02-01-2006 15:04:05", comments[i].CreationTime)
				if err != nil {
					continue
				}
				time2, err := time.Parse("02-01-2006 15:04:05", comments[j].CreationTime)
				if err != nil {
					continue
				}
				// If time of comment i is before time of comment j, swap them
				if time1.Before(time2) {
					comments[i], comments[j] = comments[j], comments[i]
					sorted = false
				}
			}
		}
	}
	return comments
}

/*
sortStringSlice is a function that takes a slice of strings as input and sorts them in
ascending order using the "bubble-sort algorithm". It returns a slice of strings.
*** NOTE THIS FUNCTION IS A QUICK FIX AS THE READS FROM THE DATABASE ARE NOT SORTED ***
*/
func sortStringSlice(s []string) []string {
	// Initialise flag that is to indicate whether the slice is sorted or not
	sorted := false

	// Sort in alphabetical order
	for !sorted {
		// Reset flag
		sorted = true
		for i := 0; i < len(s)-1; i++ {
			for j := i + 1; j < len(s); j++ {
				// If string i is after string j, swap them
				if s[i] > s[j] {
					s[i], s[j] = s[j], s[i]
					sorted = false
				}
			}
		}
	}
	return s
}

/*
getUserId takes a username string as input and returns an integer representing the userId
of the user, as well as an error value. The error value is nil if no error occurs whilst
calling other functions etc., otherwise the first error encountered is returned along with
a -1 integer as the userId.
*/
func getUserId(userName string) (int, error) {
	user, err := db.SelectDataHandler("users", "username", userName)
	if err != nil {
		return -1, errors.New("error in getting userId")
	}
	if user != nil {
		return user.(db.User).UserId, nil
	}
	return -1, errors.New("user not found error in getUserId")
}
func getUserName(userId int) (string, error) {
	user, err := db.SelectDataHandler("users", "userId", userId)
	if err != nil {
		return "", errors.New("error in getting userName" + err.Error())
	}
	if user != nil {
		return user.(db.User).Username, nil
	}
	return "", errors.New("user not found error in getUserName")
}

/*
FillUserStruct is a method for the (u *User) struct which takes a username and returns
a User struct with the username and email data from the database. It also returns an
error value, which is non-nil if there was an error in the database query.
*/
func (u *User) FillUserStruct(userName string) error {
	u.UserRank = "" // Fill this UserRank field in later task

	// Extract entire row for user from "users" table in database
	userInfo, err := db.SelectDataHandler("users", "userName", userName)
	if err != nil {
		return errors.New("error in getting user info from database: " + err.Error())
	}
	// Assign e-mail value to output struct
	u.UserEmail = userInfo.(db.User).Email

	return nil
}

/*
FillPostSummaryStruct is a method for the (ps *PostSummary) which takes a post ID
as input and returns a PostSummary struct with the post title, postID and creation
time data from the database. It also returns an error value, which is non-nil if
there was an error in the database query.
*/
func (s *Summary) FillPostSummaryStruct(postID int) error {
	s.Id = postID

	// Extract entire row for post from "posts" table in database
	postInfo, err := db.SelectDataHandler("posts", "postId", postID)
	if err != nil {
		return errors.New("error in getting post info from database:" + err.Error())
	}

	// If postInfo is em

	// Assign values to output struct
	s.Title = postInfo.(map[int]db.Post)[postID].Title
	s.CreationTime = postInfo.(map[int]db.Post)[postID].CreationTime

	return nil
}

/*
FillCommentSummaryStruct is a method for the (s *Summary) which takes a post ID
as input and returns a PostSummary struct with the post title, commentID  and
creation time data from the database. It also returns an error value, which is
non-nil if there was an error in the database query.
*/
func (s *Summary) FillCommentSummaryStruct(commentID int) error {
	// Extract entire row for post from "posts" table in database
	commentInfo, err := db.SelectDataHandler("comments", "commentId", commentID)
	if err != nil {
		return errors.New("error in getting post info from database:" + err.Error())
	}
	//s.Message += commentInfo.(map[int]db.Comment)[commentID].Content
	// Assign values to output struct
	// Get title of post that comment belongs to (THIS IS NOT EFFICIENT), should have post title in comments table
	postInfo, err := db.SelectDataHandler("posts", "postId", commentInfo.(map[int]db.Comment)[commentID].PostId)
	if err != nil {
		return errors.New("error in getting post info from database: " + err.Error())
	}
	post := postInfo.(map[int]db.Post)[commentInfo.(map[int]db.Comment)[commentID].PostId]
	s.Title = post.Title
	s.CreationTime = commentInfo.(map[int]db.Comment)[commentID].Time
	s.Id = post.PostId
	return nil
}

/*
FillEventsStruct function description...
*/
func (e *Events) FillEventsStruct() error {
	return nil
}

/*
FindReactionStatus takes a username string anda postID int as inputs and returns
an int and an error value. The int is the status of the reaction of the user on
the post. 1 means like, -1 means dislike, 0 means no reaction and -2 means we got error. The error value
is nil if no error occurs whilst calling other functions etc., otherwise the first
error encountered is returned.
*/
func findReactionStatus(username string, postId int, commentId int) (int, error) {
	userId, err := getUserId(username)

	if err != nil {
		return -2, errors.New("error in getting user id:" + err.Error())
	}

	var reaction interface{}
	if postId != -1 {
		// Get reactions for the input postID
		reaction, err = db.SelectDataHandler("reactions", "postId", postId)
	} else if commentId != -1 {
		reaction, err = db.SelectDataHandler("reactions", "commentId", commentId)
	}
	if err != nil && err.Error() != "data doesn't exist in reactions table" {
		return -2, errors.New("error in getting reactions" + err.Error())

		// Find the reaction of the input user
	} else if err == nil {
		for _, r := range reaction.(map[int]db.Reaction) {
			if r.UserId == userId {
				if r.Reaction == "like" {
					return 1, nil
				} else if r.Reaction == "dislike" {
					return -1, nil
				}
			}
		}
	}
	return 0, nil
}

/*
getTopicsOfPost takes a topicID int as input and returns a slice of integers as well
as an error value. The slice of integers contains the postIDs associated with the topic.
The error value is nil if no error occurs whilst calling other functions etc.
*/
func getPostIdsOfTopic(topicName string) ([]int, error) {
	topicMap, err := db.SelectDataHandler("topics", "topicName", topicName)
	if err != nil {
		return nil, errors.New("error in getting topic id:" + err.Error())
	}
	var topicId int
	for id, topicname := range topicMap.(map[int]string) {
		if topicname == topicName {
			topicId = id
		}
	}
	postTopicMap, err := db.SelectDataHandler("PostTopics", "topicId", topicId)
	var postIds []int
	if err != nil {
		return nil, err
	}
	for _, topic := range postTopicMap.(map[int][]int) {
		postIds = append(postIds, topic...)
	}
	return postIds, nil
}

/*
getPostSummariesOfUser takes a username string as input and returns a slice of Summary
structs as well as an error value. The slice of PostSummary structs contains the Summary
structs for all posts created by the user. The error value is nil if no error occurs
whilst calling other functions etc. Otherwise, the first error encountered is returned.
*/
func getCreatedPostSummaries(username string) ([]Summary, error) {
	userId, err := getUserId(username)
	if err != nil {
		return nil, errors.New("error in getting user id: " + err.Error())
	}

	// Get PostSummary structs for all posts created by user
	postMap, err := db.SelectDataHandler("posts", "userId", userId)
	if err != nil && err.Error() != "data doesn't exist in postTopics table" {
		return nil, errors.New("error in getting posts: " + err.Error())
	} else if postMap == nil {
		return nil, nil
	}

	// Initialise variables
	postSummaries := []Summary{}
	var newPostSummary Summary

	// Populate individual PostSummary structs and append them to the output slice
	for _, post := range postMap.(map[int]db.Post) {
		err = newPostSummary.FillPostSummaryStruct(post.PostId)
		if err != nil {
			return nil, errors.New("error in filling post summary struct: " + err.Error())
		}
		postSummaries = append(postSummaries, newPostSummary)
	}

	return postSummaries, nil
}

/*
getLikePostSummaries takes a username string as input and returns a slice of Summary
structs as well as an error value. The slice of PostSummary structs contains the Summary
structs for all posts and comments liked by the user. The error value is nil if no error
occurs whilst calling other functions etc. Otherwise, the first error encountered is
returned.
*/
func getLikedPostSummaries(username string) ([]Summary, error) {
	userId, err := getUserId(username)
	if err != nil {
		return nil, errors.New("error in getting user id: " + err.Error())
	}

	// Get PostSummary structs for all posts reacted by user
	postMap, err := db.SelectDataHandler("reactions", "userId", userId)
	if err != nil && err.Error() != "data doesn't exist in reactions table" {
		return nil, errors.New("error in getting likes \\ dislikes record: " + err.Error())
	} else if postMap == nil {
		return nil, nil
	}
	// Initialise variables
	postSummaries := []Summary{}
	var newPostSummary Summary

	// Populate individual PostSummary structs for likes and append them to the output slice
	for _, post := range postMap.(map[int]db.Reaction) {
		if post.Reaction == "like" {
			if post.PostId != -1 {
				// Differentiate for liked posts
				err = newPostSummary.FillPostSummaryStruct(post.PostId)
				if err != nil {
					return nil, errors.New("error in getting likes \\ dislikes record: " + err.Error())
				}
				newPostSummary.Message += " post liked"
				postSummaries = append(postSummaries, newPostSummary)
			} else if post.CommentId != -1 {
				// Differentiate for liked comments
				err = newPostSummary.FillCommentSummaryStruct(post.CommentId)
				if err != nil {
					return nil, errors.New("error in getting likes \\ dislikes record: " + err.Error())
				}
				newPostSummary.Message = "comment liked"
				postSummaries = append(postSummaries, newPostSummary)
			}
		}
	}
	return postSummaries, nil
}

/*
getDislikePostSummaries takes a username string as input and returns a slice of Summary
structs as well as an error value. The slice of PostSummary structs contains the Summary
structs for all posts and comments disliked by the user. The error value is nil if no error
occurs whilst calling other functions etc. Otherwise, the first error encountered is
returned.
*/
func getDislikedPostSummaries(username string) ([]Summary, error) {
	userId, err := getUserId(username)
	if err != nil {
		return nil, errors.New("error in getting user id: " + err.Error())
	}

	// Get PostSummary structs for all posts reacted by user
	postMap, err := db.SelectDataHandler("reactions", "userId", userId)
	if err != nil && err.Error() != "data doesn't exist in reactions table" {
		return nil, errors.New("error in getting likes \\ dislikes record: " + err.Error())
	} else if postMap == nil {
		return nil, nil
	}

	// Initialise variables
	postSummaries := []Summary{}
	var newPostSummary Summary

	// Populate individual PostSummary structs for likes and append them to the output slice
	for _, post := range postMap.(map[int]db.Reaction) {
		if post.Reaction == "dislike" {
			if post.PostId != -1 {
				// Differentiate for dislike posts
				err = newPostSummary.FillPostSummaryStruct(post.PostId)
				if err != nil {
					return nil, errors.New("error in getting likes \\ dislikes record: " + err.Error())
				}
				newPostSummary.Message += " post disliked"
				postSummaries = append(postSummaries, newPostSummary)
			} else if post.CommentId != -1 {
				// Differentiate for disliked comments
				err = newPostSummary.FillCommentSummaryStruct(post.CommentId)
				if err != nil {
					return nil, errors.New("error in getting likes \\ dislikes record: " + err.Error())
				}
				newPostSummary.Message += " comment disliked"
				postSummaries = append(postSummaries, newPostSummary)
			}
		}
	}

	return postSummaries, nil
}

/*
getAllPosts takes no input, and harvests all posts from the database, returning a slice
of Post structs as well as an error value. The error value is nil if no error occurs whilst
calling other functions etc., otherwise the first error encountered is returned.
*/
func getAllPosts(r *http.Request) ([]Post, error) {
	// Initialise slice of Posts structs to be returned
	allPosts := []Post{}

	// Retrieve all posts from the database
	postMap, err := db.SelectDataHandler("posts", "", nil)
	if err != nil {
		return nil, errors.New("error in getting all posts from database:" + err.Error())
	}

	// Populate individual Post structs and append them to the slice of Post structs
	for _, post := range postMap.(map[int]db.Post) {
		var newPost Post
		err = newPost.FillPostStruct(post, r)
		if err != nil {
			return nil, errors.New("error in filling post struct:" + err.Error())
		}
		allPosts = append(allPosts, newPost)
	}
	return allPosts, nil
}

/*
getAllTopics takes no input, and harvests all topics from the database, returning a slice
of strings representing all topics listed in the database, as well as an error value. The
error value is nil if no error occurs whilst calling other functions etc., otherwise the
first error encountered is returned.
*/
func getAllTopics() ([]string, error) {
	// Initialise slice of topic strings to be returned
	allTopics := []string{}

	// Retrieve all topics from the database
	topicMap, err := db.SelectDataHandler("topics", "", nil)
	if err != nil {
		return nil, errors.New("error in getting AllTopics")
	}

	// Append all topics to the slice of topic strings
	for _, topic := range topicMap.(map[int]string) {
		allTopics = append(allTopics, topic)
	}
	return allTopics, nil
}

/*
getPostsOfTopic takes a topicID int as input and returns a slice of Post structs as well
as an error value. The returned slice contains all the posts associated with the topic.
The error value is nil if no error occurs whilst calling other functions etc., otherwise
the first error encountered is returned.
*/
func getPostsOfTopic(topicName string, r *http.Request) ([]Post, error) {
	// Retrieve all posts from the database
	allPosts, err := getAllPosts(r)
	if err != nil {
		return nil, errors.New("error in getting all posts:" + err.Error())
	}

	// Retrieve all postIDs associated with the input topicID
	postIds, err := getPostIdsOfTopic(topicName)
	if err != nil {
		return nil, err
	}

	// Initialise slice of posts to be returned and append posts to it which have the
	// matching topic associated with the input topicID
	var postsOfTopic []Post
	for _, post := range allPosts {
		for _, postId := range postIds {
			if post.PostId == postId {
				postsOfTopic = append(postsOfTopic, post)
			}
		}
	}
	return sortPostSlice(postsOfTopic), nil
}

/*
insertPostToDB takes a deconstructed post as input, in the form of username, title, content
and topic strings, and returns an error value. This function inserts the post, represented by the
collection of inputs, to the database. The returned error value is nil if no error occurs whilst
calling other functions etc., otherwise the first error encountered is returned.
*/
func insertPostToDB(username string, title string, content string, topicStr string) error {
	topics := strings.Split(topicStr, " and ")
	var topicIds []int
	dt := time.Now()
	userId, err := getUserId(username)
	if err != nil {
		return errors.New("error in getting userId:" + err.Error())
	}
	topicMap, err := db.SelectDataHandler("topics", "", nil)
	if err != nil {
		return errors.New("error in getting topics map" + err.Error())
	}
	for _, topic := range topics {
		topicId := -1
		for tId, t := range topicMap.(map[int]string) {
			if t == topic {
				topicId = tId
				topicIds = append(topicIds, tId)
				break
			}
		}
		if topicId == -1 {
			return errors.New("topic is not exist!")
		}
	}
	postId, err := db.InsertData("posts", userId, title, content, dt.Format("01-02-2006 15:04:05"))
	if err != nil {
		return errors.New("error in inserting post:" + err.Error())
	}
	for _, topicId := range topicIds {
		_, errPT := db.InsertData("PostTopics", postId, topicId)
		if errPT != nil {
			return errors.New("error in inserting postTopics:" + errPT.Error())
		}
	}
	return nil
}

/*
insertComment takes a username string, postId integer and content string as input, and
returns an error value. This function inserts the comment, represented by the collection of
inputs, to the database. The returned error value is nil if no error occurs whilst calling
other functions etc., otherwise the first error encountered is returned.
*/
func insertComment(userName string, postId int, content string) error {
	userId, err := getUserId(userName)
	if err != nil {
		return errors.New("error in getting userId:" + err.Error())
	}
	dt := time.Now()
	_, err = db.InsertData("comments", userId, postId, content, dt.Format("01-02-2006 15:04:05"))
	if err != nil {
		return errors.New("error in inserting comment:" + err.Error())
	}
	return nil
}

/*
insertReaction takes a username string, postId integer and reaction string as input, and
returns an error value. This function inserts the reaction, represented by the collection of
inputs, to the database. The returned error value is nil if no error occurs whilst calling
other functions etc., otherwise the first error encountered is returned.
*/
func insertReaction(userName string, postId int, commentId int, reaction string) error {
	userId, err := getUserId(userName)
	if err != nil {
		return errors.New("error in getting userId:" + err.Error())
	}
	// check like value
	_, err = db.InsertData("reactions", userId, postId, commentId, reaction)
	if err != nil {
		return errors.New("error in inserting reaction:" + err.Error())
	}
	return nil
}

/*
FillCommentStruct is a method for the (c *Comment) struct which takes db.Comment
as input and fills the Comment struct with the comment data from the database. It
also returns an error value, which is non-nil if there was an error in the database
query.
*/
func (c *Comment) FillCommentStruct(username string, dbComment db.Comment, r *http.Request) error {
	activeUsername := sessions.ActiveSessions.GetUsername(r) // Assign values to output struct
	c.CommentId = dbComment.CommentId
	c.PostId = dbComment.PostId
	c.Content = dbComment.Content
	c.CreationTime = dbComment.Time
	c.Username = username
	// Count likes and dislikes for the input comment if it has any
	reactions, err := db.SelectDataHandler("reactions", "commentId", c.CommentId)

	if err != nil && err.Error() != "data doesn't exist in reactions table" {
		return errors.New("error in getting reactions of comment from database" + err.Error())

	} else if err == nil {
		for _, reaction := range reactions.(map[int]db.Reaction) {
			if reaction.Reaction == "like" {
				c.Likes++
			} else if reaction.Reaction == "dislike" {
				c.Dislikes++
			}
		}
	}
	if activeUsername == "" {
		c.LikeStatus = 0
	} else {
		reStatus, err := findReactionStatus(c.Username, -1, c.CommentId)
		if err != nil {
			return errors.New("error in getting reaction status" + err.Error())
		}
		c.LikeStatus = reStatus
	}
	return nil
}

/*
FillPostStruct is a method for the (p *Post) struct which takes a db.Post as input
and populates the Post struct with the data from the database, including the topics,
comments, reactions and username associated with the post. It returns an error value,
which is non-nil if there was an error in the database query.
*/
func (p *Post) FillPostStruct(dbPost db.Post, r *http.Request) error {
	// Assign values to output struct
	p.PostId = dbPost.PostId
	p.CreationTime = dbPost.CreationTime
	p.Title = dbPost.Title
	p.Content = dbPost.Content

	// Retrieve all topics from database
	topicMap, err := db.SelectDataHandler("topics", "", nil)
	if err != nil {
		return errors.New("error in getting topics from database:" + err.Error())
	}

	// Retrieve topics associated with this post
	postTopics, errPT := db.SelectDataHandler("PostTopics", "postId", p.PostId)
	if errPT != nil {
		return errors.New("error in getting topics associated with post from database:" + errPT.Error())
	}
	// topicsId is a slice of topic IDs associated with this post
	topicsId := postTopics.(map[int][]int)[dbPost.PostId]
	for _, topicId := range topicsId {
		// Append topic name to topics slice in Post struct
		p.Topics = append(p.Topics, topicMap.(map[int]string)[topicId])
	}

	// Retreieve username of post creator
	userName, err := db.SelectDataHandler("users", "userId", dbPost.UserId)
	if err != nil {
		return errors.New("error in getting user name from database" + err.Error())
	}
	p.Username = userName.(db.User).Username

	// Retrieve all comments associated with this post
	comments, err := db.SelectDataHandler("comments", "postId", p.PostId)
	if err != nil && err.Error() != "data doesn't exist" {
		return errors.New("error in getting comments from database" + err.Error())
	} else if err == nil {
		for _, comment := range comments.(map[int]db.Comment) {
			var c Comment
			err = c.FillCommentStruct(p.Username, comment, r)
			if err != nil {
				return errors.New("error in casting comment" + err.Error())
			}
			p.TotalComments++
		}
	}

	// Count all reactions (likes and dislikes) for this post
	reactions, err := db.SelectDataHandler("reactions", "postId", p.PostId)
	// if error is not nil and error is not because there are no reactions
	if err != nil && err.Error() != "data doesn't exist in reactions table" {
		return errors.New("error in getting reactions of post from database" + err.Error())
	} else if err == nil {
		for _, reaction := range reactions.(map[int]db.Reaction) {
			if reaction.Reaction == "like" {
				p.Likes++
			} else if reaction.Reaction == "dislike" {
				p.Dislikes++
			}
		}
	}
	return nil
}

/*
GetMainDataStruct is a function which takes no input and harvests all posts and
topics from the database, filling the MainData struct with the data. It also
returns an error value, which is non-nil if there was an error in the database
query. This function is called in the handler for MainPage.
*/
func GetMainDataStruct(r *http.Request, userName string) (MainData, error) {
	// Initialise output struct
	md := MainData{Username: userName}

	// Retreive all posts and topics from the database
	posts, err := getAllPosts(r)
	if err != nil {
		return md, errors.New("error in getMainData in getting all posts :" + err.Error())
	}

	topics, err := getAllTopics()
	if err != nil {
		return md, errors.New("error in getMainData in getting all topics :" + err.Error())
	}

	// Check for message-cookie, add message to output struct if it exists
	cookie, err := r.Cookie(MESSAGE_COOKIE_NAME)
	if err == nil {
		md.CookieMessage = cookie.Value
	} else {
		md.CookieMessage = ""
	}

	// Compile output struct
	md.Posts = sortPostSlice(posts)     // Bubble sort posts in descending order of creation time
	md.Topics = sortStringSlice(topics) // Bubble sort topics in alphabetical order
	return md, nil
}

/*
FillPostDataStruct is a method of PostData struct which fills the struct with all topics
from the database. It returns an error value, which is non-nil if an error occurs whilst
calling other functions etc.
*/
func GetPostDataStruct(r *http.Request, userName string) (PostData, error) {
	// Initialise output struct
	pd := PostData{Username: userName}

	// Fill output struct with data
	allTopics, err := getAllTopics()
	if err != nil {
		return pd, errors.New("error in getting AllTopics:" + err.Error())
	}
	pd.AllTopics = allTopics

	// Check for message-cookie, add message to output struct if it exists
	cookie, err := r.Cookie(MESSAGE_COOKIE_NAME)
	if err == nil {
		pd.CookieMessage = cookie.Value
	} else {
		pd.CookieMessage = ""
	}

	return pd, nil
}

/*
GetContentDataStruct takes a htpp request, username and post ID strings as input
and returns a ContentData struct along with an error value. The ContentData struct
contains all the data required for the content page, which includes the post title,
creation time, content, and topics, as well as comments and reactions. The returned
error value is non-nil if an error was encountered in any of the database queries.
*/
func GetContentDataStruct(r *http.Request, userName string, postId int) (ContentData, error) {
	// Initialise output struct
	cd := ContentData{ActiveUsername: userName}
	// Get post data
	var post Post
	dbPost, err := db.SelectDataHandler("posts", "postId", postId)
	if err != nil {
		return cd, errors.New("error in getting post from database:" + err.Error())
	}
	for _, p := range dbPost.(map[int]db.Post) {
		err = post.FillPostStruct(p, r)
		if err != nil {
			return cd, errors.New("error in casting post to Post struct:" + err.Error())
		}
	}

	// Check for message-cookie, add message to output struct if it exists
	cookie, err := r.Cookie(MESSAGE_COOKIE_NAME)
	if err == nil {
		cd.CookieMessage = cookie.Value
	} else {
		cd.CookieMessage = ""
	}
	// Assign values to output struct
	cd.CreatorUsername = post.Username
	cd.Title = post.Title
	cd.Topics = sortStringSlice(post.Topics)
	cd.CreationTime = post.CreationTime
	cd.Content = post.Content
	cd.Likes = post.Likes
	cd.Dislikes = post.Dislikes
	if userName == "guest" {
		cd.LikeStatus = 0
	} else {
		reStatus, err := findReactionStatus(userName, postId, -1)
		if err != nil {
			return cd, errors.New("error in getting reaction status" + err.Error())
		}
		cd.LikeStatus = reStatus
	}

	// Retrieve comments for the post, and fill the comments slice of the
	// ContentData struct
	comments, err := db.SelectDataHandler("comments", "postId", postId)
	if err != nil && err.Error() != "data doesn't exist" {
		return cd, errors.New("error in getting comments from database" + err.Error())
	} else if err == nil {
		for _, comment := range comments.(map[int]db.Comment) {
			var c Comment
			commentUsername, err := getUserName(comment.UserId)
			if err != nil {
				return cd, errors.New("error in getting username from getUserName" + err.Error())
			}
			err = c.FillCommentStruct(commentUsername, comment, r)
			if err != nil {
				return cd, errors.New("error in casting comment" + err.Error())
			}
			cd.Comments = append(cd.Comments, c)
		}
	}
	// Sort comments in descending order of creation time
	cd.Comments = sortCommentSlice(cd.Comments)
	return cd, nil
}

/*
FillProfileDataStruct function description...
*/
func GetProfileDataStruct(r *http.Request, activeUsername, userName string) (ProfileData, error) {
	// Initialise output struct
	profile := ProfileData{ActiveUsername: activeUsername}
	profile.UserInfo.Username = userName
	// Check for message-cookie, add message to output struct if it exists
	cookie, err := r.Cookie(MESSAGE_COOKIE_NAME)
	if err == nil {
		profile.CookieMessage = cookie.Value
	} else {
		profile.CookieMessage = ""
	}

	// Get user data
	err = profile.UserInfo.FillUserStruct(userName)
	if err != nil {
		// Not sure what to do here, as returning the error is difficult to manage in the handler
	}

	// Populate Created Posts slice
	profile.CreatedPosts, err = (getCreatedPostSummaries(userName))
	if err != nil {
		// Not sure what to do here, as returning the error is difficult to manage in the handler
	}
	profile.CreatedPosts = sortSummarySlice(profile.CreatedPosts)

	// Populate Liked posts / comments slice
	profile.LikedPosts, err = getLikedPostSummaries(userName)
	if err != nil {
		// Not sure what to do here, as returning the error is difficult to manage in the handler
	}
	profile.LikedPosts = sortSummarySlice(profile.LikedPosts)

	// Populate Disliked posts / comments slice
	profile.DislikedPosts, err = getDislikedPostSummaries(userName)
	if err != nil {
		// Not sure what to do here, as returning the error is difficult to manage in the handler
	}
	profile.DislikedPosts = sortSummarySlice(profile.DislikedPosts)

	return profile, nil
}
