package controller

import (
    "bufio"
    "fmt"
    "Library_Managment/model"
    "Library_Managment/service"
    "os"
    "strconv"
)

func RunLibrarySystem() {
    library := service.NewLibraryManager()
    scanner := bufio.NewScanner(os.Stdin)

    for {
        fmt.Println("Library Management System")
        fmt.Println("1. Add a new book")
        fmt.Println("2. Remove an existing book")
        fmt.Println("3. Borrow a book")
        fmt.Println("4. Return a book")
        fmt.Println("5. List all available books")
        fmt.Println("6. List all borrowed books by a member")
        fmt.Println("7. Add a member")
        fmt.Println("8. Exit")
        fmt.Println()
        fmt.Print("Enter your choice: ")

        if !scanner.Scan() {
            break
        }
        choice := scanner.Text()

        switch choice {
        case "1":
            addBook(library, scanner)
        case "2":
            removeBook(library, scanner)
        case "3":
            borrowBook(library, scanner)
        case "4":
            returnBook(library, scanner)
        case "5":
            listAvailableBooks(library)
        case "6":
            listBorrowedBooks(library, scanner)
        case "7":
            addmember(library, scanner)
        case "8":
            fmt.Println("Exiting the system...")
            return
        default:
            fmt.Println("Invalid choice, please try again.")
        }

        fmt.Println()

		fmt.Print("Do you want to continue? (y/n): ")
		if !scanner.Scan() {
			break
		}
		continueChoice := scanner.Text()
		if continueChoice != "y" {
			fmt.Println("Exiting the system...")
			return
		}


        fmt.Println() // Add a blank line for better readability between operations
    }
}

func addBook(library service.LibraryManager, scanner *bufio.Scanner) {
    fmt.Print("Enter book ID: ")
    scanner.Scan()
    id, _ := strconv.Atoi(scanner.Text())

    fmt.Print("Enter book title: ")
    scanner.Scan()
    title := scanner.Text()

    fmt.Print("Enter book author: ")
    scanner.Scan()
    author := scanner.Text()

    book := model.Book{
        Id:     id,
        Title:  title,
        Author: author,
        Status: "Available",
    }

    err := library.AddBook(book)
    if err != nil {
        fmt.Println("Error adding book:", err)
    }else {
        fmt.Println("Book added successfully!")
    }
}

func removeBook(library service.LibraryManager, scanner *bufio.Scanner) {
    fmt.Print("Enter book ID to remove: ")
    scanner.Scan()
    id, _ := strconv.Atoi(scanner.Text())

    err := library.RemoveBook(id)
    if err != nil {
        fmt.Println("Error removing book:", err)
    }else {
        fmt.Println("Book removed successfully!")
    }
}

func borrowBook(library service.LibraryManager, scanner *bufio.Scanner) {
    fmt.Print("Enter book ID to borrow: ")
    scanner.Scan()
    bookID, _ := strconv.Atoi(scanner.Text())

    fmt.Print("Enter member ID: ")
    scanner.Scan()
    memberID, _ := strconv.Atoi(scanner.Text())

    err := library.BorrowBook(bookID, memberID)
    if err != nil {
        fmt.Println("Error borrowing book:", err)
    } else {
        fmt.Println("Book borrowed successfully!")
    }
}

func returnBook(library service.LibraryManager, scanner *bufio.Scanner) {
    fmt.Print("Enter book ID to return: ")
    scanner.Scan()
    bookID, _ := strconv.Atoi(scanner.Text())

    fmt.Print("Enter member ID: ")
    scanner.Scan()
    memberID, _ := strconv.Atoi(scanner.Text())

    err := library.ReturnBook(bookID, memberID)
    if err != nil {
        fmt.Println("Error returning book:", err)
    } else {
        fmt.Println("Book returned successfully!")
    }
}

func listAvailableBooks(library service.LibraryManager) {
    books, _ := library.ListAvailableBooks()
    fmt.Println("Available Books:")
    for _, book := range books {
        fmt.Printf("ID: %d, Title: %s, Author: %s\n", book.Id, book.Title, book.Author)
    }
}

func listBorrowedBooks(library service.LibraryManager, scanner *bufio.Scanner) {
    fmt.Print("Enter member ID: ")
    scanner.Scan()
    memberID, _ := strconv.Atoi(scanner.Text())

    books, err := library.ListBorrowedBooks(memberID)
    if err != nil{
        fmt.Println("Error listing borrowed books:", err)

    }else{    
        fmt.Println("Borrowed Books:")
        for _, book := range books {
            fmt.Printf("ID: %d, Title: %s, Author: %s\n", book.Id, book.Title, book.Author)
        }}
}

func addmember(library service.LibraryManager, scanner *bufio.Scanner) {
    fmt.Print("Enter member ID: ")
    scanner.Scan()
    id, _ := strconv.Atoi(scanner.Text())

    fmt.Print("Enter member name: ")
    scanner.Scan()
    name := scanner.Text()

    member := model.Member{
        Id:     id,
        Name:  name,
        BorrowedBooks: []model.Book{},
    }

    err := library.AddMember(member)
    if err != nil {
        fmt.Println("Error adding member:", err)
    }else {
        fmt.Println("Member added successfully!")
    }
}
