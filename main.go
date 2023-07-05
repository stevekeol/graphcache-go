package main

import (
	"flag"
	"time"

	"github.com/CrocSwap/graphcache-go/cache"
	"github.com/CrocSwap/graphcache-go/controller"
	"github.com/CrocSwap/graphcache-go/loader"
	"github.com/CrocSwap/graphcache-go/server"
	"github.com/CrocSwap/graphcache-go/utils"
	"github.com/CrocSwap/graphcache-go/views"
)

var uniswapCandles = utils.GoDotEnvVariable("UNISWAP_CANDLES") == "true"
var testNetEnv = utils.GoDotEnvVariable("TESTNET") == "true"

func main() {
	var networkFile string
	if testNetEnv {
		networkFile = "./config/testNetwork.json"
	} else {
		networkFile = "./config/networks.json"
	}

	var netCfgPath = flag.String("netCfg", networkFile, "network config file")
	flag.Parse()

	netCfg := loader.LoadNetworkConfig(*netCfgPath)
	onChain := loader.OnChainLoader{Cfg: netCfg}

	cache := cache.New()
	cntrl := controller.New(netCfg, cache)

	for network, chainCfg := range netCfg {
		startTime := int(time.Now().Unix())
		controller.NewSubgraphSyncer(cntrl, chainCfg, network, startTime)
		if uniswapCandles {
			controller.NewUniswapSyncer(cntrl, chainCfg, network, startTime)
		}
	}

	views := views.Views{Cache: cache, OnChain: &onChain}
	apiServer := server.APIWebServer{Views: &views}
	apiServer.Serve()
}
