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
	//, grpc.WithBlock()

	conn, err := grpc.Dial("localhost:8000", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	c := pb.NewCalendarServiceClient(conn)

	ctx, _ := context.WithTimeout(context.TODO(), 10*time.Second)
	md := metadata.Pairs("authorization", "")
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
