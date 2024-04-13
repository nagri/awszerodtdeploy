# awszerodtdeploy
AWS Zero DownTime Deployment

---


Project Layout copied from https://github.com/golang-standards/project-layout 

# Design
You can Deploy applications on EC2 instances.

## Assumptions
1. Default security groups work just fine.
2. We will configure the instance with an initialization script that will install tools and packages and may configure the application from a git repo.

User will create a YAML file like `/path/of/your/yaml/file.yaml` or 

```
---
app_name: AppOne
app_version: 1.0
tag_key: Group
tag_value: AppOne
ami_id: "ami-01dad638e8f31ab9a"
instance_type: "t3.micro"
time_before_deleting_old_instance_in_min: 15 
ami_init_script: |
  #!/bin/bash
  yum update -y
  amazon-linux-extras install -y lamp-mariadb10.2-php7.2 php7.2
  yum install -y httpd mariadb-server
  systemctl start httpd
  systemctl enable httpd
  usermod -a -G apache ec2-user
  chown -R ec2-user:apache /var/www
  chmod 2775 /var/www
  find /var/www -type d -exec chmod 2775 {} \;
  find /var/www -type f -exec chmod 0664 {} \;
  echo "<?php phpinfo(); ?>" > /var/www/html/phpinfo.php

```

This will be the source of truth for the applications to be deployed. All the resources that will be created will be tagged.
These Tags will be used as identifiers for the resources for future upgrades and deployments.

- A Launch template will be created or updated based on the config from the YAML file.
- A new Instance Target Group will be created.
- An Application LB will be created if it does not exist. Target Group that was created will be attached to this ALB
- An Auto Scaling Group(ASG) will be created if one does not exist already using the Launch Template. This ASG will attach to the ALB and hence to the Target Group.

When the next Upgrade has to happen, User will create a new YAML file and run the application again.

- A new Version of Launch template will be created and marked as Default.
- A new Instance Target Group will be created.
- Target Group will be attached to this ALB
- A new Auto Scaling Group(ASG) (GreenASG) will be created.
- Once GreenASG is attached, we mark BlueASG in standby mode.
- If no errors or issues are detected, after some time(time_before_deleting_old_instance_in_min from YAML) BlueASG will be deleted.

## How to run the app?
You'll have to set AWS credentials as environment variables

```
export AWS_ACCESS_KEY_ID=ABCDS2NQIVNG6FMLYES   
export AWS_SECRET_ACCESS_KEY=q0iX9ZQub6chPID4PlT786kIQwir145f2XYZkLvy
export AWS_DEFAULT_REGION=eu-north-1
```

After you have cloned this repo, f  rom inside the root directory run

```
go run cmd/deployzerodtapp/deployzerodtapp.go
```

This will star running the application which will create/update
- AWS EC2 Launch Template
- AWS Target Group
- AWS Application Load Balancer (ALB) 
- AWS Auto Scaling Group

and stitch them together to serve an application that was created using the EC2 Launch template.

## You can rollback to older version of Launch Template with flags

Say you want to roll back to the last version of the Launch template

``` 
go run cmd/deployzerodtapp/deployzerodtapp.go --rollback`
```

Say you want to rollback to version 7

``` 
go run cmd/deployzerodtapp/deployzerodtapp.go --rollback --version=7`
```

## You can pass a config file as well
The config file contains configurations for the launch template including a script that initializes the EC2 instance to serve the application.

``` 
go run cmd/deployzerodtapp/deployzerodtapp.go --configfile=/path/of/your/yaml/file.yaml`
```

The yaml is tightly coupled with the code you can see the example at 
`cmd/deployzerodtapp/configs/configs.yaml` or this

```
---
app_name: AppOne
app_version: 1.0
tag_key: Group
tag_value: AppOne
ami_id: "ami-01dad638e8f31ab9a"
instance_type: "t3.micro"
time_before_deleting_old_instance_in_min: 15 
ami_init_script: |
  #!/bin/bash
  yum update -y
  amazon-linux-extras install -y lamp-mariadb10.2-php7.2 php7.2
  yum install -y httpd mariadb-server
  systemctl start httpd
  systemctl enable httpd
  usermod -a -G apache ec2-user
  chown -R ec2-user:apache /var/www
  chmod 2775 /var/www
  find /var/www -type d -exec chmod 2775 {} \;
  find /var/www -type f -exec chmod 0664 {} \;
  echo "<?php phpinfo(); ?>" > /var/www/html/phpinfo.php

```

## Run Test cases

``` 
go test ./...`
```

