// +build integration

package sessionmgr

import (
	"testing"
)

// func New(defaultRegion string, regions []string) (*SessionMgr, error)
func TestIntegrationNew(t *testing.T) {
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
			if s.defaultRegion != tc.defaultRegion {
				t.Fatalf("defaultRegion invalid expected: %s, got: %s", tc.defaultRegion, s.defaultRegion)
			}
			if len(s.regions) != len(tc.regions) {
				t.Fatalf("regions length invalid, expected: %d, got: %d", len(tc.regions), len(s.regions))
			}
		})
	}
}

//func (mgr *SessionMgr) All() []*session.Session
func TestIntegrationAll(t *testing.T) {
	s := New(genericRegions[0], genericRegions)
	result := len(s.All())
	if len(s.sessions) != result {
		t.Fatalf("All() length invalid, expected: %d, got: %d", len(genericRegions), result)
	}
}

//func (mgr *SessionMgr) Default() (*session.Session, error)
func TestIntegrationDefault(t *testing.T) {
	s := New(genericRegions[0], genericRegions)
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
func TestIntegrationRegion(t *testing.T) {
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
			s := New(tc.region, tc.regions)
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
