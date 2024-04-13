package awsec2

import (
	"context"

	"fmt"
	"log"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type EC2Template struct {
	Name           string
	DefaultVersion int64
	SourceVersion  int64
	TagKey         string
	TagValue       string
	AmiID          string
	InstanceType   string
	AMIInitScript  string
	AppVersion     string
	Exists         bool
}

type EC2ClientInterface interface {
	DescribeLaunchTemplates(ctx context.Context, params *ec2.DescribeLaunchTemplatesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeLaunchTemplatesOutput, error)
	CreateLaunchTemplateVersion(ctx context.Context, params *ec2.CreateLaunchTemplateVersionInput, optFns ...func(*ec2.Options)) (*ec2.CreateLaunchTemplateVersionOutput, error)
	ModifyLaunchTemplate(ctx context.Context, params *ec2.ModifyLaunchTemplateInput, optFns ...func(*ec2.Options)) (*ec2.ModifyLaunchTemplateOutput, error)
}

func chooseEc2Instance(instanceType string) ec2types.InstanceType {
	instance := ec2types.InstanceTypeT2Micro
	switch instanceType {
	case "t2.micro":
		instance = ec2types.InstanceTypeT2Micro
	case "t3.micro":
		instance = ec2types.InstanceTypeT3Micro
	case "t2.small":
		instance = ec2types.InstanceTypeT2Small
	}
	return instance
}

func CheckLaunchTemplate(client EC2ClientInterface, ec2template *EC2Template) {
	template, err := client.DescribeLaunchTemplates(context.TODO(), &ec2.DescribeLaunchTemplatesInput{
		LaunchTemplateNames: []string{ec2template.Name},
	})
	if err != nil {
		ec2template.Exists = false
		return
	}
	for _, lt := range template.LaunchTemplates {
		ec2template.Exists = true
		ec2template.SourceVersion = *lt.DefaultVersionNumber
		ec2template.DefaultVersion = *lt.DefaultVersionNumber
		return
	}
}

func (ec2template *EC2Template) checkIfTemplateExists() {
	defaultConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Println(err)
	}

	svc := ec2.NewFromConfig(defaultConfig)
	CheckLaunchTemplate(svc, ec2template)
}

// Create a Launch template if it does not exist
func (ec2template *EC2Template) CreateLaunchtemplate() {

	ec2template.checkIfTemplateExists()
	if ec2template.SourceVersion != 0 {
		fmt.Println("Upgrading the template version from SourceVersion:", ec2template.SourceVersion)
		ec2template.UpdateLaunchtemplate()
		return
	}

	defaultConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Println(err)
	}

	svc := ec2.NewFromConfig(defaultConfig)

	newLt, err := svc.CreateLaunchTemplate(context.TODO(), &ec2.CreateLaunchTemplateInput{
		LaunchTemplateData: &ec2types.RequestLaunchTemplateData{
			ImageId:      aws.String(ec2template.AmiID),
			InstanceType: chooseEc2Instance(ec2template.InstanceType),
			NetworkInterfaces: []ec2types.LaunchTemplateInstanceNetworkInterfaceSpecificationRequest{
				{
					AssociatePublicIpAddress: aws.Bool(true),
					DeviceIndex:              aws.Int32(0),
					Ipv6AddressCount:         aws.Int32(0),
				},
			},
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
			UserData: aws.String(ec2template.AMIInitScript),
		},
		LaunchTemplateName: aws.String(ec2template.Name),
		VersionDescription: aws.String("Running application version:" + ec2template.AppVersion),
	})

	if err != nil {
		fmt.Println("error while creating template", err)
	}
	fmt.Println("Created a New LaunchTemplate: ", *newLt.LaunchTemplate.LaunchTemplateName, *newLt.LaunchTemplate.LaunchTemplateId)
}

func UpdateLaunchtemplate(client EC2ClientInterface, ec2template *EC2Template) error {
	if ec2template.SourceVersion == 0 {
		return fmt.Errorf("source version cannot be 0")
	}

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
	outputCreate, err := client.CreateLaunchTemplateVersion(context.Background(), &createParams)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if outputCreate.Warning != nil {
		fmt.Printf("%v\n", outputCreate.Warning.Errors)
		for e, g := range outputCreate.Warning.Errors {
			fmt.Println(e, *g.Code, *g.Message)
		}
	}
	// set the new launch type version as the default version
	modifyParams := ec2.ModifyLaunchTemplateInput{
		DefaultVersion:     aws.String(strconv.FormatInt(*outputCreate.LaunchTemplateVersion.VersionNumber, 10)),
		LaunchTemplateName: outputCreate.LaunchTemplateVersion.LaunchTemplateName,
	}
	outputModify, err := client.ModifyLaunchTemplate(context.Background(), &modifyParams)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("default version %d\n", *outputModify.LaunchTemplate.DefaultVersionNumber)
	return nil
}

func (ec2template *EC2Template) UpdateLaunchtemplate() {

	defaultConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Println(err)
	}

	ec2client := ec2.NewFromConfig(defaultConfig)

	err = UpdateLaunchtemplate(ec2client, ec2template)
	if err != nil {
		log.Fatal(err)
	}

}

func (ec2template *EC2Template) Rollback(version int64) {

	ec2template.checkIfTemplateExists()
	if !ec2template.Exists {
		fmt.Println("Launch template", ec2template.Name, "does not exist")
		return
	}
	defaultConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Println(err)
	}

	ec2client := ec2.NewFromConfig(defaultConfig)

	var newDefaultVersion int64
	if version == 0 {
		newDefaultVersion = ec2template.DefaultVersion - 1
	} else {
		newDefaultVersion = version
	}
	modifyParams := ec2.ModifyLaunchTemplateInput{
		DefaultVersion:     aws.String(strconv.FormatInt(newDefaultVersion, 10)),
		LaunchTemplateName: aws.String(ec2template.Name),
	}
	outputModify, err := ec2client.ModifyLaunchTemplate(context.Background(), &modifyParams)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("default version %d\n", *outputModify.LaunchTemplate.DefaultVersionNumber)

}
