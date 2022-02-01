resource "linode_domain_record" "dripfile_mx_zoho_1" {
  domain_id   = linode_domain.dripfile.id
  record_type = "MX"
  target      = "mx.zoho.com"
  priority    = 10
}

resource "linode_domain_record" "dripfile_mx_zoho_2" {
  domain_id   = linode_domain.dripfile.id
  record_type = "MX"
  target      = "mx2.zoho.com"
  priority    = 20
}

resource "linode_domain_record" "dripfile_mx_zoho_3" {
  domain_id   = linode_domain.dripfile.id
  record_type = "MX"
  target      = "mx3.zoho.com"
  priority    = 50
}

resource "linode_domain_record" "dripfile_txt_zoho_spf" {
  domain_id   = linode_domain.dripfile.id
  record_type = "TXT"
  target      = "v=spf1 include:zoho.com ~all"
}

resource "linode_domain_record" "dripfile_txt_zoho_dkim" {
  domain_id   = linode_domain.dripfile.id
  record_type = "TXT"
  name        = "zmail._domainkey"
  target      = "v=DKIM1; k=rsa; p=MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDC5yvvQmOTga/oreslFFLA8OznLX8XE1hydqJNMy6CrL38lR/fWVZ48GSxIlNS+OCVMRDUb0qEzRL7tXWnJLW58uQPKWgNJZpuUBY8uKVjhdOUHWfGwneRq/q7CAjHAV32otx/0+P5mft7lyUpPkIjbUq2qOBwZReLeDLzo0nFWQIDAQAB"
}
