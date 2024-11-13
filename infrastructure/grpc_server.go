package infrastructure

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/bernardbaker/qiba.core/app"
	"github.com/bernardbaker/qiba.core/domain"
	"github.com/bernardbaker/qiba.core/proto"
)

type GameServer struct {
	proto.UnimplementedGameServiceServer
	service *app.GameService
}

func NewGameServer(service *app.GameService) *GameServer {
	return &GameServer{service: service}
}

func (s *GameServer) StartGame(ctx context.Context, req *proto.StartGameRequest) (*proto.StartGameResponse, error) {
	encryptedData, hmac, gameID, err := s.service.StartGame()
	if err != nil {
		return nil, err
	}
	return &proto.StartGameResponse{EncryptedGameData: encryptedData, Hmac: hmac, GameId: gameID}, nil
}

func (s *GameServer) Spawn(ctx context.Context, req *proto.SpawnRequest) (*proto.SpawnResponse, error) {
	data, err := s.service.Spawn(req.GameId)
	if err != nil {
		return nil, err
	}
	return &proto.SpawnResponse{Data: data}, nil
}

func (s *GameServer) Tap(ctx context.Context, req *proto.TapRequest) (*proto.TapResponse, error) {
	timestamp, _ := time.Parse(time.RFC3339, req.Timestamp)
	success, err := s.service.Tap(req.GameId, req.ObjectId, timestamp)
	if err != nil {
		return nil, err
	}
	return &proto.TapResponse{Success: success}, nil
}

func (s *GameServer) EndGame(ctx context.Context, req *proto.EndGameRequest) (*proto.EndGameResponse, error) {
	score, err := s.service.EndGame(req.GameId)
	if err != nil {
		return nil, err
	}
	return &proto.EndGameResponse{Score: score}, nil
}

type ReferralServer struct {
	proto.UnimplementedReferralServiceServer
	service *app.ReferralService
}

func NewReferralServer(service *app.ReferralService) *ReferralServer {
	return &ReferralServer{service: service}
}

func (s *ReferralServer) Referral(ctx context.Context, req *proto.ReferralRequest) (*proto.ReferralResponse, error) {
	fmt.Println("req.User", req.User)
	err := s.service.Create(req.User.UserId)
	if err != nil {
		return nil, err
	}
	// debugging
	fmt.Println(s.service.Get(strconv.FormatInt(req.User.UserId, 10)))
	//
	return &proto.ReferralResponse{Success: true}, nil
}

func (s *ReferralServer) AcceptReferral(ctx context.Context, req *proto.AcceptReferralRequest) (*proto.AcceptReferralResponse, error) {
	from := domain.User{
		UserId:       req.From.UserId,
		Username:     req.From.Username,
		FirstName:    req.From.FirstName,
		LastName:     req.From.LastName,
		LanguageCode: req.From.LanguageCode,
		IsBot:        req.From.IsBot,
	}
	to := domain.User{
		UserId:       req.To.UserId,
		Username:     req.To.Username,
		FirstName:    req.To.FirstName,
		LastName:     req.To.LastName,
		LanguageCode: req.To.LanguageCode,
		IsBot:        req.To.IsBot,
	}
	err := s.service.Update(from, to)
	// debugging
	fmt.Println(s.service.Get(strconv.FormatInt(req.From.UserId, 10)))
	if err != nil {
		return nil, err
	}
	return &proto.AcceptReferralResponse{Success: true}, nil
}
