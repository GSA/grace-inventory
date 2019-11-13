package helpers

import (
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

// Ec2Svc ... uses an SDK service iface to access SDK service client
type Ec2Svc struct {
	Client ec2iface.EC2API
}

// Instances ... pages through DescribeInstancesPages and returns all EC2 instances
func (svc *Ec2Svc) Instances() ([]*ec2.Instance, error) {
	var results []*ec2.Reservation
	err := svc.Client.DescribeInstancesPages(&ec2.DescribeInstancesInput{},
		func(page *ec2.DescribeInstancesOutput, lastPage bool) bool {
			results = append(results, page.Reservations...)
			return !lastPage
		})
	if err != nil {
		return nil, err
	}
	var instances []*ec2.Instance
	for _, r := range results {
		instances = append(instances, r.Instances...)
	}
	return instances, nil
}

// Images ... performs DescribeImages and returns all EC2 images
func (svc *Ec2Svc) Images() ([]*ec2.Image, error) {
	input := &ec2.DescribeImagesInput{Owners: self}
	result, err := svc.Client.DescribeImages(input)
	if err != nil {
		return nil, err
	}
	return result.Images, nil
}

// Volumes ... pages through DescribeVolumesPages and returns all EBS volumes
func (svc *Ec2Svc) Volumes() ([]*ec2.Volume, error) {
	var results []*ec2.Volume
	err := svc.Client.DescribeVolumesPages(&ec2.DescribeVolumesInput{},
		func(page *ec2.DescribeVolumesOutput, lastPage bool) bool {
			results = append(results, page.Volumes...)
			return !lastPage
		})
	if err != nil {
		return nil, err
	}
	return results, nil
}

// Snapshots ... pages through DescribeSnapshotsPages and returns all EBS snapshots
func (svc *Ec2Svc) Snapshots() ([]*ec2.Snapshot, error) {
	var results []*ec2.Snapshot
	input := &ec2.DescribeSnapshotsInput{OwnerIds: self}
	err := svc.Client.DescribeSnapshotsPages(input,
		func(page *ec2.DescribeSnapshotsOutput, lastPage bool) bool {
			results = append(results, page.Snapshots...)
			return !lastPage
		})
	if err != nil {
		return nil, err
	}
	return results, nil
}

// Vpcs ... pages through DescribeVpcsPages and returns all VPCs
func (svc *Ec2Svc) Vpcs() ([]*ec2.Vpc, error) {
	var results []*ec2.Vpc
	err := svc.Client.DescribeVpcsPages(&ec2.DescribeVpcsInput{},
		func(page *ec2.DescribeVpcsOutput, lastPage bool) bool {
			results = append(results, page.Vpcs...)
			return !lastPage
		})
	if err != nil {
		return nil, err
	}
	return results, nil
}

// Subnets ... pages through DescribeSubnetsPages and returns all VPC Subnets
func (svc *Ec2Svc) Subnets() ([]*ec2.Subnet, error) {
	var results []*ec2.Subnet
	err := svc.Client.DescribeSubnetsPages(&ec2.DescribeSubnetsInput{},
		func(page *ec2.DescribeSubnetsOutput, lastPage bool) bool {
			results = append(results, page.Subnets...)
			return !lastPage
		})
	if err != nil {
		return nil, err
	}
	return results, nil
}

// SecurityGroups ... pages through DescribeSecurityGroupsPages and returns all SecurityGroups
func (svc *Ec2Svc) SecurityGroups() ([]*ec2.SecurityGroup, error) {
	var results []*ec2.SecurityGroup
	err := svc.Client.DescribeSecurityGroupsPages(&ec2.DescribeSecurityGroupsInput{},
		func(page *ec2.DescribeSecurityGroupsOutput, lastPage bool) bool {
			results = append(results, page.SecurityGroups...)
			return !lastPage
		})
	if err != nil {
		return nil, err
	}
	return results, nil
}

// Addresses ... performs DescribeAddresses and returns all EC2 Addresses
func (svc *Ec2Svc) Addresses() ([]*ec2.Address, error) {
	result, err := svc.Client.DescribeAddresses(&ec2.DescribeAddressesInput{})
	if err != nil {
		return nil, err
	}
	return result.Addresses, nil
}

// KeyPairs ... performs DescribeKeyPairs and returns all EC2 KeyPairs
func (svc *Ec2Svc) KeyPairs() ([]*ec2.KeyPairInfo, error) {
	result, err := svc.Client.DescribeKeyPairs(&ec2.DescribeKeyPairsInput{})
	if err != nil {
		return nil, err
	}
	return result.KeyPairs, nil
}
