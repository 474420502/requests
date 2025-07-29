package requests

type M map[string]interface{}

// Head 请求
func Head(url string) *Request {
	return NewSession().Head(url)
}

// Get 请求
func Get(url string) *Request {
	return NewSession().Get(url)
}

// Post 请求
func Post(url string) *Request {
	return NewSession().Post(url)
}

// Put 请求
func Put(url string) *Request {
	return NewSession().Put(url)
}

// Patch 请求
func Patch(url string) *Request {
	return NewSession().Patch(url)
}

// Delete 请求
func Delete(url string) *Request {
	return NewSession().Delete(url)
}

// Connect 请求
func Connect(url string) *Request {
	return NewSession().Connect(url)
}

// Options 请求
func Options(url string) *Request {
	return NewSession().Options(url)
}

// Trace 请求
func Trace(url string) *Request {
	return NewSession().Trace(url)
}
