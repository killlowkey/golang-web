# 压力测试

## superbenchmarker 压力测试
superbenchmarker 工具在 Windows、Linux、Mac 可用
1. 安装教程：https://github.com/aliostad/SuperBenchmarker

### usage
```shell
PS F:\liunx\wrk-4.2.0> sb
SuperBenchmarker 4.5.1
Copyright (C) 2022 Ali Kheyrollahi
ERROR(S):
Required option 'u, url' is missing.

  -c, --concurrency            (Default: 1) Number of concurrent requests

  -n, --numberOfRequests       (Default: 100) Total number of requests

  -N, --numberOfSeconds        Number of seconds to run the test. If specified, -n will be ignored.

  -y, --delayInMillisecond     (Default: 0) Delay in millisecond

  -u, --url                    Required. Target URL to call. Can include placeholders.

  -m, --method                 (Default: GET) HTTP Method to use

  -t, --template               Path to request template to use

  -p, --plugin                 Name of the plugin (DLL) to provide a type for replacing placeholders or one for
                               overriding status code.Must reside in the same folder.

  -l, --logfile                Path to the log file storing run stats

  -f, --file                   Path to CSV file providing replacement values for the test

  -a, --TSV                    If you provide a tab-separated-file (TSV) with -f option instead of CSV

  -d, --dryRun                 Runs a single dry run request to make sure all is good

  -e, --timedField             Designates a datetime field in data. If set, requests will be sent according to order
                               and timing of records.

  -g, --TlsVersion             Version of TLS used. Accepted values are 0, 1, 2 and 3 for TLS 1.0, TLS 1.1 and TLS 1.2
                               and SSL3, respectively

  -v, --verbose                Provides verbose tracing information

  -b, --tokeniseBody           Tokenise the body

  -k, --cookies                Outputs cookies

  -x, --useProxy               Whether to use default browser proxy. Useful for seeing request/response in Fiddler.

  -q, --onlyRequest            In a dry-run (debug) mode shows only the request.

  -h, --headers                Displays headers for request and response.

  -z, --saveResponses          saves responses in -w parameter or if not provided in\response_<timestamp>

  -w, --responsesFolder        folder to save responses in if and only if -z parameter is set

  -?, --help                   Displays this help.

  -C, --dontcap                Don't Cap to 50 characters when Logging parameters

  -R, --responseRegex          Regex to extract from response. If it has groups, it retrieves the last group.

  -j, --jsonCount              Captures number of elements under the path e.g. root/leaf1/leaf2 finds count of leaf2
                               children - stores in the log as another parameter. If the array is at the root of the
                               JSON, use space: -j ' '

  -W, --warmUpPeriod           (Default: 0) Number of seconds to gradually increase number of concurrent users. Warm-up
                               calls do not affect stats.

  -P, --reportSliceSeconds     (Default: 3) Number of seconds as interval for reporting slices. E.g. if chosen as 5,
                               report charts have 5 second intervals.

  -F, --reportFolder           Name of the folder where report files get stored. By default it is in
                               yyyy-MM-dd_HH-mm-ss.ffffff of the start time.

  -B, --dontBrowseToReports    By default it, sb opens the browser with the report of the running test. If specified,
                               it wil not browse.

  -U, --shuffleData            If specified, shuffles the dataset provided by -f option.

  --help                       Display this help screen.

  --version                    Display version information.

Error Code: CommandLine.MissingRequiredOptionError
```

### sb 测试
`sb -u "http://example.com" -c 4 -n 1800`
1. -u：测试连接
2. -c：并发数量
3. -n：总共请求数量

`sb -u "http://example.com" -c 4 -N 1800`
1. -N：测试秒数，这里为 1800 秒

`sb -u "http://example.com/api/car/123" -n 1000 -m DELETE`
1. -m：请求方法，这里为 delete

## wrk 压力测试
### wrk 下载与安装
> wrk 工具仅在 Linux 和 Mac 上可用，以下安装在 WSL 环境下进行
1. 下载源代码：https://github.com/wg/wrk/tags
2. 解压 tar.gz：`tar -xzvf wrk-4.2.0.tar.gz`
3. 编译源代码：`make`
4. 当前目录下创建 bin 目录，并将编译后 wrk 文件移动到 bin 目录下
5. 安装 wrk：编辑 `/etc/profile` 尾部增加 `export PATH=$PATH:/mnt/f/liunx/wrk-4.2.0/bin`

### wrk usage
```shell
root@DESKTOP-HNGKMS5:/mnt/f/liunx/wrk-4.2.0# wrk
Usage: wrk <options> <url>
  Options:
    -c, --connections <N>  Connections to keep open
    -d, --duration    <T>  Duration of test
    -t, --threads     <N>  Number of threads to use

    -s, --script      <S>  Load Lua script file
    -H, --header      <H>  Add header to request
        --latency          Print latency statistics
        --timeout     <T>  Socket/request timeout
    -v, --version          Print version details

  Numeric arguments may include a SI unit (1k, 1M, 1G)
  Time arguments may include a time unit (2s, 2m, 2h)
```

### wrk 压力测试
执行该命令：`wrk -t 16 -c 100 -d 30s http://127.0.0.1:8080/hello`
1. -t：线程数量
2. -c：并发连接数
3. -d：压测时间
```shell
root@DESKTOP-HNGKMS5:/mnt/f/liunx/wrk-4.2.0# wrk -t 16 -c 100 -d 30s http://127.0.0.1:8080/hello
Running 30s test @ http://127.0.0.1:8080/hello
  16 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     3.14ms    3.65ms  49.61ms   86.44%
    Req/Sec     2.70k   355.55     4.25k    68.06%
  1290206 requests in 30.02s, 157.50MB read
Requests/sec:  42977.23
Transfer/sec:      5.25MB
```

### 参考资料
1. [wrk github](https://github.com/wg/wrk)
2. [如何通过压力测试得知系统的QPS](https://zhuanlan.zhihu.com/p/512589767)