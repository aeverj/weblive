package common

type DataRespResult struct {
	Result      *Result
	RespContent *RespContent
}

type RespContent struct {
	Target     string
	Redirect   string
	Headers    map[string][]string
	StatusCode int
	Html       string
}

type Result struct {
	Target     string `json:"url"`
	Redirect   string `json:"-"`
	Title      string `json:"title"`
	Header     string `json:"headers"`
	StatusCode int    `json:"statusCode"`
	IP         string `json:"ip"`
	CDN        string `json:"cdn"`
	Finger     string `json:"fingerPrint"`
	ScreenPath string `json:"screenshot"`
}
