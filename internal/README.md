# ErnestProvider end to end test suite

This library is provided with a test suite to test connectors against the real end provider.

## Structure

It's divided by provider and resource, so each resource provider can be tested independently.

This software is built on top of gucumber, so if you want to have specific information about it, please visit the [project page](https://github.com/gucumber/gucumber).

Gucumber basically allows you to run tests with 
```
$ gucumber
```
or specific tests with tagging support
```
$ gucumber --tags=@azure
```


## Setup

You'll first need gucumber, we've easily packaged on a Makefile, so all you need to do to get it working is run `make dev-deps`

Additionally each provider may need specific environment variables to work. See the sections bellow to know what are these environment variables.

### Azure specific setup

```
export AZURE_TENANT_ID=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
export AZURE_CLIENT_ID=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
export AZURE_CLIENT_SECRET=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
export AZURE_SUBSCRIPTION_ID=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
export AZURE_ENVIRONMENT=public
```
