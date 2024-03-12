// Copyright © 2023 Cisco Systems, Inc. and its affiliates.
// All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aws

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	cloudformationtypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func (e *AWSEnv) prepareStack(ctx context.Context) error {
	// Create a new key pair
	_, err := e.ec2Client.ImportKeyPair(ctx, &ec2.ImportKeyPairInput{
		KeyName:           aws.String(AWSKeyName),
		PublicKeyMaterial: []byte(e.publicKey),
	},
	)
	if err != nil {
		return fmt.Errorf("failed to import key pair: %w", err)
	}

	// Create a new S3 bucket
	_, err = e.s3Client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: &e.stackName,
		CreateBucketConfiguration: &s3types.CreateBucketConfiguration{
			LocationConstraint: s3types.BucketLocationConstraint(e.region),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create bucket: %w", err)
	}

	// Read template file
	f, err := os.Open("../installation/aws/VmClarity.cfn")
	if err != nil {
		return fmt.Errorf("failed to read template file: %w", err)
	}

	// Upload template file to S3 bucket
	_, err = e.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &e.stackName,
		Key:    &e.stackName,
		Body:   f,
	})
	if err != nil {
		return fmt.Errorf("failed to put object: %w", err)
	}
	e.templateURL = "https://" + e.stackName + ".s3.amazonaws.com/" + e.stackName

	return nil
}

func (e *AWSEnv) cleanupStack(ctx context.Context) error {
	// Delete the key pair
	_, err := e.ec2Client.DeleteKeyPair(ctx, &ec2.DeleteKeyPairInput{
		KeyName: aws.String(AWSKeyName),
	})
	if err != nil {
		return fmt.Errorf("failed to delete key pair: %w", err)
	}

	// Delete template file from S3 bucket
	_, err = e.s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &e.stackName,
		Key:    &e.stackName,
	})
	if err != nil {
		return fmt.Errorf("failed to delete object: %w", err)
	}

	// Delete the S3 bucket
	_, err = e.s3Client.DeleteBucket(ctx, &s3.DeleteBucketInput{
		Bucket: &e.stackName,
	})
	if err != nil {
		return fmt.Errorf("failed to delete bucket: %w", err)
	}

	return nil
}

// Get stack status by stack name.
func (e *AWSEnv) getStackStatus(ctx context.Context) (cloudformationtypes.StackStatus, error) {
	stacks, err := e.client.DescribeStacks(
		ctx,
		&cloudformation.DescribeStacksInput{
			StackName: &e.stackName,
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to describe stack: %w", err)
	}

	if len(stacks.Stacks) != 1 {
		return "", errors.New("failed to find stack")
	}

	return stacks.Stacks[0].StackStatus, nil
}

// Get EC2 instance status by instance ID.
func (e *AWSEnv) getEC2InstanceStatus(ctx context.Context, instanceID string) (ec2types.InstanceStateName, error) {
	instanceStatus, err := e.ec2Client.DescribeInstanceStatus(
		ctx,
		&ec2.DescribeInstanceStatusInput{
			InstanceIds: []string{instanceID},
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to describe instance status: %w", err)
	}

	if len(instanceStatus.InstanceStatuses) != 1 {
		return "", errors.New("failed to find instance status")
	}

	return instanceStatus.InstanceStatuses[0].InstanceState.Name, nil
}
