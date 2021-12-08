package requests

import (
	"bytes"
	"net/http"
)

func buildBodyRequest(tp *Temporary) *http.Request {
	var req *http.Request
	var err error

	if tp.Body == nil {
		req, err = http.NewRequest(tp.Method, tp.GetRawURL(), nil)
		if err != nil {
			panic(err)
		}
	} else {
		// var buf = bytes.NewBuffer(nil)
		var buf = bytes.NewBuffer(tp.Body.Bytes())
		var ct CompressType
		if tp.compressType != ContentEncodingNoCompress {
			ct = tp.compressType
		} else if tp.session.compressType != ContentEncodingNoCompress {
			ct = tp.session.compressType
		}
		switch ct {
		case ContentEncodingNoCompress:
			// if _, err = buf.Write(tp.Body.Bytes()); err != nil {
			// 	panic(err)
			// }
		case ContentEncodingGzip:
			// buf = bytes.NewBuffer(nil)
			// w := gzip.NewWriter(buf)
			// if _, err = w.Write(tp.Body.Bytes()); err != nil {
			// 	panic(err)
			// }
			// err = w.Close()
			// if err != nil {
			// 	panic(err)
			// }
			tp.Header.Add("Accept-Encoding", "gzip")
			tp.Header.Add("Content-Encoding", "gzip")
		// case ContentEncodingDeflate:
		// 	// if _, err = zlib.NewWriter(buf).Write(tp.Body.Bytes()); err != nil {
		// 	// 	panic(err)
		// 	// }
		// 	tp.Header.Add("Accept-Encoding", "deflate")
		// 	tp.Header.Add("Content-Encoding", "deflate")
		// case ContentEncodingCompress:
		// 	// if _, err = lzw.NewWriter(buf, lzw.MSB, 8).Write(tp.Body.Bytes()); err != nil {
		// 	// 	panic(err)
		// 	// }
		// 	tp.Header.Add("Accept-Encoding", "compress")
		// 	tp.Header.Add("Content-Encoding", "compress")
		// case ContentEncodingBr:
		// 	tp.Header.Add("Accept-Encoding", "br")
		// 	tp.Header.Add("Content-Encoding", "br")
		default:
			panic("compress type not support")
		}

		req, err = http.NewRequest(tp.Method, tp.GetRawURL(), buf)
		if err != nil {
			panic(err)
		}
		// req.ContentLength = int64(bodybuf.Len())
	}

	return req
}
