package workers

import (
	"fmt"

	"github.com/robfig/cron/v3"
	"shirinec.com/config"
	"shirinec.com/internal/repositories"
	"shirinec.com/internal/utils"
)

func ScheduleWorkers(mediaRepo repositories.MediaRepository) {
	c := cron.New()

	mediaCleaner := NewMediaCleanupWorker(&mediaRepo)
    mediaCleanerTimer := fmt.Sprintf("@every %s", config.AppConfig.MediaCleanerInterval)
	_, err := c.AddFunc(mediaCleanerTimer, mediaCleaner.CleanupUnusedImages)
	if err != nil {
		utils.Logger.Fatalf("ScheduleWorkers - Adding media.Cleaner.CleanupUnusedImages: %s", err.Error())
	}

    c.Start()
}
