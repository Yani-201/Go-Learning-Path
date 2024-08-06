package service

import (
	"Library_Managment/model"
	"errors"
)

type LibraryManager interface {
	AddBook(book model.Book) error
	RemoveBook(bookId int) error
	BorrowBook(bookId int, memberId int) error
	ReturnBook(bookId int, memberId int) error
	AddMember(member model.Member) error
	ListAvailableBooks() ([]model.Book, error)
	ListBorrowedBooks(memberId int) ([]model.Book, error)
}

type Library struct {
	Books   map[int]model.Book
	Members map[int]model.Member
}

func NewLibraryManager() LibraryManager {
	return &Library{
		Books:   make(map[int]model.Book),
		Members: make(map[int]model.Member),
	}
}

// AddBook implements LibraryManager.
func (l *Library) AddBook(book model.Book) error {
	if _,exists := l.Books[book.Id]; exists{
		return errors.New("book already exists")
	}
	l.Books[book.Id] = book
	return nil

}

// AddMember implements LibraryManager.
func (l *Library) AddMember(member model.Member) error {
	if _, exists := l.Members[member.Id]; exists{
		return errors.New("member already exists")
	}
	l.Members[member.Id] = member
	return nil
}

// RemoveBook implements LibraryManager.
func (l *Library) RemoveBook(bookId int) error {
	if _, exists := l.Books[bookId]; !exists{
		return errors.New("book does not exist")
	}
	delete(l.Books, bookId)
	return nil
}

// BorrowBook implements LibraryManager.
func (l *Library) BorrowBook(bookId int, memberId int) error {
	member, exist := l.Members[memberId]
	if !exist {
		return errors.New("member not Found")
	}

	book, exists := l.Books[bookId]
	if !exists{
		return errors.New("book not Found")
	}

	if book.Status == "Borrowed"{
		return errors.New("book has already been borrowed")
	}

	member.BorrowedBooks = append(member.BorrowedBooks, book)
	l.Members[memberId] = member

	book.Status = "Borrowed"
	l.Books[bookId] = book

	return nil

}

// ReturnBook implements LibraryManager.
func (l *Library) ReturnBook(bookId int, memberId int) error {
	member, exist := l.Members[memberId]
	if !exist {
		return errors.New("member not found")
	}

	book, exists := l.Books[bookId]
	if !exists {
		return errors.New("book not found")
	}

	if book.Status == "Available"{
		return errors.New("book is not borrowed")
	}

	book.Status = "Available"
	l.Books[bookId] = book

	for i, bk := range member.BorrowedBooks{
		if bk.Id == bookId {
			member.BorrowedBooks = append(member.BorrowedBooks[:i], member.BorrowedBooks[i+1:]... )
			break
		}
	}
	l.Members[memberId] = member
	return nil
}

// ListAvailableBooks implements LibraryManager.
func (l *Library) ListAvailableBooks() ([]model.Book, error) {
    availableBooks := []model.Book{}
    for _, book := range l.Books {
        if book.Status == "Available" {
            availableBooks = append(availableBooks, book)
        }
    }
    return availableBooks, nil
}

// ListBorrowedBooks implements LibraryManager.
func (l *Library) ListBorrowedBooks(memberID int) ([]model.Book, error) {
    member, exists := l.Members[memberID]
    if !exists {
        return nil, errors.New("member with this ID does not exist")
    }
    return member.BorrowedBooks, nil
}














