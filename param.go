package requests

// IParam Param接口
// Deprecated: 这个复杂的参数接口系统已废弃。
// 对于查询参数，请使用类型安全的 AddQueryInt, AddQueryBool, AddQueryFloat 等方法。
// 对于路径参数，请使用 SetPathParam 和 SetPathParams 方法。
// 这些新方法更简单、更安全，且提供更好的开发体验。
type IParam interface {
	Set(v interface{})
	Add(v interface{}) error
	ArraySet(index int, v interface{})
	ArrayAdd(index int, v interface{}) error

	IntSet(v int64)                                                // param set the value(int64)
	IntAdd(v int64) error                                          // param add the value(int64)
	IntArraySet(index int, v int64)                                // params set the values([]int64) by index
	IntArrayAdd(index int, v int64) error                          // params add the values([]int64) by index
	IntArrayDo(do func(i int, pvalue int64) interface{}) error     // range params and set value(int64)
	FloatSet(v float64)                                            // param set the value(float64)
	FloatAdd(v float64) error                                      // param add the value(float64)
	FloatArraySet(index int, v float64)                            // params set the values([]float64) by index
	FloatArrayAdd(index int, v float64) error                      // params add the values([]float64) by index
	FloatArrayDo(do func(i int, pvalue float64) interface{}) error // range params and set value(float54)
	StringSet(v string)                                            // param set the value(string)
	StringArraySet(index int, v string)                            // params set the values([]string) by index
}
