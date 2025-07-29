package requests

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"testing"

	"github.com/474420502/random"
)

func checkArrayParam(tp *Request, param string, vaild string) error {

	resp, err := tp.Execute()
	if err != nil {
		return err
	}

	data := make(map[string]interface{})
	err = json.Unmarshal(resp.Content(), &data)
	if err != nil {
		log.Println(data, string(resp.Content()))
		return err
	}

	if args, ok := data["args"]; ok {
		if page, ok := args.(map[string]interface{})[param]; ok {
			if fmt.Sprint(page)[0:len(vaild)] != vaild {
				// log.Println(data, string(resp.Content()))
				return fmt.Errorf("param: %#v", fmt.Sprint(page))
			}
		} else {
			log.Println(data, string(resp.Content()))
			return fmt.Errorf("param is %s not exists", param)
		}
	} else {
		log.Println(data, string(resp.Content()))
		return fmt.Errorf("args is not exists")
	}

	return nil
}

func checkParam(tp *Request, param string, vaild string) error {

	resp, err := tp.Execute()
	if err != nil {
		return err
	}

	data := make(map[string]interface{})
	err = json.Unmarshal(resp.Content(), &data)
	if err != nil {
		log.Println(data)
		return err
	}

	if args, ok := data["args"]; ok {
		if page, ok := args.(map[string]interface{})[param]; ok {
			if len(page.(string)) < len(vaild) {
				fmt.Errorf("")
			}

			if page.(string)[0:len(vaild)] != vaild {
				// log.Println(data)
				return fmt.Errorf("param(%s): %v != %v", param, page, vaild)
			}
		} else {
			// log.Println(data)
			return fmt.Errorf("param is %s not exists", param)
		}
	} else {
		// log.Println(data)
		return fmt.Errorf("args is not exists")
	}

	return nil
}

func TestQueryParam(t *testing.T) {

	ses := NewSession()
	tp := ses.Get("http://httpbin.org/get?page=1&name=xiaoming")

	err := checkParam(tp, "page", "1")
	if err != nil {
		t.Error(err)
	}

	p := tp.QueryParam("page")
	p.IntAdd(1)
	err = checkParam(tp, "page", "2")
	if err != nil {
		t.Error(err)
	}

	p = tp.QueryParam("page")
	p.FloatAdd(1)
	err = checkParam(tp, "page", "3")
	if err != nil {
		t.Error(err)
	}

	p = tp.QueryParam("page")
	p.IntSet(1)
	err = checkParam(tp, "page", "1")
	if err != nil {
		t.Error(err)
	}

	p = tp.QueryParam("page")
	p.FloatSet(1.5)
	err = checkParam(tp, "page", "1.5")
	if err != nil {
		t.Error(err)
	}

	p.StringSet("5")
	err = checkParam(tp, "page", "5")
	if err != nil {
		t.Error(err)
	}
}

func TestQueryArrayParam(t *testing.T) {
	var err error
	var p IParam
	ses := NewSession()
	tp := ses.Get("http://httpbin.org/get?page[]=1&page[]=2&page[]=3&name=xiaoming")
	p = tp.QueryParam("page[]")
	p.IntArrayAdd(0, 2)
	err = checkArrayParam(tp, "page[]", "[3 2 3]")
	if err != nil {
		t.Error(err)
	}

	p.IntArraySet(2, 2)
	err = checkArrayParam(tp, "page[]", "[3 2 2]")
	if err != nil {
		t.Error(err)
	}

	p.IntArrayDo(func(i int, pvalue int64) interface{} {
		if i%2 == 0 {
			pvalue++
			return pvalue
		}
		return nil
	})
	err = checkArrayParam(tp, "page[]", "[4 2 3]")
	if err != nil {
		t.Error(err)
	}

	p.FloatArrayDo(func(i int, pvalue float64) interface{} {
		if i%2 == 0 {
			pvalue += 1.5
			return pvalue
		}
		return nil
	})
	err = checkArrayParam(tp, "page[]", "[5.5 2 4.5]")
	if err != nil {
		t.Error(err)
	}

	p.FloatArrayAdd(1, 1.5)
	err = checkArrayParam(tp, "page[]", "[5.5 3.5 4.5]")
	if err != nil {
		t.Error(err)
	}

	p.FloatArraySet(0, 3.141592653589793)
	err = checkArrayParam(tp, "page[]", "[3.141592653589793 3.5 4.5]")
	if err != nil {
		t.Error(err)
	}

}

func TestParamPath(t *testing.T) {

	ses := NewSession()
	surl := "https://api.xxx.tv/oversea/xxx/api/v2/liveRoom/Page-1-30-/HK/1028/1000"
	tp := ses.Get(surl)
	param := tp.PathParam(`.+Page-(\d+)-(\d+).+`)

	param.IntAdd(1)
	purl := tp.GetURLRawPath()
	if !regexp.MustCompile("Page-2-30").MatchString(purl) {
		t.Error(purl)
	}

	param.IntArrayAdd(1, 30)
	purl = tp.GetURLRawPath()
	if !regexp.MustCompile("Page-2-60").MatchString(purl) {
		t.Error(purl)
	}

	param.IntArraySet(0, 4)
	purl = tp.GetURLRawPath()
	if !regexp.MustCompile("Page-4-60").MatchString(purl) {
		t.Error(purl)
	}

	param.IntSet(8)
	purl = tp.GetURLRawPath()
	if !regexp.MustCompile("Page-8-60").MatchString(purl) {
		t.Error(purl)
	}

	param.IntArrayDo(func(i int, pvalue int64) interface{} {
		if i == 1 {
			pvalue += 20
		}
		return pvalue
	})
	purl = tp.GetURLRawPath()
	if !regexp.MustCompile("Page-8-80").MatchString(purl) {
		t.Error(purl)
	}

	param.FloatAdd(2)
	purl = tp.GetURLRawPath()
	if !regexp.MustCompile("Page-10-80").MatchString(purl) {
		t.Error(purl)
	}

	param.FloatSet(4.5)
	purl = tp.GetURLRawPath()
	if !regexp.MustCompile("Page-4.5-80").MatchString(purl) {
		t.Error(purl)
	}

	param.FloatArrayAdd(0, 0.5)
	purl = tp.GetURLRawPath()
	if !regexp.MustCompile("Page-5-80").MatchString(purl) {
		t.Error(purl)
	}

	param.FloatArrayDo(func(i int, pvalue float64) interface{} {
		if i == 0 {
			pvalue += 2
		}
		return pvalue
	})
	purl = tp.GetURLRawPath()
	if !regexp.MustCompile("Page-7-80").MatchString(purl) {
		t.Error(purl)
	}

	param.StringSet("9")
	purl = tp.GetURLRawPath()
	if !regexp.MustCompile("Page-9-80").MatchString(purl) {
		t.Error(purl)
	}

	param.StringArraySet(1, "123")
	purl = tp.GetURLRawPath()
	if !regexp.MustCompile("Page-9-123").MatchString(purl) {
		t.Error(purl)
	}
}

func TestParamHost(t *testing.T) {
	ses := NewSession()
	surl := "https://10.api.xxx.tv/oversea/xxx/api/v2/liveRoom/Page-1-30-/HK/1028/1000"
	tp := ses.Get(surl)
	param := tp.HostParam(`(\d+).api.xx.+`)
	param.IntAdd(1)
	purl := tp.GetURLRawPath()
	if !regexp.MustCompile("11.api").MatchString(purl) {
		t.Error(purl)
	}
}

func checkBaseTypeParamSet(tp *Request, r *random.Random, t *testing.T) {
	p := tp.QueryParam("page")
	var v interface{}
	var err error

	v = r.Int63()
	p.Set(v)

	err = checkParam(tp, "page", fmt.Sprintf("%v", v))
	if err != nil {
		t.Error(err)
	}

	v = r.Int31()
	p.Set(v)

	err = checkParam(tp, "page", fmt.Sprintf("%v", v))
	if err != nil {
		t.Error(err)
	}

	v = r.Int()
	p.Set(v)

	err = checkParam(tp, "page", fmt.Sprintf("%v", v))
	if err != nil {
		t.Error(err)
	}

	v = (int16)(r.Int())
	p.Set(v)

	err = checkParam(tp, "page", fmt.Sprintf("%v", v))
	if err != nil {
		t.Error(err)
	}

	v = (int8)(r.Int())
	p.Set(v)

	err = checkParam(tp, "page", fmt.Sprintf("%v", v))
	if err != nil {
		t.Error(err)
	}

	v = r.Uint64()
	p.Set(v)

	err = checkParam(tp, "page", fmt.Sprintf("%v", v))
	if err != nil {
		t.Error(err)
	}

	v = r.Uint32()
	p.Set(v)

	err = checkParam(tp, "page", fmt.Sprintf("%v", v))
	if err != nil {
		t.Error(err)
	}

	v = uint(r.Uint64())
	p.Set(v)

	err = checkParam(tp, "page", fmt.Sprintf("%v", v))
	if err != nil {
		t.Error(err)
	}

	v = (uint16)(r.Int())
	p.Set(v)

	err = checkParam(tp, "page", fmt.Sprintf("%v", v))
	if err != nil {
		t.Error(err)
	}

	v = (uint8)(r.Int())
	p.Set(v)

	err = checkParam(tp, "page", fmt.Sprintf("%v", v))
	if err != nil {
		t.Error(err)
	}

	p = tp.QueryParam("float")
	v = r.Float64()

	p.Set(v.(float64))
	err = checkParam(tp, "float", fmt.Sprintf("%v", v.(float64)))
	if err != nil {
		t.Error(err)
	}

	v = r.Float32()
	p.Set(v.(float32))
	checkParam(tp, "float", fmt.Sprintf("%v", float64(v.(float32))))

}

func checkBaseTypeParamRegexpSet(tp *Request, r *random.Random, t *testing.T) {
	p := tp.PathParam(`Page-([0-9.]+)-([0-9.]+)`)
	var v interface{}

	var purl string

	v = r.Int63()
	p.Set(v)

	purl = tp.GetURLRawPath()
	if !regexp.MustCompile(fmt.Sprintf("Page-%v", v)).MatchString(purl) {
		t.Error(purl)
	}

	v = r.Int31()
	p.Set(v)

	purl = tp.GetURLRawPath()
	if !regexp.MustCompile(fmt.Sprintf("Page-%v", v)).MatchString(purl) {
		t.Error(purl)
	}

	v = r.Int()
	p.Set(v)

	purl = tp.GetURLRawPath()
	if !regexp.MustCompile(fmt.Sprintf("Page-%v", v)).MatchString(purl) {
		t.Error(purl)
	}

	v = (int16)(r.Int())
	p.Set(v)

	purl = tp.GetURLRawPath()
	if !regexp.MustCompile(fmt.Sprintf("Page-%v", v)).MatchString(purl) {
		t.Error(purl)
	}

	v = (int8)(r.Int())
	p.Set(v)

	purl = tp.GetURLRawPath()
	if !regexp.MustCompile(fmt.Sprintf("Page-%v", v)).MatchString(purl) {
		t.Error(purl)
	}

	v = r.Uint64()
	p.Set(v)

	purl = tp.GetURLRawPath()
	if !regexp.MustCompile(fmt.Sprintf("Page-%v", v)).MatchString(purl) {
		t.Error(purl)
	}

	v = r.Uint32()
	p.Set(v)

	purl = tp.GetURLRawPath()
	if !regexp.MustCompile(fmt.Sprintf("Page-%v", v)).MatchString(purl) {
		t.Error(purl)
	}

	v = uint(r.Uint64())
	p.Set(v)

	purl = tp.GetURLRawPath()
	if !regexp.MustCompile(fmt.Sprintf("Page-%v", v)).MatchString(purl) {
		t.Error(purl)
	}

	v = (uint16)(r.Int())
	p.Set(v)

	purl = tp.GetURLRawPath()
	if !regexp.MustCompile(fmt.Sprintf("Page-%v", v)).MatchString(purl) {
		t.Error(purl)
	}

	v = (uint8)(r.Int())
	p.Set(v)

	purl = tp.GetURLRawPath()
	if !regexp.MustCompile(fmt.Sprintf("Page-%v", v)).MatchString(purl) {
		t.Error(purl)
	}

	v = r.Float64()
	p.Set(v.(float64))
	purl = tp.GetURLRawPath()
	if !regexp.MustCompile(fmt.Sprintf("Page-%v", v)).MatchString(purl) {
		t.Error(purl, ",", fmt.Sprintf("Page-%v", v))

	}

	v = r.Float32()
	p.Set(v.(float32))
	purl = tp.GetURLRawPath()
	if !regexp.MustCompile(fmt.Sprintf("Page-%v", float64(v.(float32)))).MatchString(purl) {
		t.Error(purl, ",", fmt.Sprintf("Page-%v", v))
	}

}

func checkBaseTypeParamRegexpArraySet(tp *Request, r *random.Random, t *testing.T) {
	p := tp.PathParam(`Page-([0-9.]+)-([0-9.]+)`)
	var v interface{}

	var purl string

	v = r.Int63()
	p.ArraySet(1, v)

	purl = tp.GetURLRawPath()
	if !regexp.MustCompile(fmt.Sprintf(`Page-([0-9\.]+)-%v`, v)).MatchString(purl) {
		t.Error(purl)
	}

	v = r.Int31()
	p.ArraySet(1, v)

	purl = tp.GetURLRawPath()
	if !regexp.MustCompile(fmt.Sprintf(`Page-([0-9\.]+)-%v`, v)).MatchString(purl) {
		t.Error(purl)
	}

	v = r.Int()
	p.ArraySet(1, v)

	purl = tp.GetURLRawPath()
	if !regexp.MustCompile(fmt.Sprintf(`Page-([0-9\.]+)-%v`, v)).MatchString(purl) {
		t.Error(purl)
	}

	v = (int16)(r.Int())
	p.ArraySet(1, v)

	purl = tp.GetURLRawPath()
	if !regexp.MustCompile(fmt.Sprintf(`Page-([0-9\.]+)-%v`, v)).MatchString(purl) {
		t.Error(purl)
	}

	v = (int8)(r.Int())
	p.ArraySet(1, v)

	purl = tp.GetURLRawPath()
	if !regexp.MustCompile(fmt.Sprintf(`Page-([0-9\.]+)-%v`, v)).MatchString(purl) {
		t.Error(purl)
	}

	v = r.Uint64()
	p.ArraySet(1, v)

	purl = tp.GetURLRawPath()
	if !regexp.MustCompile(fmt.Sprintf(`Page-([0-9\.]+)-%v`, v)).MatchString(purl) {
		t.Error(purl)
	}

	v = r.Uint32()
	p.ArraySet(1, v)

	purl = tp.GetURLRawPath()
	if !regexp.MustCompile(fmt.Sprintf(`Page-([0-9\.]+)-%v`, v)).MatchString(purl) {
		t.Error(purl)
	}

	v = uint(r.Uint64())
	p.ArraySet(1, v)

	purl = tp.GetURLRawPath()
	if !regexp.MustCompile(fmt.Sprintf(`Page-([0-9\.]+)-%v`, v)).MatchString(purl) {
		t.Error(purl)
	}

	v = (uint16)(r.Int())
	p.ArraySet(1, v)

	purl = tp.GetURLRawPath()
	if !regexp.MustCompile(fmt.Sprintf(`Page-([0-9\.]+)-%v`, v)).MatchString(purl) {
		t.Error(purl)
	}

	v = (uint8)(r.Int())
	p.ArraySet(1, v)

	purl = tp.GetURLRawPath()
	if !regexp.MustCompile(fmt.Sprintf(`Page-([0-9\.]+)-%v`, v)).MatchString(purl) {
		t.Error(purl)
	}

	v = r.Float64()
	p.ArraySet(1, v.(float64))
	purl = tp.GetURLRawPath()
	if !regexp.MustCompile(fmt.Sprintf(`Page-([0-9\.]+)-%v`, v)).MatchString(purl) {
		t.Error(fmt.Sprintf(`Page-([0-9\.]+)-%v`, v), ",", purl)
	}

	v = r.Float32()
	p.ArraySet(1, v.(float32))
	purl = tp.GetURLRawPath()
	if !regexp.MustCompile(fmt.Sprintf(`Page-([0-9\.]+)-%v`, float64(v.(float32)))).MatchString(purl) {
		t.Error(fmt.Sprintf(`Page-([0-9\.]+)-%v`, v), ",", purl)
	}

}

func checkBaseTypeParamAdd(tp *Request, r *random.Random, t *testing.T) {
	p := tp.QueryParam("page")
	var v interface{}
	var err error
	var defaultint64 int64 = 1
	var defaultint32 int32 = 1
	var defaultint int = 1
	var defaultint16 int16 = 1
	var defaultint8 int8 = 1

	v = r.Int63() >> 2
	p.Set(defaultint64)
	p.Add(v)

	err = checkParam(tp, "page", fmt.Sprintf("%v", defaultint64+v.(int64)))
	if err != nil {
		t.Error(err)
	}

	v = r.Int31() >> 2
	p.Set(defaultint32)
	p.Add(v)

	err = checkParam(tp, "page", fmt.Sprintf("%v", defaultint32+v.(int32)))
	if err != nil {
		t.Error(err)
	}

	v = r.Int() >> 2
	p.Set(defaultint)
	p.Add(v)

	err = checkParam(tp, "page", fmt.Sprintf("%v", defaultint+v.(int)))
	if err != nil {
		t.Error(err)
	}

	v = (int16)(r.Int()) >> 2
	p.Set(defaultint16)
	p.Add(v)

	err = checkParam(tp, "page", fmt.Sprintf("%v", defaultint16+v.(int16)))
	if err != nil {
		t.Error(err)
	}

	v = (int8)(r.Int()) >> 2
	p.Set(defaultint8)
	p.Add(v)

	err = checkParam(tp, "page", fmt.Sprintf("%v", defaultint8+v.(int8)))
	if err != nil {
		t.Error(err)
	}

	var defaultuint64 uint64 = 1
	var defaultuint32 uint32 = 1
	var defaultuint uint = 1
	var defaultuint16 uint16 = 1
	var defaultuint8 uint8 = 1

	v = r.Uint64() >> 2
	p.Set(defaultuint64)
	p.Add(v)
	err = checkParam(tp, "page", fmt.Sprintf("%v", defaultuint64+v.(uint64)))
	if err != nil {
		t.Error(err)
	}

	v = r.Uint32() >> 2
	p.Set(defaultuint32)
	p.Add(v)
	err = checkParam(tp, "page", fmt.Sprintf("%v", defaultuint32+v.(uint32)))
	if err != nil {
		t.Error(err)
	}

	v = uint(r.Uint64()) >> 2
	p.Set(defaultuint)
	p.Add(v)
	err = checkParam(tp, "page", fmt.Sprintf("%v", defaultuint+v.(uint)))
	if err != nil {
		t.Error(err)
	}

	v = (uint16)(r.Int()) >> 2
	p.Set(defaultuint16)
	p.Add(v)
	err = checkParam(tp, "page", fmt.Sprintf("%v", defaultuint16+v.(uint16)))
	if err != nil {
		t.Error(err)
	}

	v = (uint8)(r.Int()) >> 2
	p.Set(defaultuint8)
	p.Add(v)
	err = checkParam(tp, "page", fmt.Sprintf("%v", defaultuint8+v.(uint8)))
	if err != nil {
		t.Error(err)
	}

	var defaultfloat64 float64 = 1
	var defaultfloat32 float32 = 1

	p = tp.QueryParam("float")
	v = r.Float64()

	p.Set(defaultfloat64)
	p.Add(v)
	err = checkParam(tp, "float", fmt.Sprintf("%v", defaultfloat64+v.(float64)))
	if err != nil {
		t.Error(err)
	}

	v = r.Float32()
	p.Set(defaultfloat32)
	p.Add(v)

	checkParam(tp, "float", fmt.Sprintf("%v", float32(int((1.0+v.(float32))*10000))/10000))

}

func checkBaseTypeParamRegexpAdd(tp *Request, r *random.Random, t *testing.T) {
	p := tp.PathParam(`Page-([0-9.]+)-([0-9.]+)`)

	var v interface{}

	var purl string
	var defaultint64 int64 = 1
	var defaultint32 int32 = 1
	var defaultint int = 1
	var defaultint16 int16 = 1
	var defaultint8 int8 = 1

	v = r.Int63() >> 2
	p.Set(defaultint64)
	p.Add(v)

	purl = tp.GetURLRawPath()
	if !regexp.MustCompile(fmt.Sprintf("Page-%v", defaultint64+v.(int64))).MatchString(purl) {
		t.Error(purl)
	}

	v = r.Int31() >> 2
	p.Set(defaultint32)
	p.Add(v)

	purl = tp.GetURLRawPath()
	if !regexp.MustCompile(fmt.Sprintf("Page-%v", defaultint32+v.(int32))).MatchString(purl) {
		t.Error(purl)
	}

	v = r.Int() >> 2
	p.Set(defaultint)
	p.Add(v)

	purl = tp.GetURLRawPath()
	if !regexp.MustCompile(fmt.Sprintf("Page-%v", defaultint+v.(int))).MatchString(purl) {
		t.Error(purl)
	}

	v = (int16)(r.Int()) >> 2
	p.Set(defaultint16)
	p.Add(v)

	purl = tp.GetURLRawPath()
	if !regexp.MustCompile(fmt.Sprintf("Page-%v", defaultint16+v.(int16))).MatchString(purl) {
		t.Error(purl)
	}

	v = (int8)(r.Int()) >> 2
	p.Set(defaultint8)
	p.Add(v)

	purl = tp.GetURLRawPath()
	if !regexp.MustCompile(fmt.Sprintf("Page-%v", defaultint8+v.(int8))).MatchString(purl) {
		t.Error(purl)
	}

	var defaultuint64 uint64 = 1
	var defaultuint32 uint32 = 1
	var defaultuint uint = 1
	var defaultuint16 uint16 = 1
	var defaultuint8 uint8 = 1

	v = r.Uint64() >> 2
	p.Set(defaultuint64)
	p.Add(v)
	purl = tp.GetURLRawPath()
	if !regexp.MustCompile(fmt.Sprintf("Page-%v", defaultuint64+v.(uint64))).MatchString(purl) {
		t.Error(purl)
	}

	v = r.Uint32() >> 2
	p.Set(defaultuint32)
	p.Add(v)
	purl = tp.GetURLRawPath()
	if !regexp.MustCompile(fmt.Sprintf("Page-%v", defaultuint32+v.(uint32))).MatchString(purl) {
		t.Error(purl)
	}

	v = uint(r.Uint64()) >> 2
	p.Set(defaultuint)
	p.Add(v)
	purl = tp.GetURLRawPath()
	if !regexp.MustCompile(fmt.Sprintf("Page-%v", defaultuint+v.(uint))).MatchString(purl) {
		t.Error(purl)
	}

	v = (uint16)(r.Int()) >> 2
	p.Set(defaultuint16)
	p.Add(v)
	purl = tp.GetURLRawPath()
	if !regexp.MustCompile(fmt.Sprintf("Page-%v", defaultuint16+v.(uint16))).MatchString(purl) {
		t.Error(purl)
	}

	v = (uint8)(r.Int()) >> 2
	p.Set(defaultuint8)
	p.Add(v)
	purl = tp.GetURLRawPath()
	if !regexp.MustCompile(fmt.Sprintf("Page-%v", defaultuint8+v.(uint8))).MatchString(purl) {
		t.Error(purl)
	}

	var defaultfloat64 float64 = 1
	var defaultfloat32 float32 = 1

	v = r.Float64()

	p.Set(defaultfloat64)
	p.Add(v)
	purl = tp.GetURLRawPath()
	if !regexp.MustCompile(fmt.Sprintf("Page-%v", defaultfloat64+v.(float64))).MatchString(purl) {
		t.Error(purl)
	}

	v = r.Float32()
	p.Set(defaultfloat32)
	p.Add(v)
	purl = tp.GetURLRawPath()
	if !regexp.MustCompile(fmt.Sprintf("Page-%g", defaultfloat32+v.(float32))).MatchString(purl) {
		t.Error(purl)
	}

}

func checkBaseTypeParamRegexpArrayAdd(tp *Request, r *random.Random, t *testing.T) {
	p := tp.PathParam(`Page-([0-9.]+)-([0-9.]+)`)

	var v interface{}

	var purl string
	var defaultint64 int64 = 1
	var defaultint32 int32 = 1
	var defaultint int = 1
	var defaultint16 int16 = 1
	var defaultint8 int8 = 1

	v = r.Int63() >> 2
	p.ArraySet(1, defaultint64)
	p.ArrayAdd(1, v)

	purl = tp.GetURLRawPath()
	if !regexp.MustCompile(fmt.Sprintf(`Page-([0-9\.]+)-%v`, defaultint64+v.(int64))).MatchString(purl) {
		t.Error(purl)
	}

	v = r.Int31() >> 2
	p.ArraySet(1, defaultint32)
	p.ArrayAdd(1, v)

	purl = tp.GetURLRawPath()
	if !regexp.MustCompile(fmt.Sprintf(`Page-([0-9\.]+)-%v`, defaultint32+v.(int32))).MatchString(purl) {
		t.Error(purl)
	}

	v = r.Int() >> 2
	p.ArraySet(1, defaultint)
	p.ArrayAdd(1, v)

	purl = tp.GetURLRawPath()
	if !regexp.MustCompile(fmt.Sprintf(`Page-([0-9\.]+)-%v`, defaultint+v.(int))).MatchString(purl) {
		t.Error(purl)
	}

	v = (int16)(r.Int()) >> 2
	p.ArraySet(1, defaultint16)
	p.ArrayAdd(1, v)

	purl = tp.GetURLRawPath()
	if !regexp.MustCompile(fmt.Sprintf(`Page-([0-9\.]+)-%v`, defaultint16+v.(int16))).MatchString(purl) {
		t.Error(purl)
	}

	v = (int8)(r.Int()) >> 2
	p.ArraySet(1, defaultint8)
	p.ArrayAdd(1, v)

	purl = tp.GetURLRawPath()
	if !regexp.MustCompile(fmt.Sprintf(`Page-([0-9\.]+)-%v`, defaultint8+v.(int8))).MatchString(purl) {
		t.Error(purl)
	}

	var defaultuint64 uint64 = 1
	var defaultuint32 uint32 = 1
	var defaultuint uint = 1
	var defaultuint16 uint16 = 1
	var defaultuint8 uint8 = 1

	v = r.Uint64() >> 2
	p.ArraySet(1, defaultuint64)
	p.ArrayAdd(1, v)
	purl = tp.GetURLRawPath()
	if !regexp.MustCompile(fmt.Sprintf(`Page-([0-9\.]+)-%v`, defaultuint64+v.(uint64))).MatchString(purl) {
		t.Error(purl)
	}

	v = r.Uint32() >> 2
	p.ArraySet(1, defaultuint32)
	p.ArrayAdd(1, v)
	purl = tp.GetURLRawPath()
	if !regexp.MustCompile(fmt.Sprintf(`Page-([0-9\.]+)-%v`, defaultuint32+v.(uint32))).MatchString(purl) {
		t.Error(purl)
	}

	v = uint(r.Uint64()) >> 2
	p.ArraySet(1, defaultuint)
	p.ArrayAdd(1, v)
	purl = tp.GetURLRawPath()
	if !regexp.MustCompile(fmt.Sprintf(`Page-([0-9\.]+)-%v`, defaultuint+v.(uint))).MatchString(purl) {
		t.Error(purl)
	}

	v = (uint16)(r.Int()) >> 2
	p.ArraySet(1, defaultuint16)
	p.ArrayAdd(1, v)
	purl = tp.GetURLRawPath()
	if !regexp.MustCompile(fmt.Sprintf(`Page-([0-9\.]+)-%v`, defaultuint16+v.(uint16))).MatchString(purl) {
		t.Error(purl)
	}

	v = (uint8)(r.Int()) >> 2
	p.ArraySet(1, defaultuint8)
	p.ArrayAdd(1, v)
	purl = tp.GetURLRawPath()
	if !regexp.MustCompile(fmt.Sprintf(`Page-([0-9\.]+)-%v`, defaultuint8+v.(uint8))).MatchString(purl) {
		t.Error(purl)
	}

	var defaultfloat64 float64 = 1
	var defaultfloat32 float32 = 1

	v = r.Float64()

	p.ArraySet(1, defaultfloat64)
	p.ArrayAdd(1, v)
	purl = tp.GetURLRawPath()
	if !regexp.MustCompile(fmt.Sprintf(`Page-([0-9\.]+)-%v`, defaultfloat64+v.(float64))).MatchString(purl) {
		t.Error(purl)
	}

	v = r.Float32()
	p.ArraySet(1, defaultfloat32)
	p.ArrayAdd(1, v)
	purl = tp.GetURLRawPath()
	if !regexp.MustCompile(fmt.Sprintf(`Page-([0-9\.]+)-%v`, defaultfloat32+v.(float32))).MatchString(purl) {
		t.Error(purl)
	}

}

func TestFocreParamQuery(t *testing.T) {
	r := random.New()
	ses := NewSession()
	tp := ses.Get("http://httpbin.org/get?page=1&arrpage[]=1&name=xiaoming&float=0.12")

	err := checkParam(tp, "page", "1")
	if err != nil {
		t.Error(err)
	}

	for i := 0; i < 100; i++ {
		checkBaseTypeParamSet(tp, r, t)
	}

	tp = ses.Get("http://httpbin.org/get?page=1&arrpage[]=1&name=xiaoming&float=0.12")
	for i := 0; i < 100; i++ {
		checkBaseTypeParamAdd(tp, r, t)
	}

	tp = ses.Get("http://httpbin.org/get?page[]=1&page[]=2&page[]=3&name=xiaoming")
	p := tp.QueryParam("page[]")

	p.ArraySet(2, 1)
	checkArrayParam(tp, "page[]", "[1 2 1]")

	p.ArraySet(2, int8(2))
	checkArrayParam(tp, "page[]", "[1 2 2]")

	p.ArraySet(2, int16(1))
	checkArrayParam(tp, "page[]", "[1 2 1]")

	p.ArraySet(2, int32(2))
	checkArrayParam(tp, "page[]", "[1 2 2]")

	p.ArraySet(2, int64(1))
	checkArrayParam(tp, "page[]", "[1 2 1]")

	p.ArraySet(2, uint(1))
	checkArrayParam(tp, "page[]", "[1 2 1]")

	p.ArraySet(2, uint8(2))
	checkArrayParam(tp, "page[]", "[1 2 2]")

	p.ArraySet(2, uint16(1))
	checkArrayParam(tp, "page[]", "[1 2 1]")

	p.ArraySet(2, uint32(0))
	checkArrayParam(tp, "page[]", "[0 2 1]")

	p.ArraySet(2, uint64(1))
	checkArrayParam(tp, "page[]", "[1 2 1]")

	p.ArraySet(2, 3.1)
	checkArrayParam(tp, "page[]", "[1 2 3.1]")

	p.ArraySet(2, float32(3.1))
	checkArrayParam(tp, "page[]", "[1 2 3.1]")
}

func TestFocreParamRegexp(t *testing.T) {
	r := random.New()
	ses := NewSession()
	tp := ses.Get("https://10.api.xxx.tv/oversea/xxx/api/v2/liveRoom/Page-1-30-/HK/1028/1000")

	// p := tp.PathParam(`Page-(\w+)-(\w+)`)
	for i := 0; i < 100; i++ {
		checkBaseTypeParamRegexpSet(tp, r, t)
		checkBaseTypeParamRegexpArraySet(tp, r, t)
		checkBaseTypeParamRegexpAdd(tp, r, t)
		checkBaseTypeParamRegexpArrayAdd(tp, r, t)
	}
}

// func Benchmark(b *testing.B) {
// 	var a []string = []string{"ads", "asfdf", "13123"}

// 	b.Run("+", func(b *testing.B) {
// 		for i := 0; i < b.N; i++ {
// 			concat1(a)
// 		}
// 	})

// 	b.Run("append", func(b *testing.B) {
// 		for i := 0; i < b.N; i++ {
// 			concat2(a)
// 		}
// 	})

// 	b.Run("buffer", func(b *testing.B) {
// 		for i := 0; i < b.N; i++ {
// 			concat(a)
// 		}
// 	})

// }
