package main

import (
	"bytes"
	"log"
	"text/template"
)

var tem = func() *template.Template {

	t, err := template.New("gen").Parse(stringTemplate)
	if err != nil {
		log.Fatalf("parse template error: %s", err)
	}
	return t
}()

const stringTemplate = `
{
  "annotations": {
    "list": [
      {
        "builtIn": 1,
        "datasource": "-- Grafana --",
        "enable": true,
        "hide": true,
        "iconColor": "rgba(0, 211, 255, 1)",
        "name": "Annotations & Alerts",
        "type": "dashboard"
      }
    ]
  },
  "editable": true,
  "gnetId": null,
  "graphTooltip": 0,
  "links": [],
  "panels": [
    {{.panels}}
  ],
  "refresh": "10s",
  "schemaVersion": 16,
  "style": "dark",
  "tags": [],
  "templating": {
    "list": []
  },
  "time": {
    "from": "now-1h",
    "to": "now"
  },
  "timepicker": {
    "refresh_intervals": [
      "5s",
      "10s",
      "30s",
      "1m",
      "5m",
      "15m",
      "30m",
      "1h",
      "2h",
      "1d"
    ],
    "time_options": [
      "5m",
      "15m",
      "1h",
      "6h",
      "12h",
      "24h",
      "2d",
      "7d",
      "30d"
    ]
  },
  "timezone": "",
  "title": "{{.dashBoardTitle}}",
  "uid": "null",
  "version": 0
}
`

type panelMeta struct {
	id          int
	namespace   string
	serviceName string
	dimension   string
	endpoint    string
}

func (m *panelMeta) Qps() string {
	var buf bytes.Buffer
	data := map[string]interface{}{
		"id":          m.id,
		"namespace":   m.namespace,
		"serviceName": m.serviceName,
		"dimension":   m.dimension,
		"endpoint":    m.endpoint,
	}
	qpsTem.Execute(&buf, data)
	return buf.String()
}

func (m *panelMeta) Duration() string {
	var buf bytes.Buffer
	data := map[string]interface{}{
		"id":          m.id,
		"namespace":   m.namespace,
		"serviceName": m.serviceName,
		"dimension":   m.dimension,
		"endpoint":    m.endpoint,
	}
	durationTem.Execute(&buf, data)
	return buf.String()
}

var qpsTem = func() *template.Template {

	t, err := template.New("qps").Parse(qpsPanel)
	if err != nil {
		log.Fatalf("parse qps template error: %s", err)
	}
	return t
}()

const qpsPanel = `{
	"aliasColors": {},
	"bars": false,
	"dashLength": 10,
	"dashes": false,
	"datasource": "Prometheus",
	"fill": 1,
	"gridPos": {
		"h": 9,
		"w": 12,
		"x": 0,
		"y": 0
	},
    "id": {{.id}},
	"legend": {
		"avg": false,
		"current": false,
		"max": false,
		"min": false,
		"show": true,
		"total": false,
		"values": false
	},
	"lines": true,
	"linewidth": 1,
	"links": [],
	"nullPointMode": "null as zero",
	"percentage": false,
	"pointradius": 5,
	"points": false,
	"renderer": "flot",
	"seriesOverrides": [],
	"spaceLength": 10,
	"stack": false,
	"steppedLine": false,
	"targets": [{
		"aggregator": "sum",
		"downsampleAggregator": "avg",
		"downsampleFillPolicy": "none",
		"expr": "sum(irate({{.namespace}}_{{.serviceName}}_latency_bucket{method=\"{{.dimension}}\",service=\"{{.serviceName}}\",endpoint!=\"/metrics/\"}[1m])) by (endpoint)",
		"format": "time_series",
		"intervalFactor": 1,
		"legendFormat": "{{.endpoint}}",
		"refId": "A"
	}],
	"thresholds": [],
	"timeFrom": null,
	"timeRegions": [],
	"timeShift": null,
	"title": "{{.dimension}}QPS",
	"tooltip": {
		"shared": true,
		"sort": 2,
		"value_type": "individual"
	},
	"type": "graph",
	"xaxis": {
		"buckets": null,
		"mode": "time",
		"name": null,
		"show": true,
		"values": []
	},
	"yaxes": [{
			"format": "reqps",
			"label": null,
			"logBase": 1,
			"max": null,
			"min": null,
			"show": true
		},
		{
			"format": "short",
			"label": null,
			"logBase": 1,
			"max": null,
			"min": null,
			"show": true
		}
	],
	"yaxis": {
		"align": false,
		"alignLevel": null
	}
}`

var durationTem = func() *template.Template {

	t, err := template.New("duration").Parse(durationPanel)
	if err != nil {
		log.Fatalf("parse duration template error: %s", err)
	}
	return t
}()

const durationPanel = `{
  "aliasColors": {},
  "bars": false,
  "dashLength": 10,
  "dashes": false,
  "datasource": "Prometheus",
  "fill": 1,
  "gridPos": {
	"h": 9,
	"w": 12,
	"x": 12,
	"y": 0
  },
  "id": {{.id}},
  "legend": {
	"avg": false,
	"current": false,
	"max": false,
	"min": false,
	"show": true,
	"total": false,
	"values": false
  },
  "lines": true,
  "linewidth": 1,
  "links": [],
  "nullPointMode": "null as zero",
  "percentage": false,
  "pointradius": 5,
  "points": false,
  "renderer": "flot",
  "seriesOverrides": [],
  "spaceLength": 10,
  "stack": false,
  "steppedLine": false,
  "targets": [
	{
	  "aggregator": "sum",
	  "downsampleAggregator": "avg",
	  "downsampleFillPolicy": "none",
	  "expr": "histogram_quantile(0.99,sum(rate({{.namespace}}_{{.serviceName}}_latency_bucket{method=\"{{.dimension}}\",service=\"{{.serviceName}}\",endpoint!=\"/metrics/\"}[30s]))by (le,endpoint))",
	  "format": "time_series",
	  "intervalFactor": 1,
	  "legendFormat": "{{.endpoint}}",
	  "refId": "A"
	}
  ],
  "thresholds": [],
  "timeFrom": null,
  "timeRegions": [],
  "timeShift": null,
  "title": "{{.dimension}}延时",
  "tooltip": {
	"shared": true,
	"sort": 2,
	"value_type": "individual"
  },
  "type": "graph",
  "xaxis": {
	"buckets": null,
	"mode": "time",
	"name": null,
	"show": true,
	"values": []
  },
  "yaxes": [
	{
	  "format": "s",
	  "label": null,
	  "logBase": 1,
	  "max": null,
	  "min": null,
	  "show": true
	},
	{
	  "format": "short",
	  "label": null,
	  "logBase": 1,
	  "max": null,
	  "min": null,
	  "show": true
	}
  ],
  "yaxis": {
	"align": false,
	"alignLevel": null
  }
}
`
