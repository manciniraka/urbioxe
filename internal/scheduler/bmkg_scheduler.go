package scheduler

import (
	"log"
	"time"

	"github.com/manciniraka/urbioxe/internal/logger"
	"github.com/manciniraka/urbioxe/internal/service"

	"github.com/robfig/cron/v3"
)

type BMKGScheduler struct {
	bmkgService service.BMKGService
}

func NewBMKGScheduler(
	bmkgService service.BMKGService,
) *BMKGScheduler {

	return &BMKGScheduler{
		bmkgService: bmkgService,
	}
}

func (bs *BMKGScheduler) Start() {
	log.Println("BMKG scheduler started...")

	logger.Log.Info(
		"BMKG scheduler started",
		"cron",
		"0 0 * * *",
	)

	c := cron.New(
		cron.WithLocation(
			time.Local,
		),
	)

	_, err := c.AddFunc(
		"0 0 * * *",
		func() {
			logger.Log.Info(
				"starting BMKG synchronization",
			)

			result, err := bs.bmkgService.SyncForecasts()

			if err != nil {

				logger.Log.Error(
					"BMKG scheduler failed",
					"error",
					err,
				)

				return
			}

			logger.Log.Info(
				"BMKG synchronized successfully",
				"total_district",
				result.Total,
				"success_count",
				result.Success,
				"failed_count",
				result.Failed,
			)
		},
	)

	if err != nil {

		logger.Log.Error(
			"failed to register BMKG scheduler",
			"error",
			err,
		)

		return
	}

	go func() {
		logger.Log.Info(
			"initial BMKG synchronization",
		)

		result, err := bs.bmkgService.SyncForecasts()

		if err != nil {

			logger.Log.Error(
				"initial BMKG synchronization failed",
				"error",
				err,
			)

			return
		}

		logger.Log.Info(
			"initial BMKG synchronization completed",
			"total_district",
			result.Total,
			"success_count",
			result.Success,
			"failed_count",
			result.Failed,
		)
	}()

	c.Start()
}
