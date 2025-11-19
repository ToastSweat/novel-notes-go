package main

import (
	"fmt"
	"time"
)

// GetPageByDate finds a page with the given date in a book, or nil if not found.
func GetPageByDate(bk *Book, date string) *Page {
	for i := range bk.Pages {
		if bk.Pages[i].Date == date {
			return &bk.Pages[i]
		}
	}
	return nil
}

// GetOrCreatePageForDate returns the page for a specific date, creating it if needed.
func GetOrCreatePageForDate(bk *Book, date string) *Page {
	for i := range bk.Pages {
		if bk.Pages[i].Date == date {
			return &bk.Pages[i]
		}
	}

	newPage := Page{
		Date:  date,
		Items: []Item{},
	}
	bk.Pages = append(bk.Pages, newPage)
	return &bk.Pages[len(bk.Pages)-1]
}

// GetOrCreateTodayPage returns today's page for a book.
func GetOrCreateTodayPage(bk *Book) *Page {
	today := time.Now().Format("2006-01-02")
	return GetOrCreatePageForDate(bk, today)
}

// RolloverBook moves incomplete items from the book's current page to today's page,
// leaving completed items in the old page as history, and updates CurrentDate.
func RolloverBook(bk *Book, today string) {
	// If the book is already on today's date, nothing to do.
	if bk.CurrentDate == today {
		return
	}

	// If the book has no current date yet, just set it and ensure an empty page.
	if bk.CurrentDate == "" {
		bk.CurrentDate = today
		GetOrCreatePageForDate(bk, today)
		return
	}

	// Find the old "current" page.
	oldPage := GetPageByDate(bk, bk.CurrentDate)
	if oldPage == nil {
		// No page matches the old date; just switch to today.
		bk.CurrentDate = today
		GetOrCreatePageForDate(bk, today)
		return
	}

	// Ensure today's page exists.
	newPage := GetOrCreatePageForDate(bk, today)

	// Keep completed items on the old page; move incomplete ones to today's page.
	var remaining []Item
	for _, it := range oldPage.Items {
		if it.Completed {
			// Completed items stay as history on the old page.
			remaining = append(remaining, it)
		} else {
			// Incomplete items are rolled over to today's page as new items.
			newItemID := NextItemID(bk)
			newItem := NewItem(newItemID, it.Text)
			newPage.Items = append(newPage.Items, newItem)
		}
	}
	oldPage.Items = remaining

	// Update the book's current date.
	bk.CurrentDate = today
}

// AutoRollover runs rollover once per day across the entire library.
func AutoRollover(lib *Library) {
	today := time.Now().Format("2006-01-02")

	// Already rolled over for today? Then do nothing.
	if lib.LastRollOver == today {
		return
	}

	// No bookcases? Nothing to do.
	if len(lib.Bookcases) == 0 {
		lib.LastRollOver = today
		return
	}

	// Apply rollover to every book in the library.
	for bci := range lib.Bookcases {
		for shi := range lib.Bookcases[bci].Shelves {
			for bki := range lib.Bookcases[bci].Shelves[shi].Books {
				bk := &lib.Bookcases[bci].Shelves[shi].Books[bki]
				RolloverBook(bk, today)
			}
		}
	}

	lib.LastRollOver = today
	fmt.Println("Auto rollover completed for", today)
}
