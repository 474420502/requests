package requests

import (
	"bytes"
	"strconv"
)

// ParamPath 参数
type ParamPath struct {
	Temp     *Temporary
	Key      string // 正则的Group
	Params   []string
	Selected []int
}

func concat(ss []string) string {
	buf := &bytes.Buffer{}
	for _, s := range ss {
		buf.WriteString(s)
	}
	return buf.String()
}

// IntSet 单个整型参数设置
func (p *ParamPath) IntSet(v int64) error {
	sel := p.Selected[0]
	p.Params[sel] = strconv.FormatInt(v, 10)
	p.Temp.ParsedURL.Path = concat(p.Params)
	return nil
}

// IntAdd 单个整型参数计算
func (p *ParamPath) IntAdd(v int64) error {
	sel := p.Selected[0]
	pvalue, err := strconv.ParseInt(p.Params[sel], 10, 64)
	if err != nil {
		return err
	}
	pvalue += v
	p.Params[sel] = strconv.FormatInt(pvalue, 10)
	p.Temp.ParsedURL.Path = concat(p.Params)
	return nil
}

// IntArraySet 数组整型参数计算
func (p *ParamPath) IntArraySet(index int, v int64) {
	sel := p.Selected[index]
	p.Params[sel] = strconv.FormatInt(v, 10)
	p.Temp.ParsedURL.Path = concat(p.Params)
}

// IntArrayAdd 数组整型参数计算
func (p *ParamPath) IntArrayAdd(index int, v int64) error {
	sel := p.Selected[index]
	pvalue, err := strconv.ParseInt(p.Params[sel], 10, 64)
	if err != nil {
		return err
	}
	pvalue += v
	p.Params[sel] = strconv.FormatInt(pvalue, 10)
	p.Temp.ParsedURL.Path = concat(p.Params)
	return nil
}

// IntArrayDo 数组整型参数操作 do i 数组索引 pvalue 数组值 返回值interface{} 如果nil. 则不变
func (p *ParamPath) IntArrayDo(do func(i int, pvalue int64) interface{}) error {
	for i, v := range p.Params {
		pvalue, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return err
		}
		rvalue := do(i, pvalue)
		if rvalue != nil {
			if err, ok := rvalue.(error); ok {
				return err
			}
			p.Params[i] = strconv.FormatInt(rvalue.(int64), 10)
		}
	}
	p.Temp.ParsedURL.Path = concat(p.Params)
	return nil
}

// FloatSet 单个浮点参数设置
func (p *ParamPath) FloatSet(v float64) {
	sel := p.Selected[0]
	p.Params[sel] = strconv.FormatFloat(v, 'f', -1, 64)
	p.Temp.ParsedURL.Path = concat(p.Params)
}

// FloatAdd 单个浮点参数计算
func (p *ParamPath) FloatAdd(v float64) error {
	sel := p.Selected[0]
	pvalue, err := strconv.ParseFloat(p.Params[sel], 64)
	if err != nil {
		return err
	}
	pvalue += v
	p.Params[sel] = strconv.FormatFloat(pvalue, 'f', -1, 64)
	p.Temp.ParsedURL.Path = concat(p.Params)
	return nil
}

// FloatArraySet 数组浮点参数设置
func (p *ParamPath) FloatArraySet(index int, v float64) {
	sel := p.Selected[index]
	p.Params[sel] = strconv.FormatFloat(v, 'f', -1, 64)
	p.Temp.ParsedURL.Path = concat(p.Params)
}

// FloatArrayAdd 数组浮点参数计算
func (p *ParamPath) FloatArrayAdd(index int, v float64) error {
	sel := p.Selected[index]
	pvalue, err := strconv.ParseFloat(p.Params[sel], 64)
	if err != nil {
		return err
	}
	pvalue += v
	p.Params[sel] = strconv.FormatFloat(pvalue, 'f', -1, 64)
	p.Temp.ParsedURL.Path = concat(p.Params)
	return nil
}

// FloatArrayDo 数组整型参数操作 do i 数组索引 pvalue 数组值 返回值interface{} 如果nil. 则不变
func (p *ParamPath) FloatArrayDo(do func(i int, pvalue float64) interface{}) error {
	for i, v := range p.Params {
		pvalue, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return err
		}
		rvalue := do(i, pvalue)
		if rvalue != nil {
			if err, ok := rvalue.(error); ok {
				return err
			}
			p.Params[i] = strconv.FormatFloat(rvalue.(float64), 'f', -1, 64)
		}
	}
	p.Temp.ParsedURL.Path = concat(p.Params)
	return nil
}

// StringSet 字符串参数设置
func (p *ParamPath) StringSet(v string) {
	sel := p.Selected[0]
	p.Params[sel] = v
	p.Temp.ParsedURL.Path = concat(p.Params)
}

// StringArraySet 数组字符串参数设置
func (p *ParamPath) StringArraySet(index int, v string) {
	sel := p.Selected[index]
	p.Params[sel] = v
	p.Temp.ParsedURL.Path = concat(p.Params)
}
