# Deploy a vulnerable instance on AWS

On AWS console user data needs to be defined when a new EC2 instance is created.
The user data contains the steps that will pull the vulnerable docker image.
When the docker image starts the script inside of it will copy vulnerable files to the EC2 instance.

AWS cli example:
```bash
aws ec2 run-instances --image-id ami-abcd1234 --count 1 --instance-type m3.medium \
--key-name my-key-pair --subnet-id subnet-abcd1234 --security-group-ids sg-abcd1234 \
--user-data file://installation/vulnerable_test_image/cloud-config
```