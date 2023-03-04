package cli

import (
	"fmt"
	"io/ioutil"
	"sync"
	"time"

	"memphis-load-tests-cli/handlers"

	"github.com/memphisdev/memphis.go"
	"github.com/spf13/cobra"
)

func getMsgInSize(len int64) []byte {
	bytes := make([]byte, len)
	for i := 0; i < int(len); i++ {
		bytes[i] = byte(1)
	}
	return bytes
}

var produceCmd = &cobra.Command{
	Use:     "produce",
	Aliases: []string{"prod"},
	Short:   "Ingest data into a memphis station",
	Args:    cobra.ExactArgs(0),
	Example: "produce --count 1000 --size 1024 --station xxx --replicas 1 --storage disk --concurrent 1 --host localhost --user root --token memphis --sync --json message.json",
	Run: func(cmd *cobra.Command, args []string) {
		var wg sync.WaitGroup
		c, _ := cmd.Flags().GetInt32("concurrent")
		size, _ := cmd.Flags().GetInt64("size")
		count, _ := cmd.Flags().GetInt64("count")
		station, _ := cmd.Flags().GetString("station")
		replicas, _ := cmd.Flags().GetInt32("replicas")
		storage, _ := cmd.Flags().GetString("storage")
		host, _ := cmd.Flags().GetString("host")
		user, _ := cmd.Flags().GetString("user")
		token, _ := cmd.Flags().GetString("token")
		sync, _ := cmd.Flags().GetBool("sync")
		json, _ := cmd.Flags().GetString("json")

		concurrent := int(c)
		totalBytes := int64(count) * size
		conn, err := memphis.Connect(host, user, token)
		if err != nil {
			fmt.Println("Can not connect with memphis: " + err.Error())
			return
		}

		storageType := memphis.Disk
		if storage == "memory" {
			storageType = memphis.Memory
		}
		_, err = conn.CreateStation(station, memphis.Replicas(int(replicas)), memphis.StorageTypeOpt(storageType))
		if err != nil {
			fmt.Println("Can not create station: " + err.Error())
			return
		}
		conn.Close()

		msg := getMsgInSize(size)
		if json != "" {
			msg, err = ioutil.ReadFile(json)
			if err != nil {
				fmt.Println("Can not read json files: " + err.Error())
				return
			}
		}

		wg.Add(concurrent)
		start := time.Now()
		for i := 0; i < concurrent; i++ {
			msgsCount := count / int64(concurrent)
			if i == concurrent-1 {
				msgsCount += count % int64(concurrent)
			}
			go handlers.Produce(&wg, host, user, token, station, msgsCount, msg, sync)
		}
		wg.Wait()
		latency := time.Since(start).Milliseconds()
		mbPerSec := float64(totalBytes) / float64(latency) / 1024
		msgsPerSec := (count / latency) * 1024

		fmt.Printf("Produce stats: latency: %v ms, %v msgs/ms ~ %.2f MB/sec\n", latency, msgsPerSec, mbPerSec)
	},
}

func init() {
	produceCmd.Flags().Int64("count", 10, "Number of messages to produce")
	produceCmd.Flags().Int64("size", 128, "Size of the test messages in bytes")
	produceCmd.Flags().String("station", "load-test", "station name")
	produceCmd.Flags().Int32("replicas", 1, "Number of station replicas")
	produceCmd.Flags().String("storage", "disk", "Where messages will be stored disk/memory")
	produceCmd.Flags().Int32("concurrent", 1, "Number of concurrent producers")
	produceCmd.Flags().String("host", "localhost", "memphis host")
	produceCmd.Flags().String("user", "root", "station name")
	produceCmd.Flags().String("token", "memphis", "station name")
	produceCmd.Flags().Bool("sync", false, "Synchronously publish to the station")
	produceCmd.Flags().String("json", "", "Path to a json file contains a message to send - in case you are using this option size option will be ignored")
	rootCmd.AddCommand(produceCmd)
}
