package main

import (
    "log"
    "net"

    pb "github.com/xzx1kf/consignment-service/proto/consignment"
    "golang.org/x/net/context"
    "google.golang.org/grpc"
    "google.golang.org/grpc/reflection"
)

const (
    port = ":50051"
)

type IRepository interface {
    Create(*pb.Consignment) (*pb.Consignment, error)
    GetAll() []*pb.Consignment
}

type Repository struct {
    consignments []*pb.Consignment
}

// Create appends the new consignment to the repository of consignments.
func (repo *Repository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
    repo.consignments = append(repo.consignments, consignment)
    return consignment, nil
}

func (repo *Repository) GetAll() []*pb.Consignment {
    return repo.consignments
}

// service should implement all of the methods to satisfy the service
// defined in the protobuf definition.
type service struct {
    repo IRepository
}

func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment) (*pb.Response, error) {
    // Save the consignment
    consignment, err := s.repo.Create(req)
    if err != nil {
        return nil, err
    }

    // Return the matching 'Response' message created in the protobuf
    // definition.
    return &pb.Response{Created: true, Consignment: consignment}, nil
}

func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest) (*pb.Response, error) {
    consignments := s.repo.GetAll()
    return &pb.Response{Consignments: consignments}, nil
}

func main() {
    repo := &Repository{}

    // Set up the gRPC server.
    lis, err := net.Listen("tcp", port)
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }
    s := grpc.NewServer()

    // Register the service with the gRPC server
    pb.RegisterShippingServiceServer(s, &service{repo})

    // Register reflection service on gRPC server.
    reflection.Register(s)
    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
