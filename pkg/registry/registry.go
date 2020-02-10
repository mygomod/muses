package registry

// Registry is the top abstraction which will register 'something' into register.
// We don't want to limit the
type Registry interface {
	Register(key string, value interface{})
}