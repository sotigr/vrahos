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

// hx-get="/list/projects/"
type IndexProps struct {
	ExtraHead string
}

func (p IndexPage) Template() string {
	//https://github.com/cubiclesoft/js-fileexplorer?tab=readme-ov-file
	return /*html*/ `
	{{template "DocumentStart" .}}
 
<c-box style="width: 800px; margin: 0 auto; margin-top: 30px;">
	<div style="margin-bottom: 10px"> 
		{{template "CreateProjectForm"}} 
	</div>

	<c-box-inner  
		hx-ext="sse" 
		sse-connect="/server-events/?lobbies=project_list_lobby:none" 
		sse-swap="project_list_lobby.message"
		hx-swap="beforeend" 
		class="flex flex-wrap gap-2" 
	>
		{{template "ListProject" makeprops .Props.ListProps}}
	</c-box-inner>
 
	 
</c-box> 
	{{template "DocumentEnd" .}}
	`
}

// <div class="container sample center">
// <igc-icon name="search" collection="material"></igc-icon>
// </div>
func (p IndexPage) Props(r *http.Request, meta *vrahos.MetaData) (any, map[string]string) {
	return IndexProps{
		ExtraHead: `<title>Welcome to vrahos</title>`,
	}, nil
}
