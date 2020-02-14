package helpers

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

type mockEc2Client struct {
	ec2iface.EC2API
}

func (m *mockEc2Client) DescribeInstancesPages(in *ec2.DescribeInstancesInput, fn func(*ec2.DescribeInstancesOutput, bool) bool) error {
	fn(&ec2.DescribeInstancesOutput{
		Reservations: []*ec2.Reservation{
			{
				Instances: []*ec2.Instance{{}},
			}}}, true)
	return nil
}

func (m *mockEc2Client) DescribeImages(in *ec2.DescribeImagesInput) (*ec2.DescribeImagesOutput, error) {
	return &ec2.DescribeImagesOutput{Images: []*ec2.Image{{}}}, nil
}

func (m *mockEc2Client) DescribeVolumesPages(in *ec2.DescribeVolumesInput, fn func(*ec2.DescribeVolumesOutput, bool) bool) error {
	fn(&ec2.DescribeVolumesOutput{
		Volumes: []*ec2.Volume{{}},
	}, true)
	return nil
}

func (m *mockEc2Client) DescribeSnapshotsPages(in *ec2.DescribeSnapshotsInput, fn func(*ec2.DescribeSnapshotsOutput, bool) bool) error {
	fn(&ec2.DescribeSnapshotsOutput{
		Snapshots: []*ec2.Snapshot{{}},
	}, true)
	return nil
}

func (m *mockEc2Client) DescribeVpcsPages(in *ec2.DescribeVpcsInput, fn func(*ec2.DescribeVpcsOutput, bool) bool) error {
	fn(&ec2.DescribeVpcsOutput{
		Vpcs: []*ec2.Vpc{{}},
	}, true)
	return nil
}

func (m *mockEc2Client) DescribeSubnetsPages(in *ec2.DescribeSubnetsInput, fn func(*ec2.DescribeSubnetsOutput, bool) bool) error {
	fn(&ec2.DescribeSubnetsOutput{
		Subnets: []*ec2.Subnet{{}},
	}, true)
	return nil
}

func (m *mockEc2Client) DescribeSecurityGroupsPages(in *ec2.DescribeSecurityGroupsInput, fn func(*ec2.DescribeSecurityGroupsOutput, bool) bool) error {
	fn(&ec2.DescribeSecurityGroupsOutput{
		SecurityGroups: []*ec2.SecurityGroup{{}},
	}, true)
	return nil
}

func (m *mockEc2Client) DescribeAddresses(in *ec2.DescribeAddressesInput) (*ec2.DescribeAddressesOutput, error) {
	return &ec2.DescribeAddressesOutput{Addresses: []*ec2.Address{{}}}, nil
}

func (m *mockEc2Client) DescribeKeyPairs(in *ec2.DescribeKeyPairsInput) (*ec2.DescribeKeyPairsOutput, error) {
	return &ec2.DescribeKeyPairsOutput{KeyPairs: []*ec2.KeyPairInfo{{}}}, nil
}

// func Instances() ([]*ec2.Instance, error)
func TestInstances(t *testing.T) {
	svc := Ec2Svc{Client: &mockEc2Client{}}
	expected := []*ec2.Instance{{}}
	got, err := svc.Instances()
	if err != nil {
		t.Fatalf("Instances() failed: %v", err)
	}
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("Instances() failed. Expected: %#v (%T)\nGot: %#v (%T)", expected, expected, got, got)
	}
	_, err = TypeToSheet(expected)
	if err != nil {
		t.Fatalf("TypeToSheet failed: %v", err)
	}
}

// func Images() ([]*ec2.Image, error)
func TestImages(t *testing.T) {
	svc := Ec2Svc{Client: &mockEc2Client{}}
	expected := []*ec2.Image{{}}
	got, err := svc.Images()
	if err != nil {
		t.Fatalf("Images() failed: %v", err)
	}
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("Images() failed. Expected: %#v (%T)\nGot: %#v (%T)", expected, expected, got, got)
	}
	_, err = TypeToSheet(expected)
	if err != nil {
		t.Fatalf("TypeToSheet failed: %v", err)
	}
}

// func Volumes() ([]*ec2.Volume, error)
func TestVolumes(t *testing.T) {
	svc := Ec2Svc{Client: &mockEc2Client{}}
	expected := []*ec2.Volume{{}}
	got, err := svc.Volumes()
	if err != nil {
		t.Fatalf("Volumes() failed: %v", err)
	}
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("Volumes() failed. Expected: %#v (%T)\nGot: %#v (%T)", expected, expected, got, got)
	}
	_, err = TypeToSheet(expected)
	if err != nil {
		t.Fatalf("TypeToSheet failed: %v", err)
	}
}

// func Snapshots() ([]*ec2.Snapshot, error)
func TestSnapshots(t *testing.T) {
	svc := Ec2Svc{Client: &mockEc2Client{}}
	expected := []*ec2.Snapshot{{}}
	got, err := svc.Snapshots()
	if err != nil {
		t.Fatalf("Snapshots() failed: %v", err)
	}
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("Snapshots() failed. Expected: %#v (%T)\nGot: %#v (%T)", expected, expected, got, got)
	}
	_, err = TypeToSheet(expected)
	if err != nil {
		t.Fatalf("TypeToSheet failed: %v", err)
	}
}

// func Vpcs() ([]*ec2.Vpc, error)
func TestVpcs(t *testing.T) {
	svc := Ec2Svc{Client: &mockEc2Client{}}
	expected := []*ec2.Vpc{{}}
	got, err := svc.Vpcs()
	if err != nil {
		t.Fatalf("Vpcs() failed: %v", err)
	}
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("Vpcs() failed. Expected: %#v (%T)\nGot: %#v (%T)", expected, expected, got, got)
	}
	_, err = TypeToSheet(expected)
	if err != nil {
		t.Fatalf("TypeToSheet failed: %v", err)
	}
}

// func Subnets() ([]*ec2.Subnet, error)
func TestSubnets(t *testing.T) {
	svc := Ec2Svc{Client: &mockEc2Client{}}
	expected := []*ec2.Subnet{{}}
	got, err := svc.Subnets()
	if err != nil {
		t.Fatalf("Subnets() failed: %v", err)
	}
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("Subnets() failed. Expected: %#v (%T)\nGot: %#v (%T)", expected, expected, got, got)
	}
	_, err = TypeToSheet(expected)
	if err != nil {
		t.Fatalf("TypeToSheet failed: %v", err)
	}
}

// func SecurityGroups() ([]*ec2.SecurityGroup, error)
func TestSecurityGroups(t *testing.T) {
	svc := Ec2Svc{Client: &mockEc2Client{}}
	expected := []*ec2.SecurityGroup{{}}
	got, err := svc.SecurityGroups()
	if err != nil {
		t.Fatalf("SecurityGroups() failed: %v", err)
	}
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("SecurityGroups() failed. Expected: %#v (%T)\nGot: %#v (%T)", expected, expected, got, got)
	}
	_, err = TypeToSheet(expected)
	if err != nil {
		t.Fatalf("TypeToSheet failed: %v", err)
	}
}

// func Addresses() ([]*ec2.Address, error)
func TestAddresses(t *testing.T) {
	svc := Ec2Svc{Client: &mockEc2Client{}}
	expected := []*ec2.Address{{}}
	got, err := svc.Addresses()
	if err != nil {
		t.Fatalf("Addresses() failed: %v", err)
	}
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("Addresses() failed. Expected: %#v (%T)\nGot: %#v (%T)", expected, expected, got, got)
	}
	_, err = TypeToSheet(expected)
	if err != nil {
		t.Fatalf("TypeToSheet failed: %v", err)
	}
}

// func KeyPairs() ([]*ec2.KeyPairInfo, error)
func TestKeyPairs(t *testing.T) {
	svc := Ec2Svc{Client: &mockEc2Client{}}
	expected := []*ec2.KeyPairInfo{{}}
	got, err := svc.KeyPairs()
	if err != nil {
		t.Fatalf("KeyPairs() failed: %v", err)
	}
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("KeyPairs() failed. Expected: %#v (%T)\nGot: %#v (%T)", expected, expected, got, got)
	}
	_, err = TypeToSheet(expected)
	if err != nil {
		t.Fatalf("TypeToSheet failed: %v", err)
	}
}

// func KeyPairs(sess *session.Session, cred *credentials.Credentials) ([]*ec2.KeyPairInfo, error)
