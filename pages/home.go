package pages

import (
	"net/http"
	"server/vrahos"
)

type IndexPage struct{ vrahos.BasicComponent }

func (p IndexPage) Name() string {
	return "Home"
}

func (p IndexPage) URL() string {
	return "/"
}

type IndexProps struct {
	ExtraHead string
}

func (p IndexPage) Template() string {
	return `
	{{template "DocumentStart" .}}
 
<c-box style="width: 800px; margin: 0 auto; margin-top: 30px;">
	Welcome to this sample
</c-box> 
	{{template "DocumentEnd" .}}
	`
}

func (p IndexPage) Props(r *http.Request, meta *vrahos.MetaData) (any, map[string]string) {
	return IndexProps{
		ExtraHead: `<title>Welcome to vrahos</title>`,
	}, nil
}
