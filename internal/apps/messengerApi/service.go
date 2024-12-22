package messengerApi

import (
	"insider/internal/apps/messengerApi/entities"
	"insider/internal/repositories"
	"insider/pkg"
	"time"
)

type MessengerService struct {
	messengerRepository *repositories.MessengerRepo
	validator           *entities.Validator
}

func NewMessengerService(messengerRepo *repositories.MessengerRepo) *MessengerService {
	validator := entities.NewValidator()
	return &MessengerService{
		messengerRepository: messengerRepo,
		validator:           &validator,
	}
}

func (s *MessengerService) GetAll(req *entities.GetAllRequest) (*pkg.BaseResponse, error) {
	messages, err := s.messengerRepository.GetAllMessagesMongo(req.Limit, req.Offset)
	if err != nil {
		return nil, err
	}
	return pkg.NewSuccessResponse(messages), nil
}

func (s *MessengerService) GetTwo(o int64) (*pkg.BaseResponse, error) {

	messages, err := s.messengerRepository.GetTwoPendingMessagesMongo(o)
	if err != nil {
		return nil, err
	}
	return pkg.NewSuccessResponse(messages), nil
}

func (s *MessengerService) CommitMessage(msgId string) error {
	if err := s.messengerRepository.ChangeMsgStatusMongo(msgId, "sent"); err != nil {
		return err
	}
	if err := s.messengerRepository.SetRedisKey(msgId, "sent", 30*time.Minute); err != nil {
		return err
	}
	return nil
}
