package controller_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/js-bruno/mariage-api/internal/adapter"
	"github.com/js-bruno/mariage-api/internal/controller"
)

func TestMain(m *testing.M) {
	// pool, err := dockertest.NewPool("")
	// if err != nil {
	// 	log.Fatalf("Some shitt happens in pool connection : %s", err)
	// }
	// resource, err := pool.RunWithOptions(&dockertest.RunOptions{
	// 	Repository: "mongo",
	// 	Tag:        "latest",
	// })
	// if err != nil {
	// 	log.Fatalf("container creation : %s", err)
	// }
	// defer resource.Close()
	//
	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestInsertGift(t *testing.T) {
	ctx := context.Background()
	client, close := adapter.MongoOpenConnetion()
	defer close()

	any, err := controller.InsertGift(ctx, client)

	fmt.Println("teste")
	fmt.Println(any)

	if err != nil {
		t.Error(err.Error())
	}
}
