package main

import (
	"context"
	"fmt"
	stlog "log"
	"luuk/distributed/grades"
	"luuk/distributed/registry"
	"luuk/distributed/service"
)

func main() {
	host, port := "localhost", "6000"
	r := registry.Registration{
		ServiceName:      registry.GradeService,
		ServiceURL:       fmt.Sprintf("http://%s:%s", host, port),
		RequiredService:  []registry.ServiceName{registry.LogService},
		ServiceUpdateURL: fmt.Sprintf("http://%s:%s", host, port) + "/services",
		HeartBeatURL:     fmt.Sprintf("http://%s:%s", host, port) + "/heartbeat",
	}
	ctx, err := service.Start(
		context.Background(),
		host,
		port,
		r,
		grades.RegisterHandlers,
	)
	if err != nil {
		stlog.Fatal(err)
	}

	<-ctx.Done()
	fmt.Println("Shutdown grading service")
}
