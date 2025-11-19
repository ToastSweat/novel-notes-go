package main

import "fmt"

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
