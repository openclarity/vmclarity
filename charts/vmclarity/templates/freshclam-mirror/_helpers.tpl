{{/*
Base name for the freshclamMirror components
*/}}
{{- define "vmclarity.freshclamMirror.name" -}}
{{ include  "vmclarity.names.fullname" . }}-freshclam-mirror
{{- end -}}

{{/*
Kubernetes standard labels
*/}}
{{- define "vmclarity.freshclamMirror.labels.standard" -}}
{{ include "vmclarity.labels.standard" . }}
app.kubernetes.io/component: freshclam-mirror
{{- end -}}

{{/*
Labels to use on deploy.spec.selector.matchLabels and svc.spec.selector
*/}}
{{- define "vmclarity.freshclamMirror.labels.matchLabels" -}}
{{ include "vmclarity.labels.matchLabels" . }}
app.kubernetes.io/component: freshclam-mirror
{{- end -}}
