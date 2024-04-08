package pages

import (
	"server/vrahos"
)

type Document struct{ vrahos.BasicComponent }

func (p Document) Name() string {
	return "Document"
}

type DocumentProps struct {
	ExtraHead string
}

func (p Document) Template() string {
	return `
	{{define "DocumentStart"}}
	<!DOCTYPE html>
	<html> 
		<head>
			<meta charset='utf-8' />
			<meta http-equiv='X-UA-Compatible' content='IE=edge' /> 

			<meta name='viewport' content='width=device-width, initial-scale=1' />  

			<link rel="icon" href="/static/favicon.svg" type="image/svg+xml" /> 
    
			<link rel="preload" href="/static/htmx.min.js" as="script"/>
			<link rel="preload" href="/static/.dist/main.js" as="script"/>
 
			<link rel="preload" href="/static/.dist/styles.css" as="style"/>
			<link rel="stylesheet" href="/static/.dist/styles.css" />  
 
			{{if .Props.ExtraHead}}
			{{.Props.ExtraHead}}
			{{end}} 
		</head> 
		<body>   
	{{end}}
	{{define "DocumentEnd"}}     
		<script src="/static/htmx.min.js" defer="defer"></script>   
		<script src="/static/.dist/main.js" defer="defer"></script>   
		</body> 
	</html>
	{{end}}
	`
}
