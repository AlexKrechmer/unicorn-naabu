package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

// Color codes
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Cyan   = "\033[36m"
	Purple = "\033[35m"
)

// ==== PRINT BANNERS ====
func printBanners(target string) {
	unicornArt := `
⠀⠀⠀⠀⠀⠑⢦⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠙⢷⣦⣀⠀⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⢿⣷⣿⣾⣿⣧⣄⠀⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣰⣿⣿⣿⣿⣿⣿⣿⣇⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⢀⣴⣿⣿⣿⣿⣿⣿⣿⣿⣿⣥⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠸⣿⠟⠉⠉⢹⣿⣿⣿⣿⣿⣿⣀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢠⣿⣿⣿⣿⣿⣿⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣠⣿⣿⣿⣿⣿⣿⣿⣿⡏⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣠⣿⣿⣿⣿⣿⣿⣿⣿⣿⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⢀⣠⣶⣶⣾⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣷⡄⠀⠀⠀⠀⠀⠀⠀⠀
⢀⣴⠿⠛⠉⢸⡏⠁⠉⠙⠛⠻⢿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣶⣄⡀⠀⠀⠀⠀⠀
⠉⠉⠀⠀⠀⢸⡇⠀⠀⠀⠀⠀⠀⠙⢿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣦⡀⠀⠀⠀
⠀⠀⠀⠀⠀⠈⠿⠀⠀⠀⠀⠀⠀⠀⠀⠙⢿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣧⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠉⠛⠻⢿⣿⣿⣿⣿⣿⣿⣧⡀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠸⣿⣿⣿⣿⣿⠟⢿⣷⡄
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢹⣿⣿⡟⠀⢠⣾⣿⣿
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠹⣿⣿⣀⣾⣿⡿⠃
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢠⣿⣿⣿⣿⠏
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣠⣿⣿⠻⣿⣿⡀
`
	naabuBanner := `
  ___  ___  ___ _/ /  __ __
 / _ \/ _ \/ _ \/ _ \/ // /
/_//_/\_,_/\_,_/_.__/\_,_/

                projectdiscovery.io
`
	fmt.Println(Purple + unicornArt + Reset)
	fmt.Println(Purple + naabuBanner + Reset)
	fmt.Printf("%s[*] Scanning target: %s%s\n\n", Purple, target, Reset)
}

// ==== RUN NAABU FULL TCP (JSON) ====
func runNaabuFull(target string, minRate int, useSudo bool) []string {
	fmt.Println(Cyan + "[*] Starting full Naabu sweep..." + Reset)
	openPortsFile := "open_ports.txt"
	args := []string{"-host", target, "-p-", "-json", "--rate", strconv.Itoa(minRate), "-o", openPortsFile}

	var cmd *exec.Cmd
	if useSudo {
		cmd = exec.Command("sudo", append([]string{"naabu"}, args...)...)
	} else {
		cmd = exec.Command("naabu", args...)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println(Red+"[!] Naabu scan failed:", err, Reset)
		os.Exit(1)
	}

	// Parse ports from file
	content, err := os.ReadFile(openPortsFile)
	if err != nil {
		fmt.Println(Red+"[!] Failed to read open ports file:", err, Reset)
		return nil
	}

	portSet := make(map[int]bool)
	scanner := bufio.NewScanner(strings.NewReader(string(content)))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, `"port":`) {
			parts := strings.Split(line, `"port":`)
			if len(parts) < 2 {
				continue
			}
			portPart := strings.SplitN(parts[1], ",", 2)[0]
			portPart = strings.TrimSpace(portPart)
			portNum, err := strconv.Atoi(portPart)
			if err == nil {
				portSet[portNum] = true
			}
		}
	}

	openPorts := []string{}
	for p := range portSet {
		openPorts = append(openPorts, strconv.Itoa(p))
	}
	sort.IntsFunc(openPorts, func(i, j int) bool { a, _ := strconv.Atoi(openPorts[i]); b, _ := strconv.Atoi(openPorts[j]); return a < b })
	return openPorts
}

// ==== RUN NMAP FULL SCAN ====
func runNmapFull(target string, ports []string, useSudo bool, timing int) {
	portList := "-p-"
	if len(ports) > 0 {
		portList = strings.Join(ports, ",")
		fmt.Printf(Green+"[+] Open ports found: %s%s\n", portList, Reset)
	} else {
		fmt.Println(Red + "[!] No open ports found, running full TCP scan..." + Reset)
	}

	args := []string{"-A", "-T" + strconv.Itoa(timing), portList, target}
	if useSudo {
		args = append([]string{"nmap", "--privileged"}, args...)
	} else {
		args = append([]string{"nmap"}, args...)
	}

	fmt.Println(Cyan + "[*] Running Nmap full scan..." + Reset)
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

// ==== MAIN ====
func main() {
	target := flag.String("target", "", "Target IP or hostname")
	fullScan := flag.Bool("full", false, "Run full Naabu + Nmap + OS detection scan")
	minRate := flag.Int("min-rate", 5000, "Naabu minimum rate for speed")
	useSudo := flag.Bool("sudo", true, "Use sudo for SYN scans")
	timing := flag.Int("T", 5, "Nmap timing template (0-5)")
	flag.Parse()

	if *target == "" && len(flag.Args()) > 0 {
		*target = flag.Args()[0]
	}
	if *target == "" {
		fmt.Println(Red + "[!] Please specify a target" + Reset)
		os.Exit(1)
	}

	if *useSudo && os.Geteuid() != 0 {
		fmt.Println(Red + "[!] Root required for full scan. Run with sudo." + Reset)
		os.Exit(1)
	}

	printBanners(*target)

	if *fullScan {
		openPorts := runNaabuFull(*target, *minRate, *useSudo)
		runNmapFull(*target, openPorts, *useSudo, *timing)
		fmt.Println(Green + "[+] Full scan complete." + Reset)
		return
	}

	fmt.Println(Yellow + "[*] Use -full to run Naabu + Nmap automatically" + Reset)
}
