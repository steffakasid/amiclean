package internal

import (
	"context"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/smithy-go"
	logger "github.com/sirupsen/logrus"
)

type Ec2client interface {
	DescribeInstances(ctx context.Context, params *ec2.DescribeInstancesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error)
	DescribeImages(ctx context.Context, params *ec2.DescribeImagesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeImagesOutput, error)
	DeregisterImage(ctx context.Context, params *ec2.DeregisterImageInput, optFns ...func(*ec2.Options)) (*ec2.DeregisterImageOutput, error)
	DescribeLaunchTemplateVersions(ctx context.Context, params *ec2.DescribeLaunchTemplateVersionsInput, optFns ...func(*ec2.Options)) (*ec2.DescribeLaunchTemplateVersionsOutput, error)
}

type AmiClean struct {
	ec2client      Ec2client
	olderthen      time.Duration
	awsaccount     string
	dryrun         bool
	useLaunchTpls  bool
	usedAMIs       []string
	ignorePatterns []string
}

type AWS struct {
	ec2        Ec2client
}

func NewAWSClient(conf func(ctx context.Context, optFns ...func(*config.LoadOptions) error) (cfg aws.Config, err error),
	ec2InitFunc func(cfg aws.Config, optFns ...func(*ec2.Options)) *ec2.Client,
	clouttrailInitFunc func(cfg aws.Config, optFns ...func(*cloudtrail.Options)) *cloudtrail.Client) *AWS {
	aws := &AWS{}

	cfg, err := conf(context.TODO())
	CheckError(err, logger.Fatalf)

	aws.ec2 = ec2InitFunc(cfg)
	aws.cloudtrail = clouttrailInitFunc(cfg)
	return aws
}

func (a *AWS) getUsedAMIsFromEC2() []string {
	usedImages := []string{}
	nextToken := ""
	for {
		opts := &ec2.DescribeInstancesInput{}
		if nextToken != "" {
			opts.NextToken = &nextToken
		}
		ec2Instances, err := a.ec2.DescribeInstances(context.TODO(), opts)
		CheckError(err, logger.Errorf)
		if ec2Instances != nil {
			for _, reserveration := range ec2Instances.Reservations {
				for _, instance := range reserveration.Instances {
					usedImages = uniqueAppend(usedImages, *instance.ImageId)
				}
			}
		}

		if ec2Instances == nil || ec2Instances.NextToken == nil {
			break
		}
		nextToken = *ec2Instances.NextToken
	}
	logger.Debug("UsedImages[] from EC2", usedImages)
	return usedImages
}

func (a *AWS) getUsedAMIsFromLaunchTpls() []string {
	usedImages := []string{}
	nextToken := ""
	for {
		opts := &ec2.DescribeLaunchTemplateVersionsInput{
			Versions: []string{"$Latest"},
		}
		if nextToken != "" {
			opts.NextToken = &nextToken
		}
		launchTpls, err := a.ec2.DescribeLaunchTemplateVersions(context.TODO(), opts)
		CheckError(err, logger.Errorf)
		if launchTpls != nil {
			for _, launchTplVersion := range launchTpls.LaunchTemplateVersions {
				if launchTplVersion.LaunchTemplateData.ImageId != nil {
					usedImages = append(usedImages, *launchTplVersion.LaunchTemplateData.ImageId)
				}
			}
		}
		if launchTpls == nil || launchTpls.NextToken == nil {
			break
		}
		nextToken = *launchTpls.NextToken
	}
	logger.Debug("UsedImages[] from Launch Templates", usedImages)
	return usedImages
}
}

func (a AmiClean) DeleteOlderUnusedAMIs() error {
	describeImageInput := &ec2.DescribeImagesInput{Owners: []string{"self"}}
	if a.awsaccount != "" {
		describeImageInput.Owners = append(describeImageInput.Owners, a.awsaccount)
	}
	images, err := a.ec2client.DescribeImages(context.TODO(), describeImageInput)
	if err != nil {
		return err
	}
	today := time.Now()
	olderThenDate := today.Add(a.olderthen * -1)
	for _, image := range images.Images {
		if !contains(a.usedAMIs, *image.ImageId) {
			ok, err := matchAny(*image.Name, a.ignorePatterns)
			if err != nil {
				return err
			}
			if !ok {
				creationDate, err := time.Parse("2006-01-02T15:04:05.000Z", *image.CreationDate)
				if err != nil {
					logger.Error(err)
				}
				if creationDate.Before(olderThenDate) {
					logger.Infof("Delete %s:%s as it's creationdate %s is older then %s", *image.ImageId, *image.Name, *image.CreationDate, olderThenDate.String())
					deregisterInput := &ec2.DeregisterImageInput{
						ImageId: image.ImageId,
						DryRun:  aws.Bool(a.dryrun),
					}
					_, err := a.ec2client.DeregisterImage(context.TODO(), deregisterInput)
					CheckError(err, logger.Errorf)
				} else {
					logger.Infof("Keeping %s:%s as it's creationdate %s is newer then %s", *image.ImageId, *image.Name, *image.CreationDate, olderThenDate.String())
				}
			} else {
				logger.Infof("Ignored %s\n", *image.ImageId)
			}
		} else {
			logger.Infof("Skipping %s\n", *image.ImageId)
		}
	}
	return nil
}

func CheckError(err error, logFunc func(tpl string, args ...interface{})) {
	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) {
			logFunc("code: %s, message: %s, fault: %s", ae.ErrorCode(), ae.ErrorMessage(), ae.ErrorFault().String())
		} else {
			logFunc(err.Error())
		}
	}
}
