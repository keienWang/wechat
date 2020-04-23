package gongNeng

type TianGouRiJi struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Content string `json:"content"`
	} `json:"data"`
	Author struct {
		Name string `json:"name"`
		Desc string `json:"desc"`
	} `json:"author"`
}

type ZhaManWord struct {
	Code      int      `json:"code"`
	Message   string   `json:"message"`
	ReturnObj []string `json:"returnObj"`
}

type JiTang struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Title string `json:"title"`
	} `json:"data"`
	Author struct {
		Name string `json:"name"`
		Desc string `json:"desc"`
	} `json:"author"`
}
type RenJian struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Say  string `json:"say"`
}
