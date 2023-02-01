package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/docker/docker/client"
	"github.com/rumblefrog/go-a2s"
)

func main() {
	var SRCDS_HOST string
	var SRCDS_PORT uint16
	var SRCDS_CONTAINER_NAME string
	var CHECKER_INIT uint32 = 60
	var CHECKER_TIMEOUT uint32 = 60
	var CHECKER_POLLING_INTERVAL uint32 = 10

	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	var env string
	var ok bool

	env, ok = os.LookupEnv("SRCDS_HOST")
	if !ok {
		fmt.Fprintln(os.Stderr, "environment variable 'SRCDS_HOST' is not set")
		os.Exit(1)
	} else {
		SRCDS_HOST = env
	}

	env, ok = os.LookupEnv("SRCDS_PORT")
	if !ok {
		fmt.Fprintln(os.Stderr, "environment variable 'SRCDS_PORT' is not set")
		os.Exit(1)
	} else {
		val, err := (strconv.ParseUint(env, 10, 16))
		if err != nil {
			fmt.Fprintln(os.Stderr, "environment variable 'SRCDS_PORT' can not be intepreted as uint16")
			os.Exit(1)
		} else {
			SRCDS_PORT = uint16(val)
		}
	}

	env, ok = os.LookupEnv("SRCDS_CONTAINER_NAME")
	if !ok {
		fmt.Fprintln(os.Stderr, "environment variable 'SRCDS_CONTAINER_NAME' is not set")
		os.Exit(1)
	} else {
		SRCDS_CONTAINER_NAME = env
	}

	env, ok = os.LookupEnv("CHECKER_INIT")
	if ok {
		val, err := (strconv.ParseUint(env, 10, 32))
		if err != nil {
			fmt.Fprintln(os.Stderr, "environment variable 'CHECKER_INIT' can not be intepreted as uint32")
			os.Exit(1)
		} else {
			CHECKER_INIT = uint32(val)
		}
	}

	env, ok = os.LookupEnv("CHECKER_TIMEOUT")
	if ok {
		val, err := (strconv.ParseUint(env, 10, 32))
		if err != nil {
			fmt.Fprintln(os.Stderr, "environment variable 'CHECKER_TIMEOUT' can not be intepreted as uint32")
			os.Exit(1)
		} else {
			CHECKER_TIMEOUT = uint32(val)
		}
	}

	env, ok = os.LookupEnv("CHECKER_POLLING_INTERVAL")
	if ok {
		val, err := (strconv.ParseUint(env, 10, 32))
		if err != nil {
			fmt.Fprintln(os.Stderr, "environment variable 'CHECKER_POLLING_INTERVAL' can not be intepreted as uint32")
			os.Exit(1)
		} else {
			CHECKER_POLLING_INTERVAL = uint32(val)
			if CHECKER_POLLING_INTERVAL < 3 {
				CHECKER_POLLING_INTERVAL = 3
			}
		}
	}

	fmt.Printf("SRCDS_HOST=%s\n", SRCDS_HOST)
	fmt.Printf("SRCDS_PORT=%d\n", SRCDS_PORT)
	fmt.Printf("SRCDS_CONTAINER_NAME=%s\n", SRCDS_CONTAINER_NAME)
	fmt.Printf("CHECKER_INIT=%d\n", CHECKER_INIT)
	fmt.Printf("CHECKER_TIMEOUT=%d\n", CHECKER_TIMEOUT)
	fmt.Printf("CHECKER_POLLING_INTERVAL=%d\n", CHECKER_POLLING_INTERVAL)

	json, err := cli.ContainerInspect(ctx, SRCDS_CONTAINER_NAME)
	if err != nil {
		fmt.Fprintf(os.Stderr, "container with name '%s' from environment variable 'SRCDS_CONTAINER_NAME' does not exist\n", SRCDS_CONTAINER_NAME)
		os.Exit(0)
	}

	fmt.Printf("CONTAINER ID: %s\n", json.ID)

	a2sClient, err := a2s.NewClient(fmt.Sprintf("%s:%d", SRCDS_HOST, SRCDS_PORT))
	if err != nil {
		// handle
	}
	defer a2sClient.Close()

	var counter uint32 = 0

	time.Sleep(time.Duration(CHECKER_INIT) * time.Second)

	for {
		time.Sleep(time.Duration(CHECKER_POLLING_INTERVAL) * time.Second)

		_, err := a2sClient.QueryInfo()
		if err != nil {
			counter += CHECKER_POLLING_INTERVAL
			fmt.Printf("restart counter: %d\n", counter)
			fmt.Printf("error info: %s\n", err.Error())
		} else {
			if counter != 0 {
				fmt.Printf("restart counter: %d\n", 0)
			}
			counter = 0
		}

		if counter >= CHECKER_TIMEOUT {
			fmt.Printf("restarting the container '%s'", SRCDS_CONTAINER_NAME)

			timeout := 10 * time.Second
			err := cli.ContainerRestart(ctx, SRCDS_CONTAINER_NAME, &timeout)

			fmt.Println("restarting the server...")

			if err != nil {
				fmt.Fprintf(os.Stderr, "restarting the container '%s' failed\n", SRCDS_CONTAINER_NAME)
			} else {
				fmt.Printf("restarting the container '%s' succeeded\n", SRCDS_CONTAINER_NAME)
			}

			counter = 0

			time.Sleep(time.Duration(CHECKER_INIT) * time.Second)
		}
	}
}
