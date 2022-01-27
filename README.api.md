- [API Document of Queue System](#api-document-of-queue-system)
  - [Basics](#basics)
  - [Errors](#errors)
    - [HTTP Status Codes](#http-status-codes)
    - [Error Codes](#error-codes)

---

# API Document of Queue System

## Basics
* Prefix all URLs with `/api/v1`, except the URLs for health checking.
* `v1` is a version tag and `/api` is for backend api.

## Errors

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
    <td>Too Many Requests</td>
  </tr>

  <tr>
    <td>500</td>
    <td>Internal Server Error <br> Something went wrong on server, grpc server or database. </td>
  </tr>
</table>

### Error Codes
* `Error Codes` are custom error codes of this system, and the first three digits is `HTTP Status Code`.

<table>
  <tr>
    <th>Error Code</th>
    <th>Description</th>
  </tr>

  <tr style='border-top: 2px solid #f00;'>
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

  <tr style='border-top: 2px solid #f00;'>
    <td>40101</td>
    <td></td>
  </tr>

  <tr>
    <td>40102</td>
    <td></td>
  </tr>

  <tr>
    <td>40103</td>
    <td></td>
  </tr>

  <tr>
    <td>40104</td>
    <td></td>
  </tr>

  <tr>
    <td>40105</td>
    <td></td>
  </tr>

  <tr>
    <td>40106</td>
    <td></td>
  </tr>

  <tr style='border-top: 2px solid #f00;'>
    <td>40401</td>
    <td></td>
  </tr>

  <tr>
    <td>40402</td>
    <td></td>
  </tr>

  <tr>
    <td>40403</td>
    <td></td>
  </tr>

  <tr>
    <td>40404</td>
    <td></td>
  </tr>

  <tr>
    <td>40405</td>
    <td></td>
  </tr>

  <tr style='border-top: 2px solid #f00;'>
    <td>40501</td>
    <td></td>
  </tr>

  <tr style='border-top: 2px solid #f00;'>
    <td>40901</td>
    <td></td>
  </tr>

  <tr>
    <td>40902</td>
    <td></td>
  </tr>

  <tr>
    <td>40903</td>
    <td></td>
  </tr>

  <tr>
    <td>40904</td>
    <td></td>
  </tr>

  <tr style='border-top: 2px solid #f00;'>
    <td>50001</td>
    <td></td>
  </tr>

  <tr>
    <td>50002</td>
    <td></td>
  </tr>

  <tr>
    <td>50003</td>
    <td></td>
  </tr>

  <tr>
    <td>50004</td>
    <td></td>
  </tr>

</table>