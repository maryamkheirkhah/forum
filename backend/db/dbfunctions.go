package db

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/s"
	"log"
	"path/filepath"
	"runtime"

	_ "github.com/mattn/go-sqlite3"
)

/*
selectData gets the whole query as a string and the arguments as many as desired in any type
since at each table we have a different number of arguments and different types for each column of tables
and returns the specific data rows of table info and error if exist
*/
func selectData(myQuery string, args ...any) (*sql.Rows, error) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	database, errOpen := sql.Open("sqlite3", basepath+"/forum.db")
	if errOpen != nil {
		return nil, errOpen
	}
	defer database.Close()

	rows, err := database.Query(myQuery, args...)
	if err != nil {
		return nil, err
	}
	if !rows.Next() {
		return nil, errors.New("data does not exist")
	}
	return rows, nil
}

/*SelectDataHandler gets the table name and key of desired data as any type to handle all different kinds of data.
it writes a query base on the table name and the key argument of that table.
it calls the selectData function to get the rows of the table related to the specific data
and returns data in those rows or error if exist.
*/
func SelectDataHandler(tableName string, keyName string, keyValue any) (any, error) {
	myQuery := "SELECT * FROM " + tableName + " WHERE " + keyName + "= ?"

	switch tableName {
	case "users":
		//myQuery = "SELECT userId, email, pass, creationTime FROM users WHERE " + keyName + "= ?"
		rows, err := selectData(myQuery, keyValue)
		if err != nil {
			return nil, err
		}
		u := s.User{Username: keyValue.(string)}
		for rows.Next() {
			rows.Scan(&u.UserId, &u.Email, &u.Pass, &u.Time)
		}
		return u, nil
	case "posts":
		rows, err := selectData(myQuery, keyValue)
		if err != nil {
			return nil, err
		}
		post := s.Post{}
		for rows.Next() {
			rows.Scan(&post.PostId, &post.UserId, &post.Title, &post.Content, &post.CreationTime)
		}
		return post, nil
	case "topics":
		rows, err := selectData(myQuery, keyValue)
		if err != nil {
			return nil, err
		}
		var topic string
		for rows.Next() {
			rows.Scan(&topic)
		}
		return topic, nil
	case "comments":
		rows, err := selectData(myQuery, keyValue)
		if err != nil {
			return nil, err
		}
		comment := s.Comment{}
		for rows.Next() {
			rows.Scan(&comment.UserId, &comment.PostId, &comment.Content, &comment.Time)
		}
		return comment, nil

	case "reactions":
		rows, err := selectData(myQuery, keyValue)
		if err != nil {
			return nil, err
		}
		reaction := s.Reaction{}
		for rows.Next() {
			rows.Scan(&reaction.UserId, &reaction.PostId, &reaction.CommentId, &reaction.Reaction)
		}
		return reaction, nil

	}
	return nil, nil
}

/*
CheckData, using SelectDataHandler function, check if a specific data exist in a table then it returns error otherwise it returns nil
*/
func NotExistData(tableName string, keyName string, key any) error {
	_, err := SelectDataHandler(tableName, keyName, key)
	if err != nil {
		return nil
	}
	return errors.New("data already exist")

}

/*
InsertData gets the table name and the arguments as many as desired in any type to handle all different kinds of data.
it calls the notExistData for checking the existing of this data if its not exist,
it writes a query base on the table name and the arguments of that table. insert new data to the table.
and returns error if exist.
*/
func InsertData(tableName string, args ...any) error {
	var myQuery string
	switch tableName {
	case "users":
		err := NotExistData("users", "userName", args[0])
		if err != nil {
			return errors.New("username already exist")
		}
		myQuery = "INSERT INTO users(userName, email, pass,creationTime) VALUES(?,?,?,?)"
	case "posts":
		err := NotExistData("users", "userName", args[0])
		if err == nil {
			return errors.New("user does not exist")
		}
		myQuery = "INSERT INTO posts(userId, title, content,creationTime) VALUES(?,?,?,?)"
	case "topics":
		err := NotExistData("topics", "topicName", args[0])
		if err != nil {
			return errors.New("topic already exist")
		}
		myQuery = "INSERT INTO topics(topicName) VALUES(?)"
	case "comments":
		//check user existence
		err := NotExistData("users", "userName", args[0])
		if err == nil {
			return errors.New("user does not exist")
		}
		//check post existence
		err = NotExistData("posts", "postId", args[1])
		if err == nil {
			return errors.New("post does not exist")
		}
		myQuery = "INSERT INTO comments(userId, postId, content,creationTime) VALUES(?,?,?,?)"
	case "reactions":
		err := NotExistData("users", "userName", args[0])
		if err == nil {
			return errors.New("user does not exist")
		}
		//check post existence
		err = NotExistData("posts", "postId", args[1])
		if err == nil {
			return errors.New("post does not exist")
		}
		err = NotExistData("comments", "commentId", args[2])
		if err == nil {
			return errors.New("comment does not exist")
		}
		myQuery = "INSERT INTO reactions(userId, postId, commentId,reaction) VALUES(?,?,?,?)"

	case "PostTopics":
		// add checking later!
		myQuery = "INSERT INTO PostTopics(postId, topicId) VALUES(?,?)"
	}
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	database, errOpen := sql.Open("sqlite3", basepath+"/forum.db")
	if errOpen != nil {
		log.Fatal(errOpen)
	}
	defer database.Close()
	statement, err := database.Prepare(myQuery)
	if err != nil {
		return err
	}
	statement.Exec(args...)
	return nil
}

/*
UpdateData get the table name as a string to identify which type of data is going to update
finding it in case of user want to change its username and other user's info
if user doesn't exist or any other problem return error
*/

func UpdateData(table string, key string, args ...any) error { //// should be update
	/* if CheckUserName(user.Username) {
		return errors.New("username already taken")
	} */
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	database, errOpen := sql.Open("sqlite3", basepath+"/forum.db")
	if errOpen != nil {
		return errOpen
	}
	var myQuery string
	switch table {
	case "users":
		myQuery = "UPDATE users SET userName=?, email=?, pass=? where userName=?"
	}
	statement, err := database.Prepare(myQuery)
	if err != nil {
		return err
	}
	res, err := statement.Exec(args, key)
	if err != nil {
		return err
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affect == 0 {
		return errors.New("User Not Found")
	}
	return nil
}

/*
DeleteData takes the table name as a string to identify which type of data to delete and from which table
and id to find that desired item. Then, after checking the existence of data, based on the table name,
it adds the arguments of the table in a query to delete that item.
*/
func DeleteData(tableName string, keyValue string) error {
	var key string
	switch tableName {
	case "users":
		key = "userName"
	case "posts":
		key = "postId"
	case "comments":
		key = "commentId"
	case "topics":
		key = "topicId"
	case "reactions":
		key = "reactionId"
	default:
		return errors.New("table does not exist")
	}
	err := NotExistData(tableName, key, keyValue)
	if err == nil {
		return errors.New("data does not exist")
	}
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	database, errOpen := sql.Open("sqlite3", basepath+"/forum.db")
	if errOpen != nil {
		return errOpen
	}

	myQuery := "DELETE from " + tableName + " WHERE " + key + "=?"
	statement, PrepareErr := database.Prepare(myQuery)
	if PrepareErr != nil {
		return PrepareErr
	}
	result, errSt := statement.Exec(keyValue)
	if errSt != nil {
		return errSt
	}
	affect, errRow := result.RowsAffected()
	if errRow != nil {
		return errRow
	}
	if affect == 0 {
		return errors.New("item not found")
	}
	//for testing
	fmt.Println(affect)

	return nil
}
