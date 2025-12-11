package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
)

type JobStatus string

const (
	StatusQueued     JobStatus = "queued"
	StatusProcessing JobStatus = "processing"
	StatusCompleted  JobStatus = "completed"
	StatusFailed     JobStatus = "failed"
)

type ProgressEvent struct {
	Done    int    `json:"done"`
	Total   int    `json:"total"`
	Message string `json:"message,omitempty"`
}

type JobResult struct {
	Status JobStatus `json:"status"`
	Data   any       `json:"data,omitempty"`
	Error  string    `json:"error,omitempty"`
}

type Job struct {
	ID              string
	Handler         JobHandler
	Request         *http.Request
	Context         context.Context
	Cancel          context.CancelFunc
	Progress        chan ProgressEvent
	Result          chan JobResult
	CreatedAt       time.Time
	StartedAt       *time.Time
	FinishedAt      *time.Time
	Status          JobStatus
	CompletedResult *JobResult
}

type JobHandler interface {
	Handle(
		ctx context.Context,
		r *http.Request,
		progress func(done, total int, message string),
	) (any, error)
	GetJobType() string
}

type Queue struct {
	jobs       chan *Job
	activeJobs map[string]*Job
	mutex      sync.RWMutex
	workers    int
	ctx        context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup
}

func NewQueue(workers int) *Queue {
	ctx, cancel := context.WithCancel(context.Background())
	q := &Queue{
		jobs:       make(chan *Job, workers*2),
		activeJobs: make(map[string]*Job),
		workers:    workers,
		ctx:        ctx,
		cancel:     cancel,
	}

	for i := range workers {
		q.wg.Add(1)
		go q.worker(i)
	}

	log.Printf("Queue started with %d workers", workers)
	return q
}

func (q *Queue) worker(id int) {
	defer q.wg.Done()
	log.Printf("Worker %d started", id)

	for {
		select {
		case job := <-q.jobs:
			if job == nil {
				log.Printf("Worker %d received nil job, stopping", id)
				return
			}
			q.processJob(job, id)
		case <-q.ctx.Done():
			log.Printf("Worker %d stopping due to context cancellation", id)
			return
		}
	}
}

func (q *Queue) processJob(job *Job, workerID int) {
	log.Printf("Worker %d processing job %s (%s)", workerID, job.ID, job.Handler.GetJobType())

	now := time.Now()
	job.StartedAt = &now
	job.Status = StatusProcessing

	progress := func(done, total int, message string) {
		select {
		case job.Progress <- ProgressEvent{Done: done, Total: total, Message: message}:
		case <-job.Context.Done():
		}
	}

	result, err := job.Handler.Handle(job.Context, job.Request, progress)

	finishedAt := time.Now()
	job.FinishedAt = &finishedAt

	var jobResult JobResult
	if err != nil {
		log.Printf("Worker %d error: %v", workerID, err)
		job.Status = StatusFailed
		jobResult = JobResult{
			Status: StatusFailed,
			Error:  err.Error(),
		}
	} else {
		job.Status = StatusCompleted
		jobResult = JobResult{
			Status: StatusCompleted,
			Data:   result,
		}
	}

	job.CompletedResult = &jobResult

	select {
	case job.Result <- jobResult:
	case <-job.Context.Done():
	}

	log.Printf("Worker %d completed job %s in %v", workerID, job.ID, finishedAt.Sub(*job.StartedAt))
	go func() {
		time.Sleep(30 * time.Second) // Keep job for 30 seconds
		q.mutex.Lock()
		delete(q.activeJobs, job.ID)
		q.mutex.Unlock()
		log.Printf("Cleaned up job %s", job.ID)
	}()
}

func (q *Queue) Enqueue(handler JobHandler, r *http.Request) (string, error) {
	ctx, cancel := context.WithTimeout(q.ctx, 10*time.Minute)

	jobID := uuid.New().String()
	job := &Job{
		ID:        jobID,
		Handler:   handler,
		Request:   r,
		Context:   ctx,
		Cancel:    cancel,
		Progress:  make(chan ProgressEvent, 10),
		Result:    make(chan JobResult, 1),
		CreatedAt: time.Now(),
		Status:    StatusQueued,
	}

	q.mutex.Lock()
	q.activeJobs[jobID] = job
	q.mutex.Unlock()

	select {
	case q.jobs <- job:
		log.Printf("Job %s (%s) enqueued", jobID, handler.GetJobType())
		return jobID, nil
	case <-q.ctx.Done():
		cancel()
		return "", fmt.Errorf("queue is shutting down")
	}
}

func (q *Queue) GetJob(jobID string) (*Job, bool) {
	q.mutex.RLock()
	defer q.mutex.RUnlock()
	job, exists := q.activeJobs[jobID]
	return job, exists
}

func (q *Queue) GetQueueStats() map[string]any {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	stats := map[string]any{
		"active_jobs":  len(q.activeJobs),
		"pending_jobs": len(q.jobs),
		"workers":      q.workers,
	}

	statusCounts := make(map[JobStatus]int)
	for _, job := range q.activeJobs {
		statusCounts[job.Status]++
	}
	stats["status_counts"] = statusCounts

	return stats
}

func (q *Queue) Shutdown() {
	log.Println("Shutting down queue...")

	q.mutex.Lock()
	for _, job := range q.activeJobs {
		job.Cancel()
	}
	q.mutex.Unlock()

	close(q.jobs)
	q.cancel()

	q.wg.Wait()
	log.Println("Queue shut down")
}

func (q *Queue) SSEHandler(w http.ResponseWriter, r *http.Request, jobID string) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "SSE unsupported", http.StatusInternalServerError)
		return
	}

	job, exists := q.GetJob(jobID)
	if !exists {
		fmt.Fprintf(w, "event: error\ndata: %q\n\n", "Job not found")
		flusher.Flush()
		return
	}

	if job.CompletedResult != nil {
		data, _ := json.Marshal(job.CompletedResult)
		fmt.Fprintf(w, "event: result\ndata: %s\n\n", data)
		flusher.Flush()
		return
	}

	initialData, _ := json.Marshal(map[string]any{
		"status": job.Status,
		"job_id": job.ID,
	})
	fmt.Fprintf(w, "event: status\ndata: %s\n\n", initialData)
	flusher.Flush()

	ctx := r.Context()
	for {
		select {
		case progress := <-job.Progress:
			data, _ := json.Marshal(progress)
			fmt.Fprintf(w, "event: progress\ndata: %s\n\n", data)
			flusher.Flush()

		case result := <-job.Result:
			data, _ := json.Marshal(result)
			fmt.Fprintf(w, "event: result\ndata: %s\n\n", data)
			flusher.Flush()
			return

		case <-ctx.Done():
			log.Printf("SSE connection closed for job %s", jobID)
			return

		case <-job.Context.Done():
			errorResult := JobResult{
				Status: StatusFailed,
				Error:  "Job was cancelled",
			}
			data, _ := json.Marshal(errorResult)
			fmt.Fprintf(w, "event: result\ndata: %s\n\n", data)
			flusher.Flush()
			return
		}
	}
}
