package untrackthaturl

import (
	"html/template"
	"io"
)

var indexTemplate = template.Must(template.New("index.html").Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1">
  <meta content="origin-when-cross-origin" name="referrer" />
  <title>Untrack That URL</title>
  <style type="text/css">
  .markdown-body {
      width: 95%;
	  margin: 0 auto;
	  text-align: center;
  }
  #errors { color: red; }
  </style>
  <script type="application/javascript">
    function writeError(error) {
		var errorDiv = document.getElementById("errors");
		errorDiv.innerText = "error: "+error;
	}
    function reqListener () {
      console.log(this.responseText);

      if (this.responseText === null || this.responseText === "") {
        console.error("nada");
        return;
      }

	  var info = JSON.parse(this.responseText);

	  if (info.Error !== undefined) {
		writeError(info.Error);
		return
	  }

      var resultDiv = document.getElementById("resultURL");
      resultDiv.innerText = "";

      var anchor = document.createElement("a");
      anchor.href = info.URL;
      anchor.title = "Resolved URL";
      anchor.innerText = info.URL;
	  resultDiv.appendChild(anchor);
    }

	function submitResolveURL() {
	  document.getElementById("errors").innerText = ""

	  var inputURL = document.getElementById("url").value;
      if (inputURL === undefined || inputURL === "") {
		  writeError("input URL required");
		  return
	  }
	  var params = 'url='+escape(inputURL);

	  var oReq = new XMLHttpRequest();
  	  oReq.addEventListener("load", reqListener);
	  oReq.open("POST", "/resolve.json", true);
	  oReq.setRequestHeader('Content-type', 'application/x-www-form-urlencoded');
	  oReq.send(params);
	}

    </script>
  </head>
  <body>
    <div class="markdown-body">
	  <p>Enter a URL to resolve:</p>
	  <p>
	    <div id="errors"></div>
        <label for="name">URL:</label>
	    <input type="text" id="url" name="url" required>
		<input type="button" onclick="submitResolveURL();" value="Resolve">
	  </p>
	  <p>
		<div id="resultURL"></div>
	  </p>
    </div>
  </body>
</html>
`))

func RenderHTML(w io.Writer) error {
	return indexTemplate.Execute(w, nil)
}
