{{define "layout"}}
<html>
<head>
    <title>{{.PageTitle}}</title>
    <style>
        html, textarea,input {font-family:"Helvetica Neue"; font-size:12px}

    </style>
</head>
<body>
{{template "body" .}}
</body>
</html>
{{end}}
