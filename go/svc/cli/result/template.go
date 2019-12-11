package result

const successfulStepTemplate = `
{{- define "successfulStep" -}}
{{- greenOk }}   {{ .Step.ID -}}
{{ end }}
`
const failedStepTemplate = `
{{- define "failedStep" -}}
{{- with .Step -}}
{{- redFail }} {{ .ID }}

         request: {{ .Request.Type }} {{ .Request.Endpoint }} {{ with .Request.Body }}"{{ . }}"{{ end }}
{{- end }}
         errors: 
            {{- range $key, $errString := .Valid.Errors }}
             {{ $key }}:
                {{ indent 16 $errString }}
            {{ end -}}
{{- end -}}
`

const timingsTemplate = `
{{- define "timings" }}
         DNS Lookup:         {{ time .DNSLookup }}
         TCP Connection:     {{ time .TCPConnection }}
         TLS Handshake:      {{ time .TLSHandshake }}
         Server Processing:  {{ time .ServerProcessing }}
         Content Transfer:   {{ time .ContentTransfer }}
{{ end }}
`

const closingLines = `
{{- define "closingLines" -}}
{{- $totalSteps := add .SuccessSteps .FailedSteps -}}
successful transactions: {{ printf "%d/%d" .SuccessSteps $totalSteps }}
failed transactions:     {{ printf "%d/%d" .FailedSteps $totalSteps }}
{{ end }}
`

const schema = failedStepTemplate + successfulStepTemplate + timingsTemplate + closingLines + `
{{- .Id }}:
{{- $showTimings := .ShowTimings -}}
{{- range $i, $step := .Steps -}}
    {{- if $step.OK }}
    {{ template "successfulStep" $step }}
    {{- else }}
    {{ template "failedStep" $step -}}
    {{- end -}}

    {{ if $showTimings }}
    {{ template "timings" .Timings }}
    {{ end }}
{{- end }}
`
