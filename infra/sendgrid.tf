resource "linode_domain_record" "dripfile_cname_sendgrid_1" {
  domain_id   = linode_domain.dripfile.id
  record_type = "CNAME"
  name        = "em3.dripfile.com"
  target      = "u22828247.wl248.sendgrid.net"
}

resource "linode_domain_record" "dripfile_cname_sendgrid_2" {
  domain_id   = linode_domain.dripfile.id
  record_type = "CNAME"
  name        = "s1._domainkey.dripfile.com"
  target      = "s1.domainkey.u22828247.wl248.sendgrid.net"
}

resource "linode_domain_record" "dripfile_cname_sendgrid_3" {
  domain_id   = linode_domain.dripfile.id
  record_type = "CNAME"
  name        = "s2._domainkey.dripfile.com"
  target      = "s2.domainkey.u22828247.wl248.sendgrid.net"
}

resource "linode_domain_record" "dripfile_cname_sendgrid_4" {
  domain_id   = linode_domain.dripfile.id
  record_type = "CNAME"
  name        = "url7276.dripfile.com"
  target      = "sendgrid.net"
}

resource "linode_domain_record" "dripfile_cname_sendgrid_5" {
  domain_id   = linode_domain.dripfile.id
  record_type = "CNAME"
  name        = "22828247.dripfile.com"
  target      = "sendgrid.net"
}
