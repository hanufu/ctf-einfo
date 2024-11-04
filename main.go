package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// Estrutura para representar um usuário
type Usuario struct {
	Usuario string
	Email   string
	Senha   string
	IP      string
	Pontos  int
}

func init() {
	// Conectar ao banco de dados SQLite
	var err error
	db, err = sql.Open("sqlite3", "./usuarios.db")
	if err != nil {
		panic(err)
	}

	// Criar tabela se não existir
	createTableSQL := `CREATE TABLE IF NOT EXISTS usuarios (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		usuario TEXT NOT NULL,
		email TEXT NOT NULL,
		senha TEXT NOT NULL,
		ip TEXT NOT NULL,
		pontos INTEGER DEFAULT 0
	);`
	if _, err := db.Exec(createTableSQL); err != nil {
		panic(err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		if r.Method == http.MethodGet {
			http.ServeFile(w, r, "index.html")
			return
		}

		if r.Method == http.MethodPost {
			usuario := r.FormValue("usuario")
			password := r.FormValue("password")
			userIp := r.RemoteAddr

			if login(usuario, password) {
				http.ServeFile(w, r, "profile.html")
				return
			}

			fmt.Printf("Login falhou para o usuário: %s, IP: %s\n", usuario, userIp)
			http.Redirect(w, r, "/?error=invalid_credentials", http.StatusSeeOther)
			return
		}
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}
	http.NotFound(w, r)
}

func cadastro(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.ServeFile(w, r, "cadastro.html")
		return
	}

	if r.Method == http.MethodPost {
		usuario := r.FormValue("usuario")
		email := r.FormValue("email")
		password := r.FormValue("password")
		userIp := r.RemoteAddr

		err := db.QueryRow("SELECT senha FROM usuarios WHERE usuario = ?", usuario)
		if err != nil {
			http.Error(w, "Usuário já existe", http.StatusConflict)
			return
		}
		// Cadastrar o novo usuário
		if err := cadastrarUsuario(usuario, email, password, userIp); err != nil {
			http.Error(w, "Erro ao cadastrar usuário", http.StatusInternalServerError)
			return
		}

		fmt.Printf("Usuario: %s, Email: %s, Password: %s, IP: %s\n", usuario, email, password, userIp)
		http.ServeFile(w, r, "index.html")
		return
	}
	http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
}

func cadastrarUsuario(usuario, email, senha, ip string) error {
	_, err := db.Exec("INSERT INTO usuarios (usuario, email, senha, ip) VALUES (?, ?, ?, ?)", usuario, email, senha, ip)
	return err
}

func login(usuario, password string) bool {
	var storedPassword string
	err := db.QueryRow("SELECT senha FROM usuarios WHERE usuario = ?", usuario).Scan(&storedPassword)

	if err != nil {
		if err == sql.ErrNoRows {
			// Usuário não encontrado
			return false
		}
		// Erro ao acessar o banco de dados
		fmt.Println("Erro ao acessar o banco de dados:", err)
		return false
	}

	// Comparar a senha armazenada com a senha fornecida
	return storedPassword == password
}

func serveHiddenFile(w http.ResponseWriter, r *http.Request) {
	filePath := "." + r.URL.Path
	if strings.HasSuffix(filePath, "login.txt") {
		w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(filePath))
		http.ServeFile(w, r, filePath)
		return
	}
	http.FileServer(http.Dir(".hidden")).ServeHTTP(w, r)
}

func main() {
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	http.HandleFunc("/hidden/", serveHiddenFile)
	http.HandleFunc("/", index)
	http.HandleFunc("/cadastro", cadastro)

	fmt.Println("Servidor rodando em http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
