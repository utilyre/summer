| Command | Mean [s] | Min [s] | Max [s] | Relative |
|:---|---:|---:|---:|---:|
| `./summer1 database` | 1.197 ± 0.011 | 1.184 | 1.224 | 4.57 ± 0.30 |
| `./summer2 generate -r database` | 1.098 ± 0.006 | 1.091 | 1.111 | 4.19 ± 0.27 |
| `./summer2 generate -r --read-jobs=8 --digest-jobs=8 database` | 0.262 ± 0.017 | 0.239 | 0.282 | 1.00 |
