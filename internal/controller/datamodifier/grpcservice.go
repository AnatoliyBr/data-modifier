package datamodifier

import (
	"context"

	v1 "github.com/AnatoliyBr/data-modifier/pkg/api/v1"
)

type DataModifierService struct {
	v1.UnimplementedDataModifierServer
}

func (s *DataModifierService) AddAbsenceStatus(_ context.Context, req *v1.SourceData) (*v1.ModifiedData, error) {
	userData := req.GetUserData()

	return &v1.ModifiedData{
		ModifiedUserData: &v1.UserData{
			DisplayName: userData.DisplayName,
			Email:       userData.Email,
			MobilePhone: userData.MobilePhone,
			WorkPhone:   userData.WorkPhone,
		}}, nil
}
