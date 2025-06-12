package grpc

import (
	"encoding/base64"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"strconv"
	"strings"
)

//nolint:gochecknoglobals
var simpleParsers = map[protoreflect.Kind]func(string) protoreflect.Value{
	protoreflect.StringKind: protoreflect.ValueOfString,
}

//nolint:gochecknoglobals
var parsedParsers = map[protoreflect.Kind]func(string) (protoreflect.Value, error){
	protoreflect.Int32Kind:    parseInt32,
	protoreflect.Sint32Kind:   parseInt32,
	protoreflect.Sfixed32Kind: parseInt32,
	protoreflect.Int64Kind:    parseInt64,
	protoreflect.Sint64Kind:   parseInt64,
	protoreflect.Sfixed64Kind: parseInt64,
	protoreflect.Uint32Kind:   parseUint32,
	protoreflect.Fixed32Kind:  parseUint32,
	protoreflect.Uint64Kind:   parseUint64,
	protoreflect.Fixed64Kind:  parseUint64,
	protoreflect.BoolKind:     parseBool,
	protoreflect.FloatKind:    parseFloat32,
	protoreflect.DoubleKind:   parseFloat64,
	protoreflect.BytesKind:    parseBytes,
}

func parseInt32(val string) (protoreflect.Value, error) {
	i, err := strconv.ParseInt(val, 10, 32)
	return protoreflect.ValueOfInt32(int32(i)), err
}

func parseInt64(val string) (protoreflect.Value, error) {
	i, err := strconv.ParseInt(val, 10, 64)
	return protoreflect.ValueOfInt64(i), err
}

func parseUint32(val string) (protoreflect.Value, error) {
	u, err := strconv.ParseUint(val, 10, 32)
	return protoreflect.ValueOfUint32(uint32(u)), err
}

func parseUint64(val string) (protoreflect.Value, error) {
	u, err := strconv.ParseUint(val, 10, 64)
	return protoreflect.ValueOfUint64(u), err
}

func parseBool(val string) (protoreflect.Value, error) {
	b, err := strconv.ParseBool(val)
	return protoreflect.ValueOfBool(b), err
}

func parseFloat32(val string) (protoreflect.Value, error) {
	f, err := strconv.ParseFloat(val, 32)
	return protoreflect.ValueOfFloat32(float32(f)), err
}

func parseFloat64(val string) (protoreflect.Value, error) {
	f, err := strconv.ParseFloat(val, 64)
	return protoreflect.ValueOfFloat64(f), err
}

func parseBytes(val string) (protoreflect.Value, error) {
	b, err := base64.StdEncoding.DecodeString(val)
	return protoreflect.ValueOfBytes(b), err
}

func populateProtoMessageFromParams(msg proto.Message, params map[string][]string) error {
	m := msg.ProtoReflect()
	md := m.Descriptor()

	for key, values := range params {
		fieldPath := strings.Split(key, ".")
		if err := setProtoField(m, md, fieldPath, values); err != nil {
			return err
		}
	}
	return nil
}

func setProtoField(m protoreflect.Message, md protoreflect.MessageDescriptor, path []string, values []string) error {
	if len(path) == 0 {
		return nil
	}

	field := md.Fields().ByName(protoreflect.Name(path[0]))
	if field == nil {
		return nil
	}

	if len(path) > 1 {
		subMsg := m.Mutable(field).Message()
		return setProtoField(subMsg, field.Message(), path[1:], values)
	}

	if field.IsMap() {
		return setMapField(m, field, values)
	}

	if field.IsList() {
		list := m.Mutable(field).List()
		for _, val := range values {
			parsed, err := parseProtoFieldValue(field, val)
			if err != nil {
				return errInvalidParams
			}
			list.Append(parsed)
		}
		return nil
	}

	if len(values) > 0 {
		parsed, err := parseProtoFieldValue(field, values[0])
		if err != nil {
			return errInvalidParams
		}
		m.Set(field, parsed)
	}
	return nil
}

func setMapField(m protoreflect.Message, field protoreflect.FieldDescriptor, values []string) error {
	mp := m.Mutable(field).Map()
	keyKind := field.MapKey().Kind()
	valueKind := field.MapValue().Kind()

	for _, val := range values {
		parts := strings.SplitN(val, ":", 2)
		if len(parts) != 2 {
			return errInvalidParams
		}
		keyStr, valueStr := parts[0], parts[1]

		keyVal, err := parseProtoFieldValueByKind(keyKind, keyStr)
		if err != nil {
			return errInvalidParams
		}

		if valueKind == protoreflect.MessageKind {
			return errUnsupportedKind
		}

		valueVal, err := parseProtoFieldValueByKind(valueKind, valueStr)
		if err != nil {
			return errInvalidParams
		}

		mp.Set(keyVal.MapKey(), valueVal)
	}
	return nil
}

func parseProtoFieldValue(field protoreflect.FieldDescriptor, val string) (protoreflect.Value, error) {
	kind := field.Kind()
	if kind == protoreflect.EnumKind {
		if num, err := strconv.ParseInt(val, 10, 32); err == nil {
			return protoreflect.ValueOfEnum(protoreflect.EnumNumber(num)), nil
		}
		if ev := field.Enum().Values().ByName(protoreflect.Name(val)); ev != nil {
			return protoreflect.ValueOfEnum(ev.Number()), nil
		}
		return protoreflect.Value{}, errInvalidParams
	}
	return parseProtoFieldValueByKind(kind, val)
}

func parseProtoFieldValueByKind(kind protoreflect.Kind, val string) (protoreflect.Value, error) {
	if parser, ok := simpleParsers[kind]; ok {
		return parser(val), nil
	}
	if parser, ok := parsedParsers[kind]; ok {
		return parser(val)
	}
	return protoreflect.Value{}, errUnsupportedKind
}
