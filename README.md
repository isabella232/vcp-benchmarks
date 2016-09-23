# Terraform based deploys

## Installing Terraform

Terraform is a Golang based tool and official binaries are available at https://terraform.io/downloads.html.
On OSX, Terraform can also be installed via Homebrew.

## Important information

Terraform manipulates cloud based instances, and it is imperative that you ensure that the configuration used does not alter any of your existing resources. Typically the naming chosen for the instances should avoid this, but please ensure via ``terraform plan`` before applying or destroying.

## Ansible integration

A nice dynamic inventory script for Ansible is at https://github.com/adammck/terraform-inventory
Binary releases are available at https://github.com/adammck/terraform-inventory/releases, or using homebrew ``brew install terraform-inventory``.

## Quickstart EC2

1. Copy ``terraform.tfvars`` to ``terraform.tfvars.mine`` and insert AWS credentials and sshkey.
2. Run ``./deploy`` to run start instances and provision.
3. Run ``./bench`` to run benchmarks.
3. Run ``./destroy`` to remove all AWS resources created by ``deploy``.
