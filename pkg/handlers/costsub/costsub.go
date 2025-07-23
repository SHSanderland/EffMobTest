package costsub

import (
	"log/slog"
	"net/http"
)

type CostSubscription interface {
	CostSubscription()
}

func Handler(
	l *slog.Logger, cs CostSubscription,
	w http.ResponseWriter, r *http.Request,
) {
	// const fn = "handlers.rsub.Handler"
	// log := l.With(
	// 	slog.String("fn", fn),
	// 	slog.String("requestID", middleware.GetReqID(r.Context())),
	// )

	// startPeriodStr := r.URL.Query().Get("start_period")
	// endPeriodStr := r.URL.Query().Get("end_period")
	// userIDStr := r.URL.Query().Get("user_id")
	// serviceName := r.URL.Query().Get("service_name")

	// if startPeriodStr == "" || endPeriodStr == "" {
	// 	log.Error(
	// 		"wrong period",
	// 		slog.String("start_date", startPeriodStr),
	// 		slog.String("end_date", endPeriodStr),
	// 	)
	// 	http.Error(w, "Wrong period", http.StatusBadRequest)

	// 	return
	// }

	// startDate, err := time.Parse("01-2006", startPeriodStr)
	// if err != nil {
	// 	log.Error(
	// 		"Bad period format",
	// 		slog.String("start_date", startPeriodStr),
	// 	)
	// 	http.Error(w, "Bad period format", http.StatusBadRequest)

	// 	return
	// }

	// endDate, err := time.Parse("01-2006", endPeriodStr)
	// if err != nil {
	// 	log.Error(
	// 		"Bad period format",
	// 		slog.String("end_date", endPeriodStr),
	// 	)
	// 	http.Error(w, "Bad period format", http.StatusBadRequest)

	// 	return
	// }

}
