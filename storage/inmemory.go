package storage

import (
	"errors"
	"workshop/blog"
)

// Add is a method to add new blog entries into storage
func (storage *Storage) Add(entry blog.Entry) int {

	storage.items = append(storage.items, entry)

	return len(storage.items)
}

// GetLength is a method to get storage size
func (storage *Storage) GetLength() int {

	return len(storage.items)
}

// Update is a method to update blog entry
func (storage *Storage) Update(id int, entry blog.Entry) (bool, error) {

	if len(storage.items) == 0 || len(storage.items) < id || id == 0 {
		return false, errors.New("No record found")
	}

	storage.items[id-1] = entry

	return true, nil
}

// GetByID is a method which will return item by it's ID
func (storage *Storage) GetByID(id int) (blog.Entry, error) {

	if len(storage.items) == 0 || len(storage.items) < id || id == 0 {

		return blog.Entry{}, errors.New("No record found")
	}

	return storage.items[id-1], nil
}

// GetAll is a method which will return all storage items
func (storage *Storage) GetAll() []blog.Entry {

	return storage.items
}
