package userprovide

import (
	"context"
	"playcount-monitor-backend/internal/database/repository/model"
	"playcount-monitor-backend/internal/database/txmanager"
	"playcount-monitor-backend/internal/dto"
	"playcount-monitor-backend/internal/usecase/mappers"
)

func (uc *UseCase) Get(
	ctx context.Context,
	id int,
) (*dto.User, error) {
	var user *model.User
	txErr := uc.txm.ReadOnly(ctx, func(ctx context.Context, tx txmanager.Tx) error {
		var err error
		user, err = uc.user.Get(ctx, tx, id)
		if err != nil {
			return err
		}

		return nil
	})
	if txErr != nil {
		return nil, txErr
	}

	return mappers.MapUserModelToUserDTO(user)
}

func (uc *UseCase) GetByName(
	ctx context.Context,
	name string,
) (*dto.User, error) {
	var user *model.User
	txErr := uc.txm.ReadOnly(ctx, func(ctx context.Context, tx txmanager.Tx) error {
		var err error
		user, err = uc.user.GetByName(ctx, tx, name)
		if err != nil {
			return err
		}

		return nil
	})
	if txErr != nil {
		return nil, txErr
	}

	return mappers.MapUserModelToUserDTO(user)
}

func (uc *UseCase) List(
	ctx context.Context,
) ([]*dto.User, error) {
	var users []*model.User
	txErr := uc.txm.ReadOnly(ctx, func(ctx context.Context, tx txmanager.Tx) error {
		var err error
		users, err = uc.user.List(ctx, tx)
		if err != nil {
			return err
		}

		return nil
	})
	if txErr != nil {
		return nil, txErr
	}

	return mappers.MapUserModelsToUserDTOs(users)
}
