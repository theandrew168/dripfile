resource "linode_domain" "dripfile" {
  type      = "master"
  domain    = "dripfile.com"
  soa_email = "info@shallowbrooksoftware.com"
}
