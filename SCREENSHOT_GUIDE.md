# Screenshot Guide for BubbleWebServer

This guide tells you exactly where to take screenshots and how to add them to the READMEs.

## 📁 Directory Structure

First, create an `assets` folder in the root of your project:

```
BubbleWebServer/
├── assets/                          ← Create this folder
│   ├── tui-main-screen.png
│   ├── tui-port-selection.png
│   ├── tui-initial-screen.png
│   ├── tui-request-logs.png
│   ├── tui-port-form.png
│   ├── simple-example-browser.png
│   └── api-example-responses.png
├── examples/
├── server/
└── README.md
```

## 📸 Screenshots Needed

### 1. **Main README Screenshots**

#### Screenshot 1: `assets/tui-main-screen.png`
- **Location in file:** Line 80 of `README.md`
- **What to capture:**
  - Run `go run main.go`
  - Type `/start` to start the server
  - Make 1-2 test requests (curl http://localhost:8090/hello)
  - Capture the TUI showing the banner, command prompt, and request logs
- **Size:** Full terminal window
- **Focus:** Show the TUI in action with the custom banner and live requests

#### Screenshot 2: `assets/tui-port-selection.png`
- **Location in file:** Line 85 of `README.md`
- **What to capture:**
  - Run `go run main.go`
  - Type `/port` and press Enter
  - Start typing a port number (like `3000`)
  - Capture the centered form with the rounded border
- **Size:** Full terminal window
- **Focus:** Show the port selection form centered on screen

---

### 2. **TUI Example README Screenshots**

#### Screenshot 3: `assets/tui-initial-screen.png`
- **Location in file:** Line 146 of `examples/tui/README.md`
- **What to capture:**
  - Run `cd examples/tui && go run main.go`
  - Capture immediately after launch
  - Shows welcome banner and command prompt
- **Size:** Full terminal window
- **Focus:** Clean initial state

#### Screenshot 4: `assets/tui-request-logs.png`
- **Location in file:** Line 171 of `examples/tui/README.md`
- **What to capture:**
  - Server running with `/start`
  - Make several curl requests to different endpoints
  - Capture showing 4-5 request logs in the viewport
- **Size:** Full terminal window
- **Focus:** Demonstrate real-time request monitoring

#### Screenshot 5: `assets/tui-port-form.png`
- **Location in file:** Line 212 of `examples/tui/README.md`
- **What to capture:**
  - Type `/port`
  - Start typing a port number
  - Capture the form centered on screen
- **Size:** Full terminal window
- **Focus:** Port selection form with validation

---

### 3. **Simple Example README Screenshot**

#### Screenshot 6: `assets/simple-example-browser.png`
- **Location in file:** Line 27 of `examples/simple/README.md`
- **What to capture:** Choose ONE of these options:
  - **Option A:** Browser showing http://localhost:8080/ HTML page
  - **Option B:** Terminal showing curl commands and responses
  - **Option C:** Split screen: browser on one side, terminal on other
- **Size:** Browser window or terminal
- **Focus:** Show the HTML response or JSON output

---

### 4. **API Example README Screenshot**

#### Screenshot 7: `assets/api-example-responses.png`
- **Location in file:** Line 28 of `examples/api/README.md`
- **What to capture:** Choose ONE of these options:
  - **Option A:** Postman or Insomnia showing GET/POST requests
  - **Option B:** Terminal with curl commands and their JSON responses
  - **Option C:** Multiple terminal panes showing different API calls
- **Size:** REST client window or terminal
- **Focus:** Show JSON responses from the Todo API

---

## 🖼️ How to Add Images to GitHub

### GitHub Markdown Syntax

```markdown
![Alt Text](path/to/image.png)
```

### For this project:

Since the `assets/` folder is at the root level:

**In main README.md:**
```markdown
![TUI Main Screen](assets/tui-main-screen.png)
```

**In examples/tui/README.md (nested):**
```markdown
![TUI Initial Screen](../../assets/tui-initial-screen.png)
```

### The images are ALREADY referenced in the READMEs!

All the markdown is already in place. You just need to:
1. Create the `assets/` folder
2. Take the screenshots
3. Save them with the exact filenames listed above
4. Commit and push to GitHub

---

## 🎨 Screenshot Tips

### Terminal Screenshots
- **macOS:** Cmd+Shift+4, then press Space to capture window
- **Windows:** Use Snipping Tool or Windows+Shift+S
- **Linux:** Use `gnome-screenshot` or `scrot`

### Best Practices
- Use a clean terminal with good contrast
- Remove any sensitive information
- Use a readable font size (14-16pt)
- Capture full window including title bar for context
- PNG format recommended (better quality for text)

### Terminal Theme Recommendations
- Use a dark theme with good contrast
- Default colors work fine
- Ensure text is readable at smaller sizes

---

## ✅ Checklist

Before pushing to GitHub:

- [ ] Create `assets/` folder in project root
- [ ] Screenshot 1: TUI main screen with requests
- [ ] Screenshot 2: Port selection form (main README)
- [ ] Screenshot 3: TUI initial screen
- [ ] Screenshot 4: TUI with request logs
- [ ] Screenshot 5: Port selection form (TUI example)
- [ ] Screenshot 6: Simple example in browser/terminal
- [ ] Screenshot 7: API responses
- [ ] All images are in PNG format
- [ ] All filenames match exactly (case-sensitive!)
- [ ] Test locally: render README in markdown viewer
- [ ] Commit with message: "Add screenshots to documentation"
- [ ] Push to GitHub
- [ ] Verify images display correctly on GitHub

---

## 🔧 If Images Don't Show on GitHub

Common issues:
1. **Wrong path:** Check capitalization and slashes
2. **Not committed:** Make sure to `git add assets/`
3. **Wrong branch:** Images must be in the same branch
4. **File size:** GitHub has a 25MB limit per file

---

## 📝 Alternative: Online Images

If you prefer to host images elsewhere:

```markdown
![Alt Text](https://your-url.com/image.png)
```

**Recommended hosts:**
- GitHub Issues (upload, then copy URL)
- Imgur (free image hosting)
- Your own CDN

---

## Summary

The README files are **already updated** with image placeholders. You just need to:

1. **Create:** `mkdir assets`
2. **Capture:** Take 7 screenshots following the guide above
3. **Name:** Use exact filenames (see list at top)
4. **Commit:** `git add assets/ && git commit -m "Add documentation screenshots"`
5. **Push:** `git push`
6. **Verify:** Check on GitHub that images display

All the markdown syntax is already in place! 🎉
