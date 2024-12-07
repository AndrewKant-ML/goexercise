// Package master handles master operations
package master

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"it.uniroma2.dicii/goexercise/config"
	"it.uniroma2.dicii/goexercise/log"
	"it.uniroma2.dicii/goexercise/rpc/roles"
	"sync"
)

const (
	MAPPERS  = iota
	REDUCERS = iota
)

type hostList *[]config.Host

// Map of available workers with specified role
type availableWorkers struct {
	hosts map[int]hostList
	mu    sync.Mutex
}

var (
	workers = availableWorkers{hosts: make(map[int]hostList), mu: sync.Mutex{}}
)

// Start starts the master server
func Start() {
	err := initialize()
	if err != nil {
		log.Error("Unable to initialize master", err)
		return
	}
	// TODO implement mapreduce server
}

// initialize retrieves the workers configuration needed to the master to operate
func initialize() error {
	allWorkers, err := config.GetWorkers()
	mappersNumber, err := config.GetMappersNumber()
	reducersNumber, err := config.GetReducersNumber()
	if len(*allWorkers) != mappersNumber+reducersNumber {
		return errors.New("mismatch between mappers and reducers number and defined workers")
	}
	assignWorkerRoles(allWorkers, mappersNumber, reducersNumber)
	return err
}

// assignWorkerRoles assign a role to each worker
func assignWorkerRoles(hostList hostList, mappersNumber int, reducersNumber int) {
	initWorkerStruct()
	var addr string
	for _, w := range *hostList {
		go func() {
			addr = fmt.Sprintf("%s:%d", w.Address, w.Port)
			conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.ErrorMessage(fmt.Sprintf("Unable to connect to worker at address (%s)", addr))
			}

			defer func(conn *grpc.ClientConn) {
				err := conn.Close()
				if err != nil {
					log.Error(fmt.Sprintf("Error closing grpc connection to %v", addr), err)
				}
			}(conn)

			client := roles.NewRoleServiceClient(conn)
			var res *roles.RoleAssignmentResponse
			workers.mu.Lock()
			if len(*workers.hosts[MAPPERS]) < mappersNumber {
				res, err = client.AssignRole(context.Background(), &roles.RoleAssignment{Role: roles.Role_MAPPER})
			} else if len(*workers.hosts[REDUCERS]) < reducersNumber {
				res, err = client.AssignRole(context.Background(), &roles.RoleAssignment{Role: roles.Role_REDUCER})
			} else {
				res, err = client.AssignRole(context.Background(), &roles.RoleAssignment{Role: roles.Role_ROLE_UNKNOWN})
			}
			workers.mu.Unlock()
			if err != nil {
				log.ErrorMessage(fmt.Sprintf("Unable to assign role to worker at address (%s)", addr))
				log.ErrorMessage(res.GetMessage())
			} else {
				log.Info(res.GetMessage())
			}
		}()
	}
}

// initWorkerStruct initializes mappers and reducers arrays
func initWorkerStruct() {
	workers.hosts[MAPPERS] = &[]config.Host{}
	workers.hosts[REDUCERS] = &[]config.Host{}
}
