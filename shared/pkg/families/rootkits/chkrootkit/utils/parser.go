package utils

import "fmt"
import "os"
import "log"
import "strings"
import "bufio"
import "bytes"
import "encoding/json"

func SplitFuncSeparator(sep string) func(data []byte, atEOF bool) (advance int, token []byte, err error) {
	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}

		strData := string(data)

		checkingIndex := strings.Index(strData[1:], sep)

		if checkingIndex != -1 {
			//log.Print("Checking: ", string(data[:checkingIndex+1]))
			return checkingIndex+1, data[:checkingIndex+1], nil
		}

		if atEOF {
			return len(data), data, nil
		}

		return 0, nil, nil
	}
}

// This list comes from the chkrootkit source code
var applicationRootkits = []string{
	"amd", "basename", "biff", "chfn", "chsh", "cron", "crontab", "date", "du", "dirname",
	"echo", "egrep", "env", "find", "fingerd", "gpm", "grep", "hdparm", "su", "ifconfig",
	"inetd", "inetdconf", "identd", "init", "killall", "", "ldsopreload", "login", "ls",
	"lsof", "mail", "mingetty", "netstat", "named", "passwd", "pidof", "pop2", "pop3",
	"ps", "pstree", "rpcinfo", "rlogind", "rshd", "slogin", "sendmail", "sshd", "syslogd",
	"tar", "tcpd", "tcpdump", "top", "telnetd", "timed", "traceroute", "vdir", "w", "write",
}

type Rootkit struct {
	RkType string
	Rkname string
	Message string
	Infected bool
}

func ParseChkrootkitOutput(chkrootkitOutput string) []Rootkit {
	report, err := os.Open("chkrootkit_output")
	if err != nil {
		log.Fatal(err)
	}
	defer report.Close()

	scanner := bufio.NewScanner(report)
	scanner.Split(SplitFuncSeparator("Checking"))

	rootkits := []Rootkit{}

	for scanner.Scan() {
		line := scanner.Text()

		if !strings.HasPrefix(line, "Checking `") {
			// probably should error
			continue
		}
		line = line[len("Checking `"):]

		testname, result, found := strings.Cut(line, "'... ")
		if !found {
			// probably should error
			continue
		}
		result = strings.TrimSpace(result)

		if Contains(applicationRootkits, "testname") {
			rootkits = append(rootkits, Rootkit{
				RkType: "APPLICATION",
				Rkname: "UNKNOWN",
				Message: fmt.Sprintf("Application %q %s", testname, result),
				Infected: result == "INFECTED",
			})
		} else if testname == "aliens" {
			rootkits = append(rootkits, processAliensToRootkits(result)...)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	data, err := json.MarshalIndent(rootkits, "", "    ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(data))
}

func processAliensToRootkits(result string) []Rootkit {
	scanner := bufio.NewScanner(bytes.NewBufferString(result))
	scanner.Split(SplitFuncSeparator("Searching"))

	rootkits := map[string]Rootkit{}

	for scanner.Scan() {
		line := scanner.Text()

		if !strings.HasPrefix(line, "Searching for ") {
			// probably should error
			continue
		}
		line = line[len("Searching for "):]

		name, result, found := strings.Cut(line, "...")
		if !found {
			// probably should error
			continue
		}

		name = strings.TrimSpace(name)
		name = strings.TrimSuffix(name, "default files")
		name = strings.TrimSpace(name)
		name = strings.TrimSuffix(name, "default dir")
		name = strings.TrimSpace(name)
		name = strings.TrimSuffix(name, "default files and dirs")
		name = strings.TrimSpace(name)
		name = strings.TrimSuffix(name, "default files and dir")
		name = strings.TrimSpace(name)
		name = strings.TrimSuffix(name, "files and dirs")
		name = strings.TrimSpace(name)
		name = strings.TrimSuffix(name, "modules")
		name = strings.TrimSpace(name)
		name = strings.TrimSuffix(name, "defaults")
		name = strings.TrimSpace(name)
		name = strings.TrimSuffix(name, ", it may take a while")
		name = strings.TrimSpace(name)
		name = strings.TrimSuffix(name, "logs")
		name = strings.TrimSpace(name)

		name = strings.TrimSuffix(name, "'s")
		name = strings.TrimSpace(name)

		result = strings.TrimSpace(result)

		rkType := "UNKNOWN"
		if strings.Contains(strings.ToLower(name), "lkm") {
			rkType = "KERNEL"
		}

		infected := false
		message := result
		if result != "nothing found" {
			infected = true
		}

		var rk Rootkit
		rk, ok := rootkits[name]
		if !ok {
			rk = Rootkit{
				RkType: rkType,
				Rkname: name,
				Message: message,
				Infected: infected,
			}
		} else {
			if !rk.Infected {
				rk.Infected = infected
			}

			rk.Message = fmt.Sprintf("%s %s", rootkits[name].Message, message)
		}
		rootkits[name] = rk
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	rootkitArray := []Rootkit{}
	for _, rootkit := range rootkits {
		rootkitArray = append(rootkitArray, rootkit)
	}
	return rootkitArray
}

// Index returns the index of the first occurrence of v in s,
// or -1 if not present.
func Index[E comparable](s []E, v E) int {
	for i, vs := range s {
		if v == vs {
			return i
		}
	}
	return -1
}

// Contains reports whether v is present in s.
func Contains[E comparable](s []E, v E) bool {
	return Index(s, v) >= 0
}

