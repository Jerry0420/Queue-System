import * as httpTools from './base'

const createSessionWithSSE = (sessionToken: string): EventSource => {
    const route = "/sessions/sse".concat("?session_token=", sessionToken)
    const sse = new EventSource(httpTools.generateURL(route))
    // handle sse events outside.
    // sse.onmessage = (event) => JSON.stringify(JSON.parse(event.data))
    // sse.onopen = (event) => {}
    // sse.onerror = (event) => {}
    return sse
}

const scanSession = (sessionId: string, storeId: number): Promise<any> => {
    const route = "/sessions/".concat(sessionId)
    const jsonBody: string = JSON.stringify({
        "store_id": storeId,
    })
    return fetch(
        httpTools.generateURL(route), { 
            method: httpTools.HTTPMETHOD.PUT,
            headers: {...httpTools.CONTENT_TYPE_JSON, ...httpTools.generateAuth(sessionId, false)},
            body: jsonBody
        }
    )
      .then(response => response.json())
      .then(jsonResponse => {
          console.log(jsonResponse)
          return jsonResponse
      })
      .catch(error => {
          console.error(error)
          throw new Error("scannSession error")  
      })
}

export {
    createSessionWithSSE,
    scanSession
}