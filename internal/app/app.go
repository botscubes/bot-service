package app

import (
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/goccy/go-json"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	e "github.com/botscubes/bot-service/internal/api/errors"
	resp "github.com/botscubes/bot-service/pkg/api_response"

	"github.com/botscubes/bot-service/internal/bot"
	"github.com/botscubes/bot-service/internal/config"
	"github.com/botscubes/bot-service/internal/database/pgsql"
	rdb "github.com/botscubes/bot-service/internal/database/redis"
	"github.com/botscubes/bot-service/internal/database/redisauth"
	"github.com/botscubes/user-service/pkg/token_storage"

	"github.com/gofiber/fiber/v2"
)

type App struct {
	Server         *fiber.App
	BotService     *bot.BotService
	Conf           *config.ServiceConfig
	Db             *pgsql.Db
	SessionStorage token_storage.TokenStorage
	RedisAuth      *redis.Client
	Redis          *rdb.Rdb
	Log            *zap.SugaredLogger
}

func (app *App) Run(logger *zap.SugaredLogger, c *config.ServiceConfig) error {
	var err error
	done := make(chan struct{}, 1)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	app.Log = logger
	app.Conf = c
	app.Server = fiber.New(fiber.Config{
		AppName:               "Bot API Server",
		DisableStartupMessage: true,
		JSONEncoder:           json.Marshal,
		JSONDecoder:           json.Unmarshal,
		ErrorHandler:          app.errorHandler,
	})

	// app.Server = &telego.MultiBotWebhookServer{
	// 	Server: telego.FastHTTPWebhookServer{
	// 		Server: &fasthttp.Server{
	// 			Handler: app.Router.Handler,
	// 		},
	// 		Router: app.Router,
	// 	},
	// }

	// app.BotService = bot.NewBotService(&app.Conf.Bot, app.Server)

	app.RedisAuth = redisauth.NewClient(&app.Conf.RedisAuth)
	app.SessionStorage = token_storage.NewRedisTokenStorage(app.RedisAuth)

	app.Redis = rdb.NewClient(&app.Conf.Redis)

	pgsqlUrl := "postgres://" + app.Conf.Pg.User + ":" + app.Conf.Pg.Pass + "@" + app.Conf.Pg.Host + ":" + app.Conf.Pg.Port + "/" + app.Conf.Pg.Db
	if app.Db, err = pgsql.OpenConnection(pgsqlUrl); err != nil {
		return err
	}

	defer app.Db.CloseConnection()

	app.regiterHandlers()

	go func() {
		if err := app.Server.Listen(app.Conf.ListenAddress); err != nil {
			app.Log.Error("Start server:", err)
			sigs <- syscall.SIGTERM
		}
	}()

	go func() {
		<-sigs
		app.Log.Info("Stopping...")
		// if err := app.BotService.StopBots(); err != nil {
		// 	app.Log.Info("bots stop:\n", err)
		// }
		done <- struct{}{}
	}()

	app.Log.Info("App Started")

	<-done

	return nil
}

func (app *App) errorHandler(ctx *fiber.Ctx, err error) error {
	// Status code defaults to 500
	code := fiber.StatusInternalServerError

	// Retrieve the custom status code if it's a *fiber.Error
	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		code = fiberErr.Code
	}

	app.Log.Errorf("API panic recovered: %v", err)

	return ctx.Status(code).JSON(resp.New(false, nil, e.ErrInternalServer))
}
