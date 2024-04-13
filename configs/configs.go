package configs

type ZeroDtAppConfig struct {
	AppName       string `yaml:"app_name"`
	AppVersion    string `yaml:"app_version"`
	TagKey        string `yaml:"tag_key"`
	TagValue      string `yaml:"tag_value"`
	AmiID         string `yaml:"ami_id"`
	InstanceType  string `yaml:"instance_type"`
	AmiInitScript string `yaml:"ami_init_script"`
	StandByTime   int    `yaml:"time_before_deleting_old_instance_in_min"`
}
