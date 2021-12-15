package http

import (
	"errors"
	"fmt"
	"io/ioutil"
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

const (
	_defaultHost        = "localhost"
	_defaultPort        = 3000
	_defaultProxyMode   = killgrave.ProxyNone
	_defaultStrictSlash = true
)

var (
	errGetDataFromImpostersFlag   = errors.New("error trying to get data from imposters flag")
	errGetDataFromHostFlag        = errors.New("error trying to get data from host flag")
	errGetDataFromPortFlag        = errors.New("error trying to get data from port flag")
	errGetDataFromSecureFlag      = errors.New("error trying to get data from secure flag")
	errGetDataFromCertificateFlag = errors.New("error trying to get data from certificate file flag")
	errGetDataFromPrivateKeyFlag  = errors.New("error trying to get data from private key file flag")
	errMandatoryURL               = errors.New("the field url is mandatory if you selected a proxy mode")
)

// NewHTTPCmd returns cobra.Command to run http sub command, this command will be used to run the mock server
func NewHTTPCmd() *cobra.Command {

	var cfg killgrave.Config

	cmd := &cobra.Command{
		Use:   "http",
		Short: "Configure a HTTP mock server based on your imposters",
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
		Args: cobra.NoArgs,
	}

	cmd.PersistentFlags().StringP("host", "H", _defaultHost, "Set a different host than localhost")
	cmd.PersistentFlags().IntP("port", "P", _defaultPort, "Port to run the server")
	cmd.PersistentFlags().BoolP("watcher", "w", false, "File watcher will reload the server on each file change")
	cmd.PersistentFlags().BoolP("secure", "s", false, "Run mock server using TLS (https)")
	cmd.PersistentFlags().StringP("certificate", "c", "", "TLS certificate file name")
	cmd.PersistentFlags().StringP("key", "k", "", "TLS private key file name")
	cmd.Flags().StringP("proxy", "p", _defaultProxyMode.String(), "Proxy mode, the options are all, missing or none")
	cmd.Flags().StringP("url", "u", "", "The url where the proxy will redirect to")

	return cmd
}

func runHTTP(cmd *cobra.Command, cfg killgrave.Config) error {
	done := make(chan os.Signal, 1)
	defer close(done)

	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	srv := runServer(cfg)

	watcherFlag, _ := cmd.Flags().GetBool("watcher")
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

	var cert []byte
	if cfg.PrivateKey != "" {
		cert, _ = ioutil.ReadFile(cfg.Certificate)
	}

	var key []byte
	if cfg.PrivateKey != "" {
		key, _ = ioutil.ReadFile(cfg.PrivateKey)
	}

	s := server.NewServer(
		cfg.ImpostersPath,
		router,
		&httpServer,
		proxyServer,
		cfg.Secure,
		cert,
		key,
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
		runServer(cfg)
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

	host, err := cmd.Flags().GetString("host")
	if err != nil {
		return killgrave.Config{}, fmt.Errorf("%v: %w", err, errGetDataFromHostFlag)
	}

	port, err := cmd.Flags().GetInt("port")
	if err != nil {
		return killgrave.Config{}, fmt.Errorf("%v: %w", err, errGetDataFromPortFlag)
	}

	secure, err := cmd.Flags().GetBool("secure")
	if err != nil {
		return killgrave.Config{}, fmt.Errorf("%v: %w", err, errGetDataFromSecureFlag)
	}

	certificate, err := cmd.Flags().GetString("certificate")
	if err != nil {
		return killgrave.Config{}, fmt.Errorf("%v: %w", err, errGetDataFromCertificateFlag)
	}

	privateKey, err := cmd.Flags().GetString("key")
	if err != nil {
		return killgrave.Config{}, fmt.Errorf("%v: %w", err, errGetDataFromPrivateKeyFlag)
	}

	cfg, err := killgrave.NewConfig(impostersPath, host, port, secure, certificate, privateKey)
	if err != nil {
		return killgrave.Config{}, err
	}

	return cfg, configureProxyMode(cmd, &cfg)
}

func configureProxyMode(cmd *cobra.Command, cfg *killgrave.Config) error {
	mode, err := cmd.Flags().GetString("proxy")
	if err != nil {
		return err
	}

	pMode, err := killgrave.StringToProxyMode(mode)
	if err != nil {
		return err
	}

	var url string
	if mode != killgrave.ProxyNone.String() {
		url, err = cmd.Flags().GetString("url")
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
