package main

import (
	"flag"
	log "github.com/golang/glog"
)

var (
	configFile string
)

func init() {
	flag.StringVar(&configFile, "c", "./pitchfork.conf", " set pitchfork config file path")
}

func main() {
	flag.Parse()
	defer log.Flush()

	log.Infof("bfs pitchfork start")

	var (
		config     *Config
		zk         *Zookeeper
		p          *Pitchfork
		err        error
	)

	if config, err = NewConfig(configFile); err != nil {
		log.Errorf("NewConfig(\"%s\") error(%v)", configFile, err)
		return
	}

	log.Infof("init zookeeper...")
	if zk, err = NewZookeeper(config.ZookeeperAddrs, config.ZookeeperTimeout); err != nil {
		log.Errorf("NewZookeeper() failed, Quit now")
		return
	}

	log.Infof("register pitchfork...")
	p = NewPitchfork(zk, config)
	if err = p.Register(); err != nil {
		log.Errorf("pitchfork Register() failed, Quit now")
		return
	}

	log.Infof("starts probe stores...")
	go Work(p)

	StartSignal()
	return
}
