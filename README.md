# Recon Orchestrator

![Go](https://img.shields.io/badge/Go-1.20+-00ADD8?logo=go)
![License](https://img.shields.io/badge/License-MIT-green)
![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20macOS-blue)

A lightweight Go wrapper that chains together industry-standard open-source recon tools into a single automated pipeline — subdomain enumeration, DNS resolution, HTTP probing, endpoint harvesting, subdomain takeover detection, and screenshotting — for **authorized security assessments, bug bounty engagements, and educational research**.

This tool does not perform scanning itself. It orchestrates external CLI utilities via `os/exec`, feeding the output of each phase into the next.

---

## 🛠️ Pipeline Overview

| Phase | Function | Tools Used | Output File |
|---|---|---|---|
| 1 | Subdomain Enumeration | `subfinder`, `assetfinder`, `findomain`, `amass` | `subs.txt` |
| 2 | DNS Resolution | `dnsx` | `active.txt` |
| 3 | HTTP Probing | `httpx-toolkit` | `httpx.txt` |
| 4 | Endpoint Harvesting | `gau`, `waybackurls`, `katana` | `endpoints.txt` |
| 5 | Subdomain Takeover Check | `subzy` | `takeover.txt` |
| 6 | Screenshots | `eyewitness` | `eye/` |

Each phase writes to `output/<domain>/`, and failures in one phase are logged to `debug.log` without halting the rest of the run.

---

## 🚀 Getting Started

### Prerequisites

- Go 1.20+
- Linux (Kali, Ubuntu, Debian, Parrot, Arch, Fedora/RHEL) or macOS
- Internet access (tools query third-party sources such as crt.sh and the Wayback Machine)

### Installation

Pick the section matching your OS/package manager. Every path installs the same set of tools: Go, plus `subfinder`, `assetfinder`, `findomain`, `amass`, `dnsx`, `httpx` (aliased to `httpx-toolkit`), `gau`, `waybackurls`, `katana`, `subzy`, and `eyewitness`.

#### Debian / Ubuntu / Kali / Parrot (`apt`)

```bash
# Go
sudo apt update
sudo apt install golang-go -y
go version

# Recon tools available via apt
sudo apt install amass eyewitness -y

# Everything else via go install
go install -v github.com/projectdiscovery/subfinder/v2/cmd/subfinder@latest
go install github.com/tomnomnom/assetfinder@latest
go install -v github.com/projectdiscovery/dnsx/cmd/dnsx@latest
go install -v github.com/projectdiscovery/httpx/cmd/httpx@latest
go install github.com/lc/gau/v2/cmd/gau@latest
go install github.com/tomnomnom/waybackurls@latest
go install -v github.com/projectdiscovery/katana/cmd/katana@latest
go install -v github.com/PentestPad/subzy@latest

# findomain (no apt package — prebuilt binary)
wget https://github.com/findomain/findomain/releases/latest/download/findomain-linux.zip
unzip findomain-linux.zip
chmod +x findomain
sudo mv findomain /usr/local/bin/

# symlink httpx -> httpx-toolkit (tool.go calls it by this name)
sudo ln -s ~/go/bin/httpx /usr/local/bin/httpx-toolkit

echo 'export PATH=$PATH:~/go/bin' >> ~/.bashrc
source ~/.bashrc
```

#### Arch Linux / Manjaro (`pacman`)

```bash
# Go
sudo pacman -Syu go --noconfirm
go version

# Recon tools available via pacman/AUR
sudo pacman -S amass --noconfirm
yay -S eyewitness findomain

# go install tools (identical to above)
go install -v github.com/projectdiscovery/subfinder/v2/cmd/subfinder@latest
go install github.com/tomnomnom/assetfinder@latest
go install -v github.com/projectdiscovery/dnsx/cmd/dnsx@latest
go install -v github.com/projectdiscovery/httpx/cmd/httpx@latest
go install github.com/lc/gau/v2/cmd/gau@latest
go install github.com/tomnomnom/waybackurls@latest
go install -v github.com/projectdiscovery/katana/cmd/katana@latest
go install -v github.com/PentestPad/subzy@latest

sudo ln -s ~/go/bin/httpx /usr/local/bin/httpx-toolkit
echo 'export PATH=$PATH:~/go/bin' >> ~/.bashrc
source ~/.bashrc
```

#### Fedora / RHEL / CentOS (`dnf`)

```bash
# Go
sudo dnf install golang -y
go version

# amass (may need EPEL or Go install fallback)
sudo dnf install epel-release -y
sudo dnf install amass -y
# fallback if unavailable:
go install -v github.com/owasp-amass/amass/v4/...@master

# eyewitness — install from source (not usually packaged for Fedora)
git clone https://github.com/FortyNorthSecurity/EyeWitness.git
cd EyeWitness/Python/setup
sudo ./setup.sh

# findomain
wget https://github.com/findomain/findomain/releases/latest/download/findomain-linux.zip
unzip findomain-linux.zip
chmod +x findomain
sudo mv findomain /usr/local/bin/

# go install tools (identical to above)
go install -v github.com/projectdiscovery/subfinder/v2/cmd/subfinder@latest
go install github.com/tomnomnom/assetfinder@latest
go install -v github.com/projectdiscovery/dnsx/cmd/dnsx@latest
go install -v github.com/projectdiscovery/httpx/cmd/httpx@latest
go install github.com/lc/gau/v2/cmd/gau@latest
go install github.com/tomnomnom/waybackurls@latest
go install -v github.com/projectdiscovery/katana/cmd/katana@latest
go install -v github.com/PentestPad/subzy@latest

sudo ln -s ~/go/bin/httpx /usr/local/bin/httpx-toolkit
echo 'export PATH=$PATH:~/go/bin' >> ~/.bashrc
source ~/.bashrc
```

#### macOS (`brew`)

```bash
# Go
brew install go
go version

# Recon tools available via brew
brew install amass findomain

# go install tools
go install -v github.com/projectdiscovery/subfinder/v2/cmd/subfinder@latest
go install github.com/tomnomnom/assetfinder@latest
go install -v github.com/projectdiscovery/dnsx/cmd/dnsx@latest
go install -v github.com/projectdiscovery/httpx/cmd/httpx@latest
go install github.com/lc/gau/v2/cmd/gau@latest
go install github.com/tomnomnom/waybackurls@latest
go install -v github.com/projectdiscovery/katana/cmd/katana@latest
go install -v github.com/PentestPad/subzy@latest

# eyewitness — install from source (Homebrew formula was removed)
git clone https://github.com/FortyNorthSecurity/EyeWitness.git
cd EyeWitness/Python/setup
./setup.sh

sudo ln -s ~/go/bin/httpx /usr/local/bin/httpx-toolkit
echo 'export PATH=$PATH:~/go/bin' >> ~/.zshrc
source ~/.zshrc
```

> **Note:** `amass`, `eyewitness`, and `findomain` package availability varies by distro version. If a package manager command fails, fall back to the `go install` or source-build alternative shown for that tool.

#### Verify installation (any system)

```bash
for tool in go subfinder assetfinder findomain amass dnsx httpx-toolkit gau waybackurls katana subzy eyewitness; do
  command -v $tool >/dev/null 2>&1 && echo "✔ $tool found" || echo "✘ $tool MISSING"
done
```

---

## 📖 Usage

```bash
go run tool.go example.com
```

Or build a standalone binary first:

```bash
go build -o recon tool.go
./recon example.com
```

Results are saved to `output/example.com/`, including `subs.txt`, `active.txt`, `httpx.txt`, `endpoints.txt`, `takeover.txt`, `debug.log`, and screenshots under `eye/`.

---

## ⚙️ Behavior Notes

- **Retry logic** — every external command runs through `runWithRetry`, retrying up to a set number of times before logging failure to `debug.log`.
- **Fail-soft design** — if one phase fails or a tool is missing, the pipeline logs a warning and continues rather than aborting the whole run.
- **Concurrency** — endpoint harvesting (Phase 4) runs up to 10 hosts in parallel using goroutines and a semaphore channel.
- **No credentials required** — all tools operate against public DNS/HTTP infrastructure and public archives; no API keys are hardcoded, though some tools (e.g. `amass`, `subfinder`) can use optional API keys from their own config files for higher rate limits.

---

## ⚖️ Legal Disclaimer & Responsible Use

This project is intended **strictly** for authorized security research, bug bounty engagements, and educational purposes. By using this tool, you agree to the following:

- **Get authorization first.** Only run this toolset against domains or infrastructure you own, or for which you have explicit written permission to test (e.g., a signed engagement letter, an active bug bounty program scope, or your own organization's assets).
- **No unauthorized access.** Enumerating, probing, or crawling infrastructure you don't have permission to test may violate laws such as the U.S. Computer Fraud and Abuse Act (CFAA), the UK Computer Misuse Act, the EU's GDPR/NIS2 framework, or equivalent legislation in your jurisdiction.
- **Handle discovered data responsibly.** If you find sensitive exposures during authorized testing, follow responsible disclosure practices: report to the asset owner or the relevant bug bounty program, and never publish, sell, or redistribute discovered data.
- **No warranty.** This software is provided "as is" with no guarantee of fitness for any particular purpose. The authors are not responsible for misuse, damages, or legal consequences resulting from use of this tool.
- **You are responsible for your own actions.** The maintainers of this repository do not condone and are not liable for any illegal or unethical use of this software.

If you're unsure whether you have authorization to test a target, don't — get explicit written permission first.

## 📄 License

This project is licensed under the MIT License — see the [LICENSE](LICENSE) file for full terms.
