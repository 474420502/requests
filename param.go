package requests

import "strconv"

// Param 参数
type Param struct {
	Temp *Temporary
	Key  string
}

// IntSet 单个整型参数设置
func (p *Param) IntSet(v int64) {
	values := p.Temp.GetQuery()
	values.Set(p.Key, strconv.FormatInt(v, 10))
	p.Temp.SetQuery(values)
}

// IntAdd 单个整型参数计算
func (p *Param) IntAdd(v int64) {
	values := p.Temp.GetQuery()
	vs := values[p.Key]
	pvalue, err := strconv.ParseInt(vs[0], 10, 64)
	if err != nil {
		panic(err)
	}
	pvalue += v
	values.Set(p.Key, strconv.FormatInt(pvalue, 10))
	p.Temp.SetQuery(values)
}

// IntArraySet 数组整型参数计算
func (p *Param) IntArraySet(index int, v int64) {
	values := p.Temp.GetQuery()
	vs := values[p.Key]
	vs[index] = strconv.FormatInt(v, 10)
	p.Temp.SetQuery(values)
}

// IntArrayAdd 数组整型参数计算
func (p *Param) IntArrayAdd(index int, v int64) {
	values := p.Temp.GetQuery()
	vs := values[p.Key]
	pvalue, err := strconv.ParseInt(vs[index], 10, 64)
	if err != nil {
		panic(err)
	}
	pvalue += v
	vs[index] = strconv.FormatInt(pvalue, 10)
	p.Temp.SetQuery(values)
}

// IntArrayDo 数组整型参数操作 do i 数组索引 pvalue 数组值 返回值interface{} 如果nil. 则不变
func (p *Param) IntArrayDo(do func(i int, pvalue int64) interface{}) {
	values := p.Temp.GetQuery()
	vs := values[p.Key]
	for i, v := range vs {
		pvalue, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			panic(err)
		}
		rvalue := do(i, pvalue)
		if rvalue != nil {
			vs[i] = strconv.FormatInt(rvalue.(int64), 10)
		}
	}
	p.Temp.SetQuery(values)
}

// FloatSet 单个浮点参数设置
func (p *Param) FloatSet(v float64) {
	values := p.Temp.GetQuery()
	values.Set(p.Key, strconv.FormatFloat(v, 'f', -1, 64))
	p.Temp.SetQuery(values)
}

// FloatAdd 单个浮点参数计算
func (p *Param) FloatAdd(v float64) {
	values := p.Temp.GetQuery()
	vs := values[p.Key]
	pvalue, err := strconv.ParseFloat(vs[0], 64)
	if err != nil {
		panic(err)
	}
	pvalue += v
	values.Set(p.Key, strconv.FormatFloat(pvalue, 'f', -1, 64))
	p.Temp.SetQuery(values)
}

// FloatArraySet 数组浮点参数设置
func (p *Param) FloatArraySet(index int, v float64) {
	values := p.Temp.GetQuery()
	vs := values[p.Key]
	vs[index] = strconv.FormatFloat(v, 'f', -1, 64)
	p.Temp.SetQuery(values)
}

// FloatArrayAdd 数组浮点参数计算
func (p *Param) FloatArrayAdd(index int, v float64) {
	values := p.Temp.GetQuery()
	vs := values[p.Key]
	pvalue, err := strconv.ParseFloat(vs[index], 64)
	if err != nil {
		panic(err)
	}
	pvalue += v
	vs[index] = strconv.FormatFloat(pvalue, 'f', -1, 64)
	p.Temp.SetQuery(values)
}

// FloatArrayDo 数组整型参数操作 do i 数组索引 pvalue 数组值 返回值interface{} 如果nil. 则不变
func (p *Param) FloatArrayDo(do func(i int, pvalue float64) interface{}) {
	values := p.Temp.GetQuery()
	vs := values[p.Key]
	for i, v := range vs {
		pvalue, err := strconv.ParseFloat(v, 64)
		if err != nil {
			panic(err)
		}
		rvalue := do(i, pvalue)
		if rvalue != nil {
			vs[i] = strconv.FormatFloat(rvalue.(float64), 'f', -1, 64)
		}
	}
	p.Temp.SetQuery(values)
}
