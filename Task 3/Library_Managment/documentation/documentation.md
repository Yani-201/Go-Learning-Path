# Library Management System Documentation

## Overview
The Library Management System is a simple console-based application that allows users to manage books and members. Users can perform operations such as adding, removing, borrowing, and returning books, as well as listing available and borrowed books.

## Features
- **Add a new book**: Add a new book to the library.
- **Remove an existing book**: Remove a book from the library.
- **Borrow a book**: Borrow a book for a member.
- **Return a book**: Return a borrowed book.
- **List all available books**: Display all books currently available in the library.
- **List all borrowed books by a member**: Display all books borrowed by a specific member.
- **Exit**: Exit the library management system.

## Folder Structure
library_management/ 
    ├── main.go 
    ├── controllers/ 
        │ └── library_controller.go 
    ├── models/ 
        │ ├── book.go 
        │ └── member.go 
    ├── services/ 
        │ └── library_service.go 
    ├── docs/ 
        │ └── documentation.md 
    └── go.mod


## Getting Started

### Prerequisites
- Go programming language installed on your system.

### Usage
When you run the application, you will see a menu with options to perform various operations. You can select an option by entering the corresponding number.

### Example Usage

#### Add a new book
1. Enter [`1`] to add a new book.
2. Enter the book ID, title, and author when prompted.
3. The system will confirm that the book has been added.

#### Remove an existing book
1. Enter [`2`] to remove a book.
2. Enter the book ID when prompted.
3. The system will confirm that the book has been removed.

#### Borrow a book
1. Enter [`3`]to borrow a book.
2. Enter the book ID and member ID when prompted.
3. The system will confirm that the book has been borrowed.

#### Return a book
1. Enter [`4`]to return a book.
2. Enter the book ID and member ID when prompted.
3. The system will confirm that the book has been returned.

#### List all available books
1. Enter [`5`] to list all available books.
2. The system will display a list of all books currently available in the library.

#### List all borrowed books by a member
1. Enter [`6`]to list all borrowed books by a member.
2. Enter the member ID when prompted.
3. The system will display a list of all books borrowed by the specified member.

#### Add a new member
1. Enter [`7`]to add a new member.
2. Enter the member ID and name when prompted.
3. The system will confirm that the new member has been added.

#### Exit
1. Enter [`8`] to exit the system.

## Code Explanation

### main.go
The entry point of the application. It calls the [`RunLibrarySystem`] function from the [`controllers`] package to start the library management system.

### controllers/library_controller.go
Contains the main logic for the console-based user interface. It provides functions to handle user input and perform operations on the library system.

### models/book.go
Defines the [`Book`] struct, representing a book in the library.

### models/member.go
Defines the [`Member`] struct, representing a library member.

### services/library_service.go
Implements the [`LibraryManager`] interface and provides methods to manage books and members in the library.

## Conclusion
This documentation provides an overview of the Library Management System, its features, and how to use it. The system allows users to manage books and members through a simple console-based interface.