package requests

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"testing"
)

func checkArrayParam(tp *Temporary, param string, vaild string) error {

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
			if fmt.Sprint(page) != vaild {
				log.Println(data, string(resp.Content()))
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

func checkParam(tp *Temporary, param string, vaild string) error {

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
			if page.(string) != vaild {
				log.Println(data)
				return fmt.Errorf("param: %s", param)
			}
		} else {
			log.Println(data)
			return fmt.Errorf("param is %s not exists", param)
		}
	} else {
		log.Println(data)
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
