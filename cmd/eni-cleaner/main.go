// Standard gRPC health probe
package main

import (
	"flag"
	"os"
	"time"

	"github.com/aws/amazon-vpc-cni-k8s/pkg/awsutils"
	"github.com/aws/amazon-vpc-cni-k8s/pkg/utils/logger"
	"k8s.io/apimachinery/pkg/util/wait"
)

var cleanupPeriod = defaultCleanupPeriod

var log = logger.DefaultLogger()

const (
	defaultCleanupPeriod = 1 * time.Hour
)

const (
	// StatusInvalidArguments indicates specified invalid arguments.
	StatusInvalidArguments = 1
)

func init() {
	// timeouts
	flag.DurationVar(&cleanupPeriod, "cleanup-period", cleanupPeriod, "time between cleanups")

	flag.Parse()

	argError := func(s string, v ...interface{}) {
		log.Infof("error: "+s, v...)
		os.Exit(StatusInvalidArguments)
	}

	if cleanupPeriod <= 0 {
		argError("--cleanupPeriod must be greater than zero (specified: %v)", cleanupPeriod)
	}
}

func main() {
	cache, err := awsutils.New(false)
	if err != nil {
		panic(err)
	}

	go wait.Forever(cache.CleanUpLeakedENIs, cleanupPeriod)
}
