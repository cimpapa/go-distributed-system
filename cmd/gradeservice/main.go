package main

import (
	"context"
	"fmt"
	"log"
	"luuk/distributed/grades"
	"luuk/distributed/registry"
	"luuk/distributed/service"
)


func main() {
	host, port := "localhost", "6000"
	r := registry.Registration{
		ServiceName: "Grade Service",
		ServiceURL: fmt.Sprintf("http://%s:%s", host, port),
	}
	ctx, err := service.Start(
		context.Background(),
		host,
		port,
		r,
		grades.RegisterHandlers,
	)
	if err != nil {
		log.Fatal(err)
	}
	<- ctx.Done()
	fmt.Println("Shutdown grading service")
}