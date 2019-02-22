resource "drone_activated_repository" "service-geoip" {
    name = "habx/service-geoip"
    allow_tag = true
    is_trusted = true
}

resource "drone_secret" "service-geoip_aws_access_key_id" {
    name = "aws_access_key_id"
    repository = "habx/service-geoip"
    value = "AKIAJYG5GA62VKHBGITA"
    events = [ "push", "pull_request", "tag" ]
    pull_request = "true"
}
