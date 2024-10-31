package auth

import "github.com/gofrs/uuid"

type LogSystemService interface {
	CreateLogSystem(reqFormat RequestLogSystemFormat, userId uuid.UUID, ipAddress string, userAgent string) (logSystem LogSystem, error error)
}

type LogSystemServiceImpl struct {
	LogSystemRepository LogSystemRepository
}

func ProvideLogSystemServiceImpl(LogSystemRepository LogSystemRepository) *LogSystemServiceImpl {
	s := new(LogSystemServiceImpl)
	s.LogSystemRepository = LogSystemRepository
	return s
}

func (s *LogSystemServiceImpl) CreateLogSystem(reqFormat RequestLogSystemFormat, userId uuid.UUID, ipAddress string, userAgent string) (newLogSystem LogSystem, err error) {
	if err != nil {
		return LogSystem{}, err
	}
	newLogSystem, _ = newLogSystem.NewLogSystemFormat(reqFormat, userId, ipAddress, userAgent)
	err = s.LogSystemRepository.CreateLogSystem(newLogSystem)
	if err != nil {
		return LogSystem{}, err
	}
	return newLogSystem, nil
}
