package swarm

import (
	"context"
	"runtime"
	"sync"
	"testing"
	"time"

	ci "gx/ipfs/QmUJzxQQ2kzwQubsMqBTr1NGDpLfh7pGA2E1oaJULcKDPq/go-testutil/ci"
	ma "gx/ipfs/QmWWQ2Txc2c6tqjsBpzg5Ar652cHPGNsQQp2SejkNmkUMb/go-multiaddr"
	peer "gx/ipfs/QmcJukH2sAFjY3HdBKq35WDzWoL3UUu2gt9wdfqZTUyM74/go-libp2p-peer"
	pstore "gx/ipfs/QmdeiKhUy1TVGBaKxt7y1QmBDLBdisSrLJ1x58Eoj4PXUh/go-libp2p-peerstore"
)

func TestSimultOpen(t *testing.T) {

	t.Parallel()

	ctx := context.Background()
	swarms := makeSwarms(ctx, t, 2)

	// connect everyone
	{
		var wg sync.WaitGroup
		connect := func(s *Swarm, dst peer.ID, addr ma.Multiaddr) {
			// copy for other peer
			log.Debugf("TestSimultOpen: connecting: %s --> %s (%s)", s.local, dst, addr)
			s.peers.AddAddr(dst, addr, pstore.PermanentAddrTTL)
			if _, err := s.Dial(ctx, dst); err != nil {
				t.Fatal("error swarm dialing to peer", err)
			}
			wg.Done()
		}

		log.Info("Connecting swarms simultaneously.")
		wg.Add(2)
		go connect(swarms[0], swarms[1].local, swarms[1].ListenAddresses()[0])
		go connect(swarms[1], swarms[0].local, swarms[0].ListenAddresses()[0])
		wg.Wait()
	}

	for _, s := range swarms {
		s.Close()
	}
}

func TestSimultOpenMany(t *testing.T) {
	// t.Skip("very very slow")

	addrs := 20
	rounds := 10
	if ci.IsRunning() || runtime.GOOS == "darwin" {
		// osx has a limit of 256 file descriptors
		addrs = 10
		rounds = 5
	}
	SubtestSwarm(t, addrs, rounds)
}

func TestSimultOpenFewStress(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	// t.Skip("skipping for another test")
	t.Parallel()

	msgs := 40
	swarms := 2
	rounds := 10
	// rounds := 100

	for i := 0; i < rounds; i++ {
		SubtestSwarm(t, swarms, msgs)
		<-time.After(10 * time.Millisecond)
	}
}
