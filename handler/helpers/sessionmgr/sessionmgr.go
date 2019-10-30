package sessionmgr

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

// SessionMgr ... holds one session per region provided to New()
type SessionMgr struct {
	defaultRegion string
	regions       []string
	fn            Sessioner
	sessions      []*session.Session
}

// Sessioner returns a new *session.Session and an error
type Sessioner func(cfgs ...*aws.Config) (*session.Session, error)

// New ... returns a *SessionMgr
func New(defaultRegion string, regions []string) *SessionMgr {
	return &SessionMgr{defaultRegion: defaultRegion, regions: regions, fn: session.NewSession}
}

// Init ... creates one session per region provided to New()
func (mgr *SessionMgr) Init() error {
	for _, r := range mgr.regions {
		sess, err := mgr.fn(&aws.Config{Region: aws.String(r)})
		if err != nil {
			return err
		}
		mgr.sessions = append(mgr.sessions, sess)
	}
	return nil
}

// Sessioner ... sets the method to use for creating new sessions
func (mgr *SessionMgr) Sessioner(sessioner Sessioner) {
	mgr.fn = sessioner
}

// All ... returns all sessions stored inside the *SessionMgr
func (mgr *SessionMgr) All() []*session.Session {
	return mgr.sessions
}

// Default ... returns the session matching the defaultRegion
func (mgr *SessionMgr) Default() (*session.Session, error) {
	for _, s := range mgr.sessions {
		if *s.Config.Region == mgr.defaultRegion {
			return s, nil
		}
	}
	sess, err := mgr.fn(&aws.Config{
		Region: aws.String(mgr.defaultRegion),
	})
	if err != nil {
		return nil, err
	}
	return sess, nil
}

// Region ... returns the session whose region matches 'region'
func (mgr *SessionMgr) Region(region string) (*session.Session, error) {
	for _, s := range mgr.sessions {
		if *s.Config.Region == region {
			return s, nil
		}
	}
	sess, err := mgr.fn(&aws.Config{
		Region: &region,
	})
	if err != nil {
		return nil, err
	}
	return sess, nil
}
