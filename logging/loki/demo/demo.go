package main

import (
	"math/rand"
	"strconv"
	"time"

	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	rand.Seed(time.Now().UnixNano())
	for {
		val := rand.Intn(2)
		if val == 1 {
			logger.With(zap.String("xx", strconv.Itoa(rand.Int()))).Error("=====")
		} else {
			logger.With(zap.String("xx", strconv.Itoa(rand.Int()))).Info("=====")
		}
		time.Sleep(2 * time.Second)
	}
}
