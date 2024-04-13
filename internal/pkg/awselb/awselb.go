package awselb

type AwsELB struct {
	LBARN          string `json:"lb_arn"`
	LBDNS          string `json:"lb_dns"`
	TagKey         string `json:"tag_key"`
	TagValue       string `json:"tag_value"`
	LBHostedZoneId string `json:"lb_hosted_zone_id"`
}
