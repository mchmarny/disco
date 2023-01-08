# Description: This file contains the file resources for the deployment

# data.template_file.version.rendered
data "template_file" "version" {
  template = file("../.version")
}

# data.template_file.schema_packages.rendered
data "template_file" "schema_packages" {
  template = file("schema/packages.json")
}

# data.template_file.schema_vulnerabilities.rendered
data "template_file" "schema_vulnerabilities" {
  template = file("schema/vulnerabilities.json")
}

# data.template_file.schema_licenses.rendered
data "template_file" "schema_licenses" {
  template = file("schema/licenses.json")
}


