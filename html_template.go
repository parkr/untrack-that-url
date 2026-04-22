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
    :root {
      --primary: #2563eb;
      --primary-hover: #1d4ed8;
      --bg: #f8fafc;
      --card-bg: #ffffff;
      --text: #1e293b;
      --text-muted: #64748b;
      --error: #ef4444;
      --success: #10b981;
      --border: #e2e8f0;
    }

    * { box-sizing: border-box; }

    body {
      font-family: system-ui, -apple-system, sans-serif;
      background: linear-gradient(135deg, #f0f9ff 0%, #e0f2fe 100%);
      color: var(--text);
      display: flex;
      align-items: center;
      justify-content: center;
      min-height: 100vh;
      margin: 0;
      padding: 20px;
    }

    .container {
      width: 100%;
      max-width: 500px;
      background: var(--card-bg);
      padding: 2rem;
      border-radius: 1rem;
      box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.1), 0 8px 10px -6px rgba(0, 0, 0, 0.1);
    }

    h1 {
      margin: 0 0 0.5rem 0;
      font-size: 1.5rem;
      font-weight: 700;
      text-align: center;
    }

    p.subtitle {
      color: var(--text-muted);
      text-align: center;
      margin-bottom: 2rem;
      font-size: 0.875rem;
    }

    .input-group {
      display: flex;
      flex-direction: column;
      gap: 0.5rem;
      margin-bottom: 1.5rem;
    }

    label {
      font-size: 0.875rem;
      font-weight: 600;
      color: var(--text-muted);
    }

    input[type="text"] {
      padding: 0.75rem 1rem;
      border: 1px solid var(--border);
      border-radius: 0.5rem;
      font-size: 1rem;
      width: 100%;
      transition: border-color 0.2s, ring 0.2s;
    }

    input[type="text"]:focus {
      outline: none;
      border-color: var(--primary);
      box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.1);
    }

    button {
      background: var(--primary);
      color: white;
      border: none;
      padding: 0.75rem 1.5rem;
      border-radius: 0.5rem;
      font-size: 1rem;
      font-weight: 600;
      cursor: pointer;
      transition: background 0.2s;
      width: 100%;
    }

    button:hover:not(:disabled) {
      background: var(--primary-hover);
    }

    button:disabled {
      opacity: 0.7;
      cursor: not-allowed;
    }

    #errors {
      color: var(--error);
      background: #fef2f2;
      padding: 0.75rem;
      border-radius: 0.5rem;
      font-size: 0.875rem;
      margin-bottom: 1rem;
      display: none;
      border: 1px solid #fee2e2;
    }

    #result {
      margin-top: 2rem;
      padding-top: 1.5rem;
      border-top: 1px solid var(--border);
      display: none;
    }

    #trail {
      margin-top: 1.5rem;
      display: none;
    }

    .trail-list {
      list-style: none;
      padding: 0;
      margin: 0;
      border-left: 2px solid var(--border);
      margin-left: 0.5rem;
    }

    .trail-item {
      position: relative;
      padding-left: 1.5rem;
      margin-bottom: 0.75rem;
      font-size: 0.8125rem;
      color: var(--text-muted);
      word-break: break-all;
    }

    .trail-item::before {
      content: "";
      position: absolute;
      left: -2px;
      top: 0.5rem;
      width: 8px;
      height: 8px;
      background: var(--card-bg);
      border: 2px solid var(--border);
      border-radius: 50%;
      transform: translateX(-50%);
    }

    .trail-item:last-child {
      color: var(--text);
      font-weight: 500;
    }

    .trail-item:last-child::before {
      border-color: var(--primary);
      background: var(--primary);
    }

    .result-label {
      font-size: 0.75rem;
      font-weight: 700;
      text-transform: uppercase;
      letter-spacing: 0.05em;
      color: var(--text-muted);
      margin-bottom: 0.5rem;
    }

    .result-box {
      background: var(--bg);
      padding: 0.75rem;
      border-radius: 0.5rem;
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: 0.5rem;
      border: 1px solid var(--border);
    }

    .result-url {
      font-size: 0.875rem;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
      color: var(--primary);
      text-decoration: none;
      flex: 1;
    }

    .copy-btn {
      background: white;
      color: var(--text);
      border: 1px solid var(--border);
      padding: 0.4rem 0.8rem;
      font-size: 0.75rem;
      width: auto;
    }

    .copy-btn:hover {
      background: var(--bg);
    }

    @media (max-width: 480px) {
      .container {
        padding: 1.5rem;
      }
    }
  </style>
</head>
<body>
  <div class="container">
    <h1>Untrack That URL</h1>
    <p class="subtitle">Resolve tracking links safely and privately.</p>

    <div id="errors"></div>

    <div class="input-group">
      <label for="url">Enter URL to resolve:</label>
      <input type="text" id="url" name="url" placeholder="https://t.co/..." required autofocus>
    </div>

    <button id="resolveBtn" onclick="submitResolveURL();">Resolve URL</button>

    <div id="result">
      <div id="resolvedLabel" class="result-label">Resolved URL</div>
      <div id="resolvedBox" class="result-box">
        <a id="resultURL" class="result-url" href="#" target="_blank" rel="noopener noreferrer"></a>
        <button class="copy-btn" onclick="copyResult();">Copy</button>
      </div>

      <div id="trail">
        <div class="result-label">Redirect Trail</div>
        <ul id="trailList" class="trail-list"></ul>
      </div>
    </div>
  </div>

  <script type="application/javascript">
    const urlInput = document.getElementById("url");
    const resolveBtn = document.getElementById("resolveBtn");
    const errorDiv = document.getElementById("errors");
    const resultDiv = document.getElementById("result");
    const resultURL = document.getElementById("resultURL");
    const resolvedLabel = document.getElementById("resolvedLabel");
    const resolvedBox = document.getElementById("resolvedBox");
    const trailDiv = document.getElementById("trail");
    const trailList = document.getElementById("trailList");

    function showError(msg) {
      errorDiv.innerText = msg;
      errorDiv.style.display = "block";
    }

    async function submitResolveURL() {
      const inputURL = urlInput.value.trim();
      if (!inputURL) {
        showError("Please enter a URL to resolve.");
        resultDiv.style.display = "none";
        return;
      }

      errorDiv.style.display = "none";
      resultDiv.style.display = "none";
      resolveBtn.disabled = true;
      resolveBtn.innerText = "Resolving...";

      try {
        const response = await fetch("/resolve.json", {
          method: "POST",
          headers: { "Content-Type": "application/x-www-form-urlencoded" },
          body: "url=" + encodeURIComponent(inputURL)
        });

        const info = await response.json();

        // Render trail if it exists, regardless of error
        trailList.innerHTML = "";
        if (info.trail && info.trail.length > 0) {
          info.trail.forEach(u => {
            const li = document.createElement("li");
            li.className = "trail-item";
            li.innerText = u;
            trailList.appendChild(li);
          });
          trailDiv.style.display = "block";
          resultDiv.style.display = "block";
        } else {
          trailDiv.style.display = "none";
        }

        if (info.error) {
          showError(info.error);
          resolvedLabel.style.display = "none";
          resolvedBox.style.display = "none";
        } else {
          resultURL.innerText = info.url;
          resultURL.href = info.url;
          resolvedLabel.style.display = "block";
          resolvedBox.style.display = "flex";
          resultDiv.style.display = "block";
          errorDiv.style.display = "none";
        }
      } catch (err) {
        showError("An unexpected error occurred. Please try again.");
        console.error(err);
      } finally {
        resolveBtn.disabled = false;
        resolveBtn.innerText = "Resolve URL";
      }
    }

    async function copyResult() {
      const text = resultURL.innerText;
      try {
        await navigator.clipboard.writeText(text);
        const copyBtn = document.querySelector(".copy-btn");
        const originalText = copyBtn.innerText;
        copyBtn.innerText = "Copied!";
        setTimeout(() => { copyBtn.innerText = originalText; }, 2000);
      } catch (err) {
        console.error("Failed to copy: ", err);
      }
    }

    urlInput.addEventListener("keypress", (e) => {
      if (e.key === "Enter") submitResolveURL();
    });
  </script>
</body>
</html>
`))

func RenderHTML(w io.Writer) error {
	return indexTemplate.Execute(w, nil)
}
