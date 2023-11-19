package datamodifier

import (
	"context"

	v1 "github.com/AnatoliyBr/data-modifier/pkg/api/v1"
)

type DataModifierService struct {
	WebAPI struct {
		IP       string
		Port     string
		Login    string
		Password string
	}
	v1.UnimplementedDataModifierServer
}

func (s *DataModifierService) AddAbsenceStatus(ctx context.Context, req *v1.SourceData) (*v1.ModifiedData, error) {
	userData := req.GetUserData()

	return &v1.ModifiedData{
		ModifiedUserData: &v1.UserData{
			Id:          userData.Id,
			DisplayName: userData.DisplayName,
			Email:       userData.Email,
			MobilePhone: userData.MobilePhone,
			WorkPhone:   userData.WorkPhone,
		}}, nil
}

func (s *DataModifierService) mustEmbedUnimplementedDataModifierServer() {}
