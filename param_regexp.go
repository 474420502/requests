package requests

import (
	"bytes"
	"log"
	"reflect"
	"regexp"
	"strconv"
)

// ParamPath 参数
type ParamRegexp struct {
	req      *Request // 直接引用Request而不是Temporary
	Key      string   // 正则的Group
	Params   []string
	Selected []int
}

func extractorParam(req *Request, regexpGroup string, extracted string) *ParamRegexp {
	pp := &ParamRegexp{req: req, Key: regexpGroup}
	// result := regexp.MustCompile(regexpGroup).FindAllStringSubmatch(extracted, 1)
	result := regexp.MustCompile(regexpGroup).FindAllStringSubmatchIndex(extracted, 1)

	if len(result) == 0 {
		log.Printf("Warning: regexp not find the matched: %s", extracted)
		// 返回一个空的ParamRegexp以避免panic
		return &ParamRegexp{req: req, Key: regexpGroup, Params: []string{}, Selected: []int{}}
	}

	//  = selected
	matched := result[0]
	var cur = 0
	for i := 2; i < len(matched); i += 2 {
		start := matched[i]
		end := matched[i+1]

		tok := extracted[cur:start]
		if cur != start {
			pp.Params = append(pp.Params, tok)
		}

		pp.Params = append(pp.Params, extracted[start:end])
		pp.Selected = append(pp.Selected, len(pp.Params)-1)

		cur = end
	}

	if cur < len(extracted) {
		pp.Params = append(pp.Params, extracted[cur:])
	}

	return pp
}

func concat(ss []string) string {
	buf := &bytes.Buffer{}
	for _, s := range ss {
		buf.WriteString(s)
	}
	return buf.String()
}

// Set 单个整型参数设置
func (p *ParamRegexp) Set(value interface{}) {

	sel := p.Selected[0]
	vv := reflect.ValueOf(value)
	switch k := vv.Kind(); k {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p.Params[sel] = strconv.FormatInt(vv.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		p.Params[sel] = strconv.FormatUint(vv.Uint(), 10)
	case reflect.Float64, reflect.Float32:
		p.Params[sel] = strconv.FormatFloat(vv.Float(), 'f', -1, 64)
	case reflect.String:
		p.Params[sel] = vv.String()
	}
	p.req.parsedURL.Path = concat(p.Params)
}

// Add   通用类型 参数加减 value 为通用计算类型
func (p *ParamRegexp) Add(value interface{}) error {

	sel := p.Selected[0]

	switch v := value.(type) {
	case int64:
		pvalue, err := strconv.ParseInt(p.Params[sel], 10, 64)
		if err != nil {
			return err
		}
		pvalue += v
		p.Params[sel] = strconv.FormatInt(pvalue, 10)
	case uint64:
		pvalue, err := strconv.ParseUint(p.Params[sel], 10, 64)
		if err != nil {
			return err
		}
		pvalue += v
		p.Params[sel] = strconv.FormatUint(pvalue, 10)
	case float64:
		pvalue, err := strconv.ParseFloat(p.Params[sel], 10)
		if err != nil {
			return err
		}
		pvalue += v
		p.Params[sel] = strconv.FormatFloat(pvalue, 'f', -1, 64)
	case int:
		pvalue, err := strconv.ParseInt(p.Params[sel], 10, 64)
		if err != nil {
			return err
		}
		pvalue += int64(v)
		p.Params[sel] = strconv.FormatInt(pvalue, 10)
	case int8:
		pvalue, err := strconv.ParseInt(p.Params[sel], 10, 64)
		if err != nil {
			return err
		}
		pvalue += int64(v)
		p.Params[sel] = strconv.FormatInt(pvalue, 10)
	case int16:
		pvalue, err := strconv.ParseInt(p.Params[sel], 10, 64)
		if err != nil {
			return err
		}
		pvalue += int64(v)
		p.Params[sel] = strconv.FormatInt(pvalue, 10)
	case int32:
		pvalue, err := strconv.ParseInt(p.Params[sel], 10, 64)
		if err != nil {
			return err
		}
		pvalue += int64(v)
		p.Params[sel] = strconv.FormatInt(pvalue, 10)
	case float32:
		pvalue, err := strconv.ParseFloat(p.Params[sel], 10)
		if err != nil {
			return err
		}
		pvalue += float64(v)
		p.Params[sel] = strconv.FormatFloat(pvalue, 'f', -1, 32)

	case uint:
		pvalue, err := strconv.ParseUint(p.Params[sel], 10, 64)
		if err != nil {
			return err
		}
		pvalue += uint64(v)
		p.Params[sel] = strconv.FormatUint(pvalue, 10)
	case uint8:
		pvalue, err := strconv.ParseUint(p.Params[sel], 10, 64)
		if err != nil {
			return err
		}
		pvalue += uint64(v)
		p.Params[sel] = strconv.FormatUint(pvalue, 10)
	case uint16:
		pvalue, err := strconv.ParseUint(p.Params[sel], 10, 64)
		if err != nil {
			return err
		}
		pvalue += uint64(v)
		p.Params[sel] = strconv.FormatUint(pvalue, 10)
	case uint32:
		pvalue, err := strconv.ParseUint(p.Params[sel], 10, 64)
		if err != nil {
			return err
		}
		pvalue += uint64(v)
		p.Params[sel] = strconv.FormatUint(pvalue, 10)
	}
	p.req.parsedURL.Path = concat(p.Params)
	return nil
}

// ArraySet 通用数组类型 根据 index 设置 value 为通用计算类型
func (p *ParamRegexp) ArraySet(index int, value interface{}) {
	sel := p.Selected[index]
	vv := reflect.ValueOf(value)
	switch k := vv.Kind(); k {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p.Params[sel] = strconv.FormatInt(vv.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		p.Params[sel] = strconv.FormatUint(vv.Uint(), 10)
	case reflect.Float64, reflect.Float32:
		p.Params[sel] = strconv.FormatFloat(vv.Float(), 'f', -1, 64)
	case reflect.String:
		p.Params[sel] = vv.String()
	}
	p.req.parsedURL.Path = concat(p.Params)
}

// ArrayAdd 通用数组类型 根据 index 参数加减 value 为通用计算类型
func (p *ParamRegexp) ArrayAdd(index int, value interface{}) error {

	sel := p.Selected[index]
	switch v := value.(type) {
	case int64:
		pvalue, err := strconv.ParseInt(p.Params[sel], 10, 64)
		if err != nil {
			return err
		}
		pvalue += v
		p.Params[sel] = strconv.FormatInt(pvalue, 10)
	case uint64:
		pvalue, err := strconv.ParseUint(p.Params[sel], 10, 64)
		if err != nil {
			return err
		}
		pvalue += v
		p.Params[sel] = strconv.FormatUint(pvalue, 10)
	case float64:
		pvalue, err := strconv.ParseFloat(p.Params[sel], 10)
		if err != nil {
			return err
		}
		pvalue += v
		p.Params[sel] = strconv.FormatFloat(pvalue, 'f', -1, 64)
	case int:
		pvalue, err := strconv.ParseInt(p.Params[sel], 10, 64)
		if err != nil {
			return err
		}
		pvalue += int64(v)
		p.Params[sel] = strconv.FormatInt(pvalue, 10)
	case int8:
		pvalue, err := strconv.ParseInt(p.Params[sel], 10, 64)
		if err != nil {
			return err
		}
		pvalue += int64(v)
		p.Params[sel] = strconv.FormatInt(pvalue, 10)
	case int16:
		pvalue, err := strconv.ParseInt(p.Params[sel], 10, 64)
		if err != nil {
			return err
		}
		pvalue += int64(v)
		p.Params[sel] = strconv.FormatInt(pvalue, 10)
	case int32:
		pvalue, err := strconv.ParseInt(p.Params[sel], 10, 64)
		if err != nil {
			return err
		}
		pvalue += int64(v)
		p.Params[sel] = strconv.FormatInt(pvalue, 10)
	case float32:
		pvalue, err := strconv.ParseFloat(p.Params[sel], 10)
		if err != nil {
			return err
		}
		pvalue += float64(v)
		p.Params[sel] = strconv.FormatFloat(pvalue, 'f', -1, 32)

	case uint:
		pvalue, err := strconv.ParseUint(p.Params[sel], 10, 64)
		if err != nil {
			return err
		}
		pvalue += uint64(v)
		p.Params[sel] = strconv.FormatUint(pvalue, 10)
	case uint8:
		pvalue, err := strconv.ParseUint(p.Params[sel], 10, 64)
		if err != nil {
			return err
		}
		pvalue += uint64(v)
		p.Params[sel] = strconv.FormatUint(pvalue, 10)
	case uint16:
		pvalue, err := strconv.ParseUint(p.Params[sel], 10, 64)
		if err != nil {
			return err
		}
		pvalue += uint64(v)
		p.Params[sel] = strconv.FormatUint(pvalue, 10)
	case uint32:
		pvalue, err := strconv.ParseUint(p.Params[sel], 10, 64)
		if err != nil {
			return err
		}
		pvalue += uint64(v)
		p.Params[sel] = strconv.FormatUint(pvalue, 10)
	}
	p.req.parsedURL.Path = concat(p.Params)
	return nil
}

// IntSet 单个整型参数设置
func (p *ParamRegexp) IntSet(v int64) {
	sel := p.Selected[0]
	p.Params[sel] = strconv.FormatInt(v, 10)
	p.req.parsedURL.Path = concat(p.Params)
}

// IntAdd 单个整型参数计算
func (p *ParamRegexp) IntAdd(v int64) error {
	sel := p.Selected[0]
	pvalue, err := strconv.ParseInt(p.Params[sel], 10, 64)
	if err != nil {
		return err
	}
	pvalue += v
	p.Params[sel] = strconv.FormatInt(pvalue, 10)
	p.req.parsedURL.Path = concat(p.Params)
	return nil
}

// IntArraySet 数组整型参数计算
func (p *ParamRegexp) IntArraySet(index int, v int64) {
	sel := p.Selected[index]
	p.Params[sel] = strconv.FormatInt(v, 10)
	p.req.parsedURL.Path = concat(p.Params)
}

// IntArrayAdd 数组整型参数计算
func (p *ParamRegexp) IntArrayAdd(index int, v int64) error {
	sel := p.Selected[index]
	pvalue, err := strconv.ParseInt(p.Params[sel], 10, 64)
	if err != nil {
		return err
	}
	pvalue += v
	p.Params[sel] = strconv.FormatInt(pvalue, 10)
	p.req.parsedURL.Path = concat(p.Params)
	return nil
}

// IntArrayDo 数组整型参数操作 do i 数组索引 pvalue 数组值 返回值interface{} 如果nil. 则不变
func (p *ParamRegexp) IntArrayDo(do func(i int, pvalue int64) interface{}) error {
	for i, index := range p.Selected {
		v := p.Params[index]
		pvalue, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return err
		}
		rvalue := do(i, pvalue)
		if rvalue != nil {
			if err, ok := rvalue.(error); ok {
				return err
			}
			p.Params[index] = strconv.FormatInt(rvalue.(int64), 10)
		}
	}
	p.req.parsedURL.Path = concat(p.Params)
	return nil
}

// FloatSet 单个浮点参数设置
func (p *ParamRegexp) FloatSet(v float64) {
	sel := p.Selected[0]
	p.Params[sel] = strconv.FormatFloat(v, 'f', -1, 64)
	p.req.parsedURL.Path = concat(p.Params)
}

// FloatAdd 单个浮点参数计算
func (p *ParamRegexp) FloatAdd(v float64) error {
	sel := p.Selected[0]
	pvalue, err := strconv.ParseFloat(p.Params[sel], 64)
	if err != nil {
		return err
	}
	pvalue += v
	p.Params[sel] = strconv.FormatFloat(pvalue, 'f', -1, 64)
	p.req.parsedURL.Path = concat(p.Params)
	return nil
}

// FloatArraySet 数组浮点参数设置
func (p *ParamRegexp) FloatArraySet(index int, v float64) {
	sel := p.Selected[index]
	p.Params[sel] = strconv.FormatFloat(v, 'f', -1, 64)
	p.req.parsedURL.Path = concat(p.Params)
}

// FloatArrayAdd 数组浮点参数计算
func (p *ParamRegexp) FloatArrayAdd(index int, v float64) error {
	sel := p.Selected[index]
	pvalue, err := strconv.ParseFloat(p.Params[sel], 64)
	if err != nil {
		return err
	}
	pvalue += v
	p.Params[sel] = strconv.FormatFloat(pvalue, 'f', -1, 64)
	p.req.parsedURL.Path = concat(p.Params)
	return nil
}

// FloatArrayDo 数组整型参数操作 do i 数组索引 pvalue 数组值 返回值interface{} 如果nil. 则不变
func (p *ParamRegexp) FloatArrayDo(do func(i int, pvalue float64) interface{}) error {
	for i, index := range p.Selected {
		v := p.Params[index]
		pvalue, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return err
		}
		rvalue := do(i, pvalue)
		if rvalue != nil {
			if err, ok := rvalue.(error); ok {
				return err
			}
			p.Params[index] = strconv.FormatFloat(rvalue.(float64), 'f', -1, 64)
		}
	}
	p.req.parsedURL.Path = concat(p.Params)
	return nil
}

// StringSet 字符串参数设置
func (p *ParamRegexp) StringSet(v string) {
	sel := p.Selected[0]
	p.Params[sel] = v
	p.req.parsedURL.Path = concat(p.Params)
}

// StringArraySet 数组字符串参数设置
func (p *ParamRegexp) StringArraySet(index int, v string) {
	sel := p.Selected[index]
	p.Params[sel] = v
	p.req.parsedURL.Path = concat(p.Params)
}
