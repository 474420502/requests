package requests

import (
	"reflect"
	"strconv"
)

// import "strconv"

// ParamQuery 参数
type ParamQuery struct {
	Temp *Temporary
	Key  string
}

// case string:
// 	values.Set(p.Key, v)
// case fmt.Stringer:
// 	values.Set(p.Key, v.String())
// case int, int8, int16, int32, int64,
// 	 uint, uint8, uint16, uint32, uint64:
// 	values.Set(p.Key, strconv.FormatInt(int64(v), 10))
// case float32, float64:
// 	values.Set(p.Key, strconv.FormatFloat(float64(v), 'f', -1, 64))

// Set 单个整型参数设置
func (p *ParamQuery) Set(value interface{}) {
	values := p.Temp.GetQuery()

	vv := reflect.ValueOf(value)
	switch k := vv.Kind(); k {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		values.Set(p.Key, strconv.FormatInt(vv.Int(), 10))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		values.Set(p.Key, strconv.FormatUint(vv.Uint(), 10))
	case reflect.Float32, reflect.Float64:
		values.Set(p.Key, strconv.FormatFloat(vv.Float(), 'f', -1, 64))
	case reflect.String:
		values.Set(p.Key, vv.String())
	}
	p.Temp.SetQuery(values)
}

// Add   通用类型 参数加减 value 为通用计算类型
func (p *ParamQuery) Add(value interface{}) error {
	values := p.Temp.GetQuery()
	vs := values[p.Key]

	switch v := value.(type) {
	case int64:
		pvalue, err := strconv.ParseInt(vs[0], 10, 64)
		if err != nil {
			return err
		}
		pvalue += v
		values.Set(p.Key, strconv.FormatInt(pvalue, 10))
	case uint64:
		pvalue, err := strconv.ParseUint(vs[0], 10, 64)
		if err != nil {
			return err
		}
		pvalue += v
		values.Set(p.Key, strconv.FormatUint(pvalue, 10))
	case float64:
		pvalue, err := strconv.ParseFloat(vs[0], 10)
		if err != nil {
			return err
		}
		pvalue += v
		values.Set(p.Key, strconv.FormatFloat(pvalue, 'f', -1, 64))
	case int:
		pvalue, err := strconv.ParseInt(vs[0], 10, 64)
		if err != nil {
			return err
		}
		pvalue += int64(v)
		values.Set(p.Key, strconv.FormatInt(pvalue, 10))
	case int8:
		pvalue, err := strconv.ParseInt(vs[0], 10, 64)
		if err != nil {
			return err
		}
		pvalue += int64(v)
		values.Set(p.Key, strconv.FormatInt(pvalue, 10))
	case int16:
		pvalue, err := strconv.ParseInt(vs[0], 10, 64)
		if err != nil {
			return err
		}
		pvalue += int64(v)
		values.Set(p.Key, strconv.FormatInt(pvalue, 10))
	case int32:
		pvalue, err := strconv.ParseInt(vs[0], 10, 64)
		if err != nil {
			return err
		}
		pvalue += int64(v)
		values.Set(p.Key, strconv.FormatInt(pvalue, 10))
	case float32:
		pvalue, err := strconv.ParseFloat(vs[0], 10)
		if err != nil {
			return err
		}
		pvalue += float64(v)
		values.Set(p.Key, strconv.FormatFloat(pvalue, 'f', -1, 64))

	case uint:
		pvalue, err := strconv.ParseUint(vs[0], 10, 64)
		if err != nil {
			return err
		}
		pvalue += uint64(v)
		values.Set(p.Key, strconv.FormatUint(pvalue, 10))
	case uint8:
		pvalue, err := strconv.ParseUint(vs[0], 10, 64)
		if err != nil {
			return err
		}
		pvalue += uint64(v)
		values.Set(p.Key, strconv.FormatUint(pvalue, 10))
	case uint16:
		pvalue, err := strconv.ParseUint(vs[0], 10, 64)
		if err != nil {
			return err
		}
		pvalue += uint64(v)
		values.Set(p.Key, strconv.FormatUint(pvalue, 10))
	case uint32:
		pvalue, err := strconv.ParseUint(vs[0], 10, 64)
		if err != nil {
			return err
		}
		pvalue += uint64(v)
		values.Set(p.Key, strconv.FormatUint(pvalue, 10))
	}
	p.Temp.SetQuery(values)
	return nil
}

// ArraySet 通用数组类型 根据 index 设置 value 为通用计算类型
func (p *ParamQuery) ArraySet(index int, value interface{}) {

	values := p.Temp.GetQuery()
	vs := values[p.Key]

	vv := reflect.ValueOf(value)
	switch k := vv.Kind(); k {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		vs[index] = strconv.FormatInt(vv.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		vs[index] = strconv.FormatUint(vv.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		vs[index] = strconv.FormatFloat(vv.Float(), 'f', -1, 64)
	case reflect.String:
		vs[index] = vv.String()
	}

	p.Temp.SetQuery(values)
}

// ArrayAdd 通用数组类型 根据 index 参数加减 value 为通用计算类型
func (p *ParamQuery) ArrayAdd(index int, value interface{}) error {
	values := p.Temp.GetQuery()
	vs := values[p.Key]

	switch v := value.(type) {
	case int64:
		pvalue, err := strconv.ParseInt(vs[0], 10, 64)
		if err != nil {
			return err
		}
		pvalue += v
		vs[index] = strconv.FormatInt(pvalue, 10)

	case uint64:
		pvalue, err := strconv.ParseUint(vs[0], 10, 64)
		if err != nil {
			return err
		}
		pvalue += v
		vs[index] = strconv.FormatUint(pvalue, 10)

	case float64:
		pvalue, err := strconv.ParseFloat(vs[0], 10)
		if err != nil {
			return err
		}
		pvalue += v
		vs[index] = strconv.FormatFloat(pvalue, 'f', -1, 64)

	case int:
		pvalue, err := strconv.ParseInt(vs[0], 10, 64)
		if err != nil {
			return err
		}
		pvalue += int64(v)
		vs[index] = strconv.FormatInt(pvalue, 10)

	case int8:
		pvalue, err := strconv.ParseInt(vs[0], 10, 64)
		if err != nil {
			return err
		}
		pvalue += int64(v)
		vs[index] = strconv.FormatInt(pvalue, 10)

	case int16:
		pvalue, err := strconv.ParseInt(vs[0], 10, 64)
		if err != nil {
			return err
		}
		pvalue += int64(v)
		vs[index] = strconv.FormatInt(pvalue, 10)

	case int32:
		pvalue, err := strconv.ParseInt(vs[0], 10, 64)
		if err != nil {
			return err
		}
		pvalue += int64(v)
		vs[index] = strconv.FormatInt(pvalue, 10)

	case float32:
		pvalue, err := strconv.ParseFloat(vs[0], 10)
		if err != nil {
			return err
		}
		pvalue += float64(v)
		vs[index] = strconv.FormatFloat(pvalue, 'f', -1, 64)

	case uint:
		pvalue, err := strconv.ParseUint(vs[0], 10, 64)
		if err != nil {
			return err
		}
		pvalue += uint64(v)
		vs[index] = strconv.FormatUint(pvalue, 10)

	case uint8:
		pvalue, err := strconv.ParseUint(vs[0], 10, 64)
		if err != nil {
			return err
		}
		pvalue += uint64(v)
		vs[index] = strconv.FormatUint(pvalue, 10)

	case uint16:
		pvalue, err := strconv.ParseUint(vs[0], 10, 64)
		if err != nil {
			return err
		}
		pvalue += uint64(v)
		vs[index] = strconv.FormatUint(pvalue, 10)

	case uint32:
		pvalue, err := strconv.ParseUint(vs[0], 10, 64)
		if err != nil {
			return err
		}
		pvalue += uint64(v)
		vs[index] = strconv.FormatUint(pvalue, 10)

	}
	p.Temp.SetQuery(values)
	return nil
}

// IntSet 单个整型参数设置
func (p *ParamQuery) IntSet(v int64) {
	values := p.Temp.GetQuery()
	values.Set(p.Key, strconv.FormatInt(v, 10))
	p.Temp.SetQuery(values)
}

// IntAdd 单个整型参数计算
func (p *ParamQuery) IntAdd(v int64) error {
	values := p.Temp.GetQuery()
	vs := values[p.Key]
	pvalue, err := strconv.ParseInt(vs[0], 10, 64)
	if err != nil {
		return err
	}
	pvalue += v
	values.Set(p.Key, strconv.FormatInt(pvalue, 10))
	p.Temp.SetQuery(values)
	return nil
}

// IntArraySet 数组整型参数计算
func (p *ParamQuery) IntArraySet(index int, v int64) {
	values := p.Temp.GetQuery()
	vs := values[p.Key]
	vs[index] = strconv.FormatInt(v, 10)
	p.Temp.SetQuery(values)
}

// IntArrayAdd 数组整型参数计算
func (p *ParamQuery) IntArrayAdd(index int, v int64) error {
	values := p.Temp.GetQuery()
	vs := values[p.Key]
	pvalue, err := strconv.ParseInt(vs[index], 10, 64)
	if err != nil {
		return err
	}
	pvalue += v
	vs[index] = strconv.FormatInt(pvalue, 10)
	p.Temp.SetQuery(values)
	return nil
}

// IntArrayDo 数组整型参数操作 do i 数组索引 pvalue 数组值 返回值interface{} 如果nil. 则不变
func (p *ParamQuery) IntArrayDo(do func(i int, pvalue int64) interface{}) error {
	values := p.Temp.GetQuery()
	vs := values[p.Key]
	for i, v := range vs {
		pvalue, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return err
		}
		rvalue := do(i, pvalue)
		if rvalue != nil {
			if err, ok := rvalue.(error); ok {
				return err
			}
			vs[i] = strconv.FormatInt(rvalue.(int64), 10)
		}
	}
	p.Temp.SetQuery(values)
	return nil
}

// FloatSet 单个浮点参数设置
func (p *ParamQuery) FloatSet(v float64) {
	values := p.Temp.GetQuery()
	values.Set(p.Key, strconv.FormatFloat(v, 'f', -1, 64))
	p.Temp.SetQuery(values)
}

// FloatAdd 单个浮点参数计算
func (p *ParamQuery) FloatAdd(v float64) error {
	values := p.Temp.GetQuery()
	vs := values[p.Key]
	pvalue, err := strconv.ParseFloat(vs[0], 64)
	if err != nil {
		return err
	}
	pvalue += v
	values.Set(p.Key, strconv.FormatFloat(pvalue, 'f', -1, 64))
	p.Temp.SetQuery(values)
	return nil
}

// FloatArraySet 数组浮点参数设置
func (p *ParamQuery) FloatArraySet(index int, v float64) {
	values := p.Temp.GetQuery()
	vs := values[p.Key]
	vs[index] = strconv.FormatFloat(v, 'f', -1, 64)
	p.Temp.SetQuery(values)
}

// FloatArrayAdd 数组浮点参数计算
func (p *ParamQuery) FloatArrayAdd(index int, v float64) error {
	values := p.Temp.GetQuery()
	vs := values[p.Key]
	pvalue, err := strconv.ParseFloat(vs[index], 64)
	if err != nil {
		return err
	}
	pvalue += v
	vs[index] = strconv.FormatFloat(pvalue, 'f', -1, 64)
	p.Temp.SetQuery(values)
	return nil
}

// FloatArrayDo 数组整型参数操作 do i 数组索引 pvalue 数组值 返回值interface{} 如果nil. 则不变
func (p *ParamQuery) FloatArrayDo(do func(i int, pvalue float64) interface{}) error {
	values := p.Temp.GetQuery()
	vs := values[p.Key]
	for i, v := range vs {
		pvalue, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return err
		}
		rvalue := do(i, pvalue)
		if rvalue != nil {
			if err, ok := rvalue.(error); ok {
				return err
			}
			vs[i] = strconv.FormatFloat(rvalue.(float64), 'f', -1, 64)
		}
	}
	p.Temp.SetQuery(values)
	return nil
}

// StringSet 字符串参数设置
func (p *ParamQuery) StringSet(v string) {
	values := p.Temp.GetQuery()
	values.Set(p.Key, v)
	p.Temp.SetQuery(values)
}

// StringArraySet 数组字符串参数设置
func (p *ParamQuery) StringArraySet(index int, v string) {
	values := p.Temp.GetQuery()
	vs := values[p.Key]
	vs[index] = v
	p.Temp.SetQuery(values)
}
