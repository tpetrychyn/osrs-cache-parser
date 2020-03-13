package archives

import (
	"github.com/tpetrychyn/osrs-cache-parser/pkg/cachestore"
	"log"
	"strings"
	"testing"
)

func TestInterfaceArchive_LoadInterfaces(t *testing.T) {
	store := cachestore.NewStore("../../cache")
	interfaceArchive := NewInterfaceLoader(store)

	interfaces := interfaceArchive.LoadInterfaces()
	//i := interfaces[193]
	//log.Printf("%+v", i)
	for _, v := range interfaces {
		for _, f := range v {
			if f.ContentType == 1338 {
				log.Printf("hi")
			}
			if strings.Contains(f.Text, "Advanced") {
				log.Printf("hi")
			}
		}
	}
}
