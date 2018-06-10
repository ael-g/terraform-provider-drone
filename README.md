# terraform-provider-drone
Terraform Drone provider

# Build
```
go get ./...
go build -o terraform-provider-drone
```

# Test
```
go build -o terraform-provider-drone && terraform init && TF_LOG=DEBUG terraform apply
```
