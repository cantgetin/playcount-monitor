package beatmaprepository

import (
	"context"
	"playcount-monitor-backend/internal/database/repository/model"
	"playcount-monitor-backend/internal/database/txmanager"
)

type Interface interface {
	Create(ctx context.Context, tx txmanager.Tx, beatmap *model.Beatmap) error
	Update(ctx context.Context, tx txmanager.Tx, beatmap *model.Beatmap) error
	Get(ctx context.Context, tx txmanager.Tx, id int) (*model.Beatmap, error)
	ListForMapset(ctx context.Context, tx txmanager.Tx, mapsetID int) ([]*model.Beatmap, error)
}
