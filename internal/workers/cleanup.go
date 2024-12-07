package workers

import (
	"context"

	"shirinec.com/config"
	"shirinec.com/internal/repositories"
	"shirinec.com/internal/utils"
)

type MediaCleanupWorker interface {
	CleanupUnusedImages(ctx context.Context)
}

type mediaCleanupWorker struct {
	mediaRepo repositories.MediaRepository
}

func NewMediaCleanupWorker(mediaRepo *repositories.MediaRepository) MediaCleanupWorker {
	return &mediaCleanupWorker{mediaRepo: *mediaRepo}
}

func (w *mediaCleanupWorker) CleanupUnusedImages(ctx context.Context) {
    utils.Logger.Info("Starting media cleaner worker...")
	threshold := utils.DurationToPostgresqlInterval(config.AppConfig.MediaCleanupThreshold)

	mediaList, err := w.mediaRepo.ListForCleanUp(ctx, threshold)
	if err != nil {
		utils.Logger.Errorf("mediaCleanupWorker.CleanupUnusedImages - Calling mediaRepo.ListForCleanup: %s", err.Error())
        return
	}

    utils.Logger.Info("Removing orphaned medias from disk...")
    utils.Logger.Infof("%d orphaned media found", len(mediaList))
    for _, media := range mediaList {
        if err := utils.RemoveMedia(media); err != nil {
            utils.Logger.Errorf("mediaCleanupWorker.CleanupUnusedImages - Calling utils.RemoveMedia on %s: %s", media, err.Error())
        }
    }
    utils.Logger.Info("Listed medias removed from disk")

    if err := w.mediaRepo.DeleteRemovedMedia(ctx); err != nil {
        utils.Logger.Errorf("mediaCleanupWorker.CleanupUnusedImages - Calling mediaRepo.DeleteRemovedMedia: %s", err.Error())
    }
    utils.Logger.Info("Finished media cleaner worker process")
}
