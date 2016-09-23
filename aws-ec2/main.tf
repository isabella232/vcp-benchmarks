# Configure the Amazon EC2 Provider
provider "aws" {
    access_key = "${var.aws_access_key}"
    secret_key = "${var.aws_secret_key}"
    region = "${var.aws_ec2_region}"
}

# Create a key pair
resource "aws_key_pair" "vcpbench-sshkey" {
    key_name = "vcpbench-sshkey" 
    public_key = "${var.sshkey}"
}

# Create a VPC
resource "aws_vpc" "vcpbench" {
    cidr_block           = "10.0.0.0/16"
    enable_dns_hostnames = true
}

# Create an internet gateway to give the subnet access to the world
resource "aws_internet_gateway" "vcpbench" {
    vpc_id = "${aws_vpc.vcpbench.id}"
}

# Grant VPC internet access on its main route table
resource "aws_route" "vcpbench-internet_access" {
    route_table_id         = "${aws_vpc.vcpbench.main_route_table_id}"
    destination_cidr_block = "0.0.0.0/0"
    gateway_id             = "${aws_internet_gateway.vcpbench.id}"
}

# Create a subnet to launch our instances into
resource "aws_subnet" "vcpbench" {
    vpc_id                  = "${aws_vpc.vcpbench.id}"
    cidr_block              = "10.0.1.0/24"
    map_public_ip_on_launch = true
}

# Our default security group to access
# the instances over SSH and HTTP
resource "aws_security_group" "vcpbench" {
  name        = "vcp_benchmarks"
  description = "Used for VCP Benchmarks"
  vpc_id      = "${aws_vpc.vcpbench.id}"

  # Permit ICMP
  ingress {
    from_port   = -1
    to_port     = -1
    protocol    = "icmp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # SSH access from anywhere, limit if need be
  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # Open lots of TCP ports from anywhere
  ingress {
    from_port   = 0
    to_port     = 8888
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # outbound internet access
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# Create the benchmark consumer instance
resource "aws_instance" "vcpbench-consumer" {
    ami                     = "${lookup(var.aws_ec2_ami, "consumer")}"
    instance_type           = "${lookup(var.instance_types, "consumer")}"
    key_name                = "vcpbench-sshkey"
    user_data               = "${file("aws-ec2/userdata.sh")}"
    subnet_id               = "${aws_subnet.vcpbench.id}"
    vpc_security_group_ids  = ["${aws_security_group.vcpbench.id}"]
    tags {
        name    = "${lookup(var.instance_names, "consumer")}"
    }
    provisioner "local-exec" {
        command = "echo ${self.private_ip} ${lookup(var.instance_names, "consumer")} > hosts/${lookup(var.instance_names, "consumer")}.host"
    }
    provisioner "local-exec" {
        command = "echo ansible_ssh_user: ${lookup(var.user_name, "consumer")} > host_vars/${self.public_ip}"
    }
}

# Create the benchmark consumer instances
resource "aws_instance" "vcpbench-varnish" {
    ami                     = "${lookup(var.aws_ec2_ami, "varnish")}"
    count                   = 4
    instance_type           = "${lookup(var.instance_types, "varnish")}"
    key_name                = "vcpbench-sshkey"
    user_data               = "${file("aws-ec2/userdata.sh")}"
    subnet_id               = "${aws_subnet.vcpbench.id}"
    vpc_security_group_ids  = ["${aws_security_group.vcpbench.id}"]
    tags {
        name = "${lookup(var.instance_names, "varnish")}"
    }
    provisioner "local-exec" {
        command = "echo ${self.private_ip} ${lookup(var.instance_names, "varnish")}${count.index} > hosts/${lookup(var.instance_names, "varnish")}${count.index}.host"
    }
    provisioner "local-exec" {
        command = "echo ansible_ssh_user: ${lookup(var.user_name, "varnish")} > host_vars/${self.public_ip}"
    }
    provisioner "local-exec" {
        command = "echo ${lookup(var.instance_names, "varnish")}${count.index} = ${self.private_ip}:6081 > hosts/${lookup(var.instance_names, "varnish")}${count.index}.vhahost"
    }
    provisioner "remote-exec" {
        connection {
            host = "${self.public_ip}"
            user = "${lookup(var.user_name, "varnish")}"
        }
        inline = [
            "echo ${lookup(var.instance_names, "varnish")}${count.index} > /tmp/hostname",
        ]
    }
}

# Create a new load balancer
resource "aws_elb" "vcpbench-elb" {
    name            = "vcpbench-elb"
    subnets         = ["${aws_subnet.vcpbench.id}"]
    security_groups = ["${aws_security_group.vcpbench.id}"]
    
    access_logs {
        bucket  = "foo"
        enabled = false
    }

# Dummy-API Listener
    listener {
        instance_port       = 8080
        instance_protocol   = "http"
        lb_port             = 8080
        lb_protocol         = "http"
    }

# VCP Listener
    listener {
        instance_port       = 6081
        instance_protocol   = "http"
        lb_port             = 80
        lb_protocol         = "http"
    }

# TLS Listener 
#
#  listener {
#    instance_port = 8000
#    instance_protocol = "http"
#    lb_port = 443
#    lb_protocol = "https"
#    ssl_certificate_id = "arn:aws:iam::123456789012:server-certificate/certName"
#  }

    health_check {
        healthy_threshold = 2
        unhealthy_threshold = 3
        timeout = 2
        target = "TCP:6081"
        interval = 5
    }

    instances = ["${aws_instance.vcpbench-varnish.*.id}"]
    cross_zone_load_balancing = false
    idle_timeout = 400
    connection_draining = true
    connection_draining_timeout = 400

    tags {
        Name = "VCP Benchmark ELB"
    }

    provisioner "local-exec" {
        command = "echo ${self.dns_name} > hosts/loadbalancer.cname"
    }
}

resource "null_resource" "hostsfile" {
    triggers {
        aws_instance_consumer   = "${aws_instance.vcpbench-consumer.private_ip}"
        aws_instance_varnish-0  = "${aws_instance.vcpbench-varnish.0.private_ip}"
        aws_instance_varnish-1  = "${aws_instance.vcpbench-varnish.1.private_ip}"
        aws_instance_varnish-2  = "${aws_instance.vcpbench-varnish.2.private_ip}"
        aws_instance_varnish-3  = "${aws_instance.vcpbench-varnish.3.private_ip}"
        aws_elasticloadbalancer = "${aws_elb.vcpbench-elb.dns_name}"
    }
    provisioner "file" {
        connection {
            host = "${aws_instance.vcpbench-consumer.public_ip}"
            user = "${lookup(var.user_name, "consumer")}"
        }
        source = "hosts"
        destination = "/tmp"
    }
    provisioner "remote-exec" {
        connection {
            host = "${aws_instance.vcpbench-consumer.public_ip}"
            user = "${lookup(var.user_name, "consumer")}"
        }
        inline = [
            "sudo sh -c 'cat /tmp/hosts/hostsheader /tmp/hosts/*.host > /etc/hosts'",
        ]
    }
    provisioner "file" {
        connection {
            host = "${aws_instance.vcpbench-varnish.0.public_ip}"
            user = "${lookup(var.user_name, "varnish")}"
        }
        source = "hosts"
        destination = "/tmp"
    }
    provisioner "remote-exec" {
        connection {
            host = "${aws_instance.vcpbench-varnish.0.public_ip}"
            user = "${lookup(var.user_name, "varnish")}"
        }
        inline = [
            "sudo sh -c 'cat /tmp/hosts/hostsheader /tmp/hosts/*.host > /etc/hosts'",
            "sudo sh -c 'cat /tmp/hosts/*.vhahost > /etc/vha-agent/nodes.conf'",
            "sudo sh -c 'hostname --file /tmp/hostname'",
        ]
    }
    provisioner "file" {
        connection {
            host = "${aws_instance.vcpbench-varnish.1.public_ip}"
            user = "${lookup(var.user_name, "varnish")}"
        }
        source = "hosts"
        destination = "/tmp"
    }
    provisioner "remote-exec" {
        connection {
            host = "${aws_instance.vcpbench-varnish.1.public_ip}"
            user = "${lookup(var.user_name, "varnish")}"
        }
        inline = [
            "sudo sh -c 'cat /tmp/hosts/hostsheader /tmp/hosts/*.host > /etc/hosts'",
            "sudo sh -c 'cat /tmp/hosts/*.vhahost > /etc/vha-agent/nodes.conf'",
            "sudo sh -c 'hostname --file /tmp/hostname'",
        ]
    }
    provisioner "file" {
        connection {
            host = "${aws_instance.vcpbench-varnish.2.public_ip}"
            user = "${lookup(var.user_name, "varnish")}"
        }
        source = "hosts"
        destination = "/tmp"
    }
    provisioner "remote-exec" {
        connection {
            host = "${aws_instance.vcpbench-varnish.2.public_ip}"
            user = "${lookup(var.user_name, "varnish")}"
        }
        inline = [
            "sudo sh -c 'cat /tmp/hosts/hostsheader /tmp/hosts/*.host > /etc/hosts'",
            "sudo sh -c 'cat /tmp/hosts/*.vhahost > /etc/vha-agent/nodes.conf'",
            "sudo sh -c 'hostname --file /tmp/hostname'",
        ]
    }
    provisioner "file" {
        connection {
            host = "${aws_instance.vcpbench-varnish.3.public_ip}"
            user = "${lookup(var.user_name, "varnish")}"
        }
        source = "hosts"
        destination = "/tmp"
    }
    provisioner "remote-exec" {
        connection {
            host = "${aws_instance.vcpbench-varnish.3.public_ip}"
            user = "${lookup(var.user_name, "varnish")}"
        }
        inline = [
            "sudo sh -c 'cat /tmp/hosts/hostsheader /tmp/hosts/*.host > /etc/hosts'",
            "sudo sh -c 'cat /tmp/hosts/*.vhahost > /etc/vha-agent/nodes.conf'",
            "sudo sh -c 'hostname --file /tmp/hostname'",
        ]
    }
}

