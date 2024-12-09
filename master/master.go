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
	"math/rand"
	"sort"
	"sync"
	"time"
)

type reducersList []*config.Host

type mapper struct {
	host         *config.Host
	reducersInfo []roles.ReducerInfo
}

type mappersList []*mapper

// Map of available workers with specified role
type availableWorkers struct {
	mappers        mappersList
	reducers       reducersList
	mu             sync.Mutex
	reducersNumber int
	mappersNumber  int
}

var (
	workers = availableWorkers{
		mappers:  nil,
		reducers: nil,
		mu:       sync.Mutex{},
	}
	numbers            []int
	minValue, maxValue int
)

// Start starts the master server
func Start() {
	err := initialize()
	if err != nil {
		log.Error("unable to initialize master", err)
		return
	}
}

// initialize retrieves the workers configuration needed to the master to operate
func initialize() error {
	allWorkers, err := config.GetWorkers()
	mappersNumber, err := config.GetMappersNumber()
	reducersNumber, err := config.GetReducersNumber()
	if err != nil {
		return err
	}
	if len(*allWorkers) != mappersNumber+reducersNumber {
		return errors.New("mismatch between mappers and reducers number and defined workers")
	}
	initWorkersStruct(mappersNumber, reducersNumber)
	minValue, maxValue = generateRandomIntegers(100)

	// Using a WaitGroup to synchronize multiple goroutines at once
	var wg sync.WaitGroup
	assignWorkerRoles(allWorkers, &wg)
	wg.Wait()

	sendData()
	return err
}

// initWorkersStruct initializes mappers and reducers arrays
func initWorkersStruct(mappersNumber, reducersNumber int) {
	workers.mu.Lock()
	defer workers.mu.Unlock()

	workers.reducersNumber = reducersNumber
	workers.mappersNumber = mappersNumber

	reducers := make(reducersList, 0)
	mappers := make(mappersList, 0)

	workers.reducers = reducers
	workers.mappers = mappers
}

// generateRandomIntegers generates an ordered slice of N random integers (min: 0)
func generateRandomIntegers(N int) (min, max int) {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	randomNumbers := make([]int, N)
	for i := range randomNumbers {
		randomNumbers[i] = rand.Int() // Random integer starting from 0
	}
	numbers = make([]int, N)
	copy(numbers, randomNumbers)

	sort.Slice(randomNumbers, func(i, j int) bool {
		return randomNumbers[i] < randomNumbers[j]
	})
	return randomNumbers[0], randomNumbers[N-1]
}

// assignWorkerRoles assign a role to each worker
// Master assign first reducers roles, then mappers roles. This is needed in
// order to pass reducers information to mappers
func assignWorkerRoles(hostList *[]config.Host, wg *sync.WaitGroup) {
	var addr string
	for _, w := range *hostList {
		wg.Add(1)
		go func() {
			defer wg.Done()
			addr = fmt.Sprintf("%s:%d", w.Address, w.Port)
			conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.ErrorMessage(fmt.Sprintf("Unable to connect to worker at address (%s)", addr))
				return
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
			defer workers.mu.Unlock()
			if len(workers.reducers) < workers.reducersNumber {
				res, err = client.AssignRole(context.Background(), &roles.RoleAssignment{Role: roles.Role_REDUCER})
				if err == nil {
					workers.reducers = append(workers.reducers, &w)
				}
			} else if len(workers.mappers) < workers.mappersNumber {
				res, err = client.AssignRole(context.Background(), &roles.RoleAssignment{Role: roles.Role_MAPPER})
				if err == nil {
					workers.mappers = append(workers.mappers, &mapper{
						host:         &w,
						reducersInfo: getReducersInfo(),
					})
				}
			} else {
				res, err = client.AssignRole(context.Background(), &roles.RoleAssignment{Role: roles.Role_ROLE_UNKNOWN})
			}
			if err != nil {
				log.ErrorMessage(fmt.Sprintf("Unable to assign role to worker at address (%s)", addr))
				log.ErrorMessage(res.GetMessage())
			} else {
				log.Info(res.GetMessage())
			}
		}()
	}
}

// getReducersInfo build and returns a roles.ReducerInfo
// array containing infos about ranges of concern of
// each reducer
func getReducersInfo() []roles.ReducerInfo {
	rangeSize := maxValue - minValue

	workers.mu.Lock()
	defer workers.mu.Unlock()

	reducerRangeSize := rangeSize / workers.reducersNumber
	ret := make([]roles.ReducerInfo, len(workers.reducers))
	for i, r := range workers.reducers {
		ret[i] = roles.ReducerInfo{
			Address: r.Address,
			Port:    r.Port,
			Min:     int32(reducerRangeSize * i),
			Max:     int32(reducerRangeSize*(i+1) - 1),
		}
	}
	return ret
}

// sendData sends data to mappers
func sendData() {
	// TODO send data
}

// splitSlice splits a slice into N parts
func splitSlice(slice []int, n int) [][]int {
	if n <= 0 {
		return nil
	}
	if len(slice) == 0 {
		return make([][]int, n)
	}

	chunkSize := len(slice) / n
	remainder := len(slice) % n

	var result [][]int
	start := 0

	for i := 0; i < n; i++ {
		end := start + chunkSize
		if i < remainder {
			// Distribute the remainder
			end++
		}
		result = append(result, slice[start:end])
		start = end
	}
	return result
}
