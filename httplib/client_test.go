package httplib_test

import (
	"github.com/azeroth-sha/simple/httplib"
	"testing"
)

func TestClient(t *testing.T) {
	cli := httplib.New().SetDebug(true)
	req := cli.R()
	resp, err := req.Get(`https://www.baidu.com`)
	if err != nil {
		t.Error(err)
	} else if resp.IsError() {
		t.Error(resp.Error())
	} else if buf := resp.Body(); buf == nil {
		t.Error(`body is nil`)
	} else {
		t.Log(string(buf))
	}
}
