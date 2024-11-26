package main

import (
	"context"
	"io"
	"log"
	"math/rand/v2"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/gemlab-dev/relor/gen/pb/api"
	"github.com/google/uuid"
)

func main() {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()

	c := pb.NewJobServiceClient(conn)
	l, err := c.Listen(ctx, &pb.ListenRequest{
		WorkerId: "worker-" + uuid.NewString(),
		ListenTo: pb.ListenOperation_CREATE,
	})
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	for {
		job, err := l.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("failed to receive: %v", err)
		}

		log.Printf("Doing %q (%s) tn: %q", job.Reference.WorkflowAction, job.Reference.WorkflowId, job.Reference.TransitionId)

		claimResp, err := c.Claim(ctx, &pb.ClaimRequest{
			Id: job.Id,
		})
		if err != nil {
			log.Fatalf("failed to claim: %v", err)
		}

		if len(claimResp.ResultLabels) == 0 {
			log.Fatal("no result labels")
		}
		opIdx := rand.IntN(len(claimResp.ResultLabels))

		time.Sleep(time.Duration(rand.IntN(2000)) * time.Millisecond)

		_, err = c.Complete(ctx, &pb.CompleteRequest{
			Id:          job.Id,
			ResultLabel: claimResp.ResultLabels[opIdx],
		})
		if err != nil {
			log.Fatalf("failed to complete: %v", err)
		}

		log.Printf("Result %q (%s)", claimResp.ResultLabels[opIdx], job.Reference.WorkflowId)
	}
}
