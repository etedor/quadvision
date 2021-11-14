package main

import (
	"sort"
	"time"

	"github.com/apex/log"
	qv "github.com/etedor/quadvision/internal"
)

func main() {
	first := true
	for {
		if !first {
			time.Sleep(10 * time.Minute)
		}
		first = false

		cfg := qv.LoadConfig()
		cc, id, err := qv.Login(cfg.Credentials.Username, cfg.Credentials.Secret)
		if err != nil {
			log.Error(err.Error())
			continue
		}

		cands := cfg.Streams
		sort.Slice(cands[:], func(i, j int) bool {
			return cands[i].Priority > cands[j].Priority
		})
		ss := qv.GetStreams(cands)

		err = qv.Update(cc, *id, ss)
		if err != nil {
			log.Error(err.Error())
			continue
		}
	}
}
