{{/*
Base name for the gateway components
*/}}
{{- define "vmclarity.gateway.name" -}}
{{ include  "vmclarity.names.fullname" . }}-gateway
{{- end -}}

{{/*
Kubernetes standard labels
*/}}
{{- define "vmclarity.gateway.labels.standard" -}}
{{ include "vmclarity.labels.standard" . }}
app.kubernetes.io/component: gateway
{{- end -}}

{{/*
Labels to use on deploy.spec.selector.matchLabels and svc.spec.selector
*/}}
{{- define "vmclarity.gateway.labels.matchLabels" -}}
{{ include "vmclarity.labels.matchLabels" . }}
app.kubernetes.io/component: gateway
{{- end -}}
