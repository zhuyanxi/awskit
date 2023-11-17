package main

import (
	"context"
	"flag"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/sirupsen/logrus"
)

var instanceID string

func main() {
	flag.StringVar(&instanceID, "i", "", "the instance id to re-associate")
	flag.Parse()
	if instanceID == "" {
		logrus.Errorln("instance id is empty")
		return
	}

	client, err := NewClient()
	if err != nil {
		logrus.Errorln("Error loading AWS SDK config:", err)
		return
	}

	allocateResp, err := client.allocateNewIP(context.Background())
	if err != nil {
		logrus.Errorln("Error allocating Elastic IP:", err)
		return
	}
	logrus.Infof("Allocated new Elastic IP: %s", *allocateResp.PublicIp)

	associateResp, err := client.associateNewAddress(context.Background(), allocateResp, instanceID)
	if err != nil {
		logrus.Errorln("Error associating Elastic IP:", err)
		return
	}
	logrus.Infof("Associated IP %s with instance %s successful, resp is %+v", *allocateResp.PublicIp, instanceID, associateResp)

	// Describe and release any previously associated Elastic IPs
	describeResp, err := client.ec2.DescribeAddresses(context.TODO(), &ec2.DescribeAddressesInput{})
	if err != nil {
		logrus.Errorln("Error describing Elastic IPs:", err)
		return
	}
	for _, addr := range describeResp.Addresses {
		if addr.InstanceId == nil {
			// Release the previously associated Elastic IP
			_, err := client.ec2.ReleaseAddress(context.TODO(), &ec2.ReleaseAddressInput{
				AllocationId: addr.AllocationId,
			})
			if err != nil {
				logrus.Errorln("Error releasing Elastic IP:", err)
				return
			}
			logrus.Infof("Released previous association: %s\n", *addr.PublicIp)
		} else {
			logrus.Infof("remain ip is %s", *addr.PublicIp)
		}
	}
}
