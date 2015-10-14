package bridge

import (
	"fmt"
	"github.com/dropbox/godropbox/container/set"
)

var (
	reservedBridges set.Set
)

func init() {
	reservedBridges = set.NewSet()
}

func reserveBridge() (bridge string) {
	for i := 0; ; i++ {
		bridge = fmt.Sprintf("brc%d", i)
		if !reservedBridges.Contains(bridge) {
			reservedBridges.Add(bridge)
			return
		}
	}

	return
}

func releaseBridge(bridge string) {
	reservedBridges.Remove(bridge)
}
