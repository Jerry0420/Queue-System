- [API Document of Queue System](#api-document-of-queue-system)
  - [TODO](#todo)
  - [Basics](#basics)
  - [Rate Limiting](#rate-limiting)
    - [Rules](#rules)
  - [Status Codes](#status-codes)
    - [HTTP Status Codes](#http-status-codes)
    - [Custome Status Codes](#custome-status-codes)

---

# API Document of Queue System

## TODO
* authentication
* resources
* sse

## Basics
* Prefix all URLs with `/api/v1`, except the URLs for health checking.
* `v1` is a version tag and `/api` is for backend api.

## Rate Limiting
* Rate Limiting can be used to mitigate DDoS Attacks.
### Rules
* The number of concurrent connections allowed from a single IP address is 5.
* The number of requests accepted from a given IP each second is 5.
* Return HTTP Status code 429, when breaking the rate limiting rules.

## Status Codes

### HTTP Status Codes
<table>
  <tr>
    <th>HTTP Status Code</th>
    <th>Description</th>
  </tr>

  <tr>
    <td>200</td>
    <td>OK</td>
  </tr>

  <tr>
    <td>204</td>
    <td>OK, but no content to return.</td>
  </tr>

  <tr>
    <td>400</td>
    <td>Bad Request <br> The request was unacceptable, often due to missing or invalid a required parameter. </td>
  </tr>

  <tr>
    <td>401</td>
    <td>Unauthorized <br> No valid jwt token provided. </td>
  </tr>

  <tr>
    <td>404</td>
    <td>Not Found <br> The requested resource doesn't exist. </td>
  </tr>

  <tr>
    <td>405</td>
    <td>Method Not Allowed</td>
  </tr>

  <tr>
    <td>409</td>
    <td>Conflict <br>The requested resource conflicts with another resource. </td>
  </tr>

  <tr>
    <td>429</td>
    <td>Too Many Requests <br> Check <a href="#rules">Rate Limiting Rules</a></td>
  </tr>

  <tr>
    <td>500</td>
    <td>Internal Server Error <br> Something went wrong on server, grpc server or database. </td>
  </tr>
</table>

### Custome Status Codes
* `Custome Status Codes` are custom error codes of this system, and the first three digits is `HTTP Status Code`.
* The format of all non-2xx response, will be like: 
```json
{
    "error_code": {{custom_status_code}}
}

ex: 
{
    "error_code": 40401
} 
``` 

<table>
  <tr>
    <th>Code - 400XX</th>
    <th>Description</th>
  </tr>

  <tr>
    <td>40001</td>
    <td>Lack of required params.</td>
  </tr>

  <tr>
    <td>40002</td>
    <td>Length of password is not appropriate.</td>
  </tr>

  <tr>
    <td>40003</td>
    <td>The incoming password is not equal to the original password.</td>
  </tr>

  <tr>
    <td>40004</td>
    <td>Wrong params</td>
  </tr>

  <tr>
    <td>40005</td>
    <td>The count of customers is more than 5.</td>
  </tr>

  <tr>
    <td>40006</td>
    <td>The requested timezone is not exist.</td>
  </tr>
<table>

<table>
  <tr>
    <th>Code - 401XX</th>
    <th>Description</th>
  </tr>
  <tr>
    <td>40101</td>
    <td>Fail to parse jwt token.</td>
  </tr>

  <tr>
    <td>40102</td>
    <td>Lack of jwt token.</td>
  </tr>

  <tr>
    <td>40103</td>
    <td>Other error of parsing jwt token.</td>
  </tr>

  <tr>
    <td>40104</td>
    <td>Jwt token expired.</td>
  </tr>

  <tr>
    <td>40105</td>
    <td>Lack of store session.</td>
  </tr>
<table>

<table>
  <tr>
    <th>Code - 404XX</th>
    <th>Description</th>
  </tr>
  <tr>
    <td>40401</td>
    <td>Unsupported url route.</td>
  </tr>

  <tr>
    <td>40402</td>
    <td>Store not exist.</td>
  </tr>

  <tr>
    <td>40403</td>
    <td>Sign_key not exist.</td>
  </tr>

  <tr>
    <td>40404</td>
    <td>Store_session not exist.</td>
  </tr>

  <tr>
    <td>40405</td>
    <td>Customer not exist.</td>
  </tr>
<table>

<table>
  <tr>
    <th>Code - 405XX</th>
    <th>Description</th>
  </tr>
  <tr>
    <td>40501</td>
    <td>Method Not Allowed</td>
  </tr>
<table>

<table>
  <tr>
    <th>Code - 409XX</th>
    <th>Description</th>
  </tr>
  <tr>
    <td>40901</td>
    <td>store already exist. (not exceed 24 hrs)</td>
  </tr>

  <tr>
    <td>40902</td>
    <td>Sign_key already exist.</td>
  </tr>

  <tr>
    <td>40903</td>
    <td>Store_session already exist.</td>
  </tr>
<table>

<table>
  <tr>
    <th>Code - 500XX</th>
    <th>Description</th>
  </tr>
  <tr>
    <td>50001</td>
    <td>Other internal server error.</td>
  </tr>

  <tr>
    <td>50002</td>
    <td>Unexpected database error.</td>
  </tr>

  <tr>
    <td>50003</td>
    <td>The client not support flushing.</td>
  </tr>

  <tr>
    <td>50004</td>
    <td>Unexpected grpc server error.</td>
  </tr>
</table>