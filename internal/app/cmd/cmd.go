package cmd

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	killgrave "github.com/friendsofgo/killgrave/internal"
	server "github.com/friendsofgo/killgrave/internal/server/http"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/radovskyb/watcher"
	"github.com/spf13/cobra"
)

var _version = "unknown_version"

const (
	_defaultImpostersPath = "imposters"
	_defaultConfigFile    = ""
	_defaultHost          = "localhost"
	_defaultPort          = 3000
	_defaultProxyMode     = killgrave.ProxyNone
	_defaultStrictSlash   = true

	_impostersFlag = "imposters"
	_configFlag    = "config"
	_hostFlag      = "host"
	_portFlag      = "port"
	_watcherFlag   = "watcher"
	_secureFlag    = "secure"
	_proxyModeFlag = "proxy-mode"
	_proxyURLFlag  = "proxy-url"
)

var (
	errGetDataFromImpostersFlag = errors.New("error trying to get data from imposters flag")
	errGetDataFromHostFlag      = errors.New("error trying to get data from host flag")
	errGetDataFromPortFlag      = errors.New("error trying to get data from port flag")
	errGetDataFromSecureFlag    = errors.New("error trying to get data from secure flag")
	errMandatoryURL             = errors.New("the field proxy-url is mandatory if you selected a proxy mode")
)

// NewKillgraveCmd returns cobra.Command to run killgrave command
func NewKillgraveCmd() *cobra.Command {
	var cfg killgrave.Config

	rootCmd := &cobra.Command{
		Use:           "killgrave",
		Short:         "Simple way to generate mock servers",
		SilenceErrors: true,
		SilenceUsage:  true,
		Version:       _version,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			var err error
			cfg, err = prepareConfig(cmd)
			if err != nil {
				return err
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runHTTP(cmd, cfg)
		},
	}

	rootCmd.ResetFlags()
	rootCmd.PersistentFlags().StringP(_impostersFlag, "i", _defaultImpostersPath, "Directory where your imposters are located")
	rootCmd.PersistentFlags().StringP(_configFlag, "c", _defaultConfigFile, "Path to your configuration file")
	rootCmd.Flags().StringP(_hostFlag, "H", _defaultHost, "Set a different host than localhost")
	rootCmd.Flags().IntP(_portFlag, "P", _defaultPort, "Port to run the server")
	rootCmd.Flags().BoolP(_watcherFlag, "w", false, "File watcher will reload the server on each file change")
	rootCmd.Flags().BoolP(_secureFlag, "s", false, "Run mock server using TLS (https)")
	rootCmd.Flags().StringP(_proxyModeFlag, "m", _defaultProxyMode.String(), "Proxy mode, the options are all, missing or none")
	rootCmd.Flags().StringP(_proxyURLFlag, "u", "", "The url where the proxy will redirect to")

	rootCmd.SetVersionTemplate("Killgrave version: {{.Version}}\n")

	return rootCmd
}

func runHTTP(cmd *cobra.Command, cfg killgrave.Config) error {
	done := make(chan os.Signal, 1)
	defer close(done)

	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	srv := runServer(cfg)

	watcherFlag, _ := cmd.Flags().GetBool(_watcherFlag)
	if watcherFlag || cfg.Watcher {
		w, err := runWatcher(cfg, &srv)
		if err != nil {
			return err
		}

		defer killgrave.CloseWatcher(w)
	}

	<-done
	if err := srv.Shutdown(); err != nil {
		log.Fatal(err)
	}

	return nil
}

// TODO: refactor the method NewServer of the pkg server/http should be contain how to initialize the http server
func runServer(cfg killgrave.Config) server.Server {
	router := mux.NewRouter().StrictSlash(_defaultStrictSlash)
	httpAddr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	httpServer := http.Server{
		Addr:    httpAddr,
		Handler: handlers.CORS(server.PrepareAccessControl(cfg.CORS)...)(router),
	}

	proxyServer, err := server.NewProxy(cfg.Proxy.Url, cfg.Proxy.Mode)
	if err != nil {
		log.Fatal(err)
	}

	imposterFs, err := server.NewImposterFS(cfg.ImpostersPath)
	if err != nil {
		log.Fatal(err)
	}

	s := server.NewServer(
		router,
		&httpServer,
		proxyServer,
		cfg.Secure,
		imposterFs,
	)
	if err := s.Build(); err != nil {
		log.Fatal(err)
	}

	s.Run()
	return s
}

func runWatcher(cfg killgrave.Config, currentSrv *server.Server) (*watcher.Watcher, error) {
	w, err := killgrave.InitializeWatcher(cfg.ImpostersPath)
	if err != nil {
		return nil, err
	}

	killgrave.AttachWatcher(w, func() {
		if err := currentSrv.Shutdown(); err != nil {
			log.Fatal(err)
		}
		*currentSrv = runServer(cfg)
	})
	return w, nil
}

func prepareConfig(cmd *cobra.Command) (killgrave.Config, error) {
	cfgPath, _ := cmd.Flags().GetString("config")
	if cfgPath != "" {
		return killgrave.NewConfigFromFile(cfgPath)
	}

	impostersPath, err := cmd.Flags().GetString("imposters")
	if err != nil {
		return killgrave.Config{}, fmt.Errorf("%v: %w", err, errGetDataFromImpostersFlag)
	}

	host, err := cmd.Flags().GetString(_hostFlag)
	if err != nil {
		return killgrave.Config{}, fmt.Errorf("%v: %w", err, errGetDataFromHostFlag)
	}

	port, err := cmd.Flags().GetInt(_portFlag)
	if err != nil {
		return killgrave.Config{}, fmt.Errorf("%v: %w", err, errGetDataFromPortFlag)
	}

	secure, err := cmd.Flags().GetBool(_secureFlag)
	if err != nil {
		return killgrave.Config{}, fmt.Errorf("%v: %w", err, errGetDataFromSecureFlag)
	}

	cfg, err := killgrave.NewConfig(impostersPath, host, port, secure)
	if err != nil {
		return killgrave.Config{}, err
	}

	return cfg, configureProxyMode(cmd, &cfg)
}

func configureProxyMode(cmd *cobra.Command, cfg *killgrave.Config) error {
	mode, err := cmd.Flags().GetString(_proxyModeFlag)
	if err != nil {
		return err
	}

	pMode, err := killgrave.StringToProxyMode(mode)
	if err != nil {
		return err
	}

	var url string
	if mode != killgrave.ProxyNone.String() {
		url, err = cmd.Flags().GetString(_proxyURLFlag)
		if err != nil {
			return err
		}

		if url == "" {
			return errMandatoryURL
		}
	}
	cfg.ConfigureProxy(pMode, url)
	return nil
}
