resource "drone_activated_repository" "terraform-provider-drone" {
    name = "ael-g/terraform-provider-drone"
    hooks = [ "push", "pull_request"]
    is_trusted = true
}

#resource "drone_activated_repository" "est" {
#    name = "ael-g/cicd"
#}

