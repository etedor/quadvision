package internal

import (
	"fmt"
	"math/rand"
	"net/url"
	"strings"
	"time"

	"github.com/apex/log"
)

func isDuplicate(m map[int][]string, u string) bool {
	for _, ss := range m {
		for _, s := range ss {
			p, err := url.Parse(u)
			if err != nil {
				panic(err)
			}
			prefix := fmt.Sprintf("%s://%s%s", p.Scheme, p.Host, p.Path)
			if strings.HasPrefix(s, prefix) {
				return true
			}
		}
	}
	return false
}

func GetStreams(cands []Stream) []string {
	high := -999
	m := make(map[int][]string)
	for _, s := range cands {
		ctx := log.WithFields(log.Fields{
			"priority": s.Priority,
			"url":      s.URL,
		})

		if _, ok := m[s.Priority]; !ok {
			var ss []string
			m[s.Priority] = ss
		}

		if s.Priority > high {
			high = s.Priority
		}

		// determine how many higher priority streams we currently have
		total := 0
		for p, ss := range m {
			if p == s.Priority {
				continue
			}
			total += len(ss)
		}

		if (s.Priority < high && len(m[high]) >= 4) || total >= 4 {
			// we found enough higher priority streams,
			break
		}

		u, err := streamlink(s.URL)
		if err != nil {
			ctx.Warn(err.Error())
			continue
		}

		if isDuplicate(m, u) {
			// even with a 'clean' configuration, we may see this if two candidate twitch.tv
			// streams are offline and hosting the same stream
			ctx.Warn("Duplicate stream on this URL")
			continue
		}

		ctx.Info("Found playable stream on this URL")
		m[s.Priority] = append(m[s.Priority], u)
	}

	var out []string
	for _, ss := range m {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(ss), func(i, j int) { ss[i], ss[j] = ss[j], ss[i] })
		for _, s := range ss {
			if len(out) >= 4 {
				break
			}
			out = append(out, s)
		}
	}

	for len(out) < 4 {
		out = append(out, "")
	}
	return out[:4]
}
