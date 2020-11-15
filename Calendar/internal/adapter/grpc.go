package adapter

import (
	"context"
	"errors"

	"github.com/Wilder60/ArtemisV2/Calendar/internal/logger"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
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
	db  Storage
	log *logger.Zap
}

func ProvideGRPC(s *db.Firestore, mid *middleware.GRPC, logger *logger.Zap) *grpc.Server {
	grpc_zap.ReplaceGrpcLoggerV2(logger.Log)

	server := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_zap.UnaryServerInterceptor(logger.Log),
			grpc_auth.UnaryServerInterceptor(mid.GRPCAuthFunc),
		),
	)
	service.RegisterCalendarServiceServer(server, &calendarService{
		db:  s,
		log: logger,
	})
	return server
}

func (s *calendarService) GetEventsInRange(ctx context.Context, req *pb.GetEventsInRangeRequest) (*pb.EventResponse, error) {
	userID, err := getMetadataString(ctx)
	if err != nil {
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
