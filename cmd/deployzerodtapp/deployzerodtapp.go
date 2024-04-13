package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	b64 "encoding/base64"

	"github.com/nagri/awszerodtdeploy/configs"
	"github.com/nagri/awszerodtdeploy/internal/app/deployzerodtapp"
	"gopkg.in/yaml.v3"
)

const default_config_file_path = "cmd/deployzerodtapp/configs/configs.yaml"

func main() {

	rollback := flag.Bool("rollback", false, "Rollback to previous version")
	rollback_version := flag.Int64("version", 0, "Rollback Verison, default is previous version")
	configFile := flag.String("f", "", "Config file")
	flag.Parse()

	if *configFile == "" {
		*configFile = default_config_file_path
	}
	yamlFile, err := os.ReadFile(*configFile)
	if err != nil {
		panic(err)
	}
	fmt.Println("configFile", *configFile)
	var config configs.ZeroDtAppConfig

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}

	encodedUserdata := b64.StdEncoding.EncodeToString([]byte(config.AmiInitScript))

	zdtapp := &deployzerodtapp.Zerodtapp{
		Name:          config.AppName,
		NewVersion:    config.AppVersion,
		TagKey:        config.TagKey,
		TagValue:      config.TagValue,
		AmiID:         config.AmiID,
		InstanceType:  config.InstanceType,
		AMIInitScript: encodedUserdata,
		AppVersion:    config.AppVersion,
	}

	if *rollback {

		if *rollback_version == 0 {
			zdtapp.RollbackChange(0)
		} else {
			zdtapp.RollbackChange(*rollback_version)
		}
		return
	}

	fmt.Println("Deploying", config.AppName, "version", config.AppVersion)

	// Create a new Launch template from the config file data, or update an existing one.
	template := zdtapp.CreateLaunchtemplate()

	//Create new AWS Target Group
	tgARN := zdtapp.CreateAWSTG()

	//Get AWS ALB data based on the tags mentioned in the config file
	albARN := zdtapp.GetAWSALB()
	if albARN == nil {
		// create a new ALB
		zdtapp.CreateAWSALB(tgARN)
	}

	//Get the running AWS ASG data based on the tags mentioned in the config file
	blueASGARN := zdtapp.GetBlueASG()
	if blueASGARN == nil {
		zdtapp.CreateGreenASG(template, tgARN)
		// Since the ASG didnt exist it means we are done here
		return
	} else {

		// Create a new ASG with the new template
		greenASGARN := zdtapp.CreateGreenASG(template, tgARN)

		// Now we have all the resources that we need
		// Check if the GreenASGIsHealthy

		// for {
		// 	if zdtapp.GreenASGIsHealthy(greenASGARN) {
		// 		break
		// 	} else {
		// 		time.Sleep(5 * time.Second)
		// 	}
		// }

		zdtapp.AttachGreenASGToTG(tgARN, greenASGARN)
		zdtapp.StandbyBlueASG(blueASGARN)

		// Wait here for time declared in the config
		fmt.Printf("Will wait for %d minutes before deleting the Old ASG", config.StandByTime)
		time.Sleep(time.Minute * time.Duration(config.StandByTime))
		zdtapp.DeleteBlueASG(blueASGARN)

	}

}
