package mapsetcreate

import (
	"context"
	"fmt"
	"playcount-monitor-backend/internal/database/repository/model"
	"playcount-monitor-backend/internal/database/txmanager"
)

func (uc *UseCase) Create(
	ctx context.Context,
	cmd *CreateMapsetCommand,
) error {
	txErr := uc.txm.ReadWrite(ctx, func(ctx context.Context, tx txmanager.Tx) error {
		mapsetExists, err := uc.mapset.Exists(ctx, tx, cmd.Id)
		if err != nil {
			return err
		}

		if mapsetExists {
			return fmt.Errorf("mapset with id %s already exists", cmd.Id)
		}

		// create mapset
		mapset, err := mapCommandToMapsetModel(cmd)
		if err != nil {
			return err
		}

		err = uc.mapset.Create(ctx, tx, mapset)
		if err != nil {
			return err
		}

		// create beatmaps
		for _, beatmap := range cmd.Beatmaps {
			var beatmapModel *model.Beatmap
			beatmapModel, err = mapBeatmapToBeatmapModel(&beatmap)
			if err != nil {
				return err
			}

			err = uc.beatmap.Create(ctx, tx, beatmapModel)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if txErr != nil {
		return txErr
	}

	return nil
}
