package errors

func ProvideServerHandler() *ServerHandler {
	return NewServerHandler()
}
