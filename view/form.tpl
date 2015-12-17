{{define "form"}}
<br/>
<form action="/post" method="post" ENCTYPE="multipart/form-data">
    <input type="hidden" name="id" value="{{.id}}"/>
    <input type="hidden" name="action" value="{{.action}}"/>
    <b>Title</b></br>
    <input type="text"  name="title" value="{{.title}}"/ size=100><br/><br/>
    <b>Body</b><br/>
    <textarea name="body" cols="98" rows="20">{{.body}}</textarea></br><br/>
    {{if .id}}
    {{else}}
    <b>Picture</b><br/>
    <input type="file" name="picture"/>
    {{end}}
    <input type="button" value="Cancel" onclick="window.location.href='/'">&nbsp;<input type="submit" value="submit">
</form>

{{end}}