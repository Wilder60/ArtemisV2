package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc/metadata"

	pb "github.com/Wilder60/ArtemisV2/Calendar/pkg/grpc/service"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:8000", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	c := pb.NewCalendarServiceClient(conn)

	ctx, _ := context.WithTimeout(context.TODO(), 10*time.Second)
	md := metadata.Pairs("authorization", "bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiJURVNUVVNFUiIsImlhdCI6MTYwNDIyOTQ4MCwiaXNzIjoiQXJ0ZW1pcyJ9.MG8KqiaGERU4VONXJvnVxlwV17aRlSN4JoFU57Kz5EM")
	nCtx := metautils.NiceMD(md).ToOutgoing(ctx)

	r, err := c.GetEventsInRange(nCtx, &pb.GetEventsInRangeRequest{
		Sdate: "2020-10-28",
		Edate: "2020-10-28",
	})

	if err != nil {
		log.Fatal(err)
	}

	log.Println(r.GetEvents())
}
