package delivery

import "fmt"

func V_1(route string) string {
	return fmt.Sprintf("/v1%s", route)
}