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
    "PUT": "PUT",
    "DELETE": "DELETE",
    "PATCH": "PATCH" 
}

const CONTENT_TYPE_JSON = {
    "Content-Type": "application/json"
}

const generateAuth = (token: string) => {
    token = "Bearer ".concat(token)
    return {
        "Authorization": token
    }
}

export {
    generateURL,
    HTTPMETHOD,
    CONTENT_TYPE_JSON,
    generateAuth
}