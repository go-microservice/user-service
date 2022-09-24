/**
 *
 *    ____          __
 *   / __/__ ____ _/ /__
 *  / _// _ `/ _ `/ / -_)
 * /___/\_,_/\_, /_/\__/
 *         /___/
 *
 *
 * generate by http://patorjk.com/software/taag/#p=display&f=Small%20Slant&t=Eagle
 */
package main

import (
	"encoding/json"
	"fmt"

	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	eagle "github.com/go-eagle/eagle/pkg/app"
	"github.com/go-eagle/eagle/pkg/client/consulclient"
	"github.com/go-eagle/eagle/pkg/client/etcdclient"
	"github.com/go-eagle/eagle/pkg/client/nacosclient"
	"github.com/go-eagle/eagle/pkg/config"
	logger "github.com/go-eagle/eagle/pkg/log"
	"github.com/go-eagle/eagle/pkg/registry"
	"github.com/go-eagle/eagle/pkg/registry/consul"
	"github.com/go-eagle/eagle/pkg/registry/etcd"
	"github.com/go-eagle/eagle/pkg/registry/nacos"
	"github.com/go-eagle/eagle/pkg/transport/grpc"
	v "github.com/go-eagle/eagle/pkg/version"
	"github.com/spf13/pflag"
	_ "go.uber.org/automaxprocs"

	"github.com/go-microservice/user-service/internal/server"
)

var (
	cfgDir  = pflag.StringP("config dir", "c", "config", "config path.")
	env     = pflag.StringP("env name", "e", "", "env var name.")
	version = pflag.BoolP("version", "v", false, "show version info.")
)

// @title eagle docs api
// @version 1.0
// @description eagle demo

// @host localhost:8080
// @BasePath /v1
func main() {
	pflag.Parse()
	if *version {
		ver := v.Get()
		marshaled, err := json.MarshalIndent(&ver, "", "  ")
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		fmt.Println(string(marshaled))
		return
	}

	// init config
	c := config.New(*cfgDir, config.WithEnv(*env))
	var cfg eagle.Config
	if err := c.Load("app", &cfg); err != nil {
		panic(err)
	}
	// set global
	eagle.Conf = &cfg

	// -------------- init resource -------------
	logger.Init()

	gin.SetMode(cfg.Mode)

	// init pprof server
	go func() {
		fmt.Printf("Listening and serving PProf HTTP on %s\n", cfg.PprofPort)
		if err := http.ListenAndServe(cfg.PprofPort, http.DefaultServeMux); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen ListenAndServe for PProf, err: %s", err.Error())
		}
	}()

	// start app
	app, cleanup, err := InitApp(&cfg, &cfg.GRPC)
	defer cleanup()
	if err != nil {
		panic(err)
	}
	if err := app.Run(); err != nil {
		panic(err)
	}
}

func newApp(cfg *eagle.Config, gs *grpc.Server) *eagle.App {
	return eagle.New(
		eagle.WithName(cfg.Name),
		eagle.WithVersion(cfg.Version),
		eagle.WithLogger(logger.GetLogger()),
		eagle.WithServer(
			// init HTTP server
			server.NewHTTPServer(&cfg.HTTP),
			// init gRPC server
			gs,
		),
		// eagle.WithRegistry(getConsulRegistry()),
	)
}

// create a etcd register
func getEtcdRegistry() registry.Registry {
	client, err := etcdclient.New()
	if err != nil {
		log.Fatal(err)
	}
	return etcd.New(client.Client)
}

// create a consul register
func getConsulRegistry() registry.Registry {
	client, err := consulclient.New()
	if err != nil {
		panic(err)
	}
	return consul.New(client)
}

// create a nacos register
func getNacosRegistry() registry.Registry {
	client, err := nacosclient.New()
	if err != nil {
		panic(err)
	}
	return nacos.New(client)
}
