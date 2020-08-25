package requests

// Utils 一些实用功能
type Utils struct {
}

// Param("page").ForInt(0, )
// Param("page").String
type AdjustQuery struct {
}

type AdjustQueryIterator struct {
	param string
}

func (aq *AdjustQuery) Strings(strlist ...string) {

}

func (aq *AdjustQuery) ForFloat(start, end, step float64) {

}

func (aq *AdjustQuery) ForInt(start, end, step int64) {

}

func (aq *AdjustQuery) ForChar(start, end byte) {

}
