| Command | Mean [ms] | Min [ms] | Max [ms] | Relative |
|:---|---:|---:|---:|---:|
| `./summer1 database` | 823.1 ± 94.0 | 655.0 | 925.9 | 3.17 ± 0.43 |
| `./summer2 generate -r database` | 1109.5 ± 10.1 | 1097.4 | 1130.6 | 4.27 ± 0.32 |
| `./summer2 generate -r --read-jobs=8 --digest-jobs=8 database` | 260.0 ± 19.2 | 234.9 | 291.4 | 1.00 |
