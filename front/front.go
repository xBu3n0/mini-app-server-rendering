package front

import "net/http"

func handleRoot(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./main.html")
}

func main() {
	http.HandleFunc("/", handleRoot)

	http.ListenAndServe(":3001", nil)
}
