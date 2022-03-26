package checker

import (
	"errors"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
	"strings"
)

type Checker struct {
	plugin         *protogen.Plugin
	ignorePackages []string
	ignoreMessages []string
}

func NewChecker(plugin *protogen.Plugin, ignorePackages []string, ignoreMessages []string) *Checker {
	return &Checker{
		plugin:         plugin,
		ignorePackages: ignorePackages,
		ignoreMessages: ignoreMessages,
	}
}

func (c *Checker) Check() error {
	messages := make(map[protoreflect.FullName]bool)
	for _, file := range c.plugin.Files {
		for _, service := range file.Services {
			for _, method := range service.Methods {
				c.insertRecursive(messages, method.Input)
				c.insertRecursive(messages, method.Output)
			}
		}
	}
	unused := []protoreflect.FullName{}
	for _, file := range c.plugin.Files {
		for _, message := range file.Messages {
			unused = append(unused, c.checkUnusedRecursive(messages, message)...)
		}
		for _, value := range file.Enums {
			if !messages[value.Desc.FullName()] {
				unused = append(unused, value.Desc.FullName())
			}
		}
	}
	result := []protoreflect.FullName{}
	for _, item := range unused {
		ignore := false
		for _, filter := range c.ignoreMessages {
			if string(item) == filter {
				ignore = true
			}
		}
		for _, filter := range c.ignorePackages {
			if strings.HasPrefix(string(item), filter+".") {
				ignore = true
			}
		}
		if !ignore {
			result = append(result, item)
		}
	}
	if len(result) != 0 {
		var sb strings.Builder
		for i, name := range result {
			sb.WriteString("unused message: '")
			sb.WriteString(string(name))
			sb.WriteString("'")
			if i != len(result)-1 {
				sb.WriteString("\n")
			}
		}
		return errors.New(sb.String())
	}
	return nil
}

func (c *Checker) insertRecursive(messages map[protoreflect.FullName]bool, message *protogen.Message) {
	if messages[message.Desc.FullName()] {
		return
	}
	messages[message.Desc.FullName()] = true
	for _, field := range message.Fields {
		if field.Message != nil {
			c.insertRecursive(messages, field.Message)
		}
		if field.Enum != nil {
			messages[field.Enum.Desc.FullName()] = true
		}
	}
}

func (c *Checker) checkUnusedRecursive(messages map[protoreflect.FullName]bool, message *protogen.Message) []protoreflect.FullName {
	result := []protoreflect.FullName{}
	if !messages[message.Desc.FullName()] {
		result = append(result, message.Desc.FullName())
	}
	if message.Enums != nil {
		for _, value := range message.Enums {
			if !messages[value.Desc.FullName()] {
				result = append(result, value.Desc.FullName())
			}
		}
	}
	if message.Messages != nil {
		for _, embedded := range message.Messages {
			if !embedded.Desc.IsMapEntry() {
				result = append(result, c.checkUnusedRecursive(messages, embedded)...)
			}
		}
	}
	return result
}
