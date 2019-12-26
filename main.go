package main

import (
	"context"
	"log"
	"net"
	"sync"

	pb "github.com/johnhany97/consignments-service/proto/consignment"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

type repository interface {
	Create(*pb.Consignment) (*pb.Consignment, error)
}

// Repository - simulates use of datastore for now
type Repository struct {
	mu           sync.RWMutex
	consignments []*pb.Consignment
}

// Create new consignment
func (repo *Repository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
	repo.mu.Lock()
	updated := append(repo.consignments, consignment)
	repo.consignments = updated
	repo.mu.Unlock()
	return consignment, nil
}

// Service should implement all of the methods to satisfy
// the service we defined in protobuf def. Interface can be
// checked in generated code itself for the exact method
// signatures.
type service struct {
	repo Repository
}

// CreateConsignment - created just one method for the service
// which is the create method. It takes a context & a request as
// arguments. These are entirely handled by the gRPC server
func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment) (*pb.Response, error) {
	// Save consignment
	consignment, err := s.repo.Create(req)
	if err != nil {
		return nil, err
	}

	// Return matching the `Response` message created in
	// protobuf definition
	return &pb.Response{Created: true, Consignment: consignment}, nil
}

func main() {
	repo := &Repository{}

	// Set-up gRPC server
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	// Register the service with the gRPC server, this will
	// tie our implementation with the auto-generated interface
	// for the protobuf definition
	pb.RegisterShippingServiceServer(s, &service{*repo})

	// Register reflection service on gRPC server
	reflection.Register(s)

	log.Println("Running on port:", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to server: %v", err)
	}
}
