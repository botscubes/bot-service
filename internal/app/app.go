package app

import (
	"errors"

	"github.com/goccy/go-json"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	e "github.com/botscubes/bot-service/internal/api/errors"
	h "github.com/botscubes/bot-service/internal/api/handlers"
	mb "github.com/botscubes/bot-service/internal/broker"
	resp "github.com/botscubes/bot-service/pkg/api_response"
	se "github.com/botscubes/user-service/pkg/service_error"

	"github.com/botscubes/bot-service/internal/bot"
	"github.com/botscubes/bot-service/internal/config"
	"github.com/botscubes/bot-service/internal/database/pgsql"
	rdb "github.com/botscubes/bot-service/internal/database/redis"
	"github.com/botscubes/bot-service/internal/database/redisauth"
	"github.com/botscubes/user-service/pkg/token_storage"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/gofiber/fiber/v2"
)

type App struct {
	server         *fiber.App
	botService     *bot.BotService
	conf           *config.ServiceConfig
	db             *pgsql.Db
	sessionStorage token_storage.TokenStorage
	redisAuth      *redis.Client
	redis          *rdb.Rdb
	log            *zap.SugaredLogger
	mb             mb.Broker
}

func CreateApp(logger *zap.SugaredLogger, c *config.ServiceConfig, db *pgsql.Db, b mb.Broker) *App {
	redisAuth := redisauth.NewClient(&c.RedisAuth)

	app := &App{
		log:  logger,
		conf: c,
		server: fiber.New(fiber.Config{
			AppName:               "Bot API Server",
			DisableStartupMessage: true,
			JSONEncoder:           json.Marshal,
			JSONDecoder:           json.Unmarshal,
			ErrorHandler:          errorHandler(logger),
		}),
		botService:     bot.NewBotService(c, logger),
		redisAuth:      redisAuth,
		sessionStorage: token_storage.NewRedisTokenStorage(redisAuth),
		redis:          rdb.NewClient(&c.Redis),
		db:             db,
		mb:             b,
	}

	apiHandlers := h.NewApiHandler(
		app.db,
		app.log,
		app.botService,
		app.mb,
		app.redis,
	)

	// CORS
	app.server.Use(cors.New())

	app.regiterHandlers(apiHandlers)
	return app
}

func (app *App) Run() {
	go func() {
		if err := app.server.Listen(app.conf.ListenAddress); err != nil {
			app.log.Fatalw("Start server", "error", err)
		}
	}()
}

func (app *App) Shutdown() error {
	return app.server.ShutdownWithTimeout(config.ShutdownTimeout)
}

func errorHandler(log *zap.SugaredLogger) func(ctx *fiber.Ctx, err error) error {
	return func(ctx *fiber.Ctx, err error) error {
		// Status code defaults to 500
		code := fiber.StatusInternalServerError
		errData := e.ErrInternalServer

		// Retrieve the custom status code if it's a *fiber.Error
		var fiberErr *fiber.Error
		if errors.As(err, &fiberErr) {
			code = fiberErr.Code
			errData = se.New(code, fiberErr.Message)
		}

		log.Errorf("API panic recovered: %v", err)

		return ctx.Status(code).JSON(resp.New(false, nil, errData))
	}
}
