# Testing grpc with go

Building Modern API and Microservices
with Go and Grpc

Une partie du cours est obsolete... depuis l'introduction des modules.

* Prerequis

sudo apt install -y protobuf-compiler

* Créer le projet

```
cd ~/go
mkdir src
cd src/
mkdir go-grpc
cd go-grpc/
mkdir calculator
mkdir -p greet/greetpb
```

* Installer les packages go pour grpc

OBSOLETE
```
go install google.golang.org/grpc@latest
go install github.com/golang/protobuf/protoc-gen-go@latest
export PATH="$PATH:$(go env GOPATH)/bin"
```
END OBSOLETE

Mise à jour, les commandes précédentes semblent obsoletes

```
go mod init
go mod edit -require google.golang.org/grpc@v1.55.0
go mod edit -require google.golang.org/protobuf@v1.30.0
go get -t ./...
export PATH="$PATH:$(go env GOPATH)/bin"
```

* Créer le fichier greet/greetpb/greet.proto

```
syntax = "proto3";

package greet;
option go_package="./greet/greetpb";

service GreetService{}
```

* Compiler le fichier

protoc greet/greetpb/greet.proto --go_out=plugins=grpc:.

