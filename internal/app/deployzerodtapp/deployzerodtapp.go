package deployzerodtapp

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"

	taggingtypes "github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi/types"
	"github.com/nagri/awszerodtdeploy/internal/pkg/awsec2"
	"github.com/nagri/awszerodtdeploy/internal/pkg/awsgetresource"
)

const (
	elbResource         = "elasticloadbalancing:loadbalancer"
	targetGroupResource = "elasticloadbalancing:targetgroup"
	blueASGResource     = "autoscaling:autoScalingGroup"
)

type Zerodtapp struct {
	Name           string
	AWSRoute53     []string
	AWSALB         []string
	BlueASG        []string
	GreenASG       []string
	TG             []string
	RunningVersion string
	NewVersion     string
	TagKey         string
	TagValue       string
	AmiID          string
	InstanceType   string
	AMIInitScript  string
	AppVersion     string
}

func (zdtapp *Zerodtapp) RollbackChange(version int64) {
	// This will Rollback the change to the previous version if verion is 0.
	// else it will roll back to whatever version is passed
	ec2Template := &awsec2.EC2Template{
		Name:          zdtapp.Name,
		TagKey:        zdtapp.TagKey,
		TagValue:      zdtapp.TagValue,
		AmiID:         zdtapp.AmiID,
		InstanceType:  zdtapp.InstanceType,
		AMIInitScript: zdtapp.AMIInitScript,
		AppVersion:    zdtapp.AppVersion,
	}

	if version == 0 {
		// Roll back to last good version
		fmt.Printf("Rolling back %s to version current-%d \n", zdtapp.Name, 1)
		ec2Template.Rollback(0)
	} else {
		// Roll back to the version passed a parameter
		fmt.Println("Rolling back", zdtapp.Name, "to version", version)
		ec2Template.Rollback(version)
	}

}

func (zdtapp *Zerodtapp) GetEC2LaunchTemplate() {
	// Get EC2 Launch Template resource associated for this app
}

func (zdtapp *Zerodtapp) CreateLaunchtemplate() *awsec2.EC2Template {
	// Create a Launch template for EC2 instances
	// To be used in ASG groups
	ec2Template := &awsec2.EC2Template{
		Name:          zdtapp.Name,
		TagKey:        zdtapp.TagKey,
		TagValue:      zdtapp.TagValue,
		AmiID:         zdtapp.AmiID,
		InstanceType:  zdtapp.InstanceType,
		AMIInitScript: zdtapp.AMIInitScript,
		AppVersion:    zdtapp.AppVersion,
	}
	ec2Template.CreateLaunchtemplate()
	return ec2Template

}
func (zdtapp *Zerodtapp) GetAWSRoute53() {
	// Get Route53 resource associated for this app
}

func (zdtapp *Zerodtapp) CreateAWSRoute53() {

}

func (zdtapp *Zerodtapp) DeleteAWSRoute53() {
	// Delete Route53 resource associated for this app
}

func (zdtapp *Zerodtapp) GetAWSALB() []string {
	// Get ALB resource associated for this app
	tagFilter := &taggingtypes.TagFilter{
		Key:    aws.String(zdtapp.TagKey),
		Values: []string{zdtapp.TagValue},
	}

	albARNs, err := awsgetresource.GetAWSResourceByTag(tagFilter, elbResource)
	if err != nil {
		fmt.Println(err)
	}
	zdtapp.AWSALB = albARNs
	return albARNs
}
func (zdtapp *Zerodtapp) CreateAWSALB(tgARN string) string {
	// Create AWS ALB and attach tgARN to it

	awsALBARN := ""
	return awsALBARN

}

func (zdtapp *Zerodtapp) DeleteAWSALB() {
	// Delete ALB resource associated for this app
}

func (zdtapp *Zerodtapp) CreateAWSTG() string {
	// Create AWS instance TG resource
	tgARN := ""
	return tgARN
}

func (zdtapp *Zerodtapp) DeleteAWSTG() {
	// Delete TG resource associated for this app
}

func (zdtapp *Zerodtapp) GetBlueASG() []string {
	// Get the Current ASG resource associated for this app.
	// We call this Blue ASG and this is the ASG that is currently serving the traffic.
	tagFilter := &taggingtypes.TagFilter{
		Key:    aws.String(zdtapp.TagKey),
		Values: []string{zdtapp.TagValue},
	}

	asgARN, err := awsgetresource.GetAWSResourceByTag(tagFilter, blueASGResource)
	if err != nil {
		fmt.Println(err)
	}

	zdtapp.BlueASG = asgARN

	return asgARN

}

func (zdtapp *Zerodtapp) CreateGreenASG(template *awsec2.EC2Template, tgARN string) string {
	// Create a new ASGGroup  with the new template.
	// This ASG will be called Green ASG.
	// This will replace the Blue ASG.
	greenASGARN := "ARNPlaceholder"
	return greenASGARN
}
func (zdtapp *Zerodtapp) GreenASGIsHealthy(greenASGARN string) bool {
	// Create a new ASGGroup  with the new template.
	// This ASG will be called Green ASG.
	// This will replace the Blue ASG.
	return false
}

func (zdtapp *Zerodtapp) StandbyBlueASG(arn []string) {
	// Once GreenASG is up and healthy and attached to the ALB,
	// we mark the BlueASG as standby so the if we have to revert back
	// to the BlueASG we can.
}

func (zdtapp *Zerodtapp) DeleteBlueASG(arn []string) {
	// Delete Blue ASG once the system does not show any errors for set duration of time.
}

func (zdtapp *Zerodtapp) AttachGreenASGToTG(tgARN string, greenASGARN string) {

}
