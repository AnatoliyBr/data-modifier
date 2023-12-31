// Package datamodifier implements gRPC service
// for modifying user data based on a third-party system.
package datamodifier

import (
	"context"
	"time"

	"github.com/AnatoliyBr/data-modifier/internal/entity"
	"github.com/AnatoliyBr/data-modifier/internal/usecase"
	v1 "github.com/AnatoliyBr/data-modifier/pkg/api/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	defaultTimeLayout = "2006-01-02T15:04:05"
)

// DataModifierService implements v1.DataModifierServer interface.
// It handles RPC requests and calls UseCase methods.
type DataModifierService struct {
	v1.UnimplementedDataModifierServer
	uc usecase.UseCase
}

// NewDataModifierService returns DataModifierService.
func NewDataModifierService(uc usecase.UseCase) *DataModifierService {
	return &DataModifierService{
		uc: uc,
	}
}

// AddAbsenceStatus adds an absence status to the user name.
func (s *DataModifierService) AddAbsenceStatus(_ context.Context, req *v1.SourceData) (*v1.ModifiedData, error) {
	userData := req.GetUserData()

	u := &entity.User{
		DisplayName: userData.DisplayName,
		Email:       userData.Email,
		MobilePhone: userData.MobilePhone,
		WorkPhone:   userData.WorkPhone,
	}

	if err := u.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	from, err1 := time.Parse(defaultTimeLayout, req.GetTimePeriod().DateFrom)
	to, err2 := time.Parse(defaultTimeLayout, req.GetTimePeriod().DateTo)

	if !from.Before(to) || err1 != nil || err2 != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid time period")
	}

	p := [2]entity.CustomTime{
		{Time: from},
		{Time: to},
	}

	if err := s.uc.GetUserID(u); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if err := s.uc.AddAbsenceStatus(u, p); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &v1.ModifiedData{
		ModifiedUserData: &v1.UserData{
			DisplayName: u.DisplayName,
			Email:       u.Email,
			MobilePhone: u.MobilePhone,
			WorkPhone:   u.WorkPhone,
		}}, nil
}
