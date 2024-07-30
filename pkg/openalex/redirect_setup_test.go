package openalex

import (
	"fmt"
	"log"
	"testing"
)

func TestRedirectSetup(t *testing.T) {
	book1 := Book{
		ID:          "1",
		Title:       "Learn Go with Examples",
		PublishDate: "2023-10-10",
		RedirectTo:  "",
	}

	book2 := Book{
		ID:          "2",
		Title:       "Dune 2",
		PublishDate: "2023-10-10",
		RedirectTo:  "1",
	}

	book3 := Book{
		ID:          "3",
		Title:       "Java for dummies",
		PublishDate: "2023-10-10",
		RedirectTo:  "2",
	}

	addBook(book1)
	addBook(book2)
	addBook(book3)

	id := "3"
	book, err := findBookByID(esClient, id)
	if err != nil {
		log.Fatalf("Error finding the book: %s", err)
	}

	if book != nil {
		fmt.Printf("Found Book: %+v\n", book)
		if book.RedirectTo == "" {
			fmt.Println("The 'RedirectTo' field is empty.")
		} else {
			fmt.Printf("The 'RedirectTo' field is: %s\n", book.RedirectTo)
		}
	} else {
		fmt.Println("No book found with the given title.")
	}
}
