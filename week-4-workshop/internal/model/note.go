package model

type Note struct {
	Id      int
	Title   string
	Content string
	Tags    []string
	Author  string `db:"-"`
}
