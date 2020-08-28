package requests

// IParam Param接口
type IParam interface {
	IntSet(v int64) error
	IntAdd(v int64) error
	IntArraySet(index int, v int64)
	IntArrayAdd(index int, v int64) error
	IntArrayDo(do func(i int, pvalue int64) interface{}) error
	FloatSet(v float64)
	FloatAdd(v float64) error
	FloatArraySet(index int, v float64)
	FloatArrayAdd(index int, v float64) error
	FloatArrayDo(do func(i int, pvalue float64) interface{}) error
	StringSet(v string)
	StringArraySet(index int, v string)
}
