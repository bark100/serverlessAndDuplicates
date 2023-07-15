output "apigw_url" {
  value = one(module.api_gateway[*].default_apigatewayv2_stage_invoke_url)
}
