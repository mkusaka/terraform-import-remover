resource "aws_vpc" "main" {
  cidr_block = "10.0.0.0/16"
  
  tags = {
    Name = "MainVPC"
  }
}

import {
  to = aws_vpc.main
  id = "vpc-67890"
}

resource "aws_subnet" "public" {
  vpc_id     = aws_vpc.main.id
  cidr_block = "10.0.1.0/24"
}

import {
  to = aws_subnet.public
  id = "subnet-12345"
}
