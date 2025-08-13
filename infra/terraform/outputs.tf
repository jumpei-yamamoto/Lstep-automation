output "alb_dns_name" {
  value = aws_lb.this.dns_name
}
output "ecr_backend" {
  value = aws_ecr_repository.backend.repository_url
}
output "ecr_frontend_ssr" {
  value       = aws_ecr_repository.frontend_ssr.repository_url
  description = "SSR運用時に使用"
}
output "db_address" {
  value = aws_db_instance.postgres.address
}
output "cloudfront_domain" {
  value       = var.enable_frontend_static ? aws_cloudfront_distribution.frontend[0].domain_name : null
  description = "静的サイトのCFドメイン"
}