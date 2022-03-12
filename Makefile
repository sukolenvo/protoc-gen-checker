test_proto_dir := ./checker/resources/
test_proto := $(wildcard $(test_proto_dir)test_*.proto)
test_pb := $(patsubst $(test_proto_dir)test_%.proto,$(test_proto_dir)pb/test_%.pb,$(test_proto))

.PHONY: compile_test_proto clean test

$(test_proto_dir)pb/test_%.pb: $(test_proto_dir)test_%.proto
	protoc --descriptor_set_out=$@ $^

compile_test_proto: $(test_pb)

test: compile_test_proto
	go test ./...
clean:
	rm -f $(test_pb)