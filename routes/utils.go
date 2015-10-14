package routes

import (
	"fmt"
	"github.com/dropbox/godropbox/container/set"
)

var (
	reservedTables set.Set
)

type table struct {
	Num  int
	Name string
}

func init() {
	reservedTables = set.NewSet()
}

func reserveTable() (tbl *table) {
	for i := 0; ; i++ {
		num := 9700 + i
		if !reservedTables.Contains(num) {
			reservedTables.Add(num)

			tbl = &table{
				Num:  num,
				Name: fmt.Sprintf("cortunl%d", i),
			}
			return
		}
	}

	return
}

func releaseTable(tbl *table) {
	reservedTables.Remove(tbl.Num)
}
