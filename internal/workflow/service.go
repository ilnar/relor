package workflow

import (
	pb "github.com/ilnar/wf/gen/pb/api"
)

type Server struct {
	pb.UnimplementedWorkflowServiceServer
}
