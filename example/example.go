package main

import (
	"context"
	"log"
	"os"

	ztcentral "github.com/zerotier/go-ztcentral"
)

func main() {
	c, err := ztcentral.NewClient(os.Getenv("ZEROTIER_CENTRAL_API_KEY"))
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	// get list of networks
	networks, err := c.GetNetworks(ctx)
	if err != nil {
		log.Println("error:", err.Error())
		os.Exit(1)
	}

	// print networks and members
	for _, n := range networks {
		log.Printf("%s\t%s", *n.Id, *n.Config.Name)
		members, err := c.GetMembers(ctx, *n.Id)
		if err != nil {
			log.Println("error:", err.Error())
			os.Exit(1)
		}

		for _, m := range members {
			log.Printf("\t%s\t %s", *m.NodeId, *m.Name)
		}
	}
}
