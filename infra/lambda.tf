module "playerListLambda" {
  source  = "terraform-aws-modules/lambda/aws"
  version = "7.20.1"
  function_name = "playerList"
  
}