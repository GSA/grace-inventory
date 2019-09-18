package sessionmgr

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

// SessionMgr ... holds one session per region provided to New()
type SessionMgr struct {
	defaultRegion string
	sessions      []*session.Session
}

// New ... returns a *SessionMgr after creating one session per region provided
func New(defaultRegion string, regions []string) (*SessionMgr, error) {
	mgr := &SessionMgr{defaultRegion: defaultRegion}
	for _, r := range regions {
		sess, err := session.NewSession(&aws.Config{Region: aws.String(r)})
		if err != nil {
			return nil, err
		}
		mgr.sessions = append(mgr.sessions, sess)
	}
	return mgr, nil
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
	sess, err := session.NewSession(&aws.Config{
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
	sess, err := session.NewSession(&aws.Config{
		Region: &region,
	})
	if err != nil {
		return nil, err
	}
	return sess, nil
}
