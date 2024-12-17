package cmd

import (
	"github.com/spf13/cobra"
	"it.uniroma2.dicii/goexercise/log"
	"it.uniroma2.dicii/goexercise/master"
	"it.uniroma2.dicii/goexercise/worker"
)

var rootCmd = &cobra.Command{
	Use:   "mapreduce",
	Short: "Starts the MapReduce application",
	Long:  "Starts the MapReduce application",
}

var masterCmd = &cobra.Command{
	Use:   "master",
	Short: "Starts the Master application",
	Long:  "Starts the Master application",
	Run: func(cmd *cobra.Command, args []string) {
		master.StartMaster()
	},
}

var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Starts the Master application",
	Long:  "Starts the Master application",
	Run: func(cmd *cobra.Command, args []string) {
		index, _ := cmd.Flags().GetInt("index")
		worker.Start(index)
	},
}

func init() {
	rootCmd.AddCommand(masterCmd)
	workerCmd.Flags().IntP("index", "i", 0, "Index of the worker to use")
	rootCmd.AddCommand(workerCmd)
}

func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		log.Error("Error while starting the MapReduce application", err)
		return err
	}
	return nil
}
