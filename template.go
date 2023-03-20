package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var html = `
<pre id="json"></pre>

<head>
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
</head>

<script>
    var data = %s
document.getElementById("json").innerHTML = JSON.stringify(data, undefined, 2);
</script>
`

func WrapHTML(rw http.ResponseWriter, i interface{}) {
	data, err := json.Marshal(i)
	if err != nil {
		rw.Write([]byte(err.Error()))
		return
	}
	rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	rw.Write([]byte(fmt.Sprintf(html, string(data))))
}
