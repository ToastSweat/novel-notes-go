package main

// Core data structures

type Library struct {
	Bookcases    []Bookcase `json:"bookcases"`
	TotalScore   int        `json:"total_score"`
	LastRollOver string     `json:"last_rollover"` // "YYYY-MM-DD"
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

// "Constructor" helpers

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

// ID helpers

func NextBookcaseID(lib Library) int {
	maxID := 0
	for _, bc := range lib.Bookcases {
		if bc.ID > maxID {
			maxID = bc.ID
		}
	}
	return maxID + 1
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

func NextBookID(sh *Shelf) int {
	maxID := 0
	for _, bk := range sh.Books {
		if bk.ID > maxID {
			maxID = bk.ID
		}
	}
	return maxID + 1
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

// Lookup helpers

func GetBookcaseByID(lib *Library, id int) *Bookcase {
	for i := range lib.Bookcases {
		if lib.Bookcases[i].ID == id {
			return &lib.Bookcases[i]
		}
	}
	return nil
}

func GetShelfByID(bc *Bookcase, id int) *Shelf {
	for i := range bc.Shelves {
		if bc.Shelves[i].ID == id {
			return &bc.Shelves[i]
		}
	}
	return nil
}

func GetBookByID(sh *Shelf, id int) *Book {
	for i := range sh.Books {
		if sh.Books[i].ID == id {
			return &sh.Books[i]
		}
	}
	return nil
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

// Initial sample data for first run

func initializeSampleData(lib *Library) {
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
