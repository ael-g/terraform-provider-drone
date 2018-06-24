#resource "drone_activated_repository" "terraform-provider-drone" {
#    name = "ael-g/terraform-provider-drone"
#    hooks = [ "push", "pull_request"]
#}

#resource "drone_activated_repository" "test" {
#    name = "ael-g/cicd"
#}

resource "drone_secret" "okook" {
    name = "yeah_new"
    repository = "ael-g/terraform-provider-drone"
    value = "wow woprking"
    events = [ "push", "pull_request", "tag" ]
}
