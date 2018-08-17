package main

import (
    "fmt"
    "log"

    pb "github.com/xzx1kf/consignment-service/proto/consignment"
    vesselProto "github.com/xzx1kf/vessel-service/proto/vessel"
    micro "github.com/micro/go-micro"
    "golang.org/x/net/context"
)

type Repository interface {
    Create(*pb.Consignment) (*pb.Consignment, error)
    GetAll() []*pb.Consignment
}

type ConsignmentRepository struct {
    consignments []*pb.Consignment
}

// Create appends the new consignment to the repository of consignments.
func (repo *ConsignmentRepository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
    repo.consignments = append(repo.consignments, consignment)
    return consignment, nil
}

func (repo *ConsignmentRepository) GetAll() []*pb.Consignment {
    return repo.consignments
}

// service should implement all of the methods to satisfy the service
// defined in the protobuf definition.
type service struct {
    repo Repository
    vesselClient vesselProto.VesselServiceClient
}

func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response) error {
    // Call the client instance of the vessel service with the consignemnt weight,
    // and the amoutn of containers as the capcity value.
    vesselResponse, err := s.vesselClient.FindAvailable(context.Background(), &vesselProto.Specification{
        MaxWeight: req.Weight,
        Capacity: int32(len(req.Containers)),
    })
    if err != nil {
        return err
    }
    log.Printf("Found vessel: %s \n", vesselResponse.Vessel.Name)

    req.VesselId = vesselResponse.Vessel.Id

    // Save the consignment
    consignment, err := s.repo.Create(req)
    if err != nil {
        return err
    }

    // Return the matching 'Response' message created in the protobuf
    // definition.
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
        micro.Name("go.micro.srv.consignment"),
        micro.Version("latest"),
    )

    vesselClient := vesselProto.NewVesselServiceClient("go.micro.srv.vessel", srv.Client())

    srv.Init()

    // Register the service with the gRPC server
    pb.RegisterShippingServiceHandler(srv.Server(), &service{repo, vesselClient})

    // Run the server
    if err := srv.Run(); err != nil {
        fmt.Println(err)
    }
}
