package worker

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"it.uniroma2.dicii/goexercise/config"
	"it.uniroma2.dicii/goexercise/log"
	"it.uniroma2.dicii/goexercise/rpc/barrier"
	"it.uniroma2.dicii/goexercise/rpc/mapreduce"
	"it.uniroma2.dicii/goexercise/rpc/roles"
	"net"
	"sync"
)

type workerService struct {
	roles.UnimplementedRoleServiceServer
}

type reducerInfo struct {
	address  string
	port     int64
	minValue int64
	maxValue int64
}

type role struct {
	role     roles.Role
	mu       sync.Mutex
	reducers *[]reducerInfo
}

var (
	assignedRole    = role{roles.Role_ROLE_UNKNOWN, sync.Mutex{}, nil}
	receivedNumbers []int64
)

func (w *workerService) AssignRole(_ context.Context, assignment *roles.RoleAssignment) (*roles.RoleAssignmentResponse, error) {
	assignedRole.mu.Lock()
	defer assignedRole.mu.Unlock()
	if assignedRole.role == roles.Role_ROLE_UNKNOWN {
		// Worker role has not been assigned yet
		assignedRole.role = assignment.Role
		if assignment.Role == roles.Role_MAPPER {
			if assignedRole.reducers == nil {
				// Populate reducers info array
				reducerNum, err := config.GetReducersNumber()
				if err != nil {
					log.Error("unable to retrieve reducers number from configuration", err)
					return nil, err
				}
				log.Info(fmt.Sprintf("%v", assignment.ReducerInfo))
				reducers := make([]reducerInfo, reducerNum)
				for i, ri := range (*assignment).ReducerInfo {
					reducers[i] = reducerInfo{
						address:  (*ri).Address,
						port:     (*ri).Port,
						minValue: (*ri).Min,
						maxValue: (*ri).Max,
					}
				}
				assignedRole.reducers = &reducers
				log.Info(fmt.Sprintf("assigned %d reducers successfully: %v", reducerNum, *assignedRole.reducers))
			}
		} else if assignment.Role == roles.Role_REDUCER {
			// Gets reducers info from message
			if mappersNumber == 0 {
				var err error
				mappersNumber, err = config.GetMappersNumber()
				if err != nil {
					log.Error("Unable to retrieve mappers number", err)
					return nil, err
				}
			}
			log.Info(fmt.Sprintf("mappers number: %d", mappersNumber))
			completions = &completionMessages{receivedMessages: make([]bool, mappersNumber), lastIndex: 0, mu: &sync.Mutex{}}
		}
		log.Info(fmt.Sprintf("assigned role: %v", assignment.Role))
		return &roles.RoleAssignmentResponse{Message: "Role assigned correctly"}, nil
	} else {
		return &roles.RoleAssignmentResponse{Message: "Worker role has already been assigned"}, errors.New("unable to assign worker role")
	}
}

// Start starts the worker server
func Start(index int) {
	workers, err := config.GetWorkers()
	if err != nil {
		log.Error("unable to get worker list from configuration", err)
		return
	}
	port := (*workers)[index].Port
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Error(fmt.Sprintf("unable to start worker on port %d", port), err)
		return
	}

	// Create worker server and register services
	worker := grpc.NewServer()
	roles.RegisterRoleServiceServer(worker, &workerService{})
	mapreduce.RegisterMapperServiceServer(worker, &mapperService{})
	mapreduce.RegisterReducerServiceServer(worker, &reducerService{})
	barrier.RegisterBarrierServiceServer(worker, &barrierService{})

	log.Info(fmt.Sprintf("Server is running on port %d...", port))
	if err := worker.Serve(listen); err != nil {
		log.Error("Failed to serve", err)
		return
	}
}
