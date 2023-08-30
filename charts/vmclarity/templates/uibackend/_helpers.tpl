{{/*
Base name for the uibackend components
*/}}
{{- define "vmclarity.uibackend.name" -}}
{{- printf "%s-uibackend" (include  "vmclarity.names.fullname" .) | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Kubernetes standard labels
*/}}
{{- define "vmclarity.uibackend.labels.standard" -}}
{{ include "vmclarity.labels.standard" . }}
app.kubernetes.io/component: uibackend
{{- end -}}

{{/*
Labels to use on deploy.spec.selector.matchLabels and svc.spec.selector
*/}}
{{- define "vmclarity.uibackend.labels.matchLabels" -}}
{{ include "vmclarity.labels.matchLabels" . }}
app.kubernetes.io/component: uibackend
{{- end -}}
