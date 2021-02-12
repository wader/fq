// Takes markdown on stdin and outputs same markdown with shell commands expanded
//
// ```sh (exec)
// # comment
// $ echo test
// ```
// Becomes:
// ```sh (exec)
// # comment
// $ echo test
// test
// ```
//
// [echo test]: sh-start
//
// anything here
//
// [#]: sh-end
// Becomes:
// [echo test]: sh-start
//
// test
//
// [#]: sh-end
package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	nextLine := func() (string, bool) {
		ok := scanner.Scan()
		return scanner.Text(), ok
	}

	preExecRe := regexp.MustCompile("<pre exec>")
	execRe := regexp.MustCompile("```.* \\(exec\\)")
	const nonBreakingSpace = rune(0xa0) // -> " "
	shStartRe := regexp.MustCompile(`\[(.*)\]: sh-start`)
	shEnd := "[#]: sh-end"

	for {
		l, ok := nextLine()
		if !ok {
			break
		}

		preExecReMatches := preExecRe.MatchString(l)
		execReMatches := execRe.MatchString(l)
		if preExecReMatches || execReMatches {
			fmt.Println(l)
			for {
				l, ok := nextLine()
				if !ok || ((execReMatches && l == "```") || preExecReMatches && l == "</pre>") {
					fmt.Println(l)
					break
				}

				if len(l) >= 2 && []rune(l)[len(l)-1] == nonBreakingSpace {
					fmt.Println(l)
					if strings.HasPrefix(l, "$") {
						cmd := exec.Command("sh", "-c", l[1:len(l)-2])
						o, _ := cmd.CombinedOutput()
						fmt.Print(string(o))
					}
				}
			}
		} else if sm := shStartRe.FindStringSubmatch(l); sm != nil {
			fmt.Println(l)
			fmt.Println()
			for {
				l, ok := nextLine()
				if !ok || l == shEnd {
					break
				}
			}
			cmd := exec.Command("sh", "-c", sm[1])
			o, _ := cmd.CombinedOutput()
			fmt.Print(string(o))
			fmt.Println()
			fmt.Println(shEnd)
		} else {
			fmt.Println(l)
		}
	}
}
