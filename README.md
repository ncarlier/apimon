# apimon

[![Build Status](https://travis-ci.org/ncarlier/apimon.svg?branch=master)](https://travis-ci.org/ncarlier/apimon)
[![Image size](https://images.microbadger.com/badges/image/ncarlier/apimon.svg)](https://microbadger.com/images/ncarlier/apimon)
[![Docker pulls](https://img.shields.io/docker/pulls/ncarlier/apimon.svg)](https://hub.docker.com/r/ncarlier/apimon/)

APImon is a simple tool for monitoring HTTP endpoints and sending metrics to a
robust monitoring platform (such as TICK, ELK, etc.).

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
  -v ${PWD}/configuration.yml:/etc/apimon.yml \
  ncarlier/apimon apimon -c /etc/apimon.yml
```

## Configuration

The configuration is a YAML similar to this:

```yaml
output:            # Output configuration
  type: stdout     # By default "stdout" but can also be "file://test.log" or "http://localhost:8086/write?db=test"
  format: influxdb # By default "influxdb" but can also be "json"
  meta:
    influxdb_metric_name: http_health_check
user_agent: "Mozilla/5.0 (X11; Linux x86_64; rv:10.0) Gecko/20100101 Firefox/10.0"
proxy: http://proxy-internet.localnet:3128 # Global HTTP proxy to use. By default none
healthcheck:            # Global healthcheck configuration
  interval: 5s          # By default 30s
  timeout: 2s           # By default 5s
  rules:                # By default "code: 200"
    - name: code
      spec: 200-299
monitors: # Monitors configuration
  - alias: nunux-keeper-api
    url: https://api.nunux.org/keeper/ # The URL to monitor
    unsafe: true # Don't check SSL certificate. By default false
    proxy: http://proxy-internet.localnet:3128 # Specific HTTP proxy to use (overide global). By default none
    headers: # HTTP headers to add to the request
      - "X-API-Key: xxx-xxx-xxx"
    healthcheck: # Monitor specific configuration (overide global)
      timeout: 100ms
      rules:
        - name: code
          spec: 200
  - url: https://reader.nunux.org
```

### Output configuration

This configuration section allows you to configure the target output
(i. e. your monitoring platform).

- `stdout`: Prints metrics to the STDOUT.
- `file://test.log`: Writes metrics to a log file.
- `http://...`: Post metrics to an HTTP endpoint (such as [InfluxDb][influxdb],
  [Elasticsearch][elasticsearch], ...)

> Note: When you are using `stdout` you should use the logging output flag
> (`-o`) in order to not mix metrics and logs outputs.

You can also choose the output format:

- `influxdb`: [InfluxDb line protocol][influxdb-line-protocol] (when using a
  collector compatible with [InfluxDB][influxdb])
- `json`: JSON format (when using [Elasticsearch][elasticsearch] like solution)

Some output format can be configured with additional settings.
For instance the `influxdb` format can be tuned to define the metric prefix by
setting `influxdb_metric_name`.

### Healthcheck configuration

You can configure the health check in a global way (for all monitors) by placing
this section at the root level of the configuration file.
If you configure the health check of a particular monitor, the configuration
applies upon the global one.

A health check is defined like this:

- `interval`: the waiting time between checks
- `timeout`: the maximum time allowed for a check
- `rules`: all validation rules (separated by semicolon)

A rule have the following structure:

- `name`: rule name
- `spec`: rule specification

The rule's name selects the validator to be applied.
And rule 's spec is the configuration of the validator.

Validators can be chained (it is a list).
The first failed validator stops the validation chain and the monitor is
considered as DOWN.

### Available validators

Name   | Spec 
-------|------
`code` | Validates status code (ex: `200`)<br>Validates status code in a list (ex: `200,204,205`)<br>Validates status code within an interval (ex: `200-204`)
`json-path` | Validates a JSON path of the body response (ex: `$.service[?(@.status == 'UP')]`)
`regexp:` | Validates the body response with a [regular expression][regexp-syntax] (ex: `^ok$`)

## Usage

Type `apimon -help` to get the usage.

Basically, all you have to do is to provide the configuration file either using
the `-c` parameter or the standard input of the command.

Here come examples of possible usage:

```bash
$ # Using the defaul configuration file: `./configuration.yml`
$ apimon
$ # Using a specific configuration file
$ apimon -c /etc/apimon.yml
$ # Using a configuration file provided by the shell standard input
$ cat conf.yml | apimon
$ ...
```

## FAQ

### I am using Prometheus

And that's completely fine.
You can use [InfluxDB Exporter][influxdb-exporter] to accepts InfluxDB metrics
via HTTP API and exports theme via HTTP for Prometheus.

### I am using InfluxDB with UDP

This is why `stdout` output is designed.
You can pipe the standard output to `netcat`:

```bash
$ # Sending metrics to InfluxDb using UDP
$ cat conf.yml | apimon -o output.log | nc -C -w 1 -u localhost 8125
```

### I am using Elasticsearch

You can either send the metrics (with JSON format) directly to Elasticsearch:

```bash
$ apimon -c config.yml -o output.log | curl -X POST - -d @- http://localhost:9200/index/doc
$ # You can also directly use the HTTP output provider into the APIMon configuration
```

or use [Logstash][logstash] to collect the JSON outputs:

```
input {
  pipe {
    type => "apimon"
    debug => true
    command => "/usr/local/apimon -c /etc/apimon.yml -o /var/log/apimon.log"
  }
}
```

> Note: Don't forget to use JSON format as output configuration!

---

[elasticsearch]: https://www.elastic.co/products/elasticsearch
[logstash]: https://www.elastic.co/products/logstash
[influxdb]: https://github.com/influxdata/influxdb
[influxdb-line-protocol]: https://docs.influxdata.com/influxdb/v1.4/write_protocols/line_protocol_tutorial/
[influxdb-exporter]: https://github.com/prometheus/influxdb_exporter
[regexp-syntax]: https://golang.org/pkg/regexp/syntax/
