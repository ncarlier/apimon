# apimon

[![Build Status](https://github.com/ncarlier/apimon/actions/workflows/build.yml/badge.svg)](https://github.com/ncarlier/apimon/actions/workflows/build.yml)
[![Image size](https://images.microbadger.com/badges/image/ncarlier/apimon.svg)](https://microbadger.com/images/ncarlier/apimon)
[![Docker pulls](https://img.shields.io/docker/pulls/ncarlier/apimon.svg)](https://hub.docker.com/r/ncarlier/apimon/)

APImon is a simple tool to monitor HTTP endpoints and send metrics to a
robust monitoring platform (such as TICK, Prometheus, ELK, etc.).

![Logo](apimon.svg)

## Installation

Run the following command:

```bash
$ go get -v github.com/ncarlier/apimon
```

**Or** download the binary regarding your architecture:

```bash
$ sudo curl -s https://raw.githubusercontent.com/ncarlier/apimon/master/install.sh | bash
```

**Or** use Docker:

```bash
$ docker run -d --name=apimon \
  -v ${PWD}/apimon.yml:/etc/apimon.yml \
  ncarlier/apimon apimon -c /etc/apimon.yml
```

Here a [Docker Compose example](./docker-compose.yml).

## Configuration

The configuration is a YAML file structured like this:

```yaml
output:              # Output configuration
  traget: stdout     # Output target (default: stdout)
  format: influxdb   # Output format (default: influxdb)
proxy: ${HTTP_PROXY} # Global HTTP proxy to use (default: none)
user_agent: "Mozilla/5.0 (X11; Linux x86_64; rv:10.0) Gecko/20100101 Firefox/10.0" # UserAgent used (default: Mozilla/5.0 (compatible; APImon/1.0; +https://github.com/ncarlier/apimon))
healthcheck:         # Global healthcheck configuration
  interval: 5s       # Check interval (default: 30s)
  timeout: 2s        # Request timeout (default: 5s)
  rules:             # List of rules (default: "code = 200")
    - name: code
      spec: 200-299
monitors_files: # List of external monitors configuration (default: none)
  - "monitor_*.yml" # File to load (support wildcard)
monitors: # List of monitors configuration
  - alias: nunux-keeper-api # The name used within the produced metrics (default: the URL)
    disable: false # Disable the monitor (default: false)
    url: https://api.nunux.org/keeper/ # The URL to monitor
    method: GET # HTTP method used (default: GET)
    headers: # HTTP headers to add to the request (default: none)
      - "X-API-Key: xxx-xxx-xxx"
    body: > # HTTP request body
      { "foo": "bar" }
    proxy: http://proxy-internet.localnet:3128 # Specific HTTP proxy to use. Replaces the global configuration (default: none)
    tls:
      unsafe: true # Don't verifies the server's certificate chain and host name (default: false)
      client_cert_file: ./cert.pem # PEM encoded client certificate file (default: none)
      client_key_file: ./key.pem   # PEM encoded client key file (default: none
      ca_cert_file: ./key.pem      # PEM encoded CA's certificate file (default: none)
    healthcheck: # Monitor specific configuration. Replaces the global configuration.
      timeout: 100ms
      rules:
        - name: code
          spec: 200
  - url: https://reader.nunux.org
```

Note: As the proxy example shows, you can use environment variables in your configuration file.

### Output configuration

This configuration section allows you to set the targeted output
(i. e. your monitoring platform).

- `stdout`: Prints metrics to the STDOUT.
- `file://test.log`: Writes metrics to a log file.
- `http://...`:
  - Post metrics to an HTTP endpoint such as [InfluxDb][influxdb],
  [Elasticsearch][elasticsearch], etc. (Push Mode)
  - OR serve metrics on this HTTP endpoint if `prometheus` format is used (Pull mode)

> Note: When you are using `stdout` you should use the logging output flag
> (`-o`) in order to not mix metrics and logs outputs.

You can also choose the output format:

- `influxdb`: [InfluxDb line protocol][influxdb-line-protocol] (when using a
  collector compatible with [InfluxDB][influxdb])
- `json`: JSON format (when using [Elasticsearch][elasticsearch] like solution)
- `prometheus`: when using [Prometheus][prometheus] as collector

> Note: When using `prometheus` format the output configuration must be a local
> HTTP endpoint.
> This endpoint will be mounted by APImon and you will be able to
> [configure](./prometheus.yml) Prometheus to scrap this target.

### Healthcheck configuration

You can configure the health check in a global way (for all monitors) by placing
this section at the root level of the configuration file.
If you configure the health check of a specific monitor, the configuration
applies upon the global one.

A health check is defined like this:

- `interval`: the waiting time between checks
- `timeout`: the maximum time allowed for a check
- `rules`: the list of validation rules

A rule have the following structure:

- `name`: rule name
- `spec`: rule specification

The `name` selects the validator to be applied.
And `spec` is the configuration of the validator.

Validators are chained (using list order).
The first failed validator stops the validation chain and the monitor is
considered as DOWN.

### Available validators

Name   | Spec 
-------|------
`code` | Validates status code (ex: `200`)<br>Validates status code in a list (ex: `200,204,205`)<br>Validates status code within an interval (ex: `200-204`)
`json-path` | Validates JSON response with a [JSON path][jsonpath-syntax] expression (ex: `$.service[?(@.status == 'UP')]`)
`json-expr` | Validates JSON response with an [expression][expr-syntax] (ex: `service.status == "UP" && uptime < 100)]`)
`regexp` | Validates body response with a [regular expression][regexp-syntax] (ex: `^ok$`)
`cert` | Validates certificate expiration date is before N days (ex: `30`)

### Service Discovery and dynamic configuration

APImon can work alongside [Consul](https://www.consul.io/).
To do so you must configure the [Consul environment variables](https://www.consul.io/docs/commands/index.html#environment-variables).

Ex:

```bash
$ # Configure Consul access
$ export CONSUL_HTTP_ADDR=127.0.0.1:8500
$ # Start APImon
$ apimon
```
APImon will then register as a service and you can use the K/V store for monitors configuration.

The configuration is stored in the `apimon/monitors' key in YAML:

```yaml
- alias: google
  url: https://www.google.com
```

This is the equivalent of the `monitors' section of the configuration file.

Updating this key will reload APImon monitors.

## Usage

Type `apimon -help` to get the usage.

Basically, all you have to do is to provide the configuration file either using
the `-c` parameter or the standard input of the command.

Here come examples of possible usages:

```bash
$ # Using the default configuration file: `./apimon.yml`
$ apimon
$ # Using a specific configuration file
$ apimon -c /etc/apimon.yml
$ # Using a configuration file provided by the shell standard input
$ cat conf.yml | apimon
$ ...
```

Here an example of a Grafana dashboard displaying metrics form APImon:

![screenshot](screenshot.png)

You can find dashboard examples into the `dashboard` folder:

- [Grafana dashboard using InfluxDB](./dashboard/influxdb_sample.json)
- [Grafana dashboard using Prometheus](./dashboard/prometheus_sample.json)

## FAQ

### I am using InfluxDB with UDP

This is why `stdout` output is designed.
You can pipe the standard output to `netcat`:

```bash
$ # Sending metrics to InfluxDb using UDP
$ cat conf.yml | apimon -o output.log | nc -C -w 1 -u localhost 8125
```

### I am using Elasticsearch

You can either send the metrics (using JSON format as output configuration!)
directly to Elasticsearch:

```bash
$ apimon -c config.yml -o output.log | curl -X POST - -d @- http://localhost:9200/index/doc
$ # You can also directly use the HTTP output provider into the APIMon configuration
```

or better: use [Logstash][logstash] to collect the JSON outputs:

```
input {
  pipe {
    type => "apimon"
    debug => true
    command => "/usr/local/apimon -c /etc/apimon.yml -o /var/log/apimon.log"
  }
}
```
---

This software is under MIT License (MIT)

Copyright (c) 2018 Nicolas Carlier

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

[elasticsearch]: https://www.elastic.co/products/elasticsearch
[logstash]: https://www.elastic.co/products/logstash
[influxdb]: https://github.com/influxdata/influxdb
[influxdb-line-protocol]: https://docs.influxdata.com/influxdb/v1.4/write_protocols/line_protocol_tutorial/
[prometheus]: https://prometheus.io/
[regexp-syntax]: https://golang.org/pkg/regexp/syntax/
[expr-syntax]: https://github.com/antonmedv/expr/wiki/The-Expression-Syntax
[jsonpath-syntax]: http://goessner.net/articles/JsonPath/index.html
