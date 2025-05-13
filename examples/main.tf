provider "aws" {
  region = "us-west-2"
}

resource "aws_instance" "web_server" {
  ami           = "ami-0c55b159cbfafe1f0"
  instance_type = "t2.micro"
  
  tags = {
    Name = "WebServer"
  }
}

import {
  to = aws_instance.web_server
  id = "i-abcd1234"
}

resource "aws_s3_bucket" "data" {
  bucket = "my-data-bucket"
}

import {
  to = aws_s3_bucket.data
  id = "my-data-bucket"
}
