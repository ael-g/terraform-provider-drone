resource "drone_activated_repository" "terraform-provider-drone" {
    name = "ael-g/terraform-provider-drone"
    hooks = [ "push", "pull_request"]
    is_trusted = true
}

#resource "drone_activated_repository" "test" {
#    name = "ael-g/cicd"
#}

resource "drone_secret" "okook" {
    name = "yeah"
    repository = "ael-g/terraform-provider-drone"
    value = "qsdkhqskdjqhksdjh"
    events = [ "push", "pull_request"]
}

