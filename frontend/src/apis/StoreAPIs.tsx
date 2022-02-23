const openStore = (email: string, password: string, name: string, timezone: string, queue_names: string[]) => {
    const jsonBody: string = JSON.stringify({
        "email": email,
        "password": password,
        "name": name,
        "timezone": timezone,
        "queue_names": queue_names
    })
    return fetch(
        "http://backend:8000/api/v1/stores", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
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

export {openStore}