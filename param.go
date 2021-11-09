package requests

// IParam Param接口
type IParam interface {
	IntSet(v int64) error                                          // param set the value(int64)
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
