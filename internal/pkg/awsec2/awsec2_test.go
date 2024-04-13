package awsec2

import (
	"context"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/nagri/awszerodtdeploy/internal/pkg/awsec2/mock_awsec2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCheckLaunchTemplate(t *testing.T) {
	ec2template := &EC2Template{
		Name:           "TestingLaunchTemplate",
		DefaultVersion: 1,
		SourceVersion:  1,
		TagKey:         "Group",
		TagValue:       "TheApp",
		AmiID:          "ami-01dad638e8f31ab9a",
		InstanceType:   "t3.micro",
		AMIInitScript:  "yum update -y",
		Exists:         false,
		AppVersion:     "1.0",
	}

	ctrl := gomock.NewController(t)

	mo := mock_awsec2.NewMockEC2ClientInterface(ctrl)

	template := ec2types.RequestLaunchTemplateData{
		InstanceType: chooseEc2Instance(ec2template.InstanceType),
		ImageId:      aws.String(ec2template.AmiID),
		TagSpecifications: []ec2types.LaunchTemplateTagSpecificationRequest{
			{
				ResourceType: ec2types.ResourceTypeInstance,
				Tags: []ec2types.Tag{
					{
						Key:   aws.String(ec2template.TagKey),
						Value: aws.String(ec2template.TagValue),
					},
				},
			},
		},
		UserData: aws.String(ec2template.AMIInitScript)}

	createParams := ec2.CreateLaunchTemplateVersionInput{
		LaunchTemplateData: &template,
		LaunchTemplateName: aws.String(ec2template.Name),
		SourceVersion:      aws.String(strconv.FormatInt(ec2template.SourceVersion, 10)),
		VersionDescription: aws.String("Running application version:" + ec2template.AppVersion),
	}

	modifyParams := ec2.ModifyLaunchTemplateInput{
		DefaultVersion:     aws.String(strconv.FormatInt(ec2template.DefaultVersion, 10)),
		LaunchTemplateName: aws.String(ec2template.Name),
	}
	mo.EXPECT().ModifyLaunchTemplate(context.Background(),
		&modifyParams).Return(&ec2.ModifyLaunchTemplateOutput{
		LaunchTemplate: &ec2types.LaunchTemplate{
			DefaultVersionNumber: aws.Int64(ec2template.DefaultVersion),
		},
	}, nil)

	mo.EXPECT().CreateLaunchTemplateVersion(context.Background(),
		&createParams).Return(&ec2.CreateLaunchTemplateVersionOutput{
		LaunchTemplateVersion: &ec2types.LaunchTemplateVersion{
			VersionNumber:      aws.Int64(ec2template.DefaultVersion),
			LaunchTemplateName: aws.String(ec2template.Name),
		},
	}, nil)

	err := UpdateLaunchtemplate(mo, ec2template)

	assert.Equal(t, nil, err)

}
