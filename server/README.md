```
export AWS_IDENTITY_POOL_ID= .... # your identity pool ID Terraform outputs
export AWS_LOGIN_PROVIDER=login.myapp.example.com # == 'cognito_provider_name' in secret.tfvars
export AWS_TOKEN_DURATION_SECONDS=600
export AWS_API_GATEWAY_URL= .... # your api gateway url Terraform outputs
```