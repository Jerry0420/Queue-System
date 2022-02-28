import * as httpTools from './base'

const openStore = (email: string, password: string, name: string, timezone: string, queueNames: string[]): Promise<any> => {
    const jsonBody: string = JSON.stringify({
        "email": email,
        "password": password,
        "name": name,
        "timezone": timezone,
        "queue_names": queueNames
    })
    return fetch(
        httpTools.generateURL("/stores"), {
            method: httpTools.HTTPMETHOD.POST,
            headers: httpTools.CONTENT_TYPE_JSON,
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
          throw new Error("openStore error")  
      })
}

const signInStore = (email: string, password: string): Promise<any> => {
    const jsonBody: string = JSON.stringify({
        "email": email,
        "password": password,
    })
    return fetch(
        httpTools.generateURL("/stores/signin"), {
            method: httpTools.HTTPMETHOD.POST,
            headers: httpTools.CONTENT_TYPE_JSON,
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
          throw new Error("signInStore error")  
      })
}

const refreshToken = (): Promise<any> => {
    return fetch(
        httpTools.generateURL("/stores/token"), {
            method: httpTools.HTTPMETHOD.PUT,
        }
    )
      .then(response => response.json())
      .then(jsonResponse => {
          console.log(jsonResponse)
          return jsonResponse
      })
      .catch(error => {
          console.error(error)
          throw new Error("refreshToken error")  
      })
}

const closeStore = (storeId: number, normalToken: string): Promise<any> => {
    const route = "/stores/".concat(storeId.toString())
    return fetch(
        httpTools.generateURL(route), { 
            method: httpTools.HTTPMETHOD.DELETE,
            headers: httpTools.generateAuth(normalToken)
        }
    )
      .then(response => response.json())
      .then(jsonResponse => {
          console.log(jsonResponse)
          return jsonResponse
      })
      .catch(error => {
          console.error(error)
          throw new Error("closeStore error")  
      })
}

const forgetPassword = (email: string): Promise<any> => {
    const jsonBody: string = JSON.stringify({
        "email": email,
    })
    return fetch(
        httpTools.generateURL("/stores/password/forgot"), { 
            method: httpTools.HTTPMETHOD.POST,
            headers: httpTools.CONTENT_TYPE_JSON,
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
          throw new Error("forgetPassword error")  
      })
}

const updatePassword = (storeId: number, passwordToken: string, password: string): Promise<any> => {
    const route = "/stores/".concat(storeId.toString(), "/password")
    const jsonBody: string = JSON.stringify({
        "password_token": passwordToken,
        "password": password,
    })
    return fetch(
        httpTools.generateURL(route), { 
            method: httpTools.HTTPMETHOD.PATCH,
            headers: httpTools.CONTENT_TYPE_JSON,
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
          throw new Error("updatePassword error")  
      })
}

const getStoreInfoWithSSE = (storeId: number): EventSource => {
    const route = "/stores/".concat(storeId.toString(), "/sse")
    const sse = new EventSource(httpTools.generateURL(route))
    // handle sse events outside.
    // sse.onmessage = (event) => JSON.stringify(JSON.parse(event.data))
    // sse.onopen = (event) => {}
    // sse.onerror = (event) => {}
    return sse
}

const updateStoreDescription = (storeId: number, normalToken: string, description: string): Promise<any> => {
    const route = "/stores/".concat(storeId.toString())
    const jsonBody: string = JSON.stringify({
        "description": description,
    })
    return fetch(
        httpTools.generateURL(route), { 
            method: httpTools.HTTPMETHOD.PUT,
            headers: {...httpTools.CONTENT_TYPE_JSON, ...httpTools.generateAuth(normalToken)},
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
          throw new Error("updateStoreDescription error")  
      })
}

export {
    openStore,
    signInStore,
    refreshToken,
    closeStore,
    forgetPassword,
    updatePassword,
    getStoreInfoWithSSE,
    updateStoreDescription
}