```
GOARCH=amd64 GOOS=linux go build -o ../functions/bin/sample ../functions/src/*.go
mkdir ../functions/archive && zip -j ../functions/archive/sample.zip ../functions/bin/sample
```

```
terraform init
terraform plan -var-file="secret.tfvars"
terraform apply -var-file="secret.tfvars"
```