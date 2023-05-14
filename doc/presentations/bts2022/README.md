fq presentation from Binary Tools Summit 2022 https://binary-tools.net/

[fq-bts2022-v1.pdf](fq-bts2022-v1.pdf)

Will update with link to recording when available.

Was done at the time of ~fq 0.0.5, things might have changed since.

How to build:

```
go install golang.org/x/tools/cmd/present
present -notes -content doc/presentations/bts2022 -base ~/go/pkg/mod/golang.org/x/tools@v0.1.9/cmd/present
```

```
./usage.sh | ansisvg > usage.svg
```

Export to PDF via browser.

