package notes

import (
	desc "gitlab.ozon.dev/go/classroom-14/students/week-4-workshop/pkg/api/notes/v1"
)

type Client struct {
	client desc.NotesClient
}

func NewClient(client desc.NotesClient) *Client {
	return &Client{client: client}
}
