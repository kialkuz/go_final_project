package task

import (
	"log"
	"net/http"
	"time"

	"github.com/Yandex-Practicum/final/internal/dto"
	httpService "github.com/Yandex-Practicum/final/internal/services/http"
	"github.com/Yandex-Practicum/final/internal/services/task/nextDate"
	datePkg "github.com/Yandex-Practicum/final/pkg/date"
)

func NextDayHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	now, err := time.Parse(datePkg.DateFormat, q.Get("now"))

	if err != nil {
		log.Println(err)
		httpService.WriteJsonBadRequest(w, dto.ErrorResponse{ErrorText: "ошибка получения новой даты"})
		return
	}

	nextDate, err := nextDate.NextDate(now, q.Get("date"), q.Get("repeat"))
	if err != nil {
		log.Println(err)
		httpService.WriteJsonInternalServerError(w, dto.ErrorResponse{ErrorText: "ошибка получения новой даты"})
		return
	}

	httpService.WriteJsonWithoutSerialize(w, []byte(nextDate), http.StatusOK)
}
