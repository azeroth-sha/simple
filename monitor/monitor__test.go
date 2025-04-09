package monitor_test

import (
	"github.com/azeroth-sha/simple/codec"
	"github.com/azeroth-sha/simple/monitor"
	"testing"
)

func TestSYS(t *testing.T) {
	//if stat := monitor.MustCPUStats(); stat != nil {
	//	t.Log(string(codec.MustMarshal(codec.Json, stat)))
	//} else {
	//	t.Error("stat is nil")
	//}
	//
	//if stat := monitor.MustDiskStats(); stat != nil {
	//	t.Log(string(codec.MustMarshal(codec.Json, stat)))
	//} else {
	//	t.Error("stat is nil")
	//}
	//
	//if stat := monitor.MustHostStats(); stat != nil {
	//	t.Log(string(codec.MustMarshal(codec.Json, stat)))
	//} else {
	//	t.Error("stat is nil")
	//}
	//
	//if stat := monitor.MustLoadStats(); stat != nil {
	//	t.Log(string(codec.MustMarshal(codec.Json, stat)))
	//} else {
	//	t.Error("stat is nil")
	//}
	//
	//if stat := monitor.MustMemStats(); stat != nil {
	//	t.Log(string(codec.MustMarshal(codec.Json, stat)))
	//} else {
	//	t.Error("stat is nil")
	//}
	//
	//if stat := monitor.MustNetStats(); stat != nil {
	//	t.Log(string(codec.MustMarshal(codec.Json, stat)))
	//} else {
	//	t.Error("stat is nil")
	//}

	if stat := monitor.MustSysStats(); stat != nil {
		t.Log(string(codec.MustMarshal(codec.Json, stat)))
	} else {
		t.Error("stat is nil")
	}
}
