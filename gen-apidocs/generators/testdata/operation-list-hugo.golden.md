---
content_type: "api_reference"
title: "List Pods"
weight: 10
auto_generated: true
---

### `get` List Pods

#### HTTP Request

GET /api/v1/namespaces/{namespace}/pods


#### Path Parameters

| Name | Type | Description |
|---|---|---|
| `namespace` | string | object name and auth scope, such as for teams and users |
{class="api-reference-path-parameters"}


#### Query Parameters

| Name | Type | Description |
|---|---|---|
| `watch` | boolean | Watch for changes to the described resources. |
{class="api-reference-query-parameters"}


#### Body Parameters

| Name | Type | Description |
|---|---|---|
| `body` | Pod | Pod to create. |
{class="api-reference-request-body"}


#### Response

| Status | Description | Response |
|---|---|---|
| 200 | OK | PodList |
| 401 | Unauthorized | — |
{class="api-reference-response-codes"}

