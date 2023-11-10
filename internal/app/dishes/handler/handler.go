package handler

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/PoorMercymain/GopherEats/internal/app/dishes/domain"
	"github.com/PoorMercymain/GopherEats/internal/pkg/logger"
	pb "github.com/PoorMercymain/GopherEats/pkg/api/dishes"
)

type DishesServerV1 struct {
	pb.UnimplementedDishesServiceV1Server
	service domain.DishesService
}

// NewDishesServerV1WithService returns DishesServerV1 for passed service interface.
func NewDishesServerV1WithService(service domain.DishesService) *DishesServerV1 {
	return &DishesServerV1{
		service: service,
	}
}

func (s *DishesServerV1) CreateIngredientV1(ctx context.Context, req *pb.CreateIngredientRequestV1) (*emptypb.Empty, error) {

	return &emptypb.Empty{}, nil
}

func (s *DishesServerV1) UpdateIngredientV1(ctx context.Context, req *pb.UpdateIngredientRequestV1) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *DishesServerV1) GetIngredientV1(ctx context.Context, req *pb.IngredientIdV1) (*pb.IngredientV1, error) {
	return nil, nil
}

func (s *DishesServerV1) DeleteIngredientV1(ctx context.Context, req *pb.IngredientIdV1) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *DishesServerV1) ListIngredientsV1(in *emptypb.Empty, stream pb.DishesServiceV1_ListIngredientsV1Server) (err error) {
	allIngredients, err := s.service.ListIngredients(stream.Context())
	for _, i := range allIngredients {
		err := stream.Send(i)
		if err != nil {
			logger.Logger().Infoln("Error sending message to stream: ", err)
			return
		}
	}
	return nil
}

func (s *DishesServerV1) CreateDishV1(ctx context.Context, req *pb.CreateDishRequestV1) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *DishesServerV1) UpdateDishV1(ctx context.Context, req *pb.UpdateDishRequestV1) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *DishesServerV1) GetDishV1(ctx context.Context, req *pb.DishIdV1) (*pb.DishV1, error) {
	return nil, nil
}

func (s *DishesServerV1) DeleteDishV1(ctx context.Context, req *pb.DishIdV1) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *DishesServerV1) ListDishesV1(in *emptypb.Empty, stream pb.DishesServiceV1_ListDishesV1Server) error {
	return nil
}

func (s *DishesServerV1) CreateBundleV1(ctx context.Context, req *pb.CreateBundleRequestV1) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *DishesServerV1) UpdateBundleV1(ctx context.Context, req *pb.UpdateBundleRequestV1) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *DishesServerV1) GetBundleV1(ctx context.Context, req *pb.BundleIdV1) (*pb.BundleV1, error) {
	return nil, nil
}

func (s *DishesServerV1) DeleteBundleV1(ctx context.Context, req *pb.BundleIdV1) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *DishesServerV1) ListBundlesV1(in *emptypb.Empty, stream pb.DishesServiceV1_ListBundlesV1Server) error {
	return nil
}

func (s *DishesServerV1) AddBundleWeeklyDishV1(ctx context.Context, req *pb.BundleWeeklyDishV1) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *DishesServerV1) DeleteBundleWeeklyDishV1(ctx context.Context, req *pb.BundleWeeklyDishV1) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *DishesServerV1) GetBundleWeeklyDishesV1(req *pb.GetBundleWeeklyDishesRequestV1, stream pb.DishesServiceV1_GetBundleWeeklyDishesV1Server) error {
	return nil
}

/*
func (s *DishesServerV1) SendOrderV1(ctx context.Context, req *pb.SendOrderRequestV1) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *DishesServerV1) CancelOrderV1(ctx context.Context, req *pb.CancelOrderRequestV1) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
*/
