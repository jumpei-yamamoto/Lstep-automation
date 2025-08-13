locals {
  name = var.project_name
}

# 01 VPC（シンプルにデフォルトVPCを使う場合はスキップしても可）
data "aws_vpc" "default" {
  default = true
}
data "aws_subnets" "default_private" {
  filter {
    name   = "vpc-id"
    values = [data.aws_vpc.default.id]
  }
  # 実運用では private/public を明示分離推奨。ここでは簡易に既存subnetを利用。
}

# 02 ECR（backend / frontend-ssr 用）
resource "aws_ecr_repository" "backend" {
  name                 = "${local.name}-backend"
  image_scanning_configuration { scan_on_push = true }
  force_delete         = true
}
resource "aws_ecr_repository" "frontend_ssr" {
  name                 = "${local.name}-frontend-ssr"
  image_scanning_configuration { scan_on_push = true }
  force_delete         = true
}

# 03 ECS Cluster
resource "aws_ecs_cluster" "this" {
  name = "${local.name}-cluster"
}

# 04 Security Groups
resource "aws_security_group" "alb" {
  name        = "${local.name}-alb-sg"
  description = "ALB"
  vpc_id      = data.aws_vpc.default.id
  ingress { from_port = 80  to_port = 80  protocol = "tcp" cidr_blocks = ["0.0.0.0/0"] }
  egress  { from_port = 0   to_port = 0   protocol = "-1"  cidr_blocks = ["0.0.0.0/0"] }
}

resource "aws_security_group" "ecs_service" {
  name        = "${local.name}-ecs-sg"
  description = "ECS tasks"
  vpc_id      = data.aws_vpc.default.id
  ingress {
    description       = "ALB -> ECS"
    from_port         = 8080
    to_port           = 8080
    protocol          = "tcp"
    security_groups   = [aws_security_group.alb.id]
  }
  egress  { from_port = 0 to_port = 0 protocol = "-1" cidr_blocks = ["0.0.0.0/0"] }
}

# 05 ALB
resource "aws_lb" "this" {
  name               = "${local.name}-alb"
  load_balancer_type = "application"
  security_groups    = [aws_security_group.alb.id]
  subnets            = data.aws_subnets.default_private.ids
}

resource "aws_lb_target_group" "backend" {
  name        = "${local.name}-tg"
  port        = 8080
  protocol    = "HTTP"
  vpc_id      = data.aws_vpc.default.id
  target_type = "ip"
  health_check {
    path                = "/healthz"
    matcher             = "200-399"
    interval            = 30
    timeout             = 5
    healthy_threshold   = 2
    unhealthy_threshold = 2
  }
}

resource "aws_lb_listener" "http" {
  load_balancer_arn = aws_lb.this.arn
  port              = 80
  protocol          = "HTTP"
  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.backend.arn
  }
}

# 06 IAM（ECSタスク実行ロール/タスクロール）
data "aws_iam_policy_document" "ecs_task_trust" {
  statement {
    actions = ["sts:AssumeRole"]
    principals { type = "Service" identifiers = ["ecs-tasks.amazonaws.com"] }
  }
}
resource "aws_iam_role" "ecs_task_execution" {
  name               = "${local.name}-ecs-exec"
  assume_role_policy = data.aws_iam_policy_document.ecs_task_trust.json
}
resource "aws_iam_role_policy_attachment" "ecs_exec_attach" {
  role       = aws_iam_role.ecs_task_execution.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}
resource "aws_iam_role" "ecs_task_role" {
  name               = "${local.name}-ecs-task"
  assume_role_policy = data.aws_iam_policy_document.ecs_task_trust.json
}

# 07 CloudWatch Logs
resource "aws_cloudwatch_log_group" "backend" {
  name              = "/ecs/${local.name}-backend"
  retention_in_days = 14
}

# 08 RDS(PostgreSQL) 簡易
resource "random_password" "db" {
  length  = 20
  special = true
}
resource "aws_db_subnet_group" "db" {
  name       = "${local.name}-db-subnets"
  subnet_ids = data.aws_subnets.default_private.ids
}
resource "aws_security_group" "db" {
  name        = "${local.name}-db-sg"
  description = "DB"
  vpc_id      = data.aws_vpc.default.id
  ingress {
    description     = "ECS -> DB"
    from_port       = 5432
    to_port         = 5432
    protocol        = "tcp"
    security_groups = [aws_security_group.ecs_service.id]
  }
  egress { from_port = 0 to_port = 0 protocol = "-1" cidr_blocks = ["0.0.0.0/0"] }
}
resource "aws_db_instance" "postgres" {
  identifier              = "${local.name}-pg"
  engine                  = "postgres"
  engine_version          = "15"
  instance_class          = var.db_instance_class
  allocated_storage       = var.db_allocated_storage
  db_name                 = "app"
  username                = "app"
  password                = random_password.db.result
  vpc_security_group_ids  = [aws_security_group.db.id]
  db_subnet_group_name    = aws_db_subnet_group.db.name
  skip_final_snapshot     = true
  publicly_accessible     = false
}

# 09 Secrets Manager（DB接続文字列）
resource "aws_secretsmanager_secret" "db_dsn" {
  name = "${local.name}/db_dsn"
}
resource "aws_secretsmanager_secret_version" "db_dsn_v" {
  secret_id     = aws_secretsmanager_secret.db_dsn.id
  secret_string = "postgres://app:${random_password.db.result}@${aws_db_instance.postgres.address}:5432/app?sslmode=disable"
}

# 10 ECS Task Definition & Service（backend）
resource "aws_ecs_task_definition" "backend" {
  family                   = "${local.name}-backend"
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  cpu                      = 512
  memory                   = 1024
  execution_role_arn       = aws_iam_role.ecs_task_execution.arn
  task_role_arn            = aws_iam_role.ecs_task_role.arn
  container_definitions = jsonencode([
    {
      name      = "backend",
      image     = "${aws_ecr_repository.backend.repository_url}:latest",
      essential = true,
      portMappings = [{ containerPort = 8080, hostPort = 8080, protocol = "tcp" }],
      environment = [
        { name = "PORT", value = "8080" },
        { name = "CORS_ALLOW_ORIGINS", value = "https://*.vercel.app" }
      ],
      secrets = [
        { name = "DB_DSN", valueFrom = aws_secretsmanager_secret.db_dsn.arn }
      ],
      logConfiguration = {
        logDriver = "awslogs",
        options = {
          awslogs-group         = aws_cloudwatch_log_group.backend.name,
          awslogs-region        = var.region,
          awslogs-stream-prefix = "ecs"
        }
      }
    }
  ])
  runtime_platform {
    operating_system_family = "LINUX"
    cpu_architecture        = "X86_64"
  }
}

resource "aws_ecs_service" "backend" {
  name            = "${local.name}-backend"
  cluster         = aws_ecs_cluster.this.id
  task_definition = aws_ecs_task_definition.backend.arn
  desired_count   = 1
  launch_type     = "FARGATE"
  network_configuration {
    subnets         = data.aws_subnets.default_private.ids
    security_groups = [aws_security_group.ecs_service.id]
    assign_public_ip = false
  }
  load_balancer {
    target_group_arn = aws_lb_target_group.backend.arn
    container_name   = "backend"
    container_port   = 8080
  }
  depends_on = [aws_lb_listener.http]
}

# 11 Frontend 静的（S3 + CloudFront）
resource "aws_s3_bucket" "frontend" {
  count = var.enable_frontend_static ? 1 : 0
  bucket = "${local.name}-frontend-${var.region}"
  force_destroy = true
}

resource "aws_s3_bucket_website_configuration" "frontend" {
  count  = var.enable_frontend_static ? 1 : 0
  bucket = aws_s3_bucket.frontend[0].id
  index_document { suffix = "index.html" }
  error_document { key = "404.html" }
}

# OAC (Origin Access Control) でS3をCF専用化
resource "aws_cloudfront_origin_access_control" "oac" {
  count                             = var.enable_frontend_static ? 1 : 0
  name                              = "${local.name}-oac"
  description                       = "OAC for S3 static"
  origin_access_control_origin_type = "s3"
  signing_behavior                  = "always"
  signing_protocol                  = "sigv4"
}

resource "aws_cloudfront_distribution" "frontend" {
  count = var.enable_frontend_static ? 1 : 0
  enabled             = true
  default_root_object = "index.html"

  origin {
    domain_name = aws_s3_bucket.frontend[0].bucket_regional_domain_name
    origin_id   = "s3-frontend"
    origin_access_control_id = aws_cloudfront_origin_access_control.oac[0].id
  }

  default_cache_behavior {
    target_origin_id       = "s3-frontend"
    viewer_protocol_policy = "redirect-to-https"
    allowed_methods        = ["GET", "HEAD"]
    cached_methods         = ["GET", "HEAD"]
  }

  restrictions {
    geo_restriction { restriction_type = "none" }
  }

  viewer_certificate {
    cloudfront_default_certificate = true
  }
}