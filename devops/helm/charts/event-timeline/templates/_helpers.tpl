{{- define "event-timeline.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end }}

{{- define "event-timeline.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name (default .Chart.Name .Values.nameOverride) | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}

{{- define "event-timeline.labels" -}}
helm.sh/chart: {{ .Chart.Name }}-{{ .Chart.Version | replace "+" "_" }}
{{ include "event-timeline.selectorLabels" . }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{- define "event-timeline.selectorLabels" -}}
app.kubernetes.io/name: {{ include "event-timeline.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{- define "event-timeline.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "event-timeline.fullname" .) .Values.serviceAccount.name }}
{{- else }}
default
{{- end }}
{{- end }}
