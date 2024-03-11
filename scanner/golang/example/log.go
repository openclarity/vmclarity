package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
)

func initLogger(level string, output io.Writer) error {
	logLevel, err := log.ParseLevel(level)
	if err != nil {
		return fmt.Errorf("failed to parse log level: %w", err)
	}
	log.SetLevel(logLevel)

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:          true,
		DisableTimestamp:       false,
		DisableSorting:         true,
		DisableLevelTruncation: true,
		QuoteEmptyFields:       true,
	})

	if logLevel >= log.DebugLevel {
		log.SetReportCaller(true)
	}

	log.SetOutput(output)

	return nil
}
