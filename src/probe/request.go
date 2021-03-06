package probe

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"os/exec"
	"time"
)

// Request represents abstract request
type Request struct {
	target string
}

// NewRequest create a Request Object
func NewRequest(target string) (*Request, error) {
	r := Request{
		target: target,
	}
	return &r, nil
}

// Run response request result
func (r *Request) Run() (Response, error) {
	log.Println("curl ", r.target)
	resp := Response{}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // The cancel should be deferred so resources are cleaned up

	fm := `{"NameLookUpTime":"%{time_namelookup}","ConnectTime":"%{time_connect}","AppConnectTime":"%{time_appconnect}","RedirectTime":"%{time_redirect}","PretransferTime":"%{time_pretransfer}","StarttransferTime":"%{time_starttransfer}","TotalTime":"%{time_total}","HTTPCode":"%{http_code}"}`
	ua := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.117 Safari/537.36"
	cmd := exec.CommandContext(ctx, "curl", "--connect-timeout", "3", "-m", "10", "-I", "-s", "-w", fm, "-H", "user-agent:"+ua, r.target, "-o", "/dev/null")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		return resp, err
	}

	if err := json.Unmarshal([]byte(stdout.String()), &resp); err != nil {
		return resp, err
	}
	return resp, nil
}
