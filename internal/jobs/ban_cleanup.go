package jobs

import (
	"context"
	"log"
	"time"

	"github.com/angyxys/kexel/internal/service"
)

// BanCleanupJob periodically checks and unbans expired bans
type BanCleanupJob struct {
	banService *service.BanService
	interval   time.Duration
	stopChan   chan bool
	isRunning  bool
}

func NewBanCleanupJob(banService *service.BanService, interval time.Duration) *BanCleanupJob {
	return &BanCleanupJob{
		banService: banService,
		interval:   interval,
		stopChan:   make(chan bool),
		isRunning:  false,
	}
}

// Start begins the cleanup job
func (job *BanCleanupJob) Start() {
	if job.isRunning {
		log.Println("Ban cleanup job is already running")
		return
	}

	job.isRunning = true
	log.Printf("Starting ban cleanup job with interval: %v", job.interval)

	go func() {
		ticker := time.NewTicker(job.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				job.runCleanup()
			case <-job.stopChan:
				log.Println("Ban cleanup job stopped")
				job.isRunning = false
				return
			}
		}
	}()
}

// Stop stops the cleanup job
func (job *BanCleanupJob) Stop() {
	if job.isRunning {
		job.stopChan <- true
	}
}

// runCleanup executes the cleanup logic
func (job *BanCleanupJob) runCleanup() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	count, err := job.banService.CheckExpiredBans(ctx)
	if err != nil {
		log.Printf("Error during ban cleanup: %v", err)
		return
	}

	if count > 0 {
		log.Printf("Ban cleanup completed: %d players unbanned", count)
	}
}

// IsRunning returns whether the job is currently running
func (job *BanCleanupJob) IsRunning() bool {
	return job.isRunning
}
