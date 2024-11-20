package main

import (
	"context"
	"fmt"
	stlog "log"

	"luuk/distributed/log"
	"luuk/distributed/registry"
	"luuk/distributed/service"
)

func main() {
	log.Run("./distributed.log")
	host, port := "localhost", "4000"
	reg := registry.Registration{
		ServiceName: "Log Service",
		ServiceURL:  fmt.Sprintf("http://%s:%s", host, port),
	}
	ctx, err := service.Start(
		context.Background(),
		host,
		port,
		reg,
		log.RegisterHandlers,
	)
	if err != nil {
		stlog.Fatalln(err)
	}

	<- ctx.Done()

	fmt.Println("Shutdown logging service.")
}