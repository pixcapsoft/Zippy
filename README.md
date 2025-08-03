# Zippy

> **Simple, cross-platform version control with zip storage**

---

![Zippy Banner](https://raw.githubusercontent.com/pixcapsoft/Zippy/main/banner.png)

---

## 🚀 What is Zippy?

**Zippy** is a lightweight, Git-inspired version control tool that stores your project versions as zip files. It’s perfect for simple backups, sharing, and versioning—without the complexity of full-blown VCS systems.

- **Cross-platform:** Works on Windows, Linux, and macOS  
- **No server required:** All data is local, portable, and easy to share  
- **`.zippyignore`:** Ignore files/folders just like `.gitignore`  
- **Human-readable tags:** Use version tags like `v1.0`, `release-2024`, etc.  
- **Easy restore:** Instantly restore any version—or just a specific file/folder  

---

## ✨ Features

- Simple CLI, easy to learn  
- Staging area (like `git add`)  
- Commit with message and tag  
- List, diff, and restore versions  
- Patch existing versions  
- `status` command shows what will be committed and what changed  
- Cross-platform builds (Windows, Linux, macOS, 32/64-bit)

---

## 📦 Installation

### 🔹 Pre-built Binaries

Download from the [Releases](https://github.com/pixcapsoft/Zippy/releases) page for your OS/architecture:

- `zippy.exe` for Windows  
- `zippy-linux` for Linux  
- `zippy-darwin` for macOS  

⚠️ **Note:**  
- Zippy for Linux is tested only on Ubuntu 25.0.  
- macOS builds are untested by the developer.  
- Windows builds are fully tested.  
For best results on Linux/macOS, consider building Zippy from source.

---

### 🖼️ Windows Installation

1. Download the `.msi` installer for your processor and architecture.  
2. Open the installer.  
3. Zippy is not signed with a trusted certificate, so Windows SmartScreen may show an "Unknown Publisher" warning. You can safely ignore it.  
4. Choose your installation path and proceed.  
5. The installer automatically adds Zippy to your system variables.  
6. Open Command Prompt and run:

   ```bash
   zippy -v
   ```

   If you see the version and details, Zippy is installed successfully. 💥  
   If you get an "unknown command" error, Zippy may not be installed correctly or wasn't added to your system variables. In that case, manually add the installation directory to your system `Path`.

---

### 🐧 Linux Installation

1. Extract the downloaded `.zip`  
2. Rename the binary to `zippy` (recommended)  
3. Make it executable:

   ```bash
   chmod +x zippy
   ./zippy -v
   ```

   If you see the version info, Zippy is working. To install system-wide:

   ```bash
   sudo cp zippy /usr/local/bin/zippy
   ```

---

### 🍎 macOS Installation

1. Extract the downloaded `.zip`  
2. Rename the binary to `zippy` (recommended)  
3. Make it executable:

   ```bash
   chmod +x zippy
   ./zippy -v
   ```

   If you see the version info, Zippy is working. To install system-wide:

   ```bash
   sudo cp zippy /usr/local/bin/zippy
   ```

🧠 **Notes:**

- If you see a “permission denied” error, ensure the binary is executable (`chmod +x`).  
- On macOS, you may need to allow the binary to run via **System Preferences > Security & Privacy**.  
- Zippy is unsigned by default, so it may be flagged as unverified on both macOS and Windows.

---

### 🛠️ Build from Source

1. Install [Go](https://golang.org/dl/)  
2. Clone the repo:

   ```sh
   git clone https://github.com/pixcapsoft/Zippy.git
   cd Zippy
   ```

3. Build for your platform:

   **Windows:**
   ```sh
   set GOOS=windows
   set GOARCH=amd64
   go build -o zippy.exe zippy.go
   ```

   **Linux:**
   ```sh
   set GOOS=linux
   set GOARCH=amd64
   go build -o zippy zippy.go
   ```

   **macOS:**
   ```sh
   set GOOS=darwin
   set GOARCH=amd64
   go build -o zippy zippy.go
   ```

   *(For x86 systems, use `set GOARCH=386`)*

> ⚠️ **Note:** If you build Zippy from source to create your own application, you cannot use the name "Zippy" as per the [LICENSE](LICENSE). The name is copyrighted by PixCap Soft.

---

## 🧑‍💻 Usage

### Initialize a Repository
```sh
zippy init
```

- Prompts for repo name, author, and description  
- Creates `.zippy/` metadata folder and a sample `.zippyignore`

### Ignore Files/Folders
Edit `.zippyignore` to exclude files/folders from versioning (supports globs like `*.log`, `node_modules/`, etc).

### Add Files to Staging
```sh
zippy add <file/folder>
# or add everything:
zippy add .
```

### Commit a New Version
```sh
zippy commit -m "Your message" -v "v1.0"
```

### List All Versions
```sh
zippy list
```

### Restore Files or Folders
```sh
zippy restore <version>
# or restore a specific file/folder:
zippy restore <version> <path>
```

### Compare Two Versions
```sh
zippy diff <version1> <version2>
```

### Patch (Add to Existing Version)
```sh
zippy patch <version> <file/folder>
```

### Show Status
```sh
zippy status
```

### Show Version, Help, or About
```sh
zippy version
zippy help
zippy about
```

---

## 📝 Workflow Example

```sh
zippy init
zippy add .
zippy commit -m "Initial commit" -v "v1.0"
zippy add src/main.go
zippy commit -m "Add main.go" -v "v1.1"
zippy list
zippy status
zippy restore v1.0 src/main.go
zippy patch v1.1 README.md
zippy diff v1.0 v1.1
```

---

## 📂 Repository Structure

```
.zippy/                 # Zippy metadata directory
├── config.json         # Repository configuration
├── versions/           # Version metadata files (JSON)
├── storage/            # Zip files for each version
│   ├── v1.0.zip
│   └── v1.1.zip
└── stage.json          # Staging area (auto-managed)
```

---

## 💡 Tips & Best Practices

- Always use `zippy add` before `zippy commit` to stage files  
- Use meaningful version tags and commit messages  
- Regularly check `zippy status` to see what will be committed and what changed  
- Edit `.zippyignore` to avoid archiving unwanted files (logs, build artifacts, etc)  
- Restore a single file or folder from any version—great for quick rollbacks  
- Use `zippy patch` to add files to an existing version if needed

---

## ❓ FAQ

**Q: Is Zippy a replacement for Git?**  
No. Zippy is for simple, local versioning and backup. For collaboration and advanced workflows, use Git.

**Q: Can I use Zippy on any project?**  
Yes! It works with any folder, any language, any platform.

**Q: Where are my versions stored?**  
In `.zippy/storage/` as zip files, with metadata in `.zippy/versions/`.

**Q: Can I share a Zippy repo?**  
Yes! Just share the whole project folder, including `.zippy/`.

---

## 😂 Why Zippy is (Definitely) Better Than Git

- **No merge conflicts!**  
  Because you’re the only one using it (probably).
- **No remote drama.**  
  No one can force-push to your repo but you!
- **No cryptic commands.**  
  You’ll never have to Google "how to undo git rebase" again.
- **No detached HEADs.**  
  Zippy doesn’t even know what that means.
- **No `.git` folder bloat.**  
  Just a neat `.zippy/` and some zips. Marie Kondo would approve.
- **Restore a single file in one command.**  
  No checkout gymnastics required.
- **No need to remember SHA hashes.**  
  Use human tags like `v1.0`, `final-final`, or `oh-god-why`.
- **No blame.**  
  Zippy is a judgment-free zone.
- **No pull requests.**  
  If you want to argue with yourself, that’s your business.
- **No 'git push --force'.**
  - Zippy believes in second chances (and third, and fourth...)

> **Disclaimer:** If you need to collaborate, branch, or work with a team, use Git. If you want to zip and chill, use Zippy!

---

## 🧑‍💻 Contributing

Pull requests and issues are welcome! See [CONTRIBUTION.md](CONTRIBUTION.md) for details.

---

## 📜 License

Zippy License. See [LICENSE](LICENSE) for details.

---

## 🌐 More Info

- [GitHub Repository](https://github.com/pixcapsoft/Zippy)
- [Releases & Downloads](https://github.com/pixcapsoft/Zippy/releases)
- [Issues & Feedback](https://github.com/pixcapsoft/Zippy/issues)
