package worker

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"it.uniroma2.dicii/goexercise/log"
	sn "it.uniroma2.dicii/goexercise/rpc/barrier"
	"it.uniroma2.dicii/goexercise/rpc/mapreduce"
	"sync"
)

type mapperService struct {
	mapreduce.UnimplementedMapperServiceServer
}

// Map executes map tasks
func (m *mapperService) Map(stream grpc.ClientStreamingServer[mapreduce.Number, mapreduce.Status]) error {
	var received []int64

	for {
		var req, err = stream.Recv()
		if err == io.EOF {
			// End of stream
			receivedNumbers = received
			go func() {
				err := sendCompletionMessage()
				if err != nil {
					log.Error("unable to send completion message", err)
				}
				sendDataToReducers(splitData())
			}()
			return stream.SendAndClose(&mapreduce.Status{
				Code:    0,
				Message: fmt.Sprintf("Correctly received %d numbers", len(received)),
			})
		}
		if err != nil {
			log.Error("error during mapper request ", err)
			return err
		}
		received = append(received, req.Num)
	}
}

// sendCompletionMessage sends a completion message to each reducer
func sendCompletionMessage() error {
	wg := sync.WaitGroup{}
	for _, r := range *assignedRole.reducers {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			addr := fmt.Sprintf("%s:%d", r.address, r.port)
			conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Error(fmt.Sprintf("unable to connect to reducer at address %s", addr), err)
				return
			}

			defer func(conn *grpc.ClientConn) {
				err := conn.Close()
				if err != nil {
					log.Error("unable to close reducer connection", err)
				}
			}(conn)

			client := sn.NewBarrierServiceClient(conn)
			status, err := client.SendCompletion(context.Background(), &sn.Completion{})
			if err != nil {
				log.Error(fmt.Sprintf("unable to send completion message to host %s", addr), err)
				return
			}
			log.Info(status.Message)
			defer wg.Done()
		}(&wg)
	}
	wg.Wait()
	return nil
}

// splitData divides data to sed to each reducer
func splitData() *map[reducerInfo][]int64 {
	var shuffleAndSort = make(map[reducerInfo][]int64)
	for _, r := range *assignedRole.reducers {
		//  Init map arrays
		shuffleAndSort[r] = make([]int64, 0)
	}
	// Search for the number queue
	for _, n := range receivedNumbers {
		if r := getReducerFromNumber(n); r == nil {
			// The number is not contained in any reducer range
			log.ErrorMessage(fmt.Sprintf("unable to retrieve reducer for number %d", n))
		} else {
			// Add number to reducer queue
			shuffleAndSort[*r] = append(shuffleAndSort[*r], n)
		}
	}
	log.Info(fmt.Sprintf("shuffle and sort reducers %v", shuffleAndSort))
	return &shuffleAndSort
}

// sendDataToReducers sends divided data to reducers
func sendDataToReducers(queues *map[reducerInfo][]int64) {
	for r, arr := range *queues {
		go func() {
			// Open the connection towards the reducer
			var addr = fmt.Sprintf("%s:%d", r.address, r.port)
			conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Error(fmt.Sprintf("unable to connect to reducer at address %s", addr), err)
				return
			}

			// Defer connection closing
			defer func(conn *grpc.ClientConn) {
				err := conn.Close()
				if err != nil {
					log.Error(fmt.Sprintf("unable to close connection to reducer at address %s", addr), err)
				}
			}(conn)

			client := mapreduce.NewReducerServiceClient(conn)
			stream, err := client.Reduce(context.Background())
			for _, n := range arr {
				err := stream.Send(&mapreduce.Number{Num: n})
				if err != nil {
					log.Error(fmt.Sprintf("unable to send message %d to reducer at address %s", n, addr), err)
				}
			}
		}()
	}
}

// getReducerFromNumber finds the reducer whose range include the given number
func getReducerFromNumber(n int64) *reducerInfo {
	for _, r := range *assignedRole.reducers {
		if r.minValue <= n && n <= r.maxValue {
			return &r
		}
	}
	return nil
}
