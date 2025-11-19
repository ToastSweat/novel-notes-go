package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Library struct {
	Bookcases    []Bookcase `json:"bookcases"`
	TotalScore   int        `json:"total_score"`
	LastRollOver string     `json:"last_rollover"` // "2025-11-13"
}

type Bookcase struct {
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	Shelves []Shelf `json:"shelves"`
}

type Shelf struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Books []Book `json:"books"`
}

type Book struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	CurrentDate string `json:"current_date"` // e.g. "2025-11-13"
	Pages       []Page `json:"pages"`        // all days, includes history
}

type Page struct {
	Date  string `json:"date"` // "YYYY-MM-DD"
	Items []Item `json:"items"`
}

type Item struct {
	ID          int    `json:"id"`
	Text        string `json:"text"`
	Completed   bool   `json:"completed"`
	CompletedAt string `json:"completed_at,omitempty"` // if done
}

func LoadLibrary(filename string) (Library, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		// If the file doesn't exist, return an empty Library (first run)
		if os.IsNotExist(err) {
			empty := Library{
				Bookcases:    []Bookcase{},
				TotalScore:   0,
				LastRollOver: "",
			}
			return empty, nil
		}
		// Some other error (permissions, etc.)
		return Library{}, err
	}

	var lib Library
	err = json.Unmarshal(data, &lib)
	if err != nil {
		return Library{}, err
	}

	return lib, nil
}

func NewBookcase(id int, name string) Bookcase {
	return Bookcase{
		ID:      id,
		Name:    name,
		Shelves: []Shelf{},
	}
}

func NewShelf(id int, name string) Shelf {
	return Shelf{
		ID:    id,
		Name:  name,
		Books: []Book{},
	}
}

func NewBook(id int, name string, currentDate string) Book {
	return Book{
		ID:          id,
		Name:        name,
		CurrentDate: currentDate,
		Pages:       []Page{},
	}
}

func NewItem(id int, text string) Item {
	return Item{
		ID:        id,
		Text:      text,
		Completed: false,
	}
}

func SaveLibrary(lib Library, filename string) error {
	// Turn the Library struct into JSON bytes
	data, err := json.MarshalIndent(lib, "", "  ")
	if err != nil {
		return err
	}

	// Write those bytes to a file
	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func PrintLibrarySummary(lib Library) {
	if len(lib.Bookcases) == 0 {
		fmt.Println("No bookcases yet.")
		return
	}

	for _, bc := range lib.Bookcases {
		fmt.Printf("Bookcase %d: %s\n", bc.ID, bc.Name)

		if len(bc.Shelves) == 0 {
			fmt.Println("  (no shelves)")
			continue
		}

		for _, sh := range bc.Shelves {
			fmt.Printf("  Shelf %d: %s\n", sh.ID, sh.Name)

			if len(sh.Books) == 0 {
				fmt.Println("    (no books)")
				continue
			}

			for _, bk := range sh.Books {
				fmt.Printf("    Book %d: %s (pages: %d)\n", bk.ID, bk.Name, len(bk.Pages))
			}
		}
	}
}

func NextBookcaseID(lib Library) int {
	maxID := 0
	for _, bc := range lib.Bookcases {
		if bc.ID > maxID {
			maxID = bc.ID
		}
	}
	return maxID + 1
}

func GetBookcaseByID(lib *Library, id int) *Bookcase {
	for i := range lib.Bookcases {
		if lib.Bookcases[i].ID == id {
			return &lib.Bookcases[i]
		}
	}
	return nil
}

func NextShelfID(bc *Bookcase) int {
	maxID := 0
	for _, sh := range bc.Shelves {
		if sh.ID > maxID {
			maxID = sh.ID
		}
	}
	return maxID + 1
}

func GetShelfByID(bc *Bookcase, id int) *Shelf {
	for i := range bc.Shelves {
		if bc.Shelves[i].ID == id {
			return &bc.Shelves[i]
		}
	}
	return nil
}

func NextBookID(sh *Shelf) int {
	maxID := 0
	for _, bk := range sh.Books {
		if bk.ID > maxID {
			maxID = bk.ID
		}
	}
	return maxID + 1
}

func GetBookByID(sh *Shelf, id int) *Book {
	for i := range sh.Books {
		if sh.Books[i].ID == id {
			return &sh.Books[i]
		}
	}
	return nil
}

func NextItemID(bk *Book) int {
	maxID := 0
	for _, pg := range bk.Pages {
		for _, it := range pg.Items {
			if it.ID > maxID {
				maxID = it.ID
			}
		}
	}
	return maxID + 1
}

func GetOrCreateTodayPage(bk *Book) *Page {
	today := time.Now().Format("2006-01-02")

	// Look for an existing page for today
	for i := range bk.Pages {
		if bk.Pages[i].Date == today {
			return &bk.Pages[i]
		}
	}

	// If not found, create one
	newPage := Page{
		Date:  today,
		Items: []Item{},
	}
	bk.Pages = append(bk.Pages, newPage)

	// Return pointer to the newly added page
	return &bk.Pages[len(bk.Pages)-1]
}

func PrintBookDetails(bk *Book) {
	fmt.Printf("Book %d: %s (current date: %s)\n", bk.ID, bk.Name, bk.CurrentDate)

	if len(bk.Pages) == 0 {
		fmt.Println("  (no pages)")
		return
	}

	for _, pg := range bk.Pages {
		fmt.Printf("  Page %s:\n", pg.Date)
		if len(pg.Items) == 0 {
			fmt.Println("    (no items)")
			continue
		}
		for _, it := range pg.Items {
			status := "[ ]"
			if it.Completed {
				status = "[x]"
			}
			fmt.Printf("    %s %d: %s\n", status, it.ID, it.Text)
		}
	}
}

func FindItemInBook(bk *Book, itemID int) (*Item, *Page) {
	for pi := range bk.Pages {
		for ii := range bk.Pages[pi].Items {
			if bk.Pages[pi].Items[ii].ID == itemID {
				return &bk.Pages[pi].Items[ii], &bk.Pages[pi]
			}
		}
	}
	return nil, nil
}

func PrintAllBooks(lib *Library) {
	if len(lib.Bookcases) == 0 {
		fmt.Println("No bookcases yet.")
		return
	}

	for _, bc := range lib.Bookcases {
		if len(bc.Shelves) == 0 {
			continue
		}
		for _, sh := range bc.Shelves {
			if len(sh.Books) == 0 {
				continue
			}
			for _, bk := range sh.Books {
				fmt.Printf(
					"Bookcase %d (%s) -> Shelf %d (%s) -> Book %d: %s (pages: %d)\n",
					bc.ID, bc.Name,
					sh.ID, sh.Name,
					bk.ID, bk.Name, len(bk.Pages),
				)
			}
		}
	}
}

func PrintBookHistory(bk *Book) {
	fmt.Printf("History for Book %d: %s\n", bk.ID, bk.Name)

	if len(bk.Pages) == 0 {
		fmt.Println("  (no pages yet)")
		return
	}

	currentDate := bk.CurrentDate
	hasHistory := false

	for _, pg := range bk.Pages {
		if pg.Date == currentDate {
			continue // skip current page, we only want history
		}
		hasHistory = true
		fmt.Printf("  Page %s:\n", pg.Date)
		if len(pg.Items) == 0 {
			fmt.Println("    (no items)")
			continue
		}
		for _, it := range pg.Items {
			status := "[ ]"
			if it.Completed {
				status = "[x]"
			}
			fmt.Printf("    %s %d: %s\n", status, it.ID, it.Text)
		}
	}

	if !hasHistory {
		fmt.Println("  (no history pages yet)")
	}
}

func main() {
	fmt.Println("Novel Notes", "Version 1.0")

	// 1. Load library from file (or get an empty one on first run)
	lib, err := LoadLibrary("novel_notes.json")
	if err != nil {
		fmt.Println("Error loading library:", err)
		return
	}

	// 2. If no bookcases yet, build some sample data (first run only)
	if len(lib.Bookcases) == 0 {
		fmt.Println("No bookcases found, creating sample data...")

		// Bookcase
		newBookcase := NewBookcase(1, "Test Bookcase")
		lib.Bookcases = append(lib.Bookcases, newBookcase)

		// Shelf
		newShelf := NewShelf(1, "Test Shelf")
		lib.Bookcases[0].Shelves = append(lib.Bookcases[0].Shelves, newShelf)

		// Book
		newBook := NewBook(1, "Test Book", "2025-11-13")
		lib.Bookcases[0].Shelves[0].Books = append(lib.Bookcases[0].Shelves[0].Books, newBook)

		// Items
		item1 := NewItem(1, "Write Novel Notes data model")
		item2 := NewItem(2, "Design checklist rollover logic")

		// Page
		page := Page{
			Date:  "2025-11-13",
			Items: []Item{item1, item2},
		}

		// Attach Page to the Book
		lib.Bookcases[0].Shelves[0].Books[0].Pages = append(
			lib.Bookcases[0].Shelves[0].Books[0].Pages,
			page,
		)
	}

	// 3. Read command-line arguments (skip program name)
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("Usage: go run . <command>")
		fmt.Println("Commands:")
		fmt.Println("  list    - Show bookcases, shelves, and books")
		return
	}

	switch args[0] {
	case "list":
		PrintLibrarySummary(lib)

	case "add-bookcase":
		if len(args) < 2 {
			fmt.Println("Usage: go run . add-bookcase <name>")
			return
		}

		name := strings.Join(args[1:], " ")
		newID := NextBookcaseID(lib)
		newBookcase := NewBookcase(newID, name)
		lib.Bookcases = append(lib.Bookcases, newBookcase)

		fmt.Printf("Added Bookcase %d: %s\n", newID, name)
		PrintLibrarySummary(lib)

	case "add-shelf":
		if len(args) < 3 {
			fmt.Println("Usage: go run . add-shelf <bookcase-id> <name>")
			return
		}

		// Parse bookcase ID (args[1])
		bookcaseID, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("Bookcase ID must be a number.")
			return
		}

		// Find the bookcase
		bc := GetBookcaseByID(&lib, bookcaseID)
		if bc == nil {
			fmt.Printf("No bookcase found with ID %d\n", bookcaseID)
			return
		}

		// Join the rest of the args into the shelf name
		name := strings.Join(args[2:], " ")

		newShelfID := NextShelfID(bc)
		newShelf := NewShelf(newShelfID, name)
		bc.Shelves = append(bc.Shelves, newShelf)

		fmt.Printf("Added Shelf %d to Bookcase %d: %s\n", newShelfID, bookcaseID, name)
		PrintLibrarySummary(lib)

	case "add-book":
		if len(args) < 4 {
			fmt.Println("Usage: go run . add-book <bookcase-id> <shelf-id> <name>")
			return
		}

		// Parse IDs
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

		// Find bookcase
		bc := GetBookcaseByID(&lib, bookcaseID)
		if bc == nil {
			fmt.Printf("No bookcase found with ID %d\n", bookcaseID)
			return
		}

		// Find shelf
		sh := GetShelfByID(bc, shelfID)
		if sh == nil {
			fmt.Printf("No shelf found with ID %d in bookcase %d\n", shelfID, bookcaseID)
			return
		}

		// Name is rest of the args
		name := strings.Join(args[3:], " ")

		newBookID := NextBookID(sh)
		currentDate := time.Now().Format("2006-01-02")
		newBook := NewBook(newBookID, name, currentDate)

		sh.Books = append(sh.Books, newBook)

		fmt.Printf("Added Book %d to Shelf %d in Bookcase %d: %s\n", newBookID, shelfID, bookcaseID, name)
		PrintLibrarySummary(lib)

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

		// Find bookcase
		bc := GetBookcaseByID(&lib, bookcaseID)
		if bc == nil {
			fmt.Printf("No bookcase found with ID %d\n", bookcaseID)
			return
		}

		// Find shelf
		sh := GetShelfByID(bc, shelfID)
		if sh == nil {
			fmt.Printf("No shelf found with ID %d in bookcase %d\n", shelfID, bookcaseID)
			return
		}

		// Find book
		bk := GetBookByID(sh, bookID)
		if bk == nil {
			fmt.Printf("No book found with ID %d in shelf %d\n", bookID, shelfID)
			return
		}

		// Item text (rest of args)
		text := strings.Join(args[4:], " ")

		newItemID := NextItemID(bk)
		item := NewItem(newItemID, text)

		// Make sure the book knows "today"
		today := time.Now().Format("2006-01-02")
		bk.CurrentDate = today

		// Add to today's page
		page := GetOrCreateTodayPage(bk)
		page.Items = append(page.Items, item)

		fmt.Printf("Added Item %d to Book %d: %s\n", newItemID, bookID, text)
		PrintLibrarySummary(lib)

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

		// Find bookcase
		bc := GetBookcaseByID(&lib, bookcaseID)
		if bc == nil {
			fmt.Printf("No bookcase found with ID %d\n", bookcaseID)
			return
		}

		// Find shelf
		sh := GetShelfByID(bc, shelfID)
		if sh == nil {
			fmt.Printf("No shelf found with ID %d in bookcase %d\n", shelfID, bookcaseID)
			return
		}

		// Find book
		bk := GetBookByID(sh, bookID)
		if bk == nil {
			fmt.Printf("No book found with ID %d in shelf %d\n", bookID, shelfID)
			return
		}

		// Print details
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

		// Find bookcase
		bc := GetBookcaseByID(&lib, bookcaseID)
		if bc == nil {
			fmt.Printf("No bookcase found with ID %d\n", bookcaseID)
			return
		}

		// Find shelf
		sh := GetShelfByID(bc, shelfID)
		if sh == nil {
			fmt.Printf("No shelf found with ID %d in bookcase %d\n", shelfID, bookcaseID)
			return
		}

		// Find book
		bk := GetBookByID(sh, bookID)
		if bk == nil {
			fmt.Printf("No book found with ID %d in shelf %d\n", bookID, shelfID)
			return
		}

		// Find item
		item, _ := FindItemInBook(bk, itemID)
		if item == nil {
			fmt.Printf("No item found with ID %d in book %d\n", itemID, bookID)
			return
		}

		if item.Completed {
			fmt.Printf("Item %d is already completed.\n", itemID)
			return
		}

		// Mark as completed
		item.Completed = true
		item.CompletedAt = time.Now().Format(time.RFC3339)

		// Increase total score
		lib.TotalScore++

		fmt.Printf("Completed Item %d in Book %d: %s\n", itemID, bookID, item.Text)
		fmt.Printf("New total score: %d\n", lib.TotalScore)

		// Optional: show the book again
		PrintBookDetails(bk)

	case "list-books":
		PrintAllBooks(&lib)

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

		// Find bookcase
		bc := GetBookcaseByID(&lib, bookcaseID)
		if bc == nil {
			fmt.Printf("No bookcase found with ID %d\n", bookcaseID)
			return
		}

		// Find shelf
		sh := GetShelfByID(bc, shelfID)
		if sh == nil {
			fmt.Printf("No shelf found with ID %d in bookcase %d\n", shelfID, bookcaseID)
			return
		}

		// Find book
		bk := GetBookByID(sh, bookID)
		if bk == nil {
			fmt.Printf("No book found with ID %d in shelf %d\n", bookID, shelfID)
			return
		}

		PrintBookHistory(bk)
		fmt.Printf("Total score: %d\n", lib.TotalScore)

	default:
		fmt.Println("Unknown command:", args[0])
		fmt.Println("Available commands: list, list-books, add-bookcase, add-shelf, add-book, add-item, view-book, complete-item, view-history")
	}

	// 4. Save back to file (even if nothing changed, it's fine)
	err = SaveLibrary(lib, "novel_notes.json")
	if err != nil {
		fmt.Println("Error saving library:", err)
	} else {
		fmt.Println("Library saved to novel_notes.json")
	}
}
