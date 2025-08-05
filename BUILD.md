# 🛠️ Building Zippy on Linux

> This guide walks you through building Zippy from source on a Linux system.  
> It includes step-by-step instructions, screenshots, and customization tips for authors and forks.

---

## 📦 Prerequisites

Before you begin, make sure your system has:

- **Go (Golang)** installed  
  You can check by running:
  ```bash
  go version
  ```

If Go is not installed, visit [golang.org/dl](https://golang.org/dl) and follow the instructions for your distro.

---

## 🧱 Step-by-Step Build Guide

### 🖼️ **1. Clone the Repository**

Use `git clone` to download Zippy’s source code:

![Cloning Zippy](https://raw.githubusercontent.com/ranujasanmir/Zippy-IMG/main/1.png)

```bash
git clone https://github.com/pixcapsoft/Zippy.git
```

---

### 🖼️ **2. Enter the Project Directory**

Navigate into the cloned folder:

![Navigating into Zippy](https://raw.githubusercontent.com/ranujasanmir/Zippy-IMG/main/2.png)

```bash
cd Zippy
```

---

### 🖼️ **3. Verify Go Installation and Project Files**

List the contents and check your Go version:

![Listing files and checking Go](https://raw.githubusercontent.com/ranujasanmir/Zippy-IMG/main/3.png)

```bash
ls
go version
```

You should see `zippy.go` and other project files.

---

### 🖼️ **4. Build Zippy**

Use the Go compiler to build the binary:

![Building Zippy](https://raw.githubusercontent.com/ranujasanmir/Zippy-IMG/main/4.png)

```bash
go build -o zippy zippy.go
```

This creates a binary named `zippy` in the current directory.

---

### 🖼️ **5. Run Zippy**

Test the build by running:

![Running Zippy](https://raw.githubusercontent.com/ranujasanmir/Zippy-IMG/main/5.png)

```bash
./zippy version
```

You should see Zippy’s version info and author details.

---

## ✍️ Customizing Author Name

If you're building Zippy for personal use or as a fork, you must **change the name** of the binary and update the author info, as stated in the [LICENSE](LICENSE). Attribute about original author is appreciated but not required. But we recommend to add original github project link if user want to build his own one.

---

### 🖼️ **6. Edit `zippy.go`**

Open `zippy.go` in your preferred text editor and locate the author metadata:

![Editing zippy.go](https://raw.githubusercontent.com/ranujasanmir/Zippy-IMG/main/6.png)

```go
var ZippyAuthor = "PixCap Soft"
```

Change it to your name or organization:

```go
var ZippyAuthor = "YourNameHere"
```

Also, change the binary name to something other than `Zippy` to comply with the license.

---

### 🖼️ **7. Rebuild and Verify**

Build again with your custom name:

```bash
go build -o myzippy zippy.go
./myzippy version
```

![Running customized Zippy](https://raw.githubusercontent.com/ranujasanmir/Zippy-IMG/main/7.png)

You should now see your author name in the version output.

---

## 📜 License Reminder

> The name “Zippy” is copyrighted by PixCap Soft.  
> If you modify and redistribute the tool, you must rename it and update the author metadata.

---

## 🧠 Tips

- You can move your binary to `/usr/local/bin` for system-wide access:
  ```bash
  sudo cp myzippy /usr/local/bin/
  ```

- Use `sudo chmod +x /usr/local/bin/myzippy` if the binary isn’t executable.

---

## 💬 Need Help?

Open an issue on the [GitHub repo](https://github.com/pixcapsoft/Zippy/issues) or reach out to the community for support.

Happy building! 🚀