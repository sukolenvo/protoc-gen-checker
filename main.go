package main

import (
	"fmt"
	"github.com/sukolenvo/protoc-gen-checker/checker"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	if len(os.Args) > 1 {
		fmt.Fprintf(os.Stderr, "unknown argument %q (this program should be run by protoc, not directly)\n", os.Args[1])
		os.Exit(1)
	}
	in, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unknown argument %q (this program should be run by protoc, not directly)\n", os.Args[1])
		os.Exit(1)
	}
	req := &pluginpb.CodeGeneratorRequest{}
	if err := proto.Unmarshal(in, req); err != nil {
		fmt.Fprintf(os.Stderr, "failed to unmarshal request %v\n", err)
		os.Exit(1)
	}
	params := parseParams(req.GetParameter())
	packageProvider := func(file *descriptorpb.FileDescriptorProto) string {
		return file.GetOptions().GetGoPackage()
	}
	switch params.overridePackage {
	case "":
		// Nothing to do
	case "java":
		packageProvider = func(file *descriptorpb.FileDescriptorProto) string {
			return file.GetOptions().GetJavaPackage()
		}
	case "proto":
		packageProvider = func(file *descriptorpb.FileDescriptorProto) string {
			return file.GetPackage()
		}
	default:
		fmt.Fprintf(os.Stderr, "unsupported language_package param. Should be one of: java, proto\n")
		os.Exit(1)
	}
	for _, file := range req.GetProtoFile() {
		value := packageProvider(file)
		if file.Options == nil {
			file.Options = &descriptorpb.FileOptions{}
		}
		file.Options.GoPackage = &value
		if file.GetOptions().GetGoPackage() == "" {
			fmt.Fprintf(os.Stderr, "go_package is not found for %s. Set go_package or use language_package param\n", file.GetName())
			os.Exit(1)
		}
	}
	plugin, err := protogen.Options{}.New(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to init plugin %v\n", err)
		os.Exit(1)
	}
	if err := checker.NewChecker(plugin, params.ignorePackages, params.ignoreMessages).Check(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

type Params struct {
	overridePackage string
	ignorePackages  []string
	ignoreMessages  []string
}

func parseParams(params string) *Params {
	result := &Params{
		overridePackage: "",
		ignorePackages:  []string{},
		ignoreMessages:  []string{},
	}
	for _, param := range strings.Split(params, ",") {
		var value string
		if i := strings.Index(param, "="); i >= 0 {
			value = param[i+1:]
			param = param[0:i]
		}
		switch param {
		case "":
			// Ignore.
		case "language_package":
			result.overridePackage = value
		case "ignore_package":
			result.ignorePackages = append(result.ignorePackages, value)
		case "ignore_message":
			result.ignoreMessages = append(result.ignoreMessages, value)
		}
	}
	return result
}
