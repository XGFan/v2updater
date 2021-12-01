package main

import (
	"context"
	"flag"
	"google.golang.org/grpc"
	"log"
	"os"
	"v2ray.com/core"
	"v2ray.com/core/app/proxyman/command"
	_ "v2ray.com/core/main/jsonem"
)

func main() {
	connectUrl := flag.String("c", "127.0.0.1:10085", "connect url (host:port)")
	updateFile := flag.String("f", "update.json", "update config file")
	flag.Parse()
	dial, err := grpc.Dial(*connectUrl, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	open, err := os.Open(*updateFile)
	if err != nil {
		panic(err)
	}
	config, err := core.LoadConfig("auto", *updateFile, open)
	if err != nil {
		panic(err)
	}
	hsClient := command.NewHandlerServiceClient(dial)
	for _, handlerConfig := range config.Outbound {
		_, err = hsClient.RemoveOutbound(context.Background(), &command.RemoveOutboundRequest{
			Tag: handlerConfig.Tag,
		})
		_, err = hsClient.AddOutbound(context.Background(), &command.AddOutboundRequest{
			Outbound: handlerConfig,
		})
		if err != nil {
			log.Printf("failed to call grpc command: %v", err)
		}
	}

}
