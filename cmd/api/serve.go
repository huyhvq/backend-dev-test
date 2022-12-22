package api

import (
	"github.com/huyhvq/backend-dev-test/internal/database"
	"github.com/huyhvq/backend-dev-test/internal/organizations/repositories"
	"github.com/huyhvq/backend-dev-test/internal/server"
	"github.com/huyhvq/backend-dev-test/pkg/leveledlog"
	"github.com/huyhvq/backend-dev-test/pkg/version"
	"github.com/spf13/cobra"
	"os"
)

var cfgFile string

var ServeCmd = &cobra.Command{
	Use:   "serve",
	Short: "serve APIs",
	Long:  "serve APIs server",
	Run:   serve,
}

func init() {
	ServeCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config.yaml)")
}

type application struct {
	config       config
	db           *database.DB
	logger       *leveledlog.Logger
	repositories repositories.Manager
}

func serve(cmd *cobra.Command, args []string) {
	cfg := newConfig(cfgFile)
	logger := leveledlog.NewLogger(os.Stdout, leveledlog.LevelAll, true)
	db, err := database.New(cfg.DB.DSN, cfg.DB.AutoMigrate)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	r := repositories.New(db.DB)

	app := &application{
		config:       cfg,
		db:           db,
		logger:       logger,
		repositories: r,
	}
	logger.Info("starting server on %s (version %s)", cfg.Addr, version.Get())

	if err := server.Run(cfg.Addr, app.routes()); err != nil {
		logger.Fatal(err)
	}
	logger.Info("server stopped")
}
