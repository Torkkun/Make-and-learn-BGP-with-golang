package bgp

import (
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestConfigParse(t *testing.T) {
	var configMessage = "64513 127.0.0.2 64512 127.0.0.1 passive"
	config, err := ConfigParseFromStr(configMessage)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(config)
	log.Println(string(config.Local_ip_address))
}

func TestCanTransitionToOpenConfirm(t *testing.T) {
	go func() {
		var remoteConfigMessage = "64513 127.0.0.2 64512 127.0.0.1 passive"
		remoteConfig, err := ConfigParseFromStr(remoteConfigMessage)
		if err != nil {
			log.Fatalln(err)
		}
		remoteBgpPeer := NewPeer(remoteConfig)

		remoteBgpPeer.start()

		for i := 0; i < 50; i++ {
			remoteBgpPeer.nextStep()
			t := time.NewTicker(1 * time.Second)
			select {
			case <-t.C:
				if remoteBgpPeer.Now_state == OpenConfirm {
					break
				}
			}
			t.Stop()
		}
		require.Equal(t, remoteBgpPeer.Now_state, OpenConfirm)
	}()
	var remoteConfigMessage = "64512 127.0.0.1 64513 127.0.0.2 active"
	remoteConfig, err := ConfigParseFromStr(remoteConfigMessage)
	if err != nil {
		log.Fatalln(err)
	}
	remoteBgpPeer := NewPeer(remoteConfig)

	remoteBgpPeer.start()

	for i := 0; i < 50; i++ {
		remoteBgpPeer.nextStep()
		t := time.NewTicker(1 * time.Second)
		select {
		case <-t.C:
			if remoteBgpPeer.Now_state == OpenConfirm {
				break
			}
		}
		t.Stop()
	}
	require.Equal(t, remoteBgpPeer.Now_state, OpenConfirm)
}

func TestPeerCanTransitionToOpenSentStart(t *testing.T) {
	go func() {
		var remoteConfigMessage = "64513 127.0.0.2 64512 127.0.0.1 passive"
		remoteConfig, err := ConfigParseFromStr(remoteConfigMessage)
		if err != nil {
			log.Fatalln(err)
		}
		remoteBgpPeer := NewPeer(remoteConfig)

		remoteBgpPeer.start()
		require.Equal(t, remoteBgpPeer.Now_state, Idle)

		remoteBgpPeer.nextStep()
		log.Println(remoteBgpPeer.Now_state)
		require.Equal(t, remoteBgpPeer.Now_state, Connect)

		remoteBgpPeer.nextStep()
		require.Equal(t, remoteBgpPeer.Now_state, OpenSent)
	}()

	time.Sleep(2 * time.Second)

	go func() {
		var localConfigMessage = "64512 127.0.0.1 64513 127.0.0.2 active"
		localConfig, err := ConfigParseFromStr(localConfigMessage)
		if err != nil {
			log.Fatalln(err)
		}
		localBgpPeer := NewPeer(localConfig)

		localBgpPeer.start()
		require.Equal(t, localBgpPeer.Now_state, Idle)

		localBgpPeer.nextStep()
		require.Equal(t, localBgpPeer.Now_state, Connect)

		localBgpPeer.nextStep()
		require.Equal(t, localBgpPeer.Now_state, OpenSent)
	}()
}

/* func TestPeerCanTransitionToConnectStart(t *testing.T) {
	var configMessage = "64512 127.0.0.1 64513 127.0.0.2 active"
	config, err := ConfigParseFromStr(configMessage)
	if err != nil {
		log.Fatalln(err)
	}
	bgpPeer := NewPeer(config)
	bgpPeer.start()
	bgpPeer.nextStep()
	require.Equal(t, bgpPeer.Now_state, Connect)
} */
