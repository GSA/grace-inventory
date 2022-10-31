package sessionmgr

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

var genericRegions = []string{"us-east-1", "us-west-1"}

func mockNewSession(cfgs ...*aws.Config) (*session.Session, error) {
	// server is the mock server that simply writes a 200 status back to the client
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	for _, cfg := range cfgs {
		cfg.DisableSSL = aws.Bool(true)
		cfg.Endpoint = aws.String(server.URL)
	}
	return session.NewSession(cfgs...)
}

func mockNewSessionErr(cfgs ...*aws.Config) (*session.Session, error) {
	return nil, errors.New("error")
}

// func New(defaultRegion string, regions []string) (*SessionMgr, error)
func TestNew(t *testing.T) {
	tests := []struct {
		defaultRegion string
		regions       []string
	}{
		{genericRegions[0], genericRegions},
	}
	for _, tt := range tests {
		tc := tt
		t.Run("test with 2 regions", func(t *testing.T) {
			s := New(tc.defaultRegion, tc.regions)
			s.Sessioner(mockNewSession)
			err := s.Init()
			if err != nil {
				t.Fatalf("Init() failed: %v", err)
			}
			if s.defaultRegion != tc.defaultRegion {
				t.Fatalf("defaultRegion invalid expected: %s, got: %s", tc.defaultRegion, s.defaultRegion)
			}
			if len(s.sessions) != len(tc.regions) {
				t.Fatalf("sessions length invalid, expected: %d, got: %d", len(tc.regions), len(s.sessions))
			}
		})
	}
}

// func (mgr *SessionMgr) All() []*session.Session
func TestAll(t *testing.T) {
	s := New(genericRegions[0], genericRegions)
	s.Sessioner(mockNewSession)
	err := s.Init()
	if err != nil {
		t.Fatalf("Init() failed: %v", err)
	}
	result := len(s.All())
	if len(s.sessions) != result {
		t.Fatalf("All() length invalid, expected: %d, got: %d", len(genericRegions), result)
	}
}

// func (mgr *SessionMgr) Default() (*session.Session, error)
func TestDefault(t *testing.T) {
	s := New(genericRegions[0], genericRegions)
	s.Sessioner(mockNewSession)
	err := s.Init()
	if err != nil {
		t.Fatalf("Init() failed: %v", err)
	}
	ss, err := s.Default()
	if err != nil {
		t.Fatalf("Default() failed: %v", err)
	}
	if ss.Config == nil || ss.Config.Region == nil {
		t.Fatal("ss.Config or ss.Config.Region is nil")
	}
	result := *ss.Config.Region
	if result != genericRegions[0] {
		t.Fatalf("Default() region invalid, expected: %s, got: %s", genericRegions[0], result)
	}
}

func TestDefaultAltRegion(t *testing.T) {
	s := New(genericRegions[0], genericRegions)
	s.Sessioner(mockNewSession)
	err := s.Init()
	if err != nil {
		t.Fatalf("Init() failed: %v", err)
	}
	altRegion := "us-east-2"
	s.defaultRegion = altRegion
	ss, err := s.Default()
	if err != nil {
		t.Fatalf("Default() failed: %v", err)
	}
	if ss.Config == nil || ss.Config.Region == nil {
		t.Fatal("ss.Config or ss.Config.Region is nil")
	}
	result := *ss.Config.Region
	if result != altRegion {
		t.Fatalf("Default() region invalid, expected: %s, got: %s", genericRegions[0], result)
	}
}

func TestDefaultErr(t *testing.T) {
	s := New(genericRegions[0], genericRegions)
	s.Sessioner(mockNewSessionErr)
	_, err := s.Default()
	if err == nil {
		t.Errorf("failure expected but err was nil")
	}
}

// func (mgr *SessionMgr) Region(region string) (*session.Session, error)
func TestRegion(t *testing.T) {
	tests := []struct {
		name    string
		region  string
		regions []string
	}{
		{"test empty region", "", genericRegions},
		{"test good region", genericRegions[0], genericRegions},
	}
	for _, tt := range tests {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			s := New(tc.region, tc.regions)
			s.Sessioner(mockNewSession)
			err := s.Init()
			if err != nil {
				t.Fatalf("Init() failed: %v", err)
			}
			ss, err := s.Region(tc.region)
			if err != nil {
				t.Fatalf("Region() failed: %v", err)
			}
			if ss.Config == nil || ss.Config.Region == nil {
				t.Fatal("ss.Config or ss.Config.Region is nil")
			}
			result := *ss.Config.Region
			if result != tc.region {
				t.Fatalf("Region() region invalid, expected: %s, got: %s", tc.region, result)
			}
		})
	}
}

func TestRegionErr(t *testing.T) {
	s := New(genericRegions[0], genericRegions)
	s.Sessioner(mockNewSessionErr)
	_, err := s.Region("")
	if err == nil {
		t.Errorf("failure expected but err was nil")
	}
}
