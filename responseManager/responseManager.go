package ResponseManager

import(
	"net/http"
	"encoding/json"
	"github.com/labstack/echo"
)

type APIResponse struct {
	Code    int    `json:"code"`
	Type    string `json:"type"`
	Message string `json:"message"`
}

//ResponseMapper to handle response
func ResponseMapper(code int, message string, context echo.Context) error {
	context.Response().WriteHeader(code)
	return json.NewEncoder(context.Response()).Encode(ResponseManager(code, message))
}

//ResponseManager to handle response
func ResponseManager(code int, errMessage string) APIResponse {
	response := APIResponse{}
	response.Code = code
	response.Type = http.StatusText(code)
	response.Message = errMessage
	return response
}