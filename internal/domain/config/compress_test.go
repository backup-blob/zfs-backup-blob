package config

import (
	"testing"
)

func TestNewCompressConfig(t *testing.T) {
	conf := NewCompressConfig()
	if conf == nil {
		t.Fatal("expected non-nil config")
	}
	
	cc, ok := conf.(*CompressConfig)
	if !ok {
		t.Fatal("expected *CompressConfig type")
	}
	
	// Check default level
	if cc.Level != 3 {
		t.Errorf("expected default level 3, got %d", cc.Level)
	}
	
	// Check type is Middleware
	if cc.Type() != Middleware {
		t.Errorf("expected type Middleware, got %v", cc.Type())
	}
}

func TestCompressConfig_Type(t *testing.T) {
	conf := &CompressConfig{}
	if conf.Type() != Middleware {
		t.Errorf("expected Middleware type")
	}
}

func TestCompressConfig_Remote(t *testing.T) {
	conf := &CompressConfig{
		Remote_: "s3-storage",
	}
	if conf.Remote() != "s3-storage" {
		t.Errorf("expected remote 's3-storage', got '%s'", conf.Remote())
	}
}
