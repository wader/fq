package postgres_test

import (
	"context"
	"github.com/wader/fq/internal/difftest"
	"github.com/wader/fq/internal/script"
	"os"
	"strconv"
	"testing"

	_ "github.com/wader/fq/format/all"
	"github.com/wader/fq/pkg/interp"
)

// get testdata dir from:
// git clone https://github.com/pnsafonov/testdata.git

// 1 GB postgres heap file:
// https://github.com/pnsafonov/fq_testdata_postgres14/raw/master/16397
// git clone https://github.com/pnsafonov/fq_testdata_postgres14.git
//
// time fq -d pgheap -o flavour=postgres14 ".Pages[0].PageHeaderData.pd_linp[0, 1, 2, -1] | tovalue" 16397
// Killed
//
// real    0m50.794s
// user    1m11.962s
// sys     0m8.994s
//
// sudo dmesg | tail -2
// [193541.830725] oom-kill:constraint=CONSTRAINT_NONE,nodemask=(null),cpuset=user.slice,mems_allowed=0,global_oom,task_memcg=/user.slice/user-1000.slice/session-1.scope,task=fq,pid=454783,uid=1000
// [193541.830748] Out of memory: Killed process 454783 (fq) total-vm:31508780kB, anon-rss:26629332kB, file-rss:272kB, shmem-rss:0kB, UID:1000 pgtables:58860kB oom_score_adj:0

// to make mem, cpu profiling:
// go test -cpuprofile "cpu.prof" -memprofile "mem.prof" -bench .
//
// go tool pprof mem.prof
// top20
// q
//
// go tool pprof cpu.prof
// top20
// q

// run only postgres tests
func TestFQTests(t *testing.T) {
	testPath(t, interp.DefaultRegistry)
}

func testPath(t *testing.T, registry *interp.Registry) {
	difftest.TestWithOptions(t, difftest.Options{
		Path:        "testdata",
		Pattern:     "*.fqtest",
		ColorDiff:   os.Getenv("DIFF_COLOR") != "",
		WriteOutput: os.Getenv("WRITE_ACTUAL") != "",
		Fn: func(t *testing.T, path, input string) (string, string, error) {
			t.Parallel()

			b, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			c := script.ParseCases(string(b))
			c.Path = path

			for _, p := range c.Parts {
				cr, ok := p.(*script.CaseRun)
				if !ok {
					continue
				}

				t.Run(strconv.Itoa(cr.LineNr)+"/"+cr.Command, func(t *testing.T) {
					cr.WasRun = true

					i, err := interp.New(cr, registry)
					if err != nil {
						t.Fatal(err)
					}

					err = i.Main(context.Background(), cr.Stdout(), "testversion")
					if err != nil {
						if ex, ok := err.(interp.Exiter); ok { //nolint:errorlint
							cr.ActualExitCode = ex.ExitCode()
						}
					}
				})
			}

			return path, c.ToActual(), nil
		},
	})
}
