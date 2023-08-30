{{/*
Base name for the apiserver components
*/}}
{{- define "vmclarity.apiserver.name" -}}
{{ include  "vmclarity.names.fullname" . }}-apiserver
{{- end -}}

{{/*
Kubernetes standard labels
*/}}
{{- define "vmclarity.apiserver.labels.standard" -}}
{{ include "vmclarity.labels.standard" . }}
app.kubernetes.io/component: apiserver
{{- end -}}

{{/*
Labels to use on deploy.spec.selector.matchLabels and svc.spec.selector
*/}}
{{- define "vmclarity.apiserver.labels.matchLabels" -}}
{{ include "vmclarity.labels.matchLabels" . }}
app.kubernetes.io/component: apiserver
{{- end -}}
