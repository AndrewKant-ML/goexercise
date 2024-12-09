package worker

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"it.uniroma2.dicii/goexercise/log"
	"it.uniroma2.dicii/goexercise/rpc/mapreduce"
	sn "it.uniroma2.dicii/goexercise/rpc/sync"
)

type mapperService struct {
	mapreduce.UnimplementedMapperServiceServer
}

// Map executes map tasks
func (m *mapperService) Map(stream grpc.ClientStreamingServer[mapreduce.Number, mapreduce.Status]) error {
	var received []int32

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
	for _, r := range *assignedRole.reducers {
		go func() {
			addr := fmt.Sprintf("%s:%d", r.address, r.port)
			conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Error(fmt.Sprintf("unable to connect to reducer at address %s", addr), err)
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
		}()
	}
	return nil
}

// splitData divides data to sed to each reducer
func splitData() {
	var shuffleAndSort = make(map[reducerInfo][]int32)
	for _, r := range *assignedRole.reducers {
		//  Init map arrays
		shuffleAndSort[r] = make([]int32, 0)
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
}

// sendDataToReducers sends divided data to reducers
func sendDataToReducers(queues *map[reducerInfo][]int32) {
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
func getReducerFromNumber(n int32) *reducerInfo {
	for _, r := range *assignedRole.reducers {
		if r.minValue <= n && n <= r.maxValue {
			return &r
		}
	}
	return nil
}
