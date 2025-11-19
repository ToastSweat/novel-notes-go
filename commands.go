package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func printHelp() {
	fmt.Println("Usage: go run . <command>")
	fmt.Println("Commands:")
	fmt.Println("  list                        - Show bookcases, shelves, and books")
	fmt.Println("  list-books                  - Show all books in the library")
	fmt.Println("  add-bookcase NAME")
	fmt.Println("  add-shelf BOOKCASE_ID NAME")
	fmt.Println("  add-book BOOKCASE_ID SHELF_ID NAME")
	fmt.Println("  add-item BOOKCASE_ID SHELF_ID BOOK_ID TEXT")
	fmt.Println("  view-book BOOKCASE_ID SHELF_ID BOOK_ID")
	fmt.Println("  complete-item BOOKCASE_ID SHELF_ID BOOK_ID ITEM_ID")
	fmt.Println("  view-history BOOKCASE_ID SHELF_ID BOOK_ID")
}

func handleCommand(lib *Library, args []string) {
	if len(args) == 0 {
		printHelp()
		return
	}

	switch args[0] {
	case "list":
		PrintLibrarySummary(*lib)

	case "list-books":
		PrintAllBooks(lib)

	case "add-bookcase":
		if len(args) < 2 {
			fmt.Println("Usage: go run . add-bookcase <name>")
			return
		}

		name := strings.Join(args[1:], " ")
		newID := NextBookcaseID(*lib)
		newBookcase := NewBookcase(newID, name)
		lib.Bookcases = append(lib.Bookcases, newBookcase)

		fmt.Printf("Added Bookcase %d: %s\n", newID, name)
		PrintLibrarySummary(*lib)

	case "add-shelf":
		if len(args) < 3 {
			fmt.Println("Usage: go run . add-shelf <bookcase-id> <name>")
			return
		}

		bookcaseID, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("Bookcase ID must be a number.")
			return
		}

		bc := GetBookcaseByID(lib, bookcaseID)
		if bc == nil {
			fmt.Printf("No bookcase found with ID %d\n", bookcaseID)
			return
		}

		name := strings.Join(args[2:], " ")

		newShelfID := NextShelfID(bc)
		newShelf := NewShelf(newShelfID, name)
		bc.Shelves = append(bc.Shelves, newShelf)

		fmt.Printf("Added Shelf %d to Bookcase %d: %s\n", newShelfID, bookcaseID, name)
		PrintLibrarySummary(*lib)

	case "add-book":
		if len(args) < 4 {
			fmt.Println("Usage: go run . add-book <bookcase-id> <shelf-id> <name>")
			return
		}

		bookcaseID, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("Bookcase ID must be a number.")
			return
		}

		shelfID, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Println("Shelf ID must be a number.")
			return
		}

		bc := GetBookcaseByID(lib, bookcaseID)
		if bc == nil {
			fmt.Printf("No bookcase found with ID %d\n", bookcaseID)
			return
		}

		sh := GetShelfByID(bc, shelfID)
		if sh == nil {
			fmt.Printf("No shelf found with ID %d in bookcase %d\n", shelfID, bookcaseID)
			return
		}

		name := strings.Join(args[3:], " ")

		newBookID := NextBookID(sh)
		currentDate := time.Now().Format("2006-01-02")
		newBook := NewBook(newBookID, name, currentDate)

		sh.Books = append(sh.Books, newBook)

		fmt.Printf("Added Book %d to Shelf %d in Bookcase %d: %s\n", newBookID, shelfID, bookcaseID, name)
		PrintLibrarySummary(*lib)

	case "add-item":
		if len(args) < 5 {
			fmt.Println("Usage: go run . add-item <bookcase-id> <shelf-id> <book-id> <text>")
			return
		}

		bookcaseID, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("Bookcase ID must be a number.")
			return
		}

		shelfID, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Println("Shelf ID must be a number.")
			return
		}

		bookID, err := strconv.Atoi(args[3])
		if err != nil {
			fmt.Println("Book ID must be a number.")
			return
		}

		bc := GetBookcaseByID(lib, bookcaseID)
		if bc == nil {
			fmt.Printf("No bookcase found with ID %d\n", bookcaseID)
			return
		}

		sh := GetShelfByID(bc, shelfID)
		if sh == nil {
			fmt.Printf("No shelf found with ID %d in bookcase %d\n", shelfID, bookcaseID)
			return
		}

		bk := GetBookByID(sh, bookID)
		if bk == nil {
			fmt.Printf("No book found with ID %d in shelf %d\n", bookID, shelfID)
			return
		}

		text := strings.Join(args[4:], " ")

		newItemID := NextItemID(bk)
		item := NewItem(newItemID, text)

		today := time.Now().Format("2006-01-02")
		bk.CurrentDate = today

		page := GetOrCreateTodayPage(bk)
		page.Items = append(page.Items, item)

		fmt.Printf("Added Item %d to Book %d: %s\n", newItemID, bookID, text)
		PrintLibrarySummary(*lib)

	case "view-book":
		if len(args) < 4 {
			fmt.Println("Usage: go run . view-book <bookcase-id> <shelf-id> <book-id>")
			return
		}

		bookcaseID, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("Bookcase ID must be a number.")
			return
		}

		shelfID, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Println("Shelf ID must be a number.")
			return
		}

		bookID, err := strconv.Atoi(args[3])
		if err != nil {
			fmt.Println("Book ID must be a number.")
			return
		}

		bc := GetBookcaseByID(lib, bookcaseID)
		if bc == nil {
			fmt.Printf("No bookcase found with ID %d\n", bookcaseID)
			return
		}

		sh := GetShelfByID(bc, shelfID)
		if sh == nil {
			fmt.Printf("No shelf found with ID %d in bookcase %d\n", shelfID, bookcaseID)
			return
		}

		bk := GetBookByID(sh, bookID)
		if bk == nil {
			fmt.Printf("No book found with ID %d in shelf %d\n", bookID, shelfID)
			return
		}

		PrintBookDetails(bk)
		fmt.Printf("Total score: %d\n", lib.TotalScore)

	case "complete-item":
		if len(args) < 5 {
			fmt.Println("Usage: go run . complete-item <bookcase-id> <shelf-id> <book-id> <item-id>")
			return
		}

		bookcaseID, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("Bookcase ID must be a number.")
			return
		}

		shelfID, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Println("Shelf ID must be a number.")
			return
		}

		bookID, err := strconv.Atoi(args[3])
		if err != nil {
			fmt.Println("Book ID must be a number.")
			return
		}

		itemID, err := strconv.Atoi(args[4])
		if err != nil {
			fmt.Println("Item ID must be a number.")
			return
		}

		bc := GetBookcaseByID(lib, bookcaseID)
		if bc == nil {
			fmt.Printf("No bookcase found with ID %d\n", bookcaseID)
			return
		}

		sh := GetShelfByID(bc, shelfID)
		if sh == nil {
			fmt.Printf("No shelf found with ID %d in bookcase %d\n", shelfID, bookcaseID)
			return
		}

		bk := GetBookByID(sh, bookID)
		if bk == nil {
			fmt.Printf("No book found with ID %d in shelf %d\n", bookID, shelfID)
			return
		}

		item, _ := FindItemInBook(bk, itemID)
		if item == nil {
			fmt.Printf("No item found with ID %d in book %d\n", itemID, bookID)
			return
		}

		if item.Completed {
			fmt.Printf("Item %d is already completed.\n", itemID)
			return
		}

		item.Completed = true
		item.CompletedAt = time.Now().Format(time.RFC3339)

		lib.TotalScore++

		fmt.Printf("Completed Item %d in Book %d: %s\n", itemID, bookID, item.Text)
		fmt.Printf("New total score: %d\n", lib.TotalScore)
		PrintBookDetails(bk)

	case "view-history":
		if len(args) < 4 {
			fmt.Println("Usage: go run . view-history <bookcase-id> <shelf-id> <book-id>")
			return
		}

		bookcaseID, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("Bookcase ID must be a number.")
			return
		}

		shelfID, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Println("Shelf ID must be a number.")
			return
		}

		bookID, err := strconv.Atoi(args[3])
		if err != nil {
			fmt.Println("Book ID must be a number.")
			return
		}

		bc := GetBookcaseByID(lib, bookcaseID)
		if bc == nil {
			fmt.Printf("No bookcase found with ID %d\n", bookcaseID)
			return
		}

		sh := GetShelfByID(bc, shelfID)
		if sh == nil {
			fmt.Printf("No shelf found with ID %d in bookcase %d\n", shelfID, bookcaseID)
			return
		}

		bk := GetBookByID(sh, bookID)
		if bk == nil {
			fmt.Printf("No book found with ID %d in shelf %d\n", bookID, shelfID)
			return
		}

		PrintBookHistory(bk)
		fmt.Printf("Total score: %d\n", lib.TotalScore)

	default:
		fmt.Println("Unknown command:", args[0])
		printHelp()
	}
}
