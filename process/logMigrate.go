package process




type logProcessHandler struct {
	series ProcessHandler
	summary ProcessHandler
}

func NewLog(series ProcessHandler, summary ProcessHandler) ProcessHandler {
	return &logProcessHandler{series: series, summary: summary}
}

