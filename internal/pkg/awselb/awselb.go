package awselb

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	rgtapi "github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
	"github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi/types"
)

type AwsELB struct {
	LBARN          string `json:"lb_arn"`
	LBDNS          string `json:"lb_dns"`
	TagKey         string `json:"tag_key"`
	TagValue       string `json:"tag_value"`
	LBHostedZoneId string `json:"lb_hosted_zone_id"`
}

// func (awselb *AwsELB) GetAWSELB() (*AwsELB, error) {

// 	defaultConfig, err := config.LoadDefaultConfig(context.TODO())
// 	if err != nil {
// 		fmt.Println(err)
// 		return nil, err
// 	}
// 	svc := elbv2.NewFromConfig(defaultConfig)
// 	input := &elbv2.DescribeLoadBalancersInput{
// 		LoadBalancerArns: []string{},
// 	}
// 	ctx := context.Background()
// 	result, err := svc.DescribeLoadBalancers(ctx, input)
// 	if err != nil {
// 		fmt.Println(err)
// 		return nil, err
// 	}

// 	for _, v := range result.LoadBalancers {
// 		fmt.Println("LoadBalancerName", *v.LoadBalancerName)
// 		fmt.Println("LoadBalancerArn", *v.LoadBalancerArn)
// 		fmt.Println("DNSName", *v.DNSName)
// 	}

//		return nil, nil
//	}
func (awselb *AwsELB) GetAWSELBByTag() (*AwsELB, error) {

	defaultConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	svc := rgtapi.NewFromConfig(defaultConfig)
	tagFilter := types.TagFilter{
		Key:    aws.String(awselb.TagKey),
		Values: []string{awselb.TagValue},
	}

	params := &rgtapi.GetResourcesInput{
		ResourceTypeFilters: []string{"elasticloadbalancing"},
		TagFilters: []types.TagFilter{
			tagFilter,
		},
	}
	resp, err := svc.GetResources(context.TODO(), params)
	// Build the request with its input parameters
	if err != nil {
		fmt.Printf("failed to list resources, %v", err)
	}
	for _, res := range resp.ResourceTagMappingList {
		fmt.Println(" Pointer address:", res.ResourceARN)
		fmt.Println(" Value behind pointer:", *res.ResourceARN)
		fmt.Println(" Value:", aws.ToString(res.ResourceARN))
		fmt.Println()
	}
	return nil, nil
}

// func (awselb *AwsELB) GetTargetGroups() (*AwsELB, error) {
// 	defaultConfig, err := config.LoadDefaultConfig(context.TODO())
// 	if err != nil {
// 		fmt.Println(err)
// 		return nil, err
// 	}
// 	svc := rgtapi.NewFromConfig(defaultConfig)

// 	a
// 	input := &elbv2.DescribeTargetGroupsInput{
// 		LoadBalancerArns: []string{},
// 	}

// }
