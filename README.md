# Micron Language Interpreter
Interpreter for super minimalistic Micron language written in pure Go.

Heavily inspired by Writing An Interpreter In Go book by Thorsten Ball:
* https://interpreterbook.com/
* https://amzn.com/B01N2T1VD2

### Automated builds
Automated builds with tests run after every commit to master. Status can be found below

![Build and test from latest commit on master](https://github.com/jpiechowka/micron-language-interpreter-go/workflows/Build%20and%20test%20from%20latest%20commit%20on%20master/badge.svg)

When version tags are pushed to the repository, release workflow will run and build the binaries for every operating system. After build is successful binaries will be uploaded and release will be published on the releases page: https://github.com/jpiechowka/micron-language-interpreter-go/releases

TODO Workflow + badge here

### Download release binaries
New and latest release binaries will be available to download from the releases page: https://github.com/jpiechowka/micron-language-interpreter-go/releases

### Building from source code
If you do not want to download one of the prebuilt binaries simply execute the commands below to build from source (Go needs to be installed and properly configured, see https://golang.org/doc/install)

```
git clone https://github.com/jpiechowka/micron-language-interpreter-go.git
cd micron-language-interpreter-go
go build -v -a .
```
