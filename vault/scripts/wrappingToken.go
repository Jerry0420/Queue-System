package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"net/http"
	"log"
)

// Run this server for other services to get wrapped token.
// Like vault server, every time when vault restart, it's status is set to be sealed. This server is not implementing restart mechanism too. 
func main() {
	http.HandleFunc("/wrappedToken", func (w http.ResponseWriter, r *http.Request)  {
		w.Header().Set("Content-Type", "application/json")
		
		if r.Method == http.MethodPost {
			var jsonBody map[string]string
			err := json.NewDecoder(r.Body).Decode(&jsonBody)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			if roleName, ok := jsonBody["roleName"]; ok {
				cmd := fmt.Sprintf(
					"vault write -force -wrap-ttl=%s auth/approle/role/%s/secret-id -format=json", 
					"30m", // wrap-ttl of wrapped token
					roleName,
				)
				out, err := exec.Command("sh", "-c", cmd).Output()
				var wrappedTokenResults map[string]interface{}
				json.Unmarshal(out, &wrappedTokenResults)

				if _, ok := wrappedTokenResults["wrap_info"]; !ok || err != nil{
					http.Error(w, "server is sealed.", http.StatusBadRequest)
					return
				} else {
					log.Println("Got wrappedToken!")
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(
						&map[string]interface{}{
							"wrappedToken": wrappedTokenResults["wrap_info"].(map[string]interface{})["token"],
						},
					)
				}
			} else {
				http.Error(w, "wrong parameter in body", http.StatusNotFound)
				return
			}
		} else {
			http.Error(w, "wrong http method.", http.StatusMethodNotAllowed)
			return
		}
    })

	log.Println("server start!")
	err := http.ListenAndServe(":8300", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    } 
}