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

	execRe := regexp.MustCompile("```.* \\(exec\\)")
	const nonBreakingSpace = rune(0xa0) // -> " "
	shStartRe := regexp.MustCompile(`\[(.*)\]: sh-start`)
	shEnd := "[#]: sh-end"

	for {
		l, ok := nextLine()
		if !ok {
			break
		}

		if execRe.MatchString(l) {
			fmt.Println(l)
			for {
				l, ok := nextLine()
				if !ok || l == "```" {
					fmt.Println(l)
					break
				}
				if strings.HasPrefix(l, "$") {
					fmt.Println(l)
					cmd := exec.Command("sh", "-c", l[1:])
					o, _ := cmd.CombinedOutput()
					fmt.Print(string(o))
				} else if strings.HasPrefix(l, "#") || []rune(l) == nonBreakingSpace {
					// keep comments and empty lines
					fmt.Println(l)
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
