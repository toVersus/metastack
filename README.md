# metastack

metastack is a simple helper CLI for interacting with OpenStack instance metadata (properties).

As you know, `openstack server show` command doesn't print machine friendly output like below:

```bash
$ openstack server show foo-api01 -f json | jq '.properties'
"Environment='stg', ManagedBy='terraform', Role='api'"
```

Otherwise, metastack prints server metadata as JSON format, so you can easily pass that output to [jq](https://stedolan.github.io/jq/) or other common utilities.

## Before you start

[Set environment variables using the OpenStack RC file](https://docs.openstack.org/zh_CN/user-guide/common/cli-set-environment-variables-using-openstack-rc.html)

## Usage

```bash
$ metastack foo-api01 | jq
{
  "Environment": "stg",
  "ManagedBy": "terraform",
  "Role": "api"
}
```
