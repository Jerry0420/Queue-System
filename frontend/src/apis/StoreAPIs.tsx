import * as httpTools from './base'

const openStore = (email: string, password: string, name: string, timezone: string, queue_names: string[]): Promise<any> => {
    const jsonBody: string = JSON.stringify({
        "email": email,
        "password": password,
        "name": name,
        "timezone": timezone,
        "queue_names": queue_names
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
          throw new Error("openStore error")  
      })
}

export {
    openStore,
    signInStore
}