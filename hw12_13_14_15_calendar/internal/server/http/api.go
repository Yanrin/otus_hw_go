package httpserver

import (
	"fmt"
	"net/http"

	"github.com/Yanrin/otus_hw_go/hw12_13_14_15_calendar/internal/app"
	"github.com/Yanrin/otus_hw_go/hw12_13_14_15_calendar/internal/logger"
)

type API struct {
	Log      *logger.Logger
	Calendar app.Calendar
}

func (api *API) Hello(w http.ResponseWriter, r *http.Request) {
	str := "Hello world!"
	fmt.Println(str)
	_, _ = w.Write([]byte(str))
	api.Log.Info(str)
}
