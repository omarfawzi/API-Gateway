package grpc

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
	"os"
	"strings"
)

func methodDescriptorFromFile(
	descriptorPath, serviceName, methodName string,
) (protoreflect.MethodDescriptor, error) {
	fds, err := readDescriptorSet(descriptorPath)
	if err != nil {
		return nil, err
	}

	files, err := protodesc.NewFiles(fds)
	if err != nil {
		return nil, fmt.Errorf("[gRPC] create file registry: %w", err)
	}

	svcDesc, err := findServiceDescriptor(files, fds, serviceName)
	if err != nil {
		return nil, err
	}

	methodDesc := svcDesc.Methods().ByName(protoreflect.Name(methodName))
	if methodDesc == nil {
		return nil, errInvalidMethod
	}

	return methodDesc, nil
}

func readDescriptorSet(path string) (*descriptorpb.FileDescriptorSet, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("[gRPC] read descriptor file: %w", err)
	}

	var fds descriptorpb.FileDescriptorSet
	if err := proto.Unmarshal(data, &fds); err != nil {
		return nil, fmt.Errorf("[gRPC] unmarshal descriptor: %w", err)
	}
	return &fds, nil
}

func findServiceDescriptor(
	files *protoregistry.Files,
	fds *descriptorpb.FileDescriptorSet,
	targetService string,
) (protoreflect.ServiceDescriptor, error) {
	for _, file := range fds.File {
		fd, err := files.FindFileByPath(file.GetName())
		if err != nil {
			continue
		}
		for i := 0; i < fd.Services().Len(); i++ {
			svc := fd.Services().Get(i)
			if string(svc.FullName()) == targetService {
				return svc, nil
			}
		}
	}
	return nil, errInvalidService
}

func parseServiceMethod(full string) (string, string, error) {
	parts := strings.Split(full, "/")
	if len(parts) != 2 {
		return "", "", errInvalidMethodFormat
	}
	return parts[0], parts[1], nil
}
