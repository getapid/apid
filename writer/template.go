package writer

const prelude = `
{{- define "prelude" }}
  /$$$$$$  /$$$$$$$  /$$$$$$       /$$
 /$$__  $$| $$__  $$|_  $$_/      | $$
| $$  \ $$| $$  \ $$  | $$    /$$$$$$$
| $$$$$$$$| $$$$$$$/  | $$   /$$__  $$
| $$__  $$| $$____/   | $$  | $$  | $$
| $$  | $$| $$        | $$  | $$  | $$
| $$  | $$| $$       /$$$$$$|  $$$$$$$
|__/  |__/|__/      |______/ \_______/

Booting up and running tests...

{{ end }}                                           
`

const stepTemplate = `
{{- define "step" }}
    {{ .Name }}
    {{- range $i, $check := .PassedChecks }}
        {{ green "+" }} {{ $check }}
    {{- end }}

    {{- range $i, $check := .FailedChecks }}
        {{ red "o" }} {{ $check }}
    {{- end }}  
{{ end -}}
`

const conclusion = `
{{- define "conclusion" -}}
{{- $successSteps := red .SuccessSteps }}
{{- $failedSteps := red .FailedSteps }}
specs passed: {{ green .SuccessSteps }}
specs failed: {{ red .FailedSteps }}
{{ end }}
`

const schema = prelude + stepTemplate + conclusion + `
{{ bold .Name -}}
{{- range $i, $step := .Steps }}
  {{- template "step" $step -}}
{{- end -}}
`
