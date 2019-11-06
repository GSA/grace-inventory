package helpers

import (
	"errors"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/rds/rdsiface"
)

type mockedRDS struct {
	rdsiface.RDSAPI
	mocked

	describeDBInstancesPages []*rds.DescribeDBInstancesOutput
	describeDBInstancesErr   error

	describeDBSnapshotsPages []*rds.DescribeDBSnapshotsOutput
	describeDBSnapshotsErr   error
}

func (m mockedRDS) DescribeDBInstancesPages(inp *rds.DescribeDBInstancesInput, f func(*rds.DescribeDBInstancesOutput, bool) bool) error {
	m.Called("DescribeDBInstancesPages")
	for i, p := range m.describeDBInstancesPages {
		f(p, (i == (len(m.describeDBInstancesPages) - 1)))
	}
	return m.describeDBInstancesErr
}

func (m mockedRDS) DescribeDBSnapshotsPages(inp *rds.DescribeDBSnapshotsInput, f func(*rds.DescribeDBSnapshotsOutput, bool) bool) error {
	m.Called("DescribeDBSnapshotsPages")
	for i, p := range m.describeDBSnapshotsPages {
		f(p, (i == (len(m.describeDBSnapshotsPages) - 1)))
	}
	return m.describeDBSnapshotsErr
}

func TestRDSSvc_DBInstances(t *testing.T) {
	dbname1 := "tstdbname1"
	dbname2 := "tstdbname2"
	dbname3 := "tstdbname3"

	instPage1 := []*rds.DBInstance{
		{
			DBName: &dbname1,
		},
	}
	instPage2 := []*rds.DBInstance{
		{
			DBName: &dbname2,
		},
		{
			DBName: &dbname3,
		},
	}
	pages := append(instPage1, instPage2...)

	tests := []struct {
		name       string
		dbInsPages []*rds.DescribeDBInstancesOutput
		descrErr   error

		want      []*rds.DBInstance
		wantErr   bool
		wantCalls []string
	}{
		{
			name:     "error",
			descrErr: errors.New("tst error"),

			wantErr:   true,
			wantCalls: []string{"DescribeDBInstancesPages"},
		},
		{
			name:       "ok empty",
			dbInsPages: []*rds.DescribeDBInstancesOutput{},

			want:      nil, // why not []*rds.DBInstance{}, ?
			wantCalls: []string{"DescribeDBInstancesPages"},
		},
		{
			name: "ok 1 page",
			dbInsPages: []*rds.DescribeDBInstancesOutput{
				{
					DBInstances: instPage1,
				},
			},

			want:      instPage1,
			wantCalls: []string{"DescribeDBInstancesPages"},
		},
		{
			name: "ok 2 pages",
			dbInsPages: []*rds.DescribeDBInstancesOutput{
				{
					DBInstances: instPage1,
				},
				{
					DBInstances: instPage2,
				},
			},

			want:      pages,
			wantCalls: []string{"DescribeDBInstancesPages"},
		},
	}
	for _, tt := range tests {
		tc := tt
		t.Run(tt.name, func(t *testing.T) {
			m := mockedRDS{
				describeDBInstancesPages: tc.dbInsPages,
				describeDBInstancesErr:   tc.descrErr,

				mocked: mocked{mockCalls: &mockCalls{}},
			}

			svc := RDSSvc{
				Client: &m,
			}
			got, err := svc.DBInstances()
			if (err != nil) != tc.wantErr {
				t.Errorf("RDSSvc.DBInstances() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("RDSSvc.DBInstances() = %v, want %v", got, tc.want)
			}
			if !reflect.DeepEqual(m.CallsList(), tc.wantCalls) {
				t.Errorf("Call list mismatch: %v, expected %v", m.CallsList(), tc.wantCalls)
			}
		})
	}
}

func TestRDSSvc_DBSnapshots(t *testing.T) {
	arn1 := "tstarn1"
	arn2 := "tstarn2"
	arn3 := "tstarn3"

	snapPage1 := []*rds.DBSnapshot{
		{
			DBSnapshotArn: &arn1,
		},
	}
	snapPage2 := []*rds.DBSnapshot{
		{
			DBSnapshotArn: &arn2,
		},
		{
			DBSnapshotArn: &arn3,
		},
	}
	pages := append(snapPage1, snapPage2...)

	tests := []struct {
		name        string
		dbSnapPages []*rds.DescribeDBSnapshotsOutput
		descrErr    error

		want      []*rds.DBSnapshot
		wantErr   bool
		wantCalls []string
	}{
		{
			name:     "error",
			descrErr: errors.New("tst error"),

			wantErr:   true,
			wantCalls: []string{"DescribeDBSnapshotsPages"},
		},
		{
			name:        "ok empty",
			dbSnapPages: []*rds.DescribeDBSnapshotsOutput{},

			want:      nil,
			wantCalls: []string{"DescribeDBSnapshotsPages"},
		},
		{
			name: "ok 1 page",
			dbSnapPages: []*rds.DescribeDBSnapshotsOutput{
				{
					DBSnapshots: snapPage1,
				},
			},

			want:      snapPage1,
			wantCalls: []string{"DescribeDBSnapshotsPages"},
		},
		{
			name: "ok 2 pages",
			dbSnapPages: []*rds.DescribeDBSnapshotsOutput{
				{
					DBSnapshots: snapPage1,
				},
				{
					DBSnapshots: snapPage2,
				},
			},

			want:      pages,
			wantCalls: []string{"DescribeDBSnapshotsPages"},
		},
	}
	for _, tt := range tests {
		tc := tt
		t.Run(tt.name, func(t *testing.T) {
			m := mockedRDS{
				describeDBSnapshotsPages: tc.dbSnapPages,
				describeDBSnapshotsErr:   tc.descrErr,

				mocked: mocked{mockCalls: &mockCalls{}},
			}

			svc := RDSSvc{
				Client: &m,
			}
			got, err := svc.DBSnapshots()
			if (err != nil) != tc.wantErr {
				t.Errorf("RDSSvc.DBSnapshots() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("RDSSvc.DBSnapshots() = %v, want %v", got, tc.want)
			}
			if !reflect.DeepEqual(m.CallsList(), tc.wantCalls) {
				t.Errorf("Call list mismatch: %v, expected %v", m.CallsList(), tc.wantCalls)
			}
		})
	}
}
