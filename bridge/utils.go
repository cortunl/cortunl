package bridge

import (
	"fmt"
	"github.com/dropbox/godropbox/container/set"
)

var (
	reservedBridges set.Set
)

func reserveBridge() (bridge string) {
	for i := 0; ; i++ {
		bridge = fmt.Sprintf("brc%d", i)
		if !reservedBridges.Contains(bridge) {
			return
		}
	}

	return
}

func releaseBridge(bridge string) {
	reservedBridges.Remove(bridge)
}
