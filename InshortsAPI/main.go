package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

var mut sync.Mutex

type Article struct {
	Id           string    `json:Id`
	Title        string    `json:"Title"`
	SubTitle     string    `json:"SubTitle"`
	Content      string    `json:"content"`
	CreationTime time.Time `json:"creation"`
}

var Articles []Article

func createNewArticle(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		//add new article
		mut.Lock()
		var newArticle Article
		reqBody, err := ioutil.ReadAll(r.Body)
		//if error found
		if err != nil {
			fmt.Fprintf(w, "Enter details of new article")
		}
		//calculate newId for the article
		newID, err := strconv.Atoi(Articles[len(Articles)-1].Id)

		json.Unmarshal(reqBody, &newArticle)
		newArticle.CreationTime = time.Now()
		newArticle.Id = strconv.Itoa(newID + 1)

		//add new article to existing Articles
		Articles = append(Articles, newArticle)
		fmt.Fprintf(w, "Article created!")
		json.NewEncoder(w).Encode(newArticle)
		mut.Unlock()
	} else {
		//return all articles
		fmt.Println("Endpoint Hit: returnAllArticles")
		json.NewEncoder(w).Encode(Articles)
	}

}

func getOneArticle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: getOneArticle")
	//parse the string to obtain Id
	url := (r.URL.Path)
	key := strings.Split(url, "/")[2]
	//loop over Articles to find for Id
	for _, article := range Articles {
		if article.Id == key {
			json.NewEncoder(w).Encode(article)
			fmt.Println("Article found!")
			break
		}
	}

}

func searchArticle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: searchArticle")
	key := r.URL.Query()["q"][0]
	//loop over articles to find for query

	for _, article := range Articles {
		if article.Title == key {
			json.NewEncoder(w).Encode(article)
			fmt.Println("Article found!")
		} else if article.SubTitle == key {
			json.NewEncoder(w).Encode(article)
			fmt.Println("Article found!")
		} else if article.Content == key {
			json.NewEncoder(w).Encode(article)
			fmt.Println("Article found!")
		}
	}

}

func handleRequests() {

	http.HandleFunc("/", getOneArticle)
	http.HandleFunc("/articles", createNewArticle)
	http.HandleFunc("/articles/search", searchArticle)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func main() {
	fmt.Println("API")
	Articles = []Article{
		Article{
			Id:           "1",
			Title:        "Hello",
			SubTitle:     "Hello",
			Content:      "Article content",
			CreationTime: time.Now(),
		},
		Article{
			Id:           "2",
			Title:        "Hello 2",
			SubTitle:     "Hello",
			Content:      "Article",
			CreationTime: time.Now(),
		},
	}
	handleRequests()
}
