package adapter

import (
	"context"
	"errors"
	"log"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"go.uber.org/fx"
	"google.golang.org/grpc"

	"github.com/Wilder60/ArtemisV2/Calendar/internal/db"
	"github.com/Wilder60/ArtemisV2/Calendar/internal/domain"
	"github.com/Wilder60/ArtemisV2/Calendar/internal/domain/requests"
	"github.com/Wilder60/ArtemisV2/Calendar/internal/middleware"
	"github.com/Wilder60/ArtemisV2/Calendar/pkg/grpc/service"
	pb "github.com/Wilder60/ArtemisV2/Calendar/pkg/grpc/service"
)

type calendarService struct {
	pb.UnimplementedCalendarServiceServer
	db Storage
}

func ProvideGRPC(s *db.Firestore, mid *middleware.GRPC) *grpc.Server {
	server := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_auth.UnaryServerInterceptor(mid.GRPCAuthFunc)),
	)
	service.RegisterCalendarServiceServer(server, &calendarService{
		db: s,
	})
	return server
}

func (s *calendarService) GetEventsInRange(ctx context.Context, req *pb.GetEventsInRangeRequest) (*pb.EventResponse, error) {
	log.Println("In GetEventsInRange Function")
	userID, err := getMetadataString(ctx)
	if err != nil {
		log.Println(err.Error())
		return &pb.EventResponse{}, err
	}

	rangeRequest := requests.NewProtoGetRange(userID, req)
	data, err := s.db.GetEventsInRange(rangeRequest)
	return toEventResponse(data), err
}

func (s *calendarService) GetPaginatedEvents(ctx context.Context, req *pb.GetPaginatedEventsRequest) (*pb.EventResponse, error) {
	userID, err := getMetadataString(ctx)
	if err != nil {
		return &pb.EventResponse{}, err
	}

	paginationRequest := requests.NewProtoPagination(userID, req)
	data, err := s.db.GetEventsPaginated(paginationRequest)
	return toEventResponse(data), err
}

func (s *calendarService) AddEvent(ctx context.Context, req *pb.AddEventRequest) (*pb.EmptyResponse, error) {
	userID, err := getMetadataString(ctx)
	if err != nil {
		return &pb.EmptyResponse{}, err
	}

	addRequest := requests.NewProtoAdd(userID, req)
	err = s.db.CreateEvents(addRequest)
	return &pb.EmptyResponse{}, err
}

func (s *calendarService) UpdateEvent(ctx context.Context, req *pb.UpdateEventRequest) (*pb.EmptyResponse, error) {
	userID, err := getMetadataString(ctx)
	if err != nil {
		return &pb.EmptyResponse{}, err
	}

	updateRequest := requests.NewProtoUpdate(userID, req)
	err = s.db.UpdateEvent(updateRequest)
	return &pb.EmptyResponse{}, err
}

func (s *calendarService) DeleteEvent(ctx context.Context, req *pb.DeleteRequest) (*pb.EmptyResponse, error) {
	userID, err := getMetadataString(ctx)
	if err != nil {
		return &pb.EmptyResponse{}, err
	}

	deleteRequest := requests.NewProtoDelete(userID, req)
	err = s.db.DeleteEvents(deleteRequest)
	return &pb.EmptyResponse{}, err
}

func (s *calendarService) DeleteUser(ctx context.Context, req *pb.DeleteRequest) (*pb.EmptyResponse, error) {
	userID, err := getMetadataString(ctx)
	if err != nil {
		return &pb.EmptyResponse{}, err
	}

	userRequest := requests.NewProtoDelete(userID, req)
	err = s.db.DeleteEventsForUser(userRequest)
	return &pb.EmptyResponse{}, err
}

func getMetadataString(ctx context.Context) (string, error) {
	claims, err := middleware.GetClaims(ctx)
	if err != nil {
		return "", errors.New("no metadata found for incoming request")
	}
	return claims.UserID, nil
}

func toEventResponse(arr []domain.Event) *pb.EventResponse {
	res := []*pb.Event{}
	for _, event := range arr {
		res = append(res, event.ToProto())
	}
	return &pb.EventResponse{
		Events: res,
	}
}

var GRPCModule = fx.Option(
	fx.Provide(ProvideGRPC),
)
