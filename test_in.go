package requests

import "net/http"

type ITestServer interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}
