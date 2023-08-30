{{/*
Base name for the grypeServer components
*/}}
{{- define "vmclarity.grypeServer.name" -}}
{{ include  "vmclarity.names.fullname" . }}-grype-server
{{- end -}}

{{/*
Kubernetes standard labels
*/}}
{{- define "vmclarity.grypeServer.labels.standard" -}}
{{ include "vmclarity.labels.standard" . }}
app.kubernetes.io/component: grype-server
{{- end -}}

{{/*
Labels to use on deploy.spec.selector.matchLabels and svc.spec.selector
*/}}
{{- define "vmclarity.grypeServer.labels.matchLabels" -}}
{{ include "vmclarity.labels.matchLabels" . }}
app.kubernetes.io/component: grype-server
{{- end -}}
