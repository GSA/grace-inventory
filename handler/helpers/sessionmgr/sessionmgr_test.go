package sessionmgr

import (
	"testing"
)

var genericRegions = []string{"us-east-1", "us-west-1"}

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
			s, err := New(tc.defaultRegion, tc.regions)
			if err != nil {
				t.Fatalf("New() failed: %v", err)
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

//func (mgr *SessionMgr) All() []*session.Session
func TestAll(t *testing.T) {
	s, err := New(genericRegions[0], genericRegions)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	result := len(s.All())
	if len(s.sessions) != result {
		t.Fatalf("All() length invalid, expected: %d, got: %d", len(genericRegions), result)
	}
}

//func (mgr *SessionMgr) Default() (*session.Session, error)
func TestDefault(t *testing.T) {
	s, err := New(genericRegions[0], genericRegions)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
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

//func (mgr *SessionMgr) Region(region string) (*session.Session, error)
func TestRegion(t *testing.T) {
	tests := []struct {
		name     string
		region   string
		regions  []string
		haserror bool
	}{
		{"test empty region", "", genericRegions, true},
		{"test good region", genericRegions[0], genericRegions, false},
	}
	for _, tt := range tests {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			s, err := New(tc.region, tc.regions)
			if err != nil {
				t.Fatalf("New() failed: %v", err)
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
