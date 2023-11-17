package main

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/sirupsen/logrus"
)

type Client struct {
	ec2 *ec2.Client
}

func NewClient() (*Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		logrus.Errorln("Error loading AWS SDK config:", err)
		return nil, err
	}
	return &Client{
		ec2: ec2.NewFromConfig(cfg),
	}, nil
}

func (c *Client) allocateNewIP(ctx context.Context) (*ec2.AllocateAddressOutput, error) {
	// Allocate a new Elastic IP
	allocateResp, err := c.ec2.AllocateAddress(ctx, &ec2.AllocateAddressInput{})
	if err != nil {
		logrus.Errorln("Error allocating Elastic IP:", err)
		return nil, err
	}

	return allocateResp, nil
}

func (c *Client) associateNewAddress(ctx context.Context, aaoutput *ec2.AllocateAddressOutput, insId string) (*ec2.AssociateAddressOutput, error) {
	associateResp, err := c.ec2.AssociateAddress(ctx, &ec2.AssociateAddressInput{
		AllocationId: aaoutput.AllocationId,
		InstanceId:   aws.String(insId),
	})
	if err != nil {
		logrus.Errorln("Error associating Elastic IP:", err)
		return nil, err
	}
	return associateResp, nil
}
