# AWS Credentials
aws_access_key = ""
aws_secret_key = ""
# SSH key to install and use during deployment and provisioning
sshkey = ""
# If above SSH key is not in the users keychain, specify path
#sshkey_path = "${file("path/to/private_key")}"

# Use Europe region
aws_ec2_region         = "eu-west-1"

# Use custom AMIs

# Terraform 0.6 version:
# 
aws_ec2_ami.vcp41         = "ami-876ae990"
user_name.ami-876ae990    = "ec2-user"
instance_types.benchmark  = "m4.xlarge"

# Terraform 0.7 version:
#
#aws_ec2_ami = {
#    vcp41        = "ami-876ae990"
#}
#
#user_name = {
#    ami-876ae990 = "ec2-user"
#}
#
#instance_types = {
#    benchmark  = "m4.xlarge"
#}
