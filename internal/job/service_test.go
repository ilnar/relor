package job

import (
	"context"
	"io"
	"log"
	"log/slog"
	"net"
	"sync"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/durationpb"

	pb "github.com/ilnar/wf/gen/pb/api"
	"github.com/ilnar/wf/internal/storage"
)

func TestJobListen(t *testing.T) {
	// Set up server.
	lis := bufconn.Listen(1024 * 1024)
	s := grpc.NewServer()

	js := New(slog.Default(), storage.NewJobStorage())
	pb.RegisterJobServiceServer(s, js)
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
	defer s.Stop()

	// Set up connection.
	ctx := context.Background()
	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}
	conn, err := grpc.DialContext(ctx, "bufnet",
		grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()

	// Set up listeners.
	var wg sync.WaitGroup
	client := pb.NewJobServiceClient(conn)

	var created, complete *pb.Job

	createStream, err := client.Listen(ctx, &pb.ListenRequest{
		WorkerId: "worker",
		ListenTo: pb.ListenOperation_CREATE,
	})
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		resp, err := createStream.Recv()
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Fatalf("failed to receive: %v", err)
		}
		created = resp
	}()

	completeStream, err := client.Listen(ctx, &pb.ListenRequest{
		WorkerId: "scheduler",
		ListenTo: pb.ListenOperation_COMPLETE,
	})
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		resp, err := completeStream.Recv()
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Fatalf("failed to receive: %v", err)
		}
		complete = resp
	}()

	// Create, claim, and complete job.
	const jobID = "00000000-0000-0000-0000-000000000001"
	labels := []string{"ok", "err"}

	createResp, err := client.Create(ctx, &pb.CreateRequest{
		Id: jobID,
		Reference: &pb.Reference{
			WorkflowId:     "00000000-0000-0000-0000-000000000002",
			WorkflowAction: "action",
		},
		ResultLabels: labels,
		Ttl:          durationpb.New(15 * time.Minute),
	})
	if err != nil {
		t.Fatalf("failed to create job: %v", err)
	}
	if createResp == nil {
		t.Fatal("unexpected nil response")
	}

	claimResp, err := client.Claim(ctx, &pb.ClaimRequest{
		Id: jobID,
	})
	if err != nil {
		t.Fatalf("failed to claim job: %v", err)
	}
	if claimResp == nil {
		t.Fatal("unexpected nil response")
	}

	completeResp, err := client.Complete(ctx, &pb.CompleteRequest{
		Id:          jobID,
		ResultLabel: "ok",
	})
	if err != nil {
		t.Fatalf("failed to complete job: %v", err)
	}
	if completeResp == nil {
		t.Fatal("unexpected nil response")
	}

	// Validate results.
	wg.Wait()

	if created == nil {
		t.Fatal("unexpected nil created job")
	}
	if created.Id != jobID {
		t.Errorf("unexpected created job id: %s; want %s", created.Id, jobID)
	}
	if complete == nil {
		t.Fatal("unexpected nil completed job")
	}
	if complete.Id != jobID {
		t.Errorf("unexpected completed job id: %s; want %s", complete.Id, jobID)
	}
}

func TestJobExpiration(t *testing.T) {
	// Set up server.
	lis := bufconn.Listen(1024 * 1024)
	s := grpc.NewServer()

	js := New(slog.Default(), storage.NewJobStorage())
	pb.RegisterJobServiceServer(s, js)
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
	defer s.Stop()

	// Set up connection.
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}
	conn, err := grpc.DialContext(ctx, "bufnet",
		grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()

	// Set up listeners.
	client := pb.NewJobServiceClient(conn)
	createStream, err := client.Listen(ctx, &pb.ListenRequest{
		WorkerId: "worker",
		ListenTo: pb.ListenOperation_CREATE,
	})
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}

	var created *pb.Job

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		resp, err := createStream.Recv()
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Printf("failed to receive: %v", err)
		}
		created = resp
	}()

	// Create job.
	createResp, err := client.Create(ctx, &pb.CreateRequest{
		Id: "00000000-0000-0000-0000-000000000001",
		Reference: &pb.Reference{
			WorkflowId:     "00000000-0000-0000-0000-000000000002",
			WorkflowAction: "action",
		},
		ResultLabels: []string{"ok", "err"},
	})
	if err != nil {
		t.Fatalf("failed to create job: %v", err)
	}
	if createResp == nil {
		t.Fatal("unexpected nil response")
	}

	time.Sleep(10 * time.Millisecond)

	// Ensure stream is closed by closing the server.
	cancel()

	wg.Wait()

	if created != nil {
		t.Errorf("unexpected non-nil created job %v", created)
	}
}
