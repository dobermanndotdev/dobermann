package main

import (
	"context"
	"fmt"
	"time"

	"github.com/flowck/dobermann/backend/internal/adapters/endpoint_checkers"
	"github.com/flowck/dobermann/backend/internal/common/kron"
	"github.com/flowck/dobermann/backend/internal/common/logs"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tctx, tcancel := context.WithCancel(context.Background())
	defer tcancel()

	go func() {
		time.Sleep(time.Second * 2)
		tcancel()
	}()

	go func() {
		nCtx, cancel := context.WithCancel(tctx)
		defer cancel()

		<-nCtx.Done()

		logs.Info("context done")
	}()

	httpchecker, err := endpoint_checkers.NewHttpChecker("europe", 5)
	if err != nil {
		panic(err)
	}

	job := kron.NewJob(time.Second*1, func(ctx context.Context) error {
		result, err := httpchecker.Check(ctx, "https://firmino.work")
		if err != nil {
			return fmt.Errorf("check failed: %v", err)
		}

		if result.ResponseBody != "" {
			logs.Info(result.ResponseBody[:400])
		}

		return nil
	})

	service := kron.NewService()
	service.AddJob(job)
	logs.Println(service.Start(ctx))
}
