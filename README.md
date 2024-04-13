# awszerodtdeploy
AWS Zero Down Time Deployment

---


Project Layout copied from https://github.com/golang-standards/project-layout 

## How to run the app?

From inside the root directory run

`$ go run cmd/deployzerodtapp/deployzerodtapp.go`

This will star running the application which will create/update
- AWS EC2 Launch Template
- AWS Target Group
- AWS Application Load Balancer (ALB) 
- AWS Auto Scaling Group

and stich them together to serve an application that was created using the EC2 Launch template.

## You can rollback to older version of Launch Template with flags

Say you want to roll back to the last version of the Launch template

`$ go run cmd/deployzerodtapp/deployzerodtapp.go --rollback`

Say you want to rollback to version 7

`$ go run cmd/deployzerodtapp/deployzerodtapp.go --rollback --version=7`

## You can pass a config file as well
The config file contains configurations for the launch template including a script that initializes the EC2 instance to serve the application.

`$ go run cmd/deployzerodtapp/deployzerodtapp.go --configfile=/path/of/your/yaml/file.yaml`

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

`$ go test ./...`


