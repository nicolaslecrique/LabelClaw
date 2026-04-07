package configuration

import (
	"context"
	"errors"
)

var ErrGenerateNotImplemented = errors.New("configuration generation is not implemented yet")

type Store interface {
	Load(ctx context.Context) (ActiveConfiguration, error)
	Save(ctx context.Context, configuration ActiveConfiguration) error
}

type Service struct {
	store Store
}

func NewService(store Store) Service {
	return Service{store: store}
}

func (s Service) GetCurrent(ctx context.Context) (ActiveConfiguration, error) {
	return s.store.Load(ctx)
}

func (s Service) SaveCurrent(ctx context.Context, configuration ActiveConfiguration) error {
	if err := configuration.Validate(); err != nil {
		return err
	}

	return s.store.Save(ctx, configuration)
}

func (s Service) Generate(ctx context.Context, request GenerateRequest) (ActiveConfiguration, error) {
	if err := request.Validate(); err != nil {
		return ActiveConfiguration{}, err
	}

	return ActiveConfiguration{}, ErrGenerateNotImplemented
}
