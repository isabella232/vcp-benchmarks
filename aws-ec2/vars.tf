variable "aws_access_key" {
    description = "Valid Amazon AWS access key"
}

variable "aws_secret_key" {
    description = "Valid Amazon AWS secret key"
}

# http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/using-regions-availability-zones.html
variable "aws_ec2_region" {
    description = "Amazon EC2 region to use"
    default = "eu-west-1"
}

# https://aws.amazon.com/marketplace/pp/B00O7WM7QW 
variable "aws_ec2_ami" {
    description = "The EC2 AMI to use"
    default = {
        "consumer"      = "ami-74740a07"
        "backend"       = "ami-74740a07"
        "varnish"       = "ami-74740a07"
        "loadbalancer"  = "ami-74740a07"
    }
}

variable "sshkey" {
    description = "SSH key to use for provisioning"
}

variable "instance_names" {
    description = "Base names to use for instances"
    default = {
        "consumer"      = "consumer"
        "backend"       = "backend"
        "varnish"       = "varnish"
        "loadbalancer"  = "loadbalancer"
    }
}

# https://aws.amazon.com/ec2/instance-types/ 
variable "instance_types" {
    description = "Instance types to use"
    default = {
        "consumer"      = "m3.medium"
        "backend"       = "m3.medium"
        "varnish"       = "m3.medium"
        "loadbalancer"  = "m3.medium"
    }
}

variable "user_name" {
    default = {
        "consumer"      = "ec2-user"
        "backend"       = "ec2-user"
        "varnish"       = "ec2-user"
        "loadbalancer"  = "ec2-user"
    }
}
