package main

import (
	"fmt"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.ServeFile(w, r, "index.html")
		return
	}
	fmt.Println("Post")
	matricula := r.FormValue("matricula")
	password := r.FormValue("password")

	fmt.Printf("Matricula: %s Password: %s", matricula, password)
	http.ServeFile(w, r, "index.html")
}

func cadastro(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.ServeFile(w, r, "cadastro.html")
		return
	}
	usuario := r.FormValue("usuario")
	email := r.FormValue("email")
	password := r.FormValue("password")

	fmt.Printf("Usuario: %s, Email: %s Senha: %s\n", usuario, email, password)
	http.ServeFile(w, r, "index.html")
}

// func handleHiddenDir(w http.ResponseWriter, r *http.Request) {
// 	// Caminho para o diretório "oculto"
// 	hiddenDir := ".hidden"

// 	// Verifica se o arquivo ou diretório existe
// 	path := filepath.Join(hiddenDir, filepath.Clean(r.URL.Path[len("/.hidden/"):]))
// 	_, err := os.Stat(path)
// 	if os.IsNotExist(err) {
// 		http.NotFound(w, r)
// 		return
// 	}

// 	// Serve o arquivo do diretório "oculto"
// 	http.ServeFile(w, r, path)
// }

func main() {
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	// Serve o HTML
	http.HandleFunc("/", home)
	http.HandleFunc("/cadastro", cadastro)

	fmt.Println("Server rodando em http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
