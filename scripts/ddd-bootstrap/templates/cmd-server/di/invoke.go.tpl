package di

import (
	"context"
	"fmt"

	"git.yugeeker.com/SHARED/go-lazy/config"
	"git.yugeeker.com/SHARED/go-lazy/tools/logo"
	"github.com/gofiber/fiber/v3"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var InvokeModule = fx.Options(
	fx.Invoke(func(lifecycle fx.Lifecycle, app *fiber.App, cfg *config.BasicConfiguration, log *zap.Logger) {
		app.Hooks().OnListen(func(_ fiber.ListenData) error {
			logo.NewLogoPrinter(cfg.AppName, logo.WithColor(logo.BlueColor)).Print()
			return nil
		})
		lifecycle.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				go func() {
					port := fmt.Sprintf(":%d", cfg.HttpPort)
					if err := app.Listen(port, fiber.ListenConfig{DisableStartupMessage: true}); err != nil {
						log.Error("fiber server failed", zap.Error(err))
					}
				}()
				return nil
			},
			OnStop: func(ctx context.Context) error { return app.Shutdown() },
		})
	}),
)
