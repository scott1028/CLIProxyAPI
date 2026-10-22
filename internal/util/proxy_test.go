package util

import (
	"net/http"
	"strings"
	"testing"

	"github.com/router-for-me/CLIProxyAPI/v8/sdk/config"
	log "github.com/sirupsen/logrus"
	logtest "github.com/sirupsen/logrus/hooks/test"
)

func TestSetProxyLogsRedactedProxyURLOnError(t *testing.T) {
	_, hook := logtest.NewNullLogger()
	savedHooks := make(log.LevelHooks)
	for level, hooks := range log.StandardLogger().Hooks {
		savedHooks[level] = append([]log.Hook(nil), hooks...)
	}
	log.AddHook(hook)
	t.Cleanup(func() { log.StandardLogger().ReplaceHooks(savedHooks) })

	client := &http.Client{}
	SetProxy(&config.SDKConfig{ProxyURL: "ftp://user:secret@proxy.example.com:2121"}, client)

	if client.Transport != nil {
		t.Fatalf("Transport = %T, want nil for invalid proxy", client.Transport)
	}
	entry := hook.LastEntry()
	if entry == nil || entry.Level != log.ErrorLevel {
		t.Fatalf("last log entry = %#v, want error entry", entry)
	}
	if !strings.Contains(entry.Message, "ftp://redacted@proxy.example.com:2121") {
		t.Fatalf("log message %q does not name the redacted proxy", entry.Message)
	}
	if strings.Contains(entry.Message, "secret") {
		t.Fatalf("log message %q leaks proxy password", entry.Message)
	}
}
