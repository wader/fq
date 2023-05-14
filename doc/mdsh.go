// Takes markdown on stdin and outputs same markdown with shell commands expanded
// lines inside sh sections ending with non-breaking-space " " will be kept and
// those starting with $ will also be executed as shell commands and output will
// be inserted.

// <pre sh>
// # comment
// $ echo test
// will be replaced
// </pre>
// Becomes:
// <pre sh>
// # comment
// $ echo test
// test
// </pre>
//
// ```sh (sh)
// # comment
// $ echo test
// will be replaced
// ```
// Becomes:
// ```sh (sh)
// # comment
// $ echo test
// test
// ```
//
// [echo test]: sh-start
// will be replaced
// [#]: sh-end
// Becomes:
// [echo test]: sh-start
// test
// [#]: sh-end

//nolint:gosec
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

	preShRe := regexp.MustCompile("<pre sh>")
	shRe := regexp.MustCompile("```.* \\(sh\\)")
	const nonBreakingSpace = rune(0xa0) // -> " "
	shStartRe := regexp.MustCompile(`\[(.*)\]: sh-start`)
	shEnd := "[#]: sh-end"

	for {
		l, ok := nextLine()
		if !ok {
			break
		}

		preShReMatches := preShRe.MatchString(l)
		shReMatches := shRe.MatchString(l)
		if preShReMatches || shReMatches {
			fmt.Println(l)
			for {
				l, ok := nextLine()
				if !ok || ((shReMatches && l == "```") || preShReMatches && l == "</pre>") {
					fmt.Println(l)
					break
				}

				rl := []rune(l)
				if len(rl) >= 1 && rl[len(rl)-1] == nonBreakingSpace {
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
