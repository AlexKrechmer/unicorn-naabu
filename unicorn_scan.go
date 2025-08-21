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
	"time"
)

// ==== COLOR CODES ====
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Cyan   = "\033[36m"
	Purple = "\033[35m"
	Gray   = "\033[37m"
)

func main() {
	// ==== FLAG HANDLING ====
	target := flag.String("target", "", "Target IP or hostname")
	ports := flag.String("p", "", "Ports to scan (e.g., 80,443 or 1-65535)")
	fullTCP := flag.Bool("full-tcp", false, "Scan all TCP ports (1-65535)")
	useSudo := flag.Bool("sudo", false, "Run Nmap with sudo")
	flag.Parse()

	// Positional argument fallback
	if *target == "" && len(flag.Args()) > 0 {
		*target = flag.Args()[0]
	}

	if *target == "" {
		fmt.Println(Red + "[!] Please specify a target with -target or as the first argument" + Reset)
		os.Exit(1)
	}

	// Determine port list
	portList := *ports
	if *fullTCP {
		portList = "1-65535"
	} else if portList == "" {
		portList = "80,443"
	}

	printBanners(*target)
	spinnerDone := startSpinner(*target)

	openPorts := runNaabu(*target, portList, *fullTCP)
	spinnerDone <- struct{}{} // stop spinner
	fmt.Println("\rScan complete.                 ")

	if len(openPorts) == 0 {
		fmt.Printf("%s[!] No open ports found, defaulting to %s%s\n", Red, portList, Reset)
		openPorts = strings.Split(portList, ",")
	}

	// Deduplicate & sort
	portMap := make(map[string]struct{})
	for _, p := range openPorts {
		portMap[p] = struct{}{}
	}
	var sortedPorts []string
	for p := range portMap {
		sortedPorts = append(sortedPorts, p)
	}
	sort.Slice(sortedPorts, func(i, j int) bool {
		a, _ := strconv.Atoi(sortedPorts[i])
		b, _ := strconv.Atoi(sortedPorts[j])
		return a < b
	})

	portList = strings.Join(sortedPorts, ",")
	if len(portList) > 200 {
		portList = "1-65535" // fallback for too many ports
	}

	fmt.Printf("%s[Naabu] Open ports: %s%s\n\n", Purple, portList, Reset)
	runNmap(*target, portList, *fullTCP, *useSudo)
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
⠀⠀⠀⢀⣠⣶⣶⣾⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣷⡄⠀⠀⠀⠀⠀⠀⠀⠀
⢀⣴⠿⠛⠉⢸⡏⠁⠉⠙⠛⠻⢿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣶⣄⡀⠀⠀⠀⠀⠀
⠉⠉⠀⠀⠀⢸⡇⠀⠀⠀⠀⠀⠀⠙⢿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣦⡀⠀⠀⠀
⠀⠀⠀⠀⠀⠈⠿⠀⠀⠀⠀⠀⠀⠀⠀⠙⢿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣧⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠉⠛⠻⢿⣿⣿⣿⣿⣿⣿⣧⡀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠸⣿⣿⣿⣿⣿⠟⢿⣷⡄
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢹⣿⣿⡟⠀⢠⣾⣿⣿
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠹⣿⣿⣀⣾⣿⡿⠃
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢠⣿⣿⣿⣿⠏⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣠⣿⣿⠻⣿⣿⡀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠐⠋⣹⣿⠃⠀⠈⣿⣿⣴⠇
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠠⣾⠟⠀⠀⠀⠀⠘⠉⠛⠀
`

	// Naabu banner ASCII
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

// ==== SPINNER ====
func startSpinner(target string) chan struct{} {
	done := make(chan struct{})
	go func() {
		spin := []rune{'|', '/', '-', '\\'}
		i := 0
		for {
			select {
			case <-done:
				fmt.Printf("\r%s\r", strings.Repeat(" ", 50))
				return
			default:
				fmt.Printf("\r%sScanning %s... %c%s", Yellow, target, spin[i%len(spin)], Reset)
				i++
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
	return done
}

// ==== RUN NAABU ====
func runNaabu(target, ports string, fullTCP bool) []string {
	args := []string{"-silent", "-host", target, "-no-ping"}
	if fullTCP {
		args = append(args, "-p-")
	} else {
		args = append(args, "-ports", ports)
	}

	// Corrected line: use the spread operator to unpack the args slice
	cmd := exec.Command("naabu", args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(Red + "[!] Error capturing Naabu output: " + err.Error() + Reset)
		return nil
	}
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		fmt.Println(Red + "[!] Error starting Naabu: " + err.Error() + Reset)
		return nil
	}

	openPorts := []string{}
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		port := ""
		if strings.Contains(line, ":") {
			parts := strings.Split(line, ":")
			port = strings.TrimSpace(parts[len(parts)-1])
		} else {
			port = line
		}
		if _, err := strconv.Atoi(port); err == nil {
			openPorts = append(openPorts, port)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(Red + "[!] Error reading Naabu output: " + err.Error() + Reset)
	}

	cmd.Wait()
	return openPorts
}

// ==== RUN NMAP ====
func runNmap(target, portList string, fullTCP, useSudo bool) {
	args := []string{"-sC", "-sV", "-Pn"}
	if fullTCP || portList == "1-65535" {
		args = append(args, "-p-", target)
	} else {
		args = append(args, "-p", portList, target)
	}

	fmt.Println(Cyan + "[*] Running Nmap..." + Reset)

	var nmapCmd *exec.Cmd
	if useSudo {
		allArgs := append([]string{"nmap"}, args...)
		nmapCmd = exec.Command("sudo", allArgs...)
	} else {
		// Corrected line: use the spread operator to unpack the args slice
		nmapCmd = exec.Command("nmap", args...)
	}
	nmapCmd.Stdout = os.Stdout
	nmapCmd.Stderr = os.Stderr
	if err := nmapCmd.Run(); err != nil {
		fmt.Println(Red + "[!] Error running Nmap: " + err.Error() + Reset)
	}
}
