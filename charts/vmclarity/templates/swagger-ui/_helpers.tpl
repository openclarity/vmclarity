{{/*
Base name for the swaggerUI components
*/}}
{{- define "vmclarity.swaggerUI.name" -}}
{{- printf "%s-swagger-ui" (include  "vmclarity.names.fullname" .) | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Kubernetes standard labels
*/}}
{{- define "vmclarity.swaggerUI.labels.standard" -}}
{{ include "vmclarity.labels.standard" . }}
app.kubernetes.io/component: swagger-ui
{{- end -}}

{{/*
Labels to use on deploy.spec.selector.matchLabels and svc.spec.selector
*/}}
{{- define "vmclarity.swaggerUI.labels.matchLabels" -}}
{{ include "vmclarity.labels.matchLabels" . }}
app.kubernetes.io/component: swagger-ui
{{- end -}}
