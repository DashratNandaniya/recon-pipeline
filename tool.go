
package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
)

// -------------------- HELPERS --------------------

func runWithRetry(cmd, logFile string, retries int) {
	for i := 0; i < retries; i++ {
		c := exec.Command("bash", "-c", cmd)
		out, err := c.CombinedOutput()
		if err == nil {
			return
		}
		if i == retries-1 && logFile != "" {
			f, _ := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			defer f.Close()
			fmt.Fprintf(f, "[FAILED] %s\n%s\n", cmd, string(out))
		}
	}
}

func countLines(path string) int {
	f, err := os.Open(path)
	if err != nil {
		return 0
	}
	defer f.Close()

	count := 0
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		if strings.TrimSpace(sc.Text()) != "" {
			count++
		}
	}
	return count
}

func fileHasContent(path string) bool {
	return countLines(path) > 0
}

func readLines(path string) []string {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()

	var out []string
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		if l := strings.TrimSpace(sc.Text()); l != "" {
			out = append(out, l)
		}
	}
	return out
}

// -------------------- PHASE 1 --------------------

func subdomainEnum(domain, outDir string) {
	fmt.Println("\n[1] Subdomain Enumeration")

	temp := outDir + "/temp.txt"
	final := outDir + "/subs.txt"
	log := outDir + "/debug.log"

	cmds := []string{
		fmt.Sprintf("subfinder -silent -d %s >> %s 2>/dev/null", domain, temp),
		fmt.Sprintf("assetfinder --subs-only %s >> %s 2>/dev/null", domain, temp),
		fmt.Sprintf("findomain -t %s -q >> %s 2>/dev/null", domain, temp),
		fmt.Sprintf("amass enum -passive -d %s >> %s 2>/dev/null", domain, temp),
	}

	for _, c := range cmds {
		before := countLines(temp)
		runWithRetry(c, log, 2)
		after := countLines(temp)

		if after == before {
			fmt.Println("[!] No results:", c)
		}
	}

	runWithRetry(fmt.Sprintf("sort -u %s > %s", temp, final), log, 1)
	fmt.Println("✔ Subdomains:", countLines(final))
}

// -------------------- PHASE 2 --------------------

func dnsResolve(outDir string) {
	fmt.Println("\n[2] DNS Resolution")

	input := outDir + "/subs.txt"
	output := outDir + "/active.txt"
	log := outDir + "/debug.log"

	cmd := fmt.Sprintf(
		"dnsx -l %s -silent -resp 2>/dev/null | sed 's/\\x1b\\[[0-9;]*m//g' | awk '{print $1}' | sort -u > %s",
		input, output,
	)

	runWithRetry(cmd, log, 2)
	fmt.Println("✔ Active:", countLines(output))
}

// -------------------- PHASE 3 --------------------

func httpProbe(outDir string) {
	fmt.Println("\n[3] HTTP Probing")

	input := outDir + "/active.txt"
	output := outDir + "/httpx.txt"
	log := outDir + "/debug.log"

	if countLines(input) == 0 {
		fmt.Println("[!] No input for httpx")
		return
	}

	cmd := fmt.Sprintf(
		"httpx-toolkit -l %s -silent -status-code -title -threads 30 -o %s 2>/dev/null",
		input, output,
	)

	runWithRetry(cmd, log, 2)
	fmt.Println("✔ HTTP:", countLines(output))
}

// -------------------- PHASE 4 --------------------

func endpointHarvest(outDir string) {
	fmt.Println("\n[4] Endpoint Harvest")

	input := outDir + "/active.txt"
	output := outDir + "/endpoints.txt"
	log := outDir + "/debug.log"

	hosts := readLines(input)

	var wg sync.WaitGroup
	sem := make(chan struct{}, 10)

	for _, host := range hosts {
		wg.Add(1)
		sem <- struct{}{}

		go func(h string) {
			defer wg.Done()
			defer func() { <-sem }()

			cmd := fmt.Sprintf(
				"(gau %s 2>/dev/null; echo %s | waybackurls 2>/dev/null; katana -u https://%s -silent 2>/dev/null) | sort -u >> %s",
				h, h, h, output,
			)

			runWithRetry(cmd, log, 2)
		}(host)
	}

	wg.Wait()

	runWithRetry(fmt.Sprintf("sort -u %s -o %s", output, output), log, 1)
	fmt.Println("✔ Endpoints:", countLines(output))
}

// -------------------- PHASE 5 --------------------

func runSubzy(outDir string) {
	fmt.Println("\n[5] Subdomain Takeover")

	input := outDir + "/subs.txt"
	output := outDir + "/takeover.txt"

	runWithRetry(fmt.Sprintf(
		"subzy run --targets %s --vuln --hide_fails --output %s 2>/dev/null",
		input, output,
	), "", 1)

	if fileHasContent(output) {
		fmt.Println("🔥 Takeover found!")
	} else {
		fmt.Println("✔ No takeover")
	}
}

// -------------------- PHASE 6 --------------------

func runEyewitness(outDir string) {
	fmt.Println("\n[6] Screenshots")

	input := outDir + "/active.txt"
	tmp := outDir + "/targets.txt"

	cmd1 := fmt.Sprintf(`awk '{print "https://"$0}' %s > %s`, input, tmp)
	runWithRetry(cmd1, "", 1)

	cmd2 := fmt.Sprintf(
		"eyewitness --web -f %s -d %s/eye --timeout 40 --threads 10 --no-prompt 2>/dev/null",
		tmp, outDir,
	)

	runWithRetry(cmd2, "", 1)
	fmt.Println("✔ Screenshots done")
}

// -------------------- MAIN --------------------

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run tool.go example.com")
		return
	}

	domain := os.Args[1]
	outDir := "output/" + domain
	os.MkdirAll(outDir, 0755)

	fmt.Println("Recon started for:", domain)

	subdomainEnum(domain, outDir)
	dnsResolve(outDir)
	httpProbe(outDir)
	endpointHarvest(outDir)
	runSubzy(outDir)
	runEyewitness(outDir)

	fmt.Println("\n🎯 DONE")
}

