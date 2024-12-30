package main

import (
	"context"
	"fmt"

	desc "gitlab.ozon.dev/go/classroom-14/students/week-3-workshop/pkg/api/notes/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func main() {
	conn, err := grpc.NewClient(":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	client := desc.NewNotesClient(conn)
	ctx := context.Background()
	ctx = metadata.AppendToOutgoingContext(ctx, "x-auth", "user1")

	req := &desc.SaveNoteRequest{
		Info: &desc.NoteInfo{
			Title:   "tawdad",
			Content: "cawdawdawdawdawdawd",
		},
	}

	response, err := client.SaveNote(ctx, req)
	if err != nil {
		panic(err)
	}
	fmt.Println(response.NoteId)
}
