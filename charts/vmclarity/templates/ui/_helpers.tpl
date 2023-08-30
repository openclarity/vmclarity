{{/*
Base name for the ui components
*/}}
{{- define "vmclarity.ui.name" -}}
{{- printf "%s-ui" (include  "vmclarity.names.fullname" .) | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Kubernetes standard labels
*/}}
{{- define "vmclarity.ui.labels.standard" -}}
{{ include "vmclarity.labels.standard" . }}
app.kubernetes.io/component: ui
{{- end -}}

{{/*
Labels to use on deploy.spec.selector.matchLabels and svc.spec.selector
*/}}
{{- define "vmclarity.ui.labels.matchLabels" -}}
{{ include "vmclarity.labels.matchLabels" . }}
app.kubernetes.io/component: ui
{{- end -}}
