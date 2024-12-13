package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"forum/internal/adapters"
	"forum/internal/pkg/httphelper"
	"forum/internal/ports"
	"forum/internal/service"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// ------------------------------------------------------------
	// Инициализирует конфигурацию приложения.
	config := getConfig()

	//Подключается к базе данных SQLite.
	db, err := getDB(config)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//Загружает HTML-шаблоны для генерации веб-страниц.
	tmpls, err := getTmpls()
	if err != nil {
		log.Fatal()
	}

	httphelper.InitTemplates(tmpls)

	// ------------------------------------------------------------
	// ------------------------------------------------------------

	//Инициализирует репозитории для работы с различными сущностями форума (пользователи, посты, категории, комментарии и их реакции).
	var (
		users            = adapters.NewUsersRepositorySqlite3(db)
		posts            = adapters.NewPostsRepositorySqlite3(db)
		categories       = adapters.NewCategoriesRepositorySqlite3(db)
		postCategories   = adapters.NewPostCategoriesRepositorySqlite3(db)
		sessions         = adapters.NewSessionsRepositorySqlite3(db)
		postReactions    = adapters.NewPostReactionsRepositorySqlite3(db)
		comments         = adapters.NewCommentsRepositorySqlite3(db)
		commentReactions = adapters.NewCommentReactionsRepositorySqlite3(db)
	)

	//Создаёт сервисный слой для обработки бизнес-логики.
	svc := service.NewService(
		users,
		posts,
		categories,
		postCategories,
		postReactions,
		sessions,
		comments,
		commentReactions,
		config.fileStorage,
	)
	//Инициализация HTTP-обработчиков
	handler := ports.NewHandler(svc)
	// ------------------------------------------------------------
	// ------------------------------------------------------------

	var (
		port   = ":" + config.port
		routes = handler.InitRouters()
	)

	log.Printf("server started: http://localhost:%s\n", config.port)

	if err := http.ListenAndServe(port, routes); err != nil {
		panic(err)
	}

	// ------------------------------------------------------------
}

func getDB(config *Config) (*sql.DB, error) {
	db, err := sql.Open(config.sqlite3.driver, config.sqlite3.dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func getTmpls() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {

		name := filepath.Base(page)

		ts, err := template.New(name).
			Funcs(template.FuncMap{
				"humanDate": func(t *time.Time) string {
					return t.Format("02 Jan 2006 at 15:04")
				},
			}).
			ParseFiles("./ui/html/base.tmpl")
		if err != nil {
			return nil, err
		}
		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		name = strings.ReplaceAll(name, ".tmpl", "")

		cache[name] = ts
	}

	return cache, nil
}
