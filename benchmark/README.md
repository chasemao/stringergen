Here is benchmark result.

```sh
$ go test -bench=. -benchmem -cpu=1,2
goos: linux
goarch: amd64
pkg: github.com/chasemao/stringergen/benchmark
cpu: Intel(R) Xeon(R) Gold 6231C CPU @ 3.20GHz
BenchmarkSimpleGet                                      1000000000               0.2905 ns/op          0 B/op          0 allocs/op
BenchmarkSimpleGet-2                                    1000000000               0.2795 ns/op          0 B/op          0 allocs/op
BenchmarkSimpleSprintfWithPlus                           1279094               951.4 ns/op           200 B/op          9 allocs/op
BenchmarkSimpleSprintfWithPlus-2                         1424458               835.0 ns/op           200 B/op          9 allocs/op
BenchmarkSimpleSpew                                       929484              1340 ns/op             512 B/op         16 allocs/op
BenchmarkSimpleSpew-2                                     954130              1233 ns/op             512 B/op         16 allocs/op
BenchmarkSimpleJson                                      1834711               653.8 ns/op           288 B/op          5 allocs/op
BenchmarkSimpleJson-2                                    1882478               632.9 ns/op           288 B/op          5 allocs/op
BenchmarkSimpleJsonIter                                  2362458               509.2 ns/op           200 B/op          4 allocs/op
BenchmarkSimpleJsonIter-2                                2489205               474.7 ns/op           200 B/op          4 allocs/op
BenchmarkSimpleCustomJson                                1608458               739.7 ns/op           224 B/op          6 allocs/op
BenchmarkSimpleCustomJson-2                              1797193               658.9 ns/op           224 B/op          6 allocs/op
BenchmarkSimpleCustomJsonStringBuilder                   2054240               590.9 ns/op           328 B/op          8 allocs/op
BenchmarkSimpleCustomJsonStringBuilder-2                 2077462               577.7 ns/op           328 B/op          8 allocs/op
BenchmarkSimpleStringer                                  1546446               767.3 ns/op           192 B/op          4 allocs/op
BenchmarkSimpleStringer-2                                1796024               671.6 ns/op           192 B/op          4 allocs/op
BenchmarkComplicatedGet                                   357172              3120 ns/op            2984 B/op         63 allocs/op
BenchmarkComplicatedGet-2                                 414536              2724 ns/op            2984 B/op         63 allocs/op
BenchmarkComplicatedSprintfPlus                           114205             10624 ns/op            4640 B/op         97 allocs/op
BenchmarkComplicatedSprintfPlus-2                         134042              8939 ns/op            4640 B/op         97 allocs/op
BenchmarkComplicatedSprintfPlusLog                         27451             39451 ns/op            4663 B/op         98 allocs/op
BenchmarkComplicatedSprintfPlusLog-2                       32506             39672 ns/op            4666 B/op         98 allocs/op
BenchmarkComplicatedSpew                                   17832             65997 ns/op           13696 B/op        657 allocs/op
BenchmarkComplicatedSpew-2                                 19659             60778 ns/op           13701 B/op        657 allocs/op
BenchmarkComplicatedSpewLog                                 7267            166080 ns/op           13740 B/op        658 allocs/op
BenchmarkComplicatedSpewLog-2                               7366            160207 ns/op           13755 B/op        658 allocs/op
BenchmarkComplicatedJson                                   51376             22112 ns/op           10592 B/op         82 allocs/op
BenchmarkComplicatedJson-2                                 56222             21182 ns/op           10595 B/op         82 allocs/op
BenchmarkComplicatedJsonLog                                11506            105253 ns/op           10625 B/op         83 allocs/op
BenchmarkComplicatedJsonLog-2                              10000            115040 ns/op           10670 B/op         83 allocs/op
BenchmarkComplicatedJsonIter                               50907             22759 ns/op           11489 B/op         90 allocs/op
BenchmarkComplicatedJsonIter-2                             51739             21011 ns/op           11499 B/op         90 allocs/op
BenchmarkComplicatedCustomJson                             26995             43327 ns/op           18810 B/op        298 allocs/op
BenchmarkComplicatedCustomJson-2                           31198             38279 ns/op           18817 B/op        298 allocs/op
BenchmarkComplicatedCustomJsonStringBuilder                27650             41771 ns/op           37100 B/op        465 allocs/op
BenchmarkComplicatedCustomJsonStringBuilder-2              27621             40257 ns/op           37126 B/op        465 allocs/op
BenchmarkComplicatedStringer                               19887             58549 ns/op           18202 B/op        380 allocs/op
BenchmarkComplicatedStringer-2                             23569             50527 ns/op           18208 B/op        380 allocs/op
BenchmarkComplicatedStringerRetype                         17168             68573 ns/op           19986 B/op        540 allocs/op
BenchmarkComplicatedStringerRetype-2                       20367             59711 ns/op           19991 B/op        540 allocs/op
PASS
ok      github.com/chasemao/stringergen/benchmark       64.309s
```