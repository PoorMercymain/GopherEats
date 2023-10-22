package main

import (
	"GopherEats/internal/pkg/logger"
	"context"
	"fmt"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main() {
	di := fx.New(
		fx.Provide(logger.ProvideLogger),
		fx.Invoke(func (zap *zap.SugaredLogger)  {
			zap.Infoln("cringe")
		}),
	)

	err := di.Start(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func ()  {
		err = di.Stop(context.Background())
		if err != nil {
			fmt.Println(err)
			return
		}
	}()
}