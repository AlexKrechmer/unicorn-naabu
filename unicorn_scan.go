package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"sync"
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

// ==== STRUCTS FOR JSON PARSING ====
type NaabuResult struct {
	Ip   string `json:"ip"`
	Port struct {
		Port  int    `json:"port"`
		Proto string `json:"proto"`
	} `json:"port"`
}

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

// ==== MAIN ====
func main() {
	target := flag.String("target", "", "Target IP or hostname")
	fullTCP := flag.Bool("full-tcp", true, "Scan all TCP ports (1-65535) by default")
	fullScan := flag.Bool("full", false, "Run full Naabu + Nmap + OS detection scan")
	useSudo := flag.Bool("sudo", true, "Use sudo for SYN scans by default")
	timing := flag.Int("T", 5, "Nmap timing template (0-5)")
	minRate := flag.Int("min-rate", 5000, "Naabu minimum rate for speed")
	saveFile := flag.String("save", "", "Optional file to save Naabu/Nmap outputs")
	flag.Parse()

	if *target == "" && len(flag.Args()) > 0 {
		*target = flag.Args()[0]
	}
	if *target == "" {
		fmt.Println(Red + "[!] Please specify a target with -target or as the first argument" + Reset)
		os.Exit(1)
	}

	if *useSudo && os.Geteuid() != 0 {
		fmt.Println(Red + "[!] SYN scans require sudo/root. Please run this script with sudo." + Reset)
		os.Exit(1)
	}

	printBanners(*target)

	if *fullScan {
		runFullScan(*target, *fullTCP, *useSudo, *minRate, *timing, *saveFile)
		return
	}

	openPorts := runNaabuLiveAsync(*target, *fullTCP, *useSudo, *minRate, *saveFile)

	if len(openPorts) > 0 {
		sortPorts(openPorts)
		portList := strings.Join(openPorts, ",")
		fmt.Printf("%s[Naabu] Open ports found: %s%s\n\n", Purple, portList, Reset)
		runNmapColor(*target, portList, *useSudo, *timing, *saveFile)
	} else {
		fmt.Printf("%s[!] No open ports found, escalating straight to full TCP Nmap scan...%s\n", Red, Reset)
		runNmapFull(*target, *useSudo, *timing, *saveFile)
	}

	fmt.Println(Green + "[+] Scan summary complete." + Reset)
}

// ==== NAABU LIVE ASYNC (JSON PARSE + live printing) ====
func runNaabuLiveAsync(target string, fullTCP, useSudo bool, minRate int, saveFile string) []string {
	args := []string{"-target", target, "-json", "--rate", strconv.Itoa(minRate)}
	if fullTCP {
		args = append(args, "-p-")
	}

	var cmd *exec.Cmd
	if useSudo {
		cmd = exec.Command("sudo", append([]string{"naabu"}, args...)...)
	} else {
		cmd = exec.Command("naabu", args...)
	}

	fmt.Println(Cyan + "[*] Starting Naabu sweep..." + Reset)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(Red+"[!] Failed to get stdout pipe:", err, Reset)
		return nil
	}
	cmd.Stderr = cmd.Stdout

	if err := cmd.Start(); err != nil {
		fmt.Println(Red+"[!] Failed to start Naabu:", err, Reset)
		return nil
	}

	scanner := bufio.NewScanner(stdout)
	openPorts := []string{}
	portSet := make(map[int]bool)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var result map[string]interface{}
		if err := json.Unmarshal([]byte(line), &result); err != nil {
			// ignore non-JSON lines
			continue
		}

		// check if port info exists
		if portObj, ok := result["port"].(map[string]interface{}); ok {
			if portVal, ok := portObj["port"].(float64); ok {
				port := int(portVal)
				if !portSet[port] {
					portSet[port] = true
					openPorts = append(openPorts, strconv.Itoa(port))
					fmt.Printf(Green+"[+] Open port found: %d%s\n", port, Reset)
				}
			}
		}
	}

	cmd.Wait()
	return openPorts
}

// ==== SORT PORTS ====
func sortPorts(ports []string) {
	sort.Slice(ports, func(i, j int) bool {
		a, _ := strconv.Atoi(ports[i])
		b, _ := strconv.Atoi(ports[j])
		return a < b
	})
}

// ==== RUN NMAP WITH COLORS (for found ports) ====
func runNmapColor(target, ports string, useSudo bool, timing int, saveFile string) {
	args := []string{"-sV", "-T" + strconv.Itoa(timing), "-p", ports, target}
	if useSudo {
		args = append([]string{"nmap", "--privileged"}, args...)
	} else {
		args = append([]string{"nmap"}, args...)
	}
	runCommand(args, saveFile, "[Nmap] Service/Version scan results:")
}

// ==== RUN FULL NMAP TCP (fallback or -full mode) ====
func runNmapFull(target string, useSudo bool, timing int, saveFile string) {
	args := []string{"-sV", "-T" + strconv.Itoa(timing), "-p-", target}
	if useSudo {
		args = append([]string{"nmap", "--privileged"}, args...)
	} else {
		args = append([]string{"nmap"}, args...)
	}
	runCommand(args, saveFile, "[Nmap] Full TCP scan results:")
}

// ==== RUN FULL NAABU + NMAP + OS DETECTION ====
func runFullScan(target string, fullTCP, useSudo bool, minRate, timing int, saveFile string) {
	fmt.Println(Yellow + "[*] Running full scan mode: Naabu + Nmap + OS detection" + Reset)

	openPorts := runNaabuLiveAsync(target, fullTCP, useSudo, minRate, saveFile)

	if len(openPorts) > 0 {
		sortPorts(openPorts)
		portList := strings.Join(openPorts, ",")
		fmt.Printf("%s[Naabu] Open ports found: %s%s\n\n", Purple, portList, Reset)
		// Add -A for OS detection
		args := []string{"-A", "-T" + strconv.Itoa(timing), "-p", portList, target}
		if useSudo {
			args = append([]string{"nmap", "--privileged"}, args...)
		} else {
			args = append([]string{"nmap"}, args...)
		}
		runCommand(args, saveFile, "[Nmap] Aggressive scan with OS detection:")
	} else {
		fmt.Printf("%s[!] No open ports from Naabu. Running full Nmap -A scan...%s\n", Red, Reset)
		args := []string{"-A", "-T" + strconv.Itoa(timing), "-p-", target}
		if useSudo {
			args = append([]string{"nmap", "--privileged"}, args...)
		} else {
			args = append([]string{"nmap"}, args...)
		}
		runCommand(args, saveFile, "[Nmap] Full -A TCP scan results:")
	}
}

// ==== GENERIC COMMAND RUNNER (shared by Nmap funcs) ====
func runCommand(args []string, saveFile, banner string) {
	fmt.Println(Cyan + banner + Reset)

	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(Red+"[!] Command failed:", err, Reset)
		return
	}

	fmt.Println(string(output))

	if saveFile != "" {
		f, err := os.OpenFile(saveFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err == nil {
			defer f.Close()
			f.WriteString("\n" + banner + "\n")
			f.Write(output)
		} else {
			fmt.Println(Red+"[!] Failed to write to save file:", err, Reset)
		}
	}
}
