package main

//import (
//	// _ "github.com/coredns/coredns/plugin/log"
//	_ "github.com/coredns/coredns/plugin/whoami"
//	_ "gophers.dev/cmds/donutdns/plugins/donutdns"
//
//	"github.com/coredns/coredns/core/dnsserver"
//	"github.com/coredns/coredns/coremain"
//)
//
//var directives = []string{
//	"donutdns",
//	"whoami",
//	"startup",
//	"shutdown",
//}
//
//func init() {
//	dnsserver.Directives = directives
//}
//
//func main() {
//	coremain.Run()
//}
//
//

import (
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/coremain"
	_ "github.com/coredns/coredns/plugin/debug"
	_ "github.com/coredns/coredns/plugin/log"
	_ "github.com/coredns/coredns/plugin/whoami"
	_ "github.com/coredns/example"
	_ "gophers.dev/cmds/donutdns/plugins/donutdns"
)

var directives = []string{
	"donutdns",
	"debug",
	"example",
	"log",
	"whoami",
	"startup",
	"shutdown",
}

func init() {
	dnsserver.Directives = directives
}

func main() {
	coremain.Run()
}
