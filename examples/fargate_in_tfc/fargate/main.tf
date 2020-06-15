provider aws {
  region = "us-west-1"
}

resource "random_id" "project_tag" {
  byte_length = 4
}

data "aws_availability_zones" "available" {
  state = "available"
}

module "vpc" {
  source = "terraform-aws-modules/vpc/aws"

  name = "${random_id.project_tag.hex}-vpc"
  cidr = "172.16.0.0/16"
  azs  = data.aws_availability_zones.available.names
  public_subnets = [
    for num in data.aws_availability_zones.available.names :
    cidrsubnet("172.16.0.0/16", 8, 100 + index(data.aws_availability_zones.available.names, num))
  ]
}

resource "aws_default_security_group" "default" {
  vpc_id = module.vpc.vpc_id

  ingress {
    protocol  = -1
    from_port = 0
    to_port   = 0
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_ecs_cluster" "ecs" {
  name = "${random_id.project_tag.hex}-ecs"
}

resource "aws_ecs_task_definition" "nginx" {
  family                   = "nginx"
  container_definitions    = file("files/nginx.json")
  requires_compatibilities = ["FARGATE"]
  cpu                      = 256
  memory                   = 512
  network_mode             = "awsvpc"
  execution_role_arn       = aws_iam_role.execution.arn
}

resource "aws_iam_role" "execution" {
  name = "execution_role"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "ecs-tasks.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "execution-attach" {
  role       = aws_iam_role.execution.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

resource "aws_ecs_service" "nginx" {
  name            = "nginx"
  cluster         = aws_ecs_cluster.ecs.id
  task_definition = aws_ecs_task_definition.nginx.arn
  desired_count   = 1
  launch_type     = "FARGATE"
  network_configuration {
    subnets          = module.vpc.public_subnets
    assign_public_ip = true
  }
}

data "aws_network_interface" "interface" {
  filter {
    name   = "subnet-id"
    values = module.vpc.public_subnets
  }
}
