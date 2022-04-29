package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

const (
	DB_USER     = "postgres"
	DB_PASSWORD = "123456."
	DB_NAME     = "dbTest"
)

// DB set up
func setupDB() *sql.DB {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)

	checkErr(err)

	return db
}

type Post struct {
	Postid      int    `json:"id"`
	Postsub     string `json:"subject"`
	Postcontent string `json:"content"`
}
type JsonResponse struct {
	Type    string `json:"type"`
	Data    []Post `json:"data"`
	Message string `json:"message"`
}

// Main function
func main() {

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {

			LoginCheck(w, r)
		}

	})

	http.HandleFunc("/posts", func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "GET" {

			GetPosts(w, r)

		} else if r.Method == "POST" {

			CreatePost(w, r)

		} else if r.Method == "DELETE" {

			RemovePost(w, r)

		} else if r.Method == "PUT" {

			UpdatePost(w, r)

		}

	})

	//http.HandleFunc("/posts/remove", RemovePost)
	//http.HandleFunc("/posts/update", UpdatePost)

	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

// Function for handling messages
func printMessage(message string) {
	fmt.Println("")
	fmt.Println(message)
	fmt.Println("")
}

// Function for handling errors
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

// Get all movies

// response and request handlers
func GetPosts(w http.ResponseWriter, r *http.Request) {

	db := setupDB()

	printMessage("Getting posts...")

	// Get all movies from movies table that don't have movieID = "1"
	rows, err := db.Query("SELECT * FROM blog")

	// check errors
	checkErr(err)

	// var response []JsonResponse
	var posts []Post

	// Foreach movie
	for rows.Next() {
		var id int
		var subject string
		var content string

		err = rows.Scan(&id, &subject, &content)

		// check errors
		checkErr(err)

		posts = append(posts, Post{Postid: id, Postsub: subject, Postcontent: content})
	}

	var response = JsonResponse{Type: "success", Data: posts}

	json.NewEncoder(w).Encode(response)

}

func CreatePost(w http.ResponseWriter, r *http.Request) {

	postID := r.FormValue("postid")
	postsub := r.FormValue("subject")
	postcontent := r.FormValue("content")

	var response = JsonResponse{}

	if postID == "" || postsub == "" || postcontent == "" {
		response = JsonResponse{Type: "error", Message: "You are missing postid or subject or content parameter."}
		//response = JsonResponse{Type: "error", Message: postID}
	} else {
		db := setupDB()

		printMessage("Inserting post into DB")

		fmt.Println("Inserting new post with ID: " + postID + " and subject: " + postsub)

		var lastInsertID int
		err := db.QueryRow("INSERT INTO blog(id, subject, content) VALUES($1, $2, $3) returning id;", postID, postsub, postcontent).Scan(&lastInsertID)

		// check errors
		checkErr(err)

		response = JsonResponse{Type: "success", Message: "The post has been inserted successfully!"}
	}

	json.NewEncoder(w).Encode(response)

}

// Delete a post

// response and request handlers
func RemovePost(w http.ResponseWriter, r *http.Request) {
	postID := r.FormValue("postid")

	var response = JsonResponse{}

	if postID == "" {
		response = JsonResponse{Type: "error", Message: "You are missing postID parameter."}
	} else {
		db := setupDB()

		printMessage("Deleting movie from DB")

		_, err := db.Exec("DELETE FROM blog where id = $1", postID)

		// check errors
		checkErr(err)

		response = JsonResponse{Type: "success", Message: "The movie has been deleted successfully!"}
	}

	json.NewEncoder(w).Encode(response)
}

func UpdatePost(w http.ResponseWriter, r *http.Request) {

	postID := r.FormValue("postid")
	postsub := r.FormValue("subject")
	postcontent := r.FormValue("content")

	var response = JsonResponse{}

	if postID == "" || postsub == "" || postcontent == "" {
		response = JsonResponse{Type: "error", Message: "You are missing postid or subject or content parameter."}
	} else {
		db := setupDB()

		printMessage("Updating post from DB")

		_, err := db.Exec("UPDATE blog SET subject = $2, content = $3 where (id = $1);", postID, postsub, postcontent)

		// check errors
		checkErr(err)

		response = JsonResponse{Type: "success", Message: "The post has been update successfully!"}
	}

	json.NewEncoder(w).Encode(response)
}

func LoginCheck(w http.ResponseWriter, r *http.Request) {

	userName := r.FormValue("username")
	passWord := r.FormValue("password")

	var response = JsonResponse{}

	if userName == "" || passWord == "" {
		response = JsonResponse{Type: "error", Message: "You are missing username or password parameter."}
		//response = JsonResponse{Type: "error", Message: postID}
	} else {
		db := setupDB()

		printMessage("login Cheking ...")

		rows, err := db.Query("SELECT password FROM users WHERE username = $1;", userName)

		for rows.Next() {
			var passwd string

			err = rows.Scan(&passwd)

			// check errors
			checkErr(err)

			if passwd == passWord {

				response = JsonResponse{Type: "success", Message: "successfully loged in"}
			} else {

				response = JsonResponse{Type: "no success", Message: "The username or Password is incorrect!!!"}

			}
		}

		// check errors
		checkErr(err)

	}

	json.NewEncoder(w).Encode(response)

}
