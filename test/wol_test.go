package test

import (
	"flag"
	"github.com/DGHeroin/wol"
	"testing"
)
var (
	testSubnet = ""
	testMAC    = ""
)

func init() {
	flag.StringVar(&testSubnet, "subnet", "", "")
	flag.StringVar(&testMAC, "mac", "", "")
	flag.Parse()
}

func TestWakeupWithSubNetwork(t *testing.T) {
	pkg, err := wol.NewMagicPacket(testMAC)
	if err != nil {
		t.Error(err)
		return
	}

	err = pkg.Send(testSubnet)
	if err != nil {
		t.Error(err)
	}
}

func TestWakeupWithSubNetworkPort(t *testing.T) {
	pkg, err := wol.NewMagicPacket(testMAC)
	if err != nil {
		t.Error(err)
		return
	}
	err = pkg.SendPort(testSubnet, "9")
	if err != nil {
		t.Error(err)
	}
}

func TestWakeupBroadcast(t *testing.T) {
	pkg, err := wol.NewMagicPacket(testMAC)
	if err != nil {
		t.Error(err)
		return
	}
	pkg.Broadcast()
}

