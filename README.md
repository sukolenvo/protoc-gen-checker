## About

Helper protoc plugin to clean unused messages from protobuf. 
Prints out orphan messages and enums that are not used by services.

Sample usage:
```bash
protoc --checker_out="." -I . --checker_opt=language_package=java proto/events.proto
```

Go package is used to resolve file to package mapping. Alternatively you can use java package or proto package 
by using option `langauge_package`. E.g.:
```bash
protoc --checker_out="." -I . --checker_opt=language_package=java proto/events.proto
```

Packages can be excluded from with `ignore_package` param. Messages can be excluded with `ignore_message` parame:
```bash
protoc --checker_out="." -I . --checker_opt=ignore_package=google,ignore_message=foo.Bar proto/events.proto
```

Check all proto files in directory oneliner:
```bash
protoc --checker_out="." -I . $(find . -name "*.proto" |tr '\n' ' ')
```

## Install
```bash
wget https://github.com/sukolenvo/protoc-gen-checker/releases/download/v1.0.0/protoc-gen-checker
chmod +x protoc-gen-checker
# move to the PATH e.g mv protoc-gen-checker /usr/local/bin/
# or use -plugin param e.g protoc --plugin="./protoc-gen-checker" ...
```