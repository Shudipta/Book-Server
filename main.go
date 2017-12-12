package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	//"io/ioutil"
)

type Book struct {
	Id		int `json:"Id"`
	Title   string `json:"Title"`
	Author	string `json:"Author"`
	Edition	int `json:"Edition"`
}

var books []Book

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the \"Book Server\"")
	fmt.Println("Endpoint Hit: \"hello\" page")
}

func showBookList(w http.ResponseWriter, r *http.Request) {
	response, err := json.MarshalIndent(books, "", " ")
	if err != nil {
		fmt.Fprintln(w, "Error occured in converting into json is %s", err)
		fmt.Println("Error occured in converting into json is", err)
		return
	}
	fmt.Fprintln(w, string(response))
}

func addBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	if r.Method == "GET" {
		//book, titleError := r.URL.Query()
		fmt.Println(r.URL)
		fmt.Println(r.URL.Query())
		fmt.Println(r.URL.Query()["Edition"] == nil)
		//fmt.Println(r.URL.Query()["Edition"][0])
		data := r.URL.Query()
		if data["Title"] == nil || data["Author"] == nil || data["Edition"] == nil {
			fmt.Fprintf(w, "contains empty field")
			fmt.Println("contains empty field")
			return
		}

		ed, edErr := strconv.Atoi(data["Edition"][0])
		if edErr != nil {
			fmt.Fprintf(w, "error while converting edition")
			return
		}

		id := len(books) + 1
		book = Book{Id: id, Title: data["Title"][0], Author: data["Author"][0], Edition: ed}
		//books = append(books, book)

	} else if r.Method == "POST" {
		err := json.NewDecoder(r.Body).Decode(&book)
		defer r.Body.Close()

		//fmt.Println(book)
		if err != nil {
			fmt.Println(err, "error getting json data in POST method")
			fmt.Fprintf(w,"error getting json data in POST method")
			return
		}
		if book.Title == "" || book.Author == "" {
			fmt.Fprintf(w, "contains empty field")
			fmt.Println("contains empty field")
			return
		}
		book.Id = len(books) + 1
		//fmt.Println(book)

	}

	books = append(books, book)
	fmt.Println("Endpoint Hit: \"addBook\" page")

	addResponsErr := json.NewEncoder(w).Encode(http.Response{StatusCode: http.StatusAccepted, Status: "added succesfully"})
	if addResponsErr != nil {
		fmt.Fprintf(w, "response error in adding new book")
		fmt.Println("response error in adding new book")
		return
	}
}

func editBook(w http.ResponseWriter, r *http.Request) {
	var book Book

	if r.Method == "PUT" {
		newDecoder := json.NewDecoder(r.Body)
		defer r.Body.Close()
		fmt.Println(newDecoder)
		err := newDecoder.Decode(&book)

		fmt.Println(book)
		if err != nil {
			fmt.Println(err, "error getting json data in PUT method")
			fmt.Fprintf(w,"error getting json data in PUT method")
			return
		}


		for i, _ := range books{
			if book.Id == books[i].Id {
				books[i] = book
				//log.Fatal(json.NewEncoder(w).Encode(http.Response{StatusCode: http.StatusAccepted, Status: "updated succesfully"}))
				updateResponsErr := json.NewEncoder(w).Encode(http.Response{StatusCode: http.StatusAccepted, Status: "updated succesfully"})
				if updateResponsErr != nil {
					fmt.Fprintf(w, "response error in updating new book")
					fmt.Println("response error in updating new book")
					return
				}
				return
			}
		}
		json.NewEncoder(w).Encode(http.Response{
			StatusCode: http.StatusBadRequest,
			Status: "Bad request",
		})

		//fmt.Fprintf(w, "response error in updating book")
		//fmt.Println("response error in updating book")
	}
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	var book Book

	if r.Method == "DELETE" {
		newDecoder := json.NewDecoder(r.Body)
		defer r.Body.Close()
		fmt.Println(newDecoder)
		err := newDecoder.Decode(&book)

		fmt.Println(book)
		if err != nil {
			fmt.Println(err, "error getting json data in PUT method")
			fmt.Fprintf(w,"error getting json data in PUT method")
			return
		}


		for i, _ := range books{
			if book.Id == books[i].Id {
				tmpBook := books[i]
				books = append(books[:i], books[i+1:]...)
				//log.Fatal(json.NewEncoder(w).Encode(http.Response{StatusCode: http.StatusAccepted, Status: "updated succesfully"}))
				updateResponsErr := json.NewEncoder(w).Encode(http.Response{StatusCode: http.StatusAccepted, Status: "deleted succesfully"})
				if updateResponsErr != nil {
					fmt.Fprintf(w, "response error in deleting book", tmpBook)
					fmt.Println("response error in deleting book", tmpBook)
				}
				return
			}
		}
		json.NewEncoder(w).Encode(http.Response{
			StatusCode: http.StatusBadRequest,
			Status: "Bad request",
		})

		//fmt.Fprintf(w, "response error in updating book")
		//fmt.Println("response error in updating book")
	}
}



func handleRequests() {
	http.HandleFunc("/", showBookList)
	http.HandleFunc("/addBook", addBook)
	http.HandleFunc("/editBook", editBook)
	http.HandleFunc("/deleteBook", deleteBook)
	log.Fatal(http.ListenAndServe(":10000", nil))
}

func main() {
	handleRequests()
}