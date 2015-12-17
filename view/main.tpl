{{define "body"}}
<div><img src="/img/gdg-riga.png"><h1>Workshop blog</h1></div>
<a href="/post">Create new post</a>
{{ range $index, $element := .storage }}
    <div style="margin-top:20px">
        <div><h3><a href="/edit/{{.ID}}">{{.Title}}</a></h3>{{.CreateAt}}</div>
        <div>{{.Body}}</div>
        <div>{{if .Picture}}<img src="/img/3_{{.Picture}}/">{{end}}</div>
    </div>
{{end}}
{{end}}
