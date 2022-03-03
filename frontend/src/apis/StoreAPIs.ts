import * as httpTools from './base'

const openStore = (email: string, password: string, name: string, timezone: string, queueNames: string[]): httpTools.RequestParams => {
    const jsonBody: string = JSON.stringify({
        "email": email,
        "password": password,
        "name": name,
        "timezone": timezone,
        "queue_names": queueNames
    })
    return {
        endpoint: httpTools.generateURL("/stores"),
        method: httpTools.HTTPMETHOD.POST,
        headers: httpTools.CONTENT_TYPE_JSON,
        body: jsonBody 
    }
}

const signInStore = (email: string, password: string): httpTools.RequestParams => {
    const jsonBody: string = JSON.stringify({
        "email": email,
        "password": password,
    })
    return {
        endpoint: httpTools.generateURL("/stores/signin"),
        method: httpTools.HTTPMETHOD.POST,
        headers: httpTools.CONTENT_TYPE_JSON,
        body: jsonBody 
    }
}

const refreshToken = (): httpTools.RequestParams => {
    return {
        endpoint: httpTools.generateURL("/stores/token"),
        method: httpTools.HTTPMETHOD.PUT 
    }
}

const closeStore = (storeId: number, normalToken: string): httpTools.RequestParams => {
    const route = "/stores/".concat(storeId.toString())
    return { 
        endpoint: httpTools.generateURL(route),
        method: httpTools.HTTPMETHOD.DELETE,
        headers: httpTools.generateAuth(normalToken)
    }
}

const forgetPassword = (email: string): httpTools.RequestParams => {
    const jsonBody: string = JSON.stringify({
        "email": email,
    })
    return { 
        endpoint: httpTools.generateURL("/stores/password/forgot"),
        method: httpTools.HTTPMETHOD.POST,
        headers: httpTools.CONTENT_TYPE_JSON,
        body: jsonBody
    }
}

const updatePassword = (storeId: number, passwordToken: string, password: string): httpTools.RequestParams => {
    const route = "/stores/".concat(storeId.toString(), "/password")
    const jsonBody: string = JSON.stringify({
        "password_token": passwordToken,
        "password": password,
    })
    return { 
        endpoint: httpTools.generateURL(route), 
        method: httpTools.HTTPMETHOD.PATCH,
        headers: httpTools.CONTENT_TYPE_JSON,
        body: jsonBody
    }
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

const updateStoreDescription = (storeId: number, normalToken: string, description: string): httpTools.RequestParams => {
    const route = "/stores/".concat(storeId.toString())
    const jsonBody: string = JSON.stringify({
        "description": description,
    })
    return { 
        endpoint: httpTools.generateURL(route), 
        method: httpTools.HTTPMETHOD.PUT,
        headers: {...httpTools.CONTENT_TYPE_JSON, ...httpTools.generateAuth(normalToken)},
        body: jsonBody
    }
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