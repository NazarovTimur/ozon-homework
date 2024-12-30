package main

import (
	"context"
	"fmt"

	"gitlab.ozon.dev/go/classroom-14/students/week-4-workshop/internal/client/notes"
	"gitlab.ozon.dev/go/classroom-14/students/week-4-workshop/internal/model"

	desc "gitlab.ozon.dev/go/classroom-14/students/week-4-workshop/pkg/api/notes/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial(":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	client := desc.NewNotesClient(conn)

	wrappedClient := notes.NewClient(client)

	for i := range 100_000 {
		_, err = wrappedClient.SaveNote(context.Background(), &model.Note{
			Title:   fmt.Sprintf("title_%d", i),
			Content: fmt.Sprintf("content_%d", i),
			Tags:    []string{"tag1", fmt.Sprintf("unique_tag_%d", i)},
		})
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("complete")
}
