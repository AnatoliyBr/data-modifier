package usecase_test

import (
	"testing"

	"github.com/AnatoliyBr/data-modifier/internal/entity"
	"github.com/AnatoliyBr/data-modifier/internal/usecase"
	"github.com/AnatoliyBr/data-modifier/internal/webapi"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAppUseCase_GetUserID(t *testing.T) {
	u1 := entity.TestUser()
	u1.ID = 0
	u2 := entity.TestUser()

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockWebAPI := webapi.NewMockWebAPI(ctl)
	mockWebAPI.EXPECT().GetUserID(gomock.Eq(u1)).DoAndReturn(func(u *entity.User) error {
		u.ID = 1234
		return nil
	})

	uc := usecase.NewAppUseCase(mockWebAPI)

	assert.NoError(t, uc.GetUserID(u1))
	assert.Equal(t, u1.ID, u2.ID)
}

func TestAppUseCase_AddAbsenceStatus(t *testing.T) {
	u1 := entity.TestUser()
	u2 := entity.TestUser()
	absenceData := entity.TestUserAbsenceData()
	p := [2]entity.CustomTime{absenceData.DateFrom, absenceData.DateTo}

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockWebAPI := webapi.NewMockWebAPI(ctl)
	mockWebAPI.EXPECT().AddAbsenceStatus(
		gomock.Eq(u1),
		gomock.Eq(p)).DoAndReturn(
		func(u *entity.User, p [2]entity.CustomTime) error {
			e := entity.ReasonList[absenceData.ReasonID].Emoji
			u.DisplayName = u.DisplayName + " " + e
			return nil
		})

	uc := usecase.NewAppUseCase(mockWebAPI)

	assert.NoError(t, uc.AddAbsenceStatus(u1, p))
	assert.NotEqual(t, u1.DisplayName, u2.DisplayName)
}
