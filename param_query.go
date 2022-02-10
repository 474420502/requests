package requests

import "strconv"

// import "strconv"

// ParamQuery 参数
type ParamQuery struct {
	Temp *Temporary
	Key  string
}

// Set 单个整型参数设置
func (p *ParamQuery) Set(value interface{}) {
	values := p.Temp.GetQuery()
	switch v := value.(type) {
	case int64:
		values.Set(p.Key, strconv.FormatInt(v, 10))
	case uint64:
		values.Set(p.Key, strconv.FormatUint(v, 10))
	case float64:
		values.Set(p.Key, strconv.FormatFloat(v, 'f', -1, 64))
	case int:
		values.Set(p.Key, strconv.FormatInt(int64(v), 10))
	case int8:
		values.Set(p.Key, strconv.FormatInt(int64(v), 10))
	case int16:
		values.Set(p.Key, strconv.FormatInt(int64(v), 10))
	case int32:
		values.Set(p.Key, strconv.FormatInt(int64(v), 10))
	case uint:
		values.Set(p.Key, strconv.FormatUint(uint64(v), 10))
	case uint8:
		values.Set(p.Key, strconv.FormatUint(uint64(v), 10))
	case uint16:
		values.Set(p.Key, strconv.FormatUint(uint64(v), 10))
	case uint32:
		values.Set(p.Key, strconv.FormatUint(uint64(v), 10))
	case float32:
		values.Set(p.Key, strconv.FormatFloat(float64(v), 'f', -1, 64))
	case string:
		values.Set(p.Key, v)
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

	switch v := value.(type) {
	case int64:
		vs[index] = strconv.FormatInt(v, 10)
	case uint64:
		vs[index] = strconv.FormatUint(v, 10)
	case float64:
		vs[index] = strconv.FormatFloat(v, 'f', -1, 64)
	case int:
		vs[index] = strconv.FormatInt(int64(v), 10)
	case int8:
		vs[index] = strconv.FormatInt(int64(v), 10)
	case int16:
		vs[index] = strconv.FormatInt(int64(v), 10)
	case int32:
		vs[index] = strconv.FormatInt(int64(v), 10)
	case uint:
		vs[index] = strconv.FormatUint(uint64(v), 10)
	case uint8:
		vs[index] = strconv.FormatUint(uint64(v), 10)
	case uint16:
		vs[index] = strconv.FormatUint(uint64(v), 10)
	case uint32:
		vs[index] = strconv.FormatUint(uint64(v), 10)
	case float32:
		vs[index] = strconv.FormatFloat(float64(v), 'f', -1, 64)
	case string:
		vs[index] = v
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
