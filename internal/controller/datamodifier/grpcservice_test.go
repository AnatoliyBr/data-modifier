package datamodifier_test

import (
	"errors"
	"testing"

	"github.com/AnatoliyBr/data-modifier/internal/controller/datamodifier"
	"github.com/AnatoliyBr/data-modifier/internal/entity"
	"github.com/AnatoliyBr/data-modifier/internal/usecase"
	"github.com/AnatoliyBr/data-modifier/internal/webapi"
	v1 "github.com/AnatoliyBr/data-modifier/pkg/api/v1"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

const (
	defaultTimeLayout = "2006-01-02T15:04:05"
)

func TestDataModifierService_AddAbsenceStatus(t *testing.T) {
	u := entity.TestUser()
	absenceData := entity.TestUserAbsenceData()
	p := [2]entity.CustomTime{absenceData.DateFrom, absenceData.DateTo}

	testSourceData := func() *v1.SourceData {
		return &v1.SourceData{
			UserData: &v1.UserData{
				DisplayName: u.DisplayName,
				Email:       u.Email,
				MobilePhone: u.MobilePhone,
				WorkPhone:   u.WorkPhone},
			TimePeriod: &v1.TimePeriod{
				DateFrom: p[0].Format(defaultTimeLayout),
				DateTo:   p[1].Format(defaultTimeLayout),
			},
		}
	}

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	testCases := []struct {
		name    string
		sd      func() *v1.SourceData
		wa      func() *webapi.MockWebAPI
		isValid bool
	}{
		{
			name: "valid",
			sd: func() *v1.SourceData {
				return testSourceData()
			},
			wa: func() *webapi.MockWebAPI {
				u1 := entity.TestUser()
				u1.ID = 0
				mockWebAPI := webapi.NewMockWebAPI(ctl)
				mockWebAPI.EXPECT().GetUserID(gomock.Eq(u1)).DoAndReturn(func(u *entity.User) error {
					u.ID = 1234
					return nil
				})

				u2 := entity.TestUser()
				mockWebAPI.EXPECT().AddAbsenceStatus(
					gomock.Eq(u2),
					gomock.Eq(p)).DoAndReturn(
					func(u *entity.User, p [2]entity.CustomTime) error {
						e := entity.ReasonList[absenceData.ReasonID].Emoji
						u.DisplayName = u.DisplayName + " " + e
						return nil
					})

				return mockWebAPI
			},
			isValid: true,
		},
		{
			name: "invalid argument",
			sd: func() *v1.SourceData {
				sd := testSourceData()
				sd.UserData.Email = "invalid"
				return sd
			},
			wa: func() *webapi.MockWebAPI {
				mockWebAPI := webapi.NewMockWebAPI(ctl)
				mockWebAPI.EXPECT().GetUserID(gomock.Any()).Return(nil).AnyTimes()
				mockWebAPI.EXPECT().AddAbsenceStatus(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
				return mockWebAPI
			},
			isValid: false,
		},
		{
			name: "invalid time period",
			sd: func() *v1.SourceData {
				sd := testSourceData()
				sd.TimePeriod.DateFrom = "invalid"
				return sd
			},
			wa: func() *webapi.MockWebAPI {
				mockWebAPI := webapi.NewMockWebAPI(ctl)
				mockWebAPI.EXPECT().GetUserID(gomock.Any()).Return(nil).AnyTimes()
				mockWebAPI.EXPECT().AddAbsenceStatus(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
				return mockWebAPI
			},
			isValid: false,
		},
		{
			name: "user id not found",
			sd: func() *v1.SourceData {
				return testSourceData()
			},
			wa: func() *webapi.MockWebAPI {
				u1 := entity.TestUser()
				u1.ID = 0
				mockWebAPI := webapi.NewMockWebAPI(ctl)
				mockWebAPI.EXPECT().GetUserID(gomock.Eq(u1)).Return(errors.New("not found"))
				mockWebAPI.EXPECT().AddAbsenceStatus(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
				return mockWebAPI
			},
			isValid: false,
		},
		{
			name: "status not found",
			sd: func() *v1.SourceData {
				sd := testSourceData()
				sd.UserData.Email = "mihalich@mail.ru"
				return sd
			},
			wa: func() *webapi.MockWebAPI {
				u1 := entity.TestUser()
				u1.ID = 0
				u1.Email = "mihalich@mail.ru"
				mockWebAPI := webapi.NewMockWebAPI(ctl)
				mockWebAPI.EXPECT().GetUserID(gomock.Eq(u1)).DoAndReturn(func(u *entity.User) error {
					u.ID = 5678
					return nil
				})

				u2 := entity.TestUser()
				u2.ID = 5678
				u2.Email = "mihalich@mail.ru"
				mockWebAPI.EXPECT().AddAbsenceStatus(gomock.Eq(u2), gomock.Any()).Return(errors.New("not found"))

				return mockWebAPI
			},
			isValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			uc := usecase.NewAppUseCase(tc.wa())
			dm := datamodifier.NewDataModifierService(uc)
			if tc.isValid {
				_, err := dm.AddAbsenceStatus(nil, tc.sd())
				assert.NoError(t, err)
			} else {
				_, err := dm.AddAbsenceStatus(nil, tc.sd())
				assert.Error(t, err)
			}
		})
	}
}
