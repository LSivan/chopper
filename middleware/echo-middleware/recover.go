package echo_middleware

import (
	"fmt"
	xerror "github.com/LSivan/hatchet/x-error"
	"github.com/labstack/echo"
	"net/http"
	"time"
)

type Handler func(err error) error

type panicError struct {
	Time time.Time
	Err  error
}

func (p panicError) Error() string {
	return fmt.Sprintf("panic on %s, %v", p.Time.Format(time.RFC3339), p.Err)
}
func newPanicError(err error) panicError {
	return panicError{
		Time: time.Now(),
		Err:  err,
	}
}

func Recover(handlerFunc echo.HandlerFunc, handler Handler) echo.HandlerFunc {

	return func(c echo.Context) error {

		defer func() {
			if handler == nil {
				handler = func(err error) error {
					c.Logger().Errorf("%s", newPanicError(err).Error())
					return nil
				}
			}
			err := recover()
			switch e := err.(type) {
			case nil:
				return
			case error:
				xerror.DoIfErrorNotNil(handler(e), func(err2 error) {
					c.Logger().Errorf("fail to handle error:%+v", err2)
				})
			default:
				xerror.DoIfErrorNotNil(handler(fmt.Errorf("%v", e)), func(err2 error) {
					c.Logger().Errorf("fail to handle error:%+v", err2)
				})
			}
			_ = c.JSON(http.StatusOK, struct {
				Msg string `json:"msg"`
			}{"Internal Server Error"})
		}()
		return handlerFunc(c)
	}
}
