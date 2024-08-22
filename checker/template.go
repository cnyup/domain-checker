package checker

var DingTemplate = `{
"msgtype": "markdown",
"markdown": {
"title": "nginx domain check cert alert",
"text": "###  **主机名**: {{ .EcsInfo.Name }}  
###  **内网IP**:  {{ .EcsInfo.LanIp }}  
###  **外网IP**:  {{ .EcsInfo.WanIp }}  
{{ if not .ThresholdDomain }}{{ else }}
___________________________  
#### **触发告警阈值域名**:  
{{ range $val := .ThresholdDomain -}}  
- {{ $val.DomainName }}  还有 <font color=FF0000> {{ $val.ExpiredDays }} </font> 天过期  
{{ end -}}  
##### 上述域名请提前更换证书{{ end }}  
{{ if not .ExpireDomain }}
{{ else }}  
___________________________  
#### **失效域名**:  
{{ range $val := .ExpireDomain -}}> **{{ $val.DomainName }}**
{{ end -}}  
> ##### <font color=FF0000> 上述域名已经过期，请确认并进行后续处理  </font> {{ end }} 
"}}`
