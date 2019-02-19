package routing

import (
	"Project/doit/app"
	"Project/doit/code"
	"fmt"
	"github.com/go-ozzo/ozzo-routing"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
	"net/http"
	"os"
	"runtime/debug"
)

func errorHandler(c *routing.Context) (err error) {

	/*-----捕获panic异常，处理相关错误-----*/
	defer func() {
		if e := recover(); e != nil {
			//格式化堆栈跟踪
			stack := string(debug.Stack())
			app.Logger.Error().Str("stacktrace", stack).Msgf("recovered from panic: %v", e)
			fmt.Fprintln(os.Stderr, stack)
			http.Error(c.Response, "Internal Server Error.", http.StatusInternalServerError)
			c.Abort()
			err = nil
		}
	}()

	/*处理关联程序并处理对应错误*/
	err = c.Next()
	if err != nil {
		c.Abort()
		switch e := errors.Cause(err).(type) {
		//验证类错误
		case validation.Errors:
			err := code.New(http.StatusBadRequest, code.CodeBadRequest)
			for k, v := range e {
				err.Err(fmt.Sprintf("\"%s\" %s", k, v))
			}
			c.Response.WriteHeader(err.Status)
			return c.Write(err)
		//HTTP类错误
		case routing.HTTPError:
			http.Error(c.Response, e.Error(), e.StatusCode())
		//错误状态码
		case *code.Error:
			c.Response.WriteHeader(e.Status)
			return c.Write(e)
		//处理函数的其他错误
		default:
			app.Logger.Error().Err(err).Msg("handler error.")
			http.Error(c.Response, "Internal Server Error.", http.StatusInternalServerError)
			return nil
		}
	}
	return nil
}
