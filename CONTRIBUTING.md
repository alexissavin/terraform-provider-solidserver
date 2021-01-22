# Setup Go and Workspace

Before contributing, make sure you have a working Go Workspace

* Install Go

* Set a persistent GO PATH by adding these lines into your .bashrc
```
GO111MODULE=on
GOPATH="$HOME/go"
```

* Create a Go workspace and set GO PATH
```
mkdir -p $HOME/go
export GOPATH="$HOME/go"

# Folder contains your golang source codes
mkdir -p $GOPATH/src

# Folder contains the binaries when you install an go based executable
mkdir -p $GOPATH/bin

# Folder contains the Go packages you install
mkdir -p $GOPATH/pkg

# Folder contains the Github Source code for the repos you cloned
mkdir -p $GOPATH/src/github.com

```
All your imports will be resolved from this GO PATH only

For more info: https://golang.org/cmd/go/#hdr-GOPATH_environment_variable

# Cloning the Project Repository

Then, clone the repository. The folder structure to maintain the local repo is *$GOPATH/src/github.com/<org-name>/<project-name>*

For this document, let's assume you're cloning the below repo.
https://github.com/alexissavin/terraform-provider-solidserver

```
mkdir -p $GOPATH/src/github.com/alexissavin/terraform-provider-solidserver
cd $GOPATH/src/github.com/alexissavin/terraform-provider-solidserver

git clone git@github.com:alexissavin/terraform-provider-solidserver
```

# Installing Go Packages

To install a go package, you can use *go get* command.

```
go get
```

Note: Depends on the package, it may install a binary under $GOPATH/bin or a package in $GOPATH/pkg

# Building the Provider for Testing Purpose

* Go into the project's folder and run the make command.
```
cd $GOPATH/src/github.com/alexissavin/terraform-provider-solidserver
make
```

* Create a testing folder and copy the sample test file:
```
mkdir -p _tests
cp tests/tests.tf _tests/tests.tf
cd _tests/tests.tf 
```

* Then create a "variables.tf" file with the following content:
```
variable "solidserver_user" {
  type = string
  default = "<SOLIDserver User>"
}

variable "solidserver_password" {
  type = string
  default = "<SOLIDserver User>"
}

variable "solidserver_host" {
  type = string
  default = "<SOLIDserver IP Address>"
}
```

Keep in mind that the provider built for testing is different from the one published on the terraform registry. The locally generated provider is generated in the following folder:
```
$HOME/.terraform.d/plugins/terraform.efficientip.com/efficientip/${PKG_NAME}/${RELEASE}/${OS_ARCH}
```

Then you must leverage it in your test files:
```
terraform {
  required_providers {
    solidserver = {
      source  = "terraform.efficientip.com/efficientip/solidserver"
      version = ">= 99999.9"
    }
  }
}
```