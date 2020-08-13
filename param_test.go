package requests

import (
	"encoding/json"
	"fmt"
	"log"
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
}

func TestQueryArrayParam(t *testing.T) {
	var err error
	var p *Param
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
