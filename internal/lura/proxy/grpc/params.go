package grpc

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"strconv"
)

func populateProtoMessageFromParams(msg proto.Message, params map[string]string) error {
	m := msg.ProtoReflect()
	md := m.Descriptor()

	for key, val := range params {
		field := md.Fields().ByName(protoreflect.Name(key))
		if field == nil {
			continue
		}

		// Skip repeated or map fields for simplicity
		if field.IsList() || field.IsMap() {
			continue
		}

		var parsed protoreflect.Value
		var err error

		switch field.Kind() {
		case protoreflect.StringKind:
			parsed = protoreflect.ValueOfString(val)

		case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
			var i int64
			i, err = strconv.ParseInt(val, 10, 32)
			parsed = protoreflect.ValueOfInt32(int32(i))

		case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
			var i int64
			i, err = strconv.ParseInt(val, 10, 64)
			parsed = protoreflect.ValueOfInt64(i)

		case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
			var u uint64
			u, err = strconv.ParseUint(val, 10, 32)
			parsed = protoreflect.ValueOfUint32(uint32(u))

		case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
			var u uint64
			u, err = strconv.ParseUint(val, 10, 64)
			parsed = protoreflect.ValueOfUint64(u)

		case protoreflect.BoolKind:
			var b bool
			b, err = strconv.ParseBool(val)
			parsed = protoreflect.ValueOfBool(b)

		case protoreflect.FloatKind:
			var f float64
			f, err = strconv.ParseFloat(val, 32)
			parsed = protoreflect.ValueOfFloat32(float32(f))

		case protoreflect.DoubleKind:
			var f float64
			f, err = strconv.ParseFloat(val, 64)
			parsed = protoreflect.ValueOfFloat64(f)

		default:
			continue // skip unsupported kinds like enums, bytes, messages
		}

		if err != nil {
			return errInvalidParams
		}

		m.Set(field, parsed)
	}
	return nil
}
