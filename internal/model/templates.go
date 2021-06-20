package model

import "text/template"

var OrgResultTemplate = template.Must(template.New("orgResultTemplate").Parse(orgResultTemplate))

const orgResultTemplate = `
_id                 {{ index .Organization "_id" }}
url                 {{ index .Organization "url" }}
external_id         {{ index .Organization "external_id" }}
name                {{ index .Organization "name" }}
domain_names        {{ index .Organization "domain_names" }}
created_at          {{ index .Organization "created_at" }}
details             {{ index .Organization "details" }}
shared_tickets      {{ index .Organization "shared_tickets" }}
tags                {{ index .Organization "tags" }}
{{ range $index, $element := .UserNames }}user_{{ $index }}              {{ $element }}
{{ end }}{{ range $index, $element := .TicketSubjects }}ticket_{{ $index }}            {{ $element }}
{{ end }}`

var UserResultTemplate = template.Must(template.New("userResultTemplate").Parse(userResultTemplate))

const userResultTemplate = `
_id                 {{ index .User "_id" }}
url                 {{ index .User "url" }}
external_id         {{ index .User "external_id" }}
name                {{ index .User "name" }}
alias               {{ index .User "alias" }}
created_at          {{ index .User "created_at" }}
active              {{ index .User "active" }}
verified            {{ index .User "verified" }}
shared              {{ index .User "shared" }}
locale              {{ index .User "locale" }}
timezone            {{ index .User "timezone" }}
last_login_at       {{ index .User "last_login_at" }}
email               {{ index .User "email" }}
phone               {{ index .User "phone" }}
signature           {{ index .User "signature" }}
organization_id     {{ index .User "organization_id" }}
tags                {{ index .User "tags" }}
suspended           {{ index .User "suspended" }}
role                {{ index .User "role" }}
organization_name   {{ .OrganizationName }}
{{ range $index, $element := .TicketSubjects }}ticket_{{ $index }}            {{ $element }}
{{ end }}`

var TicketResultTemplate = template.Must(template.New("ticketResultTemplate").Parse(ticketResultTemplate))

const ticketResultTemplate = `
_id                 {{ index .Ticket "_id" }}
url                 {{ index .Ticket "url" }}
external_id         {{ index .Ticket "external_id" }}
created_at          {{ index .Ticket "created_at" }}
type                {{ index .Ticket "type" }}
subject             {{ index .Ticket "subject" }}
description         {{ index .Ticket "description" }}
priority            {{ index .Ticket "priority" }}
status              {{ index .Ticket "status" }}
submitter_id        {{ index .Ticket "submitter_id" }}
assignee_id         {{ index .Ticket "assignee_id" }}
organization_id     {{ index .Ticket "organization_id" }}
tags                {{ index .Ticket "tags" }}
has_incidents       {{ index .Ticket "has_incidents" }}
due_at              {{ index .Ticket "due_at" }}
via                 {{ index .Ticket "via" }}
organization_name	{{ .OrganizationName }}
`
