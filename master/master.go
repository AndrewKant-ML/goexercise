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
	"it.uniroma2.dicii/goexercise/rpc/mapreduce"
	"it.uniroma2.dicii/goexercise/rpc/roles"
	"math/rand"
	"sort"
	"sync"
	"time"
)

type reducersList []*config.Host

type mapper struct {
	host         *config.Host
	reducersInfo []*roles.ReducerInfo
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
	numbers            []int64
	minValue, maxValue int64
)

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

	assignWorkerRoles(allWorkers)
	log.Info("Master successfully initialized")
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
func generateRandomIntegers(N int) (min, max int64) {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	randomNumbers := make([]int64, N)
	for i := range randomNumbers {
		randomNumbers[i] = rand.Int63() // Random integer starting from 0
	}
	numbers = make([]int64, N)
	copy(numbers, randomNumbers)

	sort.Slice(randomNumbers, func(i, j int) bool {
		return randomNumbers[i] < randomNumbers[j]
	})
	return randomNumbers[0], randomNumbers[N-1]
}

// assignWorkerRoles assign a role to each worker
// Master assign first reducers roles, then mappers roles. This is needed in
// order to pass reducers information to mappers
func assignWorkerRoles(hostList *[]config.Host) {
	var wg sync.WaitGroup
	for i, w := range (*hostList)[:workers.reducersNumber] {
		wg.Add(1)
		go func(w config.Host) {
			defer wg.Done()
			log.Info(fmt.Sprintf("Assigning role to %s:%d", w.Address, w.Port))
			addr := fmt.Sprintf("%s:%d", w.Address, w.Port)
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

			res, err = client.AssignRole(context.Background(), &roles.RoleAssignment{Role: roles.Role_REDUCER, Index: int64(i)})
			if err == nil {
				workers.reducers = append(workers.reducers, &w)
			}

			if err != nil {
				log.ErrorMessage(fmt.Sprintf("Unable to assign role to worker at address (%s)", addr))
				log.ErrorMessage(res.GetMessage())
			} else {
				log.Info(res.GetMessage())
			}
		}(w)
	}
	wg.Wait()
	log.Info("Successfully assigned reducers roles")

	reducersInfo := getReducersInfo()
	log.Info(fmt.Sprintf("Reducers info: %v", reducersInfo))
	for _, w := range (*hostList)[workers.reducersNumber:] {
		log.Info(fmt.Sprintf("Assigning role to %s:%d", w.Address, w.Port))
		addr := fmt.Sprintf("%s:%d", w.Address, w.Port)
		wg.Add(1)
		go func(w config.Host, addr string) {
			defer wg.Done()
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

			res, err = client.AssignRole(context.Background(), &roles.RoleAssignment{Role: roles.Role_MAPPER, ReducerInfo: reducersInfo})
			if err == nil {
				workers.mappers = append(workers.mappers, &mapper{host: &w, reducersInfo: reducersInfo})
			}

			if err != nil {
				log.ErrorMessage(fmt.Sprintf("Unable to assign role to worker at address (%s)", addr))
				log.ErrorMessage(res.GetMessage())
			} else {
				log.Info(res.GetMessage())
			}
		}(w, addr)
	}
	wg.Wait()
	log.Info("Successfully assigned mappers roles")

	workers.mu.Lock()
	defer workers.mu.Unlock()
	log.Info(fmt.Sprintf("Mappers: %#v", workers.mappers))
	log.Info(fmt.Sprintf("Reducers: %#v", workers.reducers))
}

// getReducersInfo build and returns a roles.ReducerInfo
// array containing infos about ranges of concern of
// each reducer
func getReducersInfo() []*roles.ReducerInfo {
	rangeSize := maxValue - minValue

	workers.mu.Lock()
	defer workers.mu.Unlock()

	reducerRangeSize := rangeSize / int64(workers.reducersNumber)
	ret := make([]*roles.ReducerInfo, len(workers.reducers))
	for i, r := range workers.reducers {
		if i == len(workers.reducers)-1 {
			ret[i] = &roles.ReducerInfo{
				Address: r.Address,
				Port:    r.Port,
				Min:     reducerRangeSize * int64(i),
				Max:     maxValue + 1,
			}
		} else {
			ret[i] = &roles.ReducerInfo{
				Address: r.Address,
				Port:    r.Port,
				Min:     reducerRangeSize * int64(i),
				Max:     reducerRangeSize*int64(i+1) - 1,
			}
		}
	}
	return ret
}

// sendData sends data to mappers
func sendData() {
	slices := splitSlice(numbers, workers.mappersNumber)
	log.Info(fmt.Sprintf("Sending %d slices: %v", len(slices), slices))
	wg := sync.WaitGroup{}
	for i, slice := range slices {
		wg.Add(1)
		w := workers.mappers[i]
		addr := fmt.Sprintf("%s:%d", w.host.Address, w.host.Port)
		go func(slice []int64, addr string) {
			defer wg.Done()
			log.Info(fmt.Sprintf("Sending data to %s", addr))
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

			client := mapreduce.NewMapperServiceClient(conn)
			stream, err := client.Map(context.Background())
			if err != nil {
				log.Error(fmt.Sprintf("Unable to open stream to %s", addr), err)
				return
			}

			for _, value := range slice {
				err := stream.Send(&mapreduce.Number{Num: value})
				if err != nil {
					log.Error(fmt.Sprintf("Unable to send data to %s", addr), err)
				}
			}

			// Close the stream to signal that the client is done sending
			resp, err := stream.CloseAndRecv()
			if err != nil {
				log.Error("failed to receive response: %v", err)
			}
			log.Info(resp.GetMessage())
		}(slice, addr)
	}
	wg.Wait()
	log.Info("Data sent correctly")
}

// splitSlice splits a slice into N parts
func splitSlice(slice []int64, n int) [][]int64 {
	if n <= 0 {
		return nil
	}
	if len(slice) == 0 {
		return make([][]int64, n)
	}

	chunkSize := len(slice) / n
	remainder := len(slice) % n

	var result [][]int64
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

// StartMaster starts the master server
func StartMaster() {
	err := initialize()
	if err != nil {
		log.Error("unable to initialize master", err)
		return
	}
	sendData()
}
