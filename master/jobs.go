package master

import (
	"context"
	"io"
	mr "it.uniroma2.dicii/goexercise/rpc/mapreduce"
	"sync"
)

type jobScheduler struct {
	mr.UnimplementedJobSchedulerServer
	jobs      map[int32]*mr.JobStatusMessage
	nextJobId int32
	mu        sync.Mutex
}

// CreateJob adds a new job to the list of active jobs
func (j *jobScheduler) CreateJob(_ context.Context, _ *mr.Job) (*mr.Job, error) {
	j.mu.Lock()
	defer func() {
		j.nextJobId++
		j.mu.Unlock()
	}()
	j.jobs[j.nextJobId] = &mr.JobStatusMessage{
		JobId:  j.nextJobId,
		Status: mr.JobStatus_CREATED,
	}
	return &mr.Job{JobId: j.nextJobId}, nil
}

// SendJob handles the reception of a mr job
func (j *jobScheduler) SendJob(stream mr.JobScheduler_SendDataServer) error {
	for {
		num, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&mr.JobStatusMessage{})
		}
		j.mu.Lock()
		if j, ok := j.jobs[num.JobId]; !ok {
			// New job
			j.jobs[job.JobId] = &mr.JobStatusMessage{}
		}
	}
}
