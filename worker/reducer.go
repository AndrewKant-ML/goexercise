package worker

import (
	"fmt"
	"google.golang.org/grpc"
	"io"
	"it.uniroma2.dicii/goexercise/log"
	"it.uniroma2.dicii/goexercise/rpc/mapreduce"
	"os"
	"sort"
	"strconv"
	"strings"
)

type reducerService struct {
	mapreduce.UnimplementedReducerServiceServer
}

// Reduce executes reduce tasks
func (r *reducerService) Reduce(stream grpc.ClientStreamingServer[mapreduce.Number, mapreduce.Status]) error {
	for {
		num, err := stream.Recv()
		if err == io.EOF {
			// End of stream
			go orderNumbers()
			return stream.SendAndClose(&mapreduce.Status{
				Message: fmt.Sprintf("Correctly received %d numbers", len(receivedNumbers)),
			})
		}
		if err != nil {
			log.Error("error during reducer request handling", err)
			return err
		}
		receivedNumbers = append(receivedNumbers, num.Num)
	}
}

// orderNumbers reorder received numbers and save them in a local file
func orderNumbers() {
	sort.Slice(receivedNumbers, func(i, j int) bool {
		return receivedNumbers[i] < receivedNumbers[j]
	})

	file, err := os.OpenFile("output.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Error("unable to create output file", err)
		return
	}

	// Deferring file closing
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Error("unable to close output file", err)
			return
		}
	}(file)

	stringSlice := make([]string, len(receivedNumbers))
	for i, num := range receivedNumbers {
		stringSlice[i] = strconv.Itoa(int(num))
	}

	// Join the string slice with a separator
	result := strings.Join(stringSlice, ", ")
	i, err := file.WriteString(result)
	log.Info(fmt.Sprintf("Successfully wrote %d bytes", i))
}
