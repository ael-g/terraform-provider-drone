#resource "drone_activated_repository" "terraform-provider-drone" {
#    name = "ael-g/terraform-provider-drone"
#}

resource "drone_activated_repository" "test" {
    name = "ael-g/terraform-provider-drone"
    allow_tag = true
}

resource "drone_secret" "okook" {
    name = "yeah_new"
    repository = "ael-g/terraform-provider-drone"
    value = "wow dfsdfsdf"
    events = [ "push", "pull_request", "tag" ]
}
