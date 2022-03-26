package checker

import (
	"errors"
	"fmt"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestChecker(t *testing.T) {
	tests := []struct {
		testFile       string
		expected       string
		ignorePackages []string
		ignoreMessages []string
	}{
		{"test_unused_message.proto", "unused message: 'com.sukolenvo.checker.test.Unused'\nunused message: 'com.sukolenvo.checker.test.Nested'", []string{}, []string{}},
		{"test_unused_enum.proto", "unused message: 'com.sukolenvo.checker.test.UnusedEnum'", []string{}, []string{}},
		{"test_unused_embedded.proto", "unused message: 'com.sukolenvo.checker.test.UnaryRequest.UnaryRequestStatus'", []string{}, []string{}},
		{"test_nounused.proto", "", []string{}, []string{}},
		{"test_unused_enum.proto", "", []string{"com"}, []string{}},
		{"test_unused_embedded.proto", "", []string{"com"}, []string{}},
		{"test_unused_enum.proto", "", []string{"com.sukolenvo.checker.test"}, []string{}},
		{"test_unused_enum.proto", "", []string{}, []string{"com.sukolenvo.checker.test.UnusedEnum"}},
		{"test_unused_enum.proto", "unused message: 'com.sukolenvo.checker.test.UnusedEnum'", []string{"google"}, []string{"google.UnusedEnum"}},
		{"test_map.proto", "unused message: 'com.sukolenvo.checker.test.Unused'", []string{}, []string{}},
	}
	for _, test := range tests {
		t.Run(test.testFile, func(t *testing.T) {
			descriptor, err := ReadTestProto(test.testFile)
			if err != nil {
				t.Fatal("failed to read file", err)
			}
			request := pluginpb.CodeGeneratorRequest{}
			params := fmt.Sprintf("M%s=test.test_message", *descriptor.Name)
			request.Parameter = &params
			request.ProtoFile = []*descriptorpb.FileDescriptorProto{descriptor}
			plugin, err := protogen.Options{}.New(&request)
			if err != nil {
				t.Fatal("failed to prepare plugin", err)
			}
			result := NewChecker(plugin, test.ignorePackages, test.ignoreMessages).Check()
			message := ""
			if result != nil {
				message = result.Error()
			}
			if message != test.expected {
				t.Fatalf("Expected '%s', but got '%v'", test.expected, result)
			}
		})
	}
}

func ReadTestProto(protoFile string) (*descriptorpb.FileDescriptorProto, error) {
	if !strings.HasSuffix(protoFile, ".proto") {
		return nil, errors.New("Expecting proto file but got: " + protoFile)
	}
	if checkFileExists("resources/"+protoFile) != nil {
		return nil, errors.New("Proto file not found: " + protoFile)
	}
	descriptor := protoFile[0:len(protoFile)-len(".proto")] + ".pb"
	descriptorPath := "resources/pb/" + descriptor
	if checkFileExists(descriptorPath) != nil {
		return nil, errors.New(`test file descriptor is missing, run "make compile_test_proto" first to generate ` + descriptorPath)
	}
	data, err := ioutil.ReadFile(descriptorPath)
	if err != nil {
		return nil, err
	}
	fileDescriptorSet := descriptorpb.FileDescriptorSet{}
	err = proto.Unmarshal(data, &fileDescriptorSet)
	if err != nil {
		return nil, err
	}
	if len(fileDescriptorSet.GetFile()) != 1 {
		return nil, errors.New("file descriptor set contains unexpected number of files: " + string(rune(len(fileDescriptorSet.GetFile()))))
	}
	return fileDescriptorSet.GetFile()[0], nil
}

func checkFileExists(filePath string) error {
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			return errors.New("file not exists: " + filePath)
		}
	}
	return nil
}
