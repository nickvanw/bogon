package main

import (
	"os"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/nickvanw/bogon"
	"github.com/nickvanw/bogon/commands/config"
	"github.com/nickvanw/bogon/commands/config/viperprovider"
	"github.com/nickvanw/bogon/commands/plugins"
	"github.com/nickvanw/bogon/commands/plugins/bookmarks"
	"github.com/nickvanw/bogon/commands/plugins/bookmarks/boltdb"
	"github.com/nickvanw/bogon/commands/plugins/twitter"
	"github.com/nickvanw/ircx"
	"github.com/spf13/viper"
)

func main() {
	app := cli.NewApp()
	app.Name = "bogon"
	app.Action = realMain
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "name, n",
			Value:  "ircx",
			Usage:  "IRC Name to use",
			EnvVar: "BOGON_NAME",
		},
		cli.StringFlag{
			Name:   "server, s",
			Value:  "chat.freenode.org:6667",
			Usage:  "IRC host:port to connect to",
			EnvVar: "BOGON_SERVER",
		},
		cli.StringFlag{
			Name:   "user, u",
			Value:  "",
			Usage:  "User in User:Pass to use when connecting, must have both",
			EnvVar: "BOGON_USER",
		},
		cli.StringFlag{
			Name:   "pass, p",
			Value:  "",
			Usage:  "Password in User:Pass to use when connecting, must have both",
			EnvVar: "BOGON_PASSWORD",
		},
		cli.StringFlag{
			Name:   "channels, c",
			Usage:  "channels to join on startup",
			Value:  "#test",
			EnvVar: "BOGON_CHANNELS",
		},
		cli.StringFlag{
			Name:   "config, cfg",
			Usage:  "config file",
			Value:  "",
			EnvVar: "BOGON_CONFIG",
		},
		cli.StringFlag{
			Name:   "bookmark, bm",
			Usage:  "location of bookmark db.",
			Value:  "bm.bdb",
			EnvVar: "BOGON_BOOKMARK_FILE",
		},
		cli.StringFlag{
			Name:   "bm.serve",
			Usage:  "IP:Port to serve the bookmark viewer",
			Value:  ":9001",
			EnvVar: "BOGON_BOOKMARK_SERVER",
		},
		cli.StringFlag{
			Name:   "admin, a",
			Usage:  "location of admin socket",
			Value:  "",
			EnvVar: "BOGON_ADMIN_SOCKET",
		},
		cli.BoolFlag{
			Name:   "debug",
			Usage:  "enable debug logging",
			EnvVar: "BOGON_DEBUG",
		},
	}
	app.Run(os.Args)
}

func realMain(c *cli.Context) {
	// Create a logger
	logger := log.NewLogfmtLogger(os.Stderr)
	if c.Bool("debug") {
		logger = level.NewFilter(logger, level.AllowAll())
	} else {
		logger = level.NewFilter(logger, level.AllowInfo())
	}
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	// Create the underlying connection
	bot := ircx.Classic(c.String("server"), c.String("name"))
	bot.SetLogger(logger)
	bot.Config.MaxRetries = 10
	if usr := c.String("user"); usr != "" {
		bot.Config.User = usr
	}
	bot.Config.Password = c.String("pass")

	channels := strings.Split(c.String("channels"), ",")

	// Create a new bogon
	bogon, err := bogon.New(bot, c.String("name"), channels)
	if err != nil {
		level.Error(bot.Logger()).Log("action", "create", "error", err)
		os.Exit(1)
	}

	// Setup & add config
	viper.AutomaticEnv()

	// Add config file if provided
	if cfg := c.String("config"); cfg != "" {
		viper.SetConfigFile(cfg)
		viper.WatchConfig()
		if err := viper.ReadInConfig(); err != nil {
			level.Error(bot.Logger()).Log("action", "config", "error", err)
			os.Exit(2)
		}
	}

	// Add viper provider
	config.RegisterProvider(viperprovider.V{})

	// Set up commands
	if err := commandSetup(bogon, c); err != nil {
		level.Error(bot.Logger()).Log("action", "command_setup", "error", err)
		os.Exit(3)
	}

	if cfg := c.String("admin"); cfg != "" {
		bogon.AdminSocket(cfg)
	}

	// Connect!
	if err := bogon.Connect(); err != nil {
		level.Error(bot.Logger()).Log("action", "connect", "error", err)
		os.Exit(4)
	}
	level.Info(bot.Logger()).Log("action", "connected", "server", c.String("server"))
	bogon.Start()
}

func commandSetup(bogon *bogon.Client, c *cli.Context) error {
	// Register basic plugins
	bogon.AddCommands(plugins.Exports()...)

	// Register twitter command
	api, err := twitter.NewFromEnv()
	if err != nil && err != twitter.ErrMissingTokens {
		return err
	}
	bogon.AddCommands(api.TwitterHandler())
	bogon.AddCommands(api.RawTwitterHandler())

	// Register bookmark handler
	if bmdb := c.String("bookmark"); bmdb != "" {
		db, err := boltdb.New(bmdb)
		if err != nil {
			return err
		}
		bm, err := bookmarks.New(db)
		if err != nil {
			return err
		}
		bogon.AddCommands(bm.Exports()...)
		bm.Block(bogon.ListCommands())
		if srv := c.String("bm.serve"); srv != "" {
			go bm.Serve(srv)
		}
	}

	return nil
}
