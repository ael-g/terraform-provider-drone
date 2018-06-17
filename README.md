# terraform-provider-drone
Terraform Drone provider

# Build
```
Create const file with:
const (
        token = ""
        host  = ""
)
go get ./...
go build -o terraform-provider-drone
```

# Test
```
go build -o terraform-provider-drone && terraform init && TF_LOG=DEBUG terraform apply
```
