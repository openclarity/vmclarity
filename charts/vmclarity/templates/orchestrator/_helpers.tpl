{{/*
Base name for the orchestrator components
*/}}
{{- define "vmclarity.orchestrator.name" -}}
{{- printf "%s-orchestrator" (include  "vmclarity.names.fullname" .) | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Kubernetes standard labels
*/}}
{{- define "vmclarity.orchestrator.labels.standard" -}}
{{ include "vmclarity.labels.standard" . }}
app.kubernetes.io/component: orchestrator
{{- end -}}

{{/*
Labels to use on deploy.spec.selector.matchLabels and svc.spec.selector
*/}}
{{- define "vmclarity.orchestrator.labels.matchLabels" -}}
{{ include "vmclarity.labels.matchLabels" . }}
app.kubernetes.io/component: orchestrator
{{- end -}}
