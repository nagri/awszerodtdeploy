package awsgetresource

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	rgtapi "github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
	"github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi/types"
)

func GetAWSResourceByTag(tagFilter *types.TagFilter, resourceTypeFilter string) ([]string, error) {

	defaultConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	svc := rgtapi.NewFromConfig(defaultConfig)

	params := &rgtapi.GetResourcesInput{
		ResourceTypeFilters: []string{resourceTypeFilter},
		TagFilters: []types.TagFilter{
			*tagFilter,
		},
	}
	resp, err := svc.GetResources(context.TODO(), params)
	// Build the request with its input parameters
	if err != nil {
		fmt.Printf("failed to list resources, %v", err)
	}
	resourceARN := []string{}
	for _, res := range resp.ResourceTagMappingList {
		resourceARN = append(resourceARN, *res.ResourceARN)
	}
	return resourceARN, nil
}
