{{/*
Expand the name of the chart.
*/}}
{{- define "release-tracker.name" -}}
{{- .Chart.Name }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "release-tracker.labels" -}}
app.kubernetes.io/name: {{ include "release-tracker.name" . }}
app.kubernetes.io/version: {{ .Chart.AppVersion }}
{{- end }}
