package app

import (
	"context"
	phoneBalanceV1Pb "github.com/Geniuskaa/Bank-system/pkg/gen/proto/v1"
	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"strings"
)

type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) CreatePattern(ctx context.Context, request *phoneBalanceV1Pb.Pattern) (*phoneBalanceV1Pb.PatternResponse, error) {
	log.Print(request)
	if !strings.EqualFold(request.GetPhoneNumber(), "+79644272738") {
		return &phoneBalanceV1Pb.PatternResponse{CompletedSuccesfully: true}, nil
	}

	return nil, status.Errorf(codes.NotFound, "pattern with phoneNumber %d is already created", request.GetPhoneNumber())
}

func (s *Server) GetAllPatterns(ctx context.Context, request *phoneBalanceV1Pb.EmptyRequest) (*phoneBalanceV1Pb.AllPatternsResponse, error) {
	log.Print(request)
	//TODO: запрос в бд с получением всех шаблонов
	return nil, status.Errorf(codes.NotFound, "there aren`t any patterns")
}

func (s *Server) GetPatternById(ctx context.Context, request *phoneBalanceV1Pb.PatternId) (*phoneBalanceV1Pb.Pattern, error) {
	log.Print(request)
	if request.GetId() == 1 {
		return &phoneBalanceV1Pb.Pattern{
			Id:          1,
			Title:       "Мамин телефон",
			PhoneNumber: "+79644272738",
			Created:     &timestamp.Timestamp{Seconds: 1656338658},
			Updated:     nil,
		}, nil
	}
	return nil, status.Errorf(codes.NotFound, "we didn`t find pattern with id %d", request.GetId())
}

func (s *Server) EditPatterById(ctx context.Context, request *phoneBalanceV1Pb.Pattern) (*phoneBalanceV1Pb.PatternResponse, error) {
	log.Print(request)
	//TODO: достать из бд шаблон и заменить его на новый

	return nil, status.Errorf(codes.NotFound, "some problems with editing pattern with id %d", request.GetId())
}

func (s *Server) DeleteById(ctx context.Context, request *phoneBalanceV1Pb.PatternId) (*phoneBalanceV1Pb.PatternResponse, error) {
	log.Print(request)
	//TODO: удалить из бд шаблон с заданным ID

	return nil, status.Errorf(codes.NotFound, "we couldn`t delete pattern with id %d", request.GetId())
}

//func (s *Server) FindByUserId(ctx context.Context, request *phoneBalanceV1Pb.FinesRequest) (*fineV1Pb.FinesResponse, error) {
//	log.Print(request)
//	if request.UserId == 1 {
//		return &fineV1Pb.FinesResponse{
//			UserId: 1,
//			Items:  []*fineV1Pb.Fine{},
//		}, nil
//	}
//
//	return nil, status.Errorf(codes.NotFound, "user with id %d not found", request.GetUserId())
//}
