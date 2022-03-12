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