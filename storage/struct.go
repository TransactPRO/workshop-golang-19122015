package storage

import "workshop/blog"

type Storage struct {
	id    int
	items []blog.Entry
}
