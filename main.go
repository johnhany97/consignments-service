package main

import (
	"context"
	"fmt"
	"log"

	pb "github.com/johnhany97/consignments-service/proto/consignment"
	vesselProto "github.com/johnhany97/vessels-service/proto/vessel"
	"github.com/micro/go-micro"
)

// Repository - used to simulate datastore
type Repository interface {
	Create(*pb.Consignment) (*pb.Consignment, error)
	GetAll() []*pb.Consignment
}

// ConsignmentRepository - implementation of datastore
type ConsignmentRepository struct {
	consignments []*pb.Consignment
}

// Create new consignment
func (repo *ConsignmentRepository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
	updated := append(repo.consignments, consignment)
	repo.consignments = updated
	return consignment, nil
}

// GetAll consignments
func (repo *ConsignmentRepository) GetAll() []*pb.Consignment {
	return repo.consignments
}

// Service should implement all of the methods to satisfy
// the service we defined in protobuf def. Interface can be
// checked in generated code itself for the exact method
// signatures.
type service struct {
	repo         Repository
	vesselClient vesselProto.VesselServiceClient
}

// CreateConsignment - created just one method for the service
// which is the create method. It takes a context & a request as
// arguments. These are entirely handled by the gRPC server
func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response) error {

	// Call a client instance of vessel service with consignment weight
	// and amount of containers as the capacity value
	vesselResponse, err := s.vesselClient.FindAvailable(
		context.Background(),
		&vesselProto.Specification{
			MaxWeight: req.Weight,
			Capacity:  int32(leng(req.Containers)),
		})
	log.Printf("Found vessel: %s \n", vesselResponse.Vessel.Name)
	if err != nil {
		return err
	}

	req.VesselId = vesselReponse.Vessel.Id

	// Save consignment
	consignment, err := s.repo.Create(req)
	if err != nil {
		return err
	}

	res.Created = true
	res.Consignment = consignment
	return nil
}

func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest, res *pb.Response) error {
	consignments := s.repo.GetAll()

	res.Consignments = consignments
	return nil
}

func main() {
	repo := &ConsignmentRepository{}

	srv := micro.NewService(
		micro.Name("service.consignments"),
	)

	srv.Init()

	vesselClient := vesselProto.NewVesselServiceClient("service.vessels", srv.Client())

	// Register handler
	pb.RegisterShippingServiceHandler(srv.Server(), &service{repo, vesselClient})

	// Run the server
	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
