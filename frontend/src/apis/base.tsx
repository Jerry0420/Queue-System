const version: string = "/api/v1"

const generateURL = (route: string): string => {
    const backendHost: string = (process.env.BACKEND_HOST as string)
    const backendPort: string = (process.env.BACKEND_PORT as string)
    let url = "http://".concat(backendHost, ":", backendPort, version, route)
    return url
}

const HTTPMETHOD = {
    "GET": "GET",
    "POST": "POST",
    "PUT": "PUT" 
}

const CONTENT_TYPE_JSON = {
    "Content-Type": "application/json"
}

const AUTHORIZATION = "Authorization"

export {
    generateURL,
    HTTPMETHOD,
    CONTENT_TYPE_JSON,
    AUTHORIZATION
}