# Terraform based deploys

## Installing Terraform

Terraform is a Golang based tool and official binaries are available at https://terraform.io/downloads.html.
On OSX, Terraform can also be installed via Homebrew.

## Important information

Terraform manipulates cloud based instances, and it is imperative that you ensure that the configuration used does not alter any of your existing resources. Typically the naming chosen for the instances should avoid this, but please ensure via ``terraform plan`` before applying or destroying.

## Ansible integration

A nice dynamic inventory script for Ansible is at https://github.com/adammck/terraform-inventory
Binary releases are available at https://github.com/adammck/terraform-inventory/releases
Copy the ``terraform-inventory`` binary that works for your platform to the ``bin/`` subdirectory.

## Amazon EC2

### Quickstart EC2

1. Edit the terraform.tfvars file and insert AWS credentials and sshkey, or add them to a separate terraform.tfvars.mine file.
2. Run ``./deploy`` to run start instances, provision and run tests
3. Run ``./destroy`` to remove all AWS resources created by ``deploy``

### Manual EC2 run

1. Edit the terraform.tfvars file and insert AWS credentials and sshkey
2. Run ``terraform plan aws-ec2 [-var-file=terraform.tfvars.mine]`` to verify the plan
3. Run ``terraform apply aws-ec2 [-var-file=terraform.tfvars.mine]`` to spin up instances
4. Run ``ansible-playbook --inventory-file=bin/terraform-inventory playbook.yml``
5. Run ``terraform destroy aws-ec2 [-var-file=terraform.tfvars.mine]`` to remove all resources

### Custom setups

The ec2 region, instance types and ami to use is configured with variables, default is **m3.medium** and the official Centos7 AMI:

    aws_ec2_region = "us-east-1"
    aws_ec2_ami = "ami-6d1c2007" (Centos7 64bit)
    instance_types = "m3.medium" (Can be configured per node)

