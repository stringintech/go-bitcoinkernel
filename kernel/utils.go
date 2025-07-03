package kernel

func ReverseBytes(data []byte) []byte {
	result := make([]byte, len(data))
	for i, b := range data {
		result[len(data)-1-i] = b
	}
	return result
}

type cResource interface {
	isReady() bool
	uninitializedError() error
}

type cManagedResource interface {
	cResource
	Destroy()
}

func validateReady(r cResource) error {
	if !r.isReady() {
		return r.uninitializedError()
	}
	return nil
}

func checkReady(r cResource) {
	if !r.isReady() {
		panic(r.uninitializedError())
	}
}
