package cli

import (
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/briandowns/spinner"
	"github.com/memphisdev/memphis.go"
	"github.com/spf13/cobra"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func generateRandomJSON(size int) []byte {
	if size < 2 {
		return []byte("{}")
	}

	// Generating JSON-like string
	json := randStringBytes(size - 12)

	// Ensuring it starts and ends with curly braces
	jsonStr := "{\"message\":\"" + json + "\"}"
	return []byte(jsonStr)
}

func createStation(host, user, pass, station string, accountId int) error {
	c, err := memphis.Connect(host,
		user,
		memphis.Password(pass),
		memphis.AccountId(accountId),
	)
	if err != nil {
		return fmt.Errorf("failed to connect to Memphis server: %v", err)
	}

	_, err = c.CreateStation(station)
	if err != nil {
		return fmt.Errorf("failed to create station: %v", err)
	}
	c.Close()
	return nil
}

func generateClientsAndProduceMessages(host, user, pass, station, pName, partitionKey, messagePath, cName, cGroup string, mSize, count, partitionNumber, accountId, concurrency, batchSize, batchMaxWaitTime int, syncProduce bool, consumers []*memphis.Consumer, loader *spinner.Spinner) (int64, error) {
	fmt.Println("Initializing station and clients...")
	loader.Start()
	err := createStation(host, user, pass, station, accountId)
	if err != nil {
		return 0, fmt.Errorf(err.Error())
	}

	// creating separate conns and producers for each goroutine
	conns := make([]*memphis.Conn, concurrency)
	producers := make([]*memphis.Producer, concurrency)
	for i := 0; i < concurrency; i++ {
		conn, err := memphis.Connect(host,
			user,
			memphis.Password(pass),
			memphis.AccountId(accountId),
		)
		if err != nil {
			return 0, fmt.Errorf("failed to connect to Memphis server: %v", err)
		}
		conns[i] = conn

		p, err := conn.CreateProducer(station, pName)
		if err != nil {
			return 0, fmt.Errorf("failed to create producer: %v", err)
		}
		producers[i] = p

		if consumers != nil {
			c, err := conn.CreateConsumer(station, cName, memphis.ConsumerGroup(cGroup), memphis.BatchSize(batchSize), memphis.BatchMaxWaitTime(time.Duration(batchMaxWaitTime)*time.Millisecond))
			if err != nil {
				return 0, fmt.Errorf("failed to create consumer: %v", err)
			}
			consumers[i] = c
		}
	}

	var messageBytes []byte
	if messagePath != "" {
		messageBytes, err = os.ReadFile(messagePath)
		if err != nil {
			return 0, err
		}
	} else {
		messageBytes = generateRandomJSON(mSize)
	}
	produceSyncOpts := memphis.AsyncProduce()
	if syncProduce {
		produceSyncOpts = memphis.SyncProduce()
	}

	var producePartitionOpts memphis.ProduceOpt
	if partitionNumber > 0 {
		producePartitionOpts = memphis.ProducerPartitionNumber(partitionNumber)
	} else if partitionKey != "" {
		producePartitionOpts = memphis.ProducerPartitionKey(partitionKey)
	}
	loader.Stop()

	fmt.Println("Producing messages...")
	loader.Start()
	wg := sync.WaitGroup{}
	wg.Add(concurrency)
	start := time.Now()
	// send messages concurrently
	for i := 0; i < concurrency; i++ {
		go func(i int) {
			defer wg.Done()
			p := producers[i]
			messagesCount := count / concurrency
			// last producer will send the remaining messages, will handle cases in which the count is not evenly divisible by the concurrency
			if i == concurrency-1 {
				messagesCount = count - (count/concurrency)*(concurrency-1)
			}
			for j := 0; j < messagesCount; j++ {
				err := p.Produce(messageBytes, produceSyncOpts, producePartitionOpts)
				if err != nil {
					fmt.Printf("Produce failed: %v\n", err)
				}
			}
		}(i)
	}

	wg.Wait() // wait for all producers to finish
	loader.Stop()
	duration := time.Since(start).Milliseconds()

	return duration, nil
}

var benchProduceCmd = &cobra.Command{
	Use:     "producer",
	Aliases: []string{"produce"},
	Short:   "Produce messages to a station",
	Example: "mem bench produce --message-size 1024 --message-count 5",
	Run: func(cmd *cobra.Command, args []string) {
		host, err := cmd.Flags().GetString("host")
		if host == "" || err != nil {
			host = "localhost"
		}
		accountId, err := cmd.Flags().GetInt("account-id")
		if accountId < 1 || err != nil {
			accountId = 1
		}
		user, err := cmd.Flags().GetString("user")
		if user == "" || err != nil {
			fmt.Println("Please provide a user name")
			return
		}
		pass, err := cmd.Flags().GetString("password")
		if pass == "" || err != nil {
			fmt.Println("Please provide a password")
			return
		}
		station, err := cmd.Flags().GetString("station")
		if station == "" || err != nil {
			station = "benchmark-station"
		}
		pName, err := cmd.Flags().GetString("producer-name")
		if pName == "" || err != nil {
			pName = "p-bench"
		}
		mSize, err := cmd.Flags().GetInt("message-size")
		if mSize < 128 || mSize > 8388608 || err != nil {
			mSize = 128
		}
		count, err := cmd.Flags().GetInt("count")
		if count < 1 || err != nil {
			count = 1
		}
		partitionNumber, err := cmd.Flags().GetInt("partition-number")
		if err != nil {
			partitionNumber = 0
		}
		partitionKey, _ := cmd.Flags().GetString("partition-key")
		message, _ := cmd.Flags().GetString("message")        // default is ""
		syncProduce, _ := cmd.Flags().GetBool("sync-produce") // default is false
		concurrency, err := cmd.Flags().GetInt("concurrency")
		if concurrency < 1 || err != nil {
			concurrency = 1
		}

		loader := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		duration, err := generateClientsAndProduceMessages(host, user, pass, station, pName, partitionKey, message, "", "", mSize, count, partitionNumber, accountId, concurrency, 0, 0, syncProduce, nil, loader)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		durationForPrint := float64(duration)
		totalDurUnits := "ms"
		if duration >= 1000 {
			durationForPrint = durationForPrint / 1000
			totalDurUnits = "sec"
		}

		avgDurForPrint := float64(duration) / float64(count)
		avgDurUnits := "ms"
		if avgDurForPrint >= 1000 {
			avgDurForPrint = avgDurForPrint / 1000
			avgDurUnits = "sec"
		}

		fmt.Printf("%v message have been produced, total latency: %v%s, latency for a single message: %v%s\n", count, durationForPrint, totalDurUnits, avgDurForPrint, avgDurUnits)
		time.Sleep(time.Duration(count)*time.Microsecond + 1*time.Second) // wait for all messages to arrive
	},
}

var benchConsumeCmd = &cobra.Command{
	Use:     "consumer",
	Aliases: []string{"consume"},
	Short:   "Consume messages from a station",
	Example: "mem bench consume --batch-size 500 --concurrency 2",
	Run: func(cmd *cobra.Command, args []string) {
		host, err := cmd.Flags().GetString("host")
		if host == "" || err != nil {
			host = "localhost"
		}
		accountId, err := cmd.Flags().GetInt("account-id")
		if accountId < 1 || err != nil {
			accountId = 1
		}
		user, err := cmd.Flags().GetString("user")
		if user == "" || err != nil {
			fmt.Println("Please provide a user name")
			return
		}
		pass, err := cmd.Flags().GetString("password")
		if pass == "" || err != nil {
			fmt.Println("Please provide a password")
			return
		}
		station, err := cmd.Flags().GetString("station")
		if station == "" || err != nil {
			station = "benchmark-station"
		}
		pName, err := cmd.Flags().GetString("producer-name")
		if pName == "" || err != nil {
			pName = "p-bench"
		}
		mSize, err := cmd.Flags().GetInt("message-size")
		if mSize < 128 || mSize > 8388608 || err != nil {
			mSize = 128
		}
		count, err := cmd.Flags().GetInt("count")
		if count < 1 || err != nil {
			count = 1
		}
		cName, err := cmd.Flags().GetString("consumer-name")
		if cName == "" || err != nil {
			cName = "c-bench"
		}
		cGroup, err := cmd.Flags().GetString("group")
		if cGroup == "" || err != nil {
			cGroup = "cg-bench"
		}
		batchSize, err := cmd.Flags().GetInt("batch-size")
		if batchSize < 1 || err != nil {
			batchSize = 10
		}
		batchMaxWaitTime, err := cmd.Flags().GetInt("batch-max-wait-time")
		if batchMaxWaitTime < 5000 || err != nil {
			batchMaxWaitTime = 5000
		}
		partitionKey, _ := cmd.Flags().GetString("partition-key")
		message, _ := cmd.Flags().GetString("message") // default is ""
		concurrency, err := cmd.Flags().GetInt("concurrency")
		if concurrency < 1 || err != nil {
			concurrency = 1
		}

		loader := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		consumers := make([]*memphis.Consumer, concurrency)
		_, err = generateClientsAndProduceMessages(host, user, pass, station, pName, partitionKey, message, cName, cGroup, mSize, count, 1, accountId, concurrency, batchSize, batchMaxWaitTime, false, consumers, loader)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fetchPartitionOpts := []memphis.ConsumingOpt{}
		if partitionKey != "" {
			fetchPartitionOpts = append(fetchPartitionOpts, memphis.ConsumerPartitionKey(partitionKey))
		}

		ch := make(chan int, count/batchSize)
		totalConsumed := 0

		fmt.Println("Start consuming...")
		loader.Start()
		time.Sleep(time.Duration(count)*time.Microsecond + 1*time.Second) // wait for all messages to arrive
		// consume messages concurrently
		for i := 0; i < concurrency; i++ {
			go func(index int, chann chan int) {
				c := consumers[index]
				for totalConsumed < count {
					msgs, err := c.Fetch(batchSize, false, fetchPartitionOpts...)
					if err != nil {
						fmt.Printf("Fetch failed: %v\n", err)
					}
					if len(msgs) > 0 {
						chann <- len(msgs)
						for _, msg := range msgs {
							msg.Ack()
						}
					}
				}
			}(i, ch)
		}

		// wait for all messages to arrive
		start := time.Now()
		for totalConsumed < count {
			amount := <-ch
			totalConsumed += amount
		}
		duration := time.Since(start).Milliseconds()

		durationForPrint := float64(duration)
		totalDurUnits := "ms"
		if duration >= 1000 {
			durationForPrint = durationForPrint / 1000
			totalDurUnits = "sec"
		}

		avgDurForPrint := float64(duration) / float64(count)
		avgDurUnits := "ms"
		if avgDurForPrint >= 1000 {
			avgDurForPrint = avgDurForPrint / 1000
			avgDurUnits = "sec"
		}

		loader.Stop()
		fmt.Printf("%v message have been consumed, total latency: %v%s, latency for a single message: %v%s\n", count, durationForPrint, totalDurUnits, avgDurForPrint, avgDurUnits)
		time.Sleep(1 * time.Second) // wait for all messages to be acked
	},
}

var benchCmd = &cobra.Command{
	Use:     "benchmark",
	Aliases: []string{"bench"},
	Short:   "",
}

func init() {
	benchProduceCmd.Flags().String("station", "benchmark-station", "The desired station to which the messages will be produced, default is benchmark-station")
	benchProduceCmd.Flags().String("partition-key", "", "The desired partition key with which the messages will be produced, this will take priority in case partition-number flag is also provided")
	benchProduceCmd.Flags().Int("partition-number", 0, "The desired partition number to which the messages will be produced, default is 1")
	benchProduceCmd.Flags().String("producer-name", "p-bench", "The desired name of the producer, default is p-bench")
	benchProduceCmd.Flags().Int("message-size", 128, "The desired message size in bytes, default is 128, min is 128, max is 8,388,608(8MB). In case message flag is empty this will cause random data to be created")
	benchProduceCmd.Flags().Int("count", 1, "The desired amount of messages to be produced, default is 1")
	benchProduceCmd.Flags().String("message", "", "Path for a file holds the desired message to be produced, default is empty. In case this flag is empty this will cause random data to be created")
	benchProduceCmd.Flags().Bool("sync-produce", false, "Whether to wait for an acknowledgement for every message, default is false")
	benchProduceCmd.Flags().Int("concurrency", 1, "The desired amount of concurrent producers, default is 1")
	benchProduceCmd.Flags().String("host", "localhost", "Memphis host, default is localhost")
	benchProduceCmd.Flags().Int("account-id", 1, "The account id to use when connecting to the Memphis server, default is 1 (no need to pass when using the open-source edition)")
	benchProduceCmd.Flags().String("user", "", "The user name to use when connecting to the Memphis server")
	benchProduceCmd.Flags().String("password", "", "The password to use when connecting to the Memphis server")

	benchConsumeCmd.Flags().String("station", "benchmark-station", "The desired station to which the messages will be produced, default is benchmark-station")
	benchConsumeCmd.Flags().String("partition-key", "", "The desired partition key with which the messages will be consumed")
	benchConsumeCmd.Flags().String("consumer-name", "c-bench", "The desired name of the consumer, default is c-bench")
	benchConsumeCmd.Flags().String("group", "cg-bench", "The desired name of the consumers group, default is cg-bench")
	benchConsumeCmd.Flags().Int("batch-size", 10, "The desired batch size, default is 10")
	benchConsumeCmd.Flags().Int("batch-max-wait-time", 5000, "The desired max wait time (in millis) for a batch, default and minimum is 5000")
	benchConsumeCmd.Flags().Int("concurrency", 1, "The desired amount of concurrent producers, default is 1")
	benchConsumeCmd.Flags().String("producer-name", "p-bench", "The desired name of the producer, default is p-bench")
	benchConsumeCmd.Flags().Int("message-size", 128, "The desired message size in bytes, default is 128, min is 128, max is 8,388,608(8MB). In case message flag is empty this will cause random data to be created")
	benchConsumeCmd.Flags().Int("count", 1, "The desired amount of messages to be produced, default is 1")
	benchConsumeCmd.Flags().String("message", "", "Path for a file holds the desired message to be produced, default is empty. In case this flag is empty this will cause random data to be created")
	benchConsumeCmd.Flags().String("host", "localhost", "Memphis host, default is localhost")
	benchConsumeCmd.Flags().Int("account-id", 1, "The account id to use when connecting to the Memphis server (for open source users)")
	benchConsumeCmd.Flags().String("user", "", "The user name to use when connecting to the Memphis server")
	benchConsumeCmd.Flags().String("password", "", "The password to use when connecting to the Memphis server")

	benchCmd.AddCommand(benchProduceCmd)
	benchCmd.AddCommand(benchConsumeCmd)
	rootCmd.AddCommand(benchCmd)
}
