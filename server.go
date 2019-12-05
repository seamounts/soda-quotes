package main

import (
	"io"
	"net/http"
	"strconv"
	"sync"
	"text/template"
	"time"

	"github.com/labstack/echo"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var db *gorm.DB

func init() {
	var err error
	db, err = gorm.Open("sqlite3", "db/quotes.db")
	if err != nil {
		panic("failed to connect database")
	}
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

type ChickenSoul struct {
	ID      uint   `gorm:"primary_key"`
	Content string `gorm:"type:text"`
	Hits    string `gorm:"type:varchar(100)"`
}

type (
	Stats struct {
		Uptime       time.Time           `json:"uptime"`
		RequestCount uint64              `json:"requestCount"`
		Statuses     map[string]int      `json:"statuses"`
		ApiStats     map[string]*ApiStat `json:"api_stats"`
		mutex        sync.RWMutex
	}

	ApiStat struct {
		Path         string         `json:"path"`
		Statuses     map[string]int `json:"statuses"`
		RequestCount uint64         `json:"requestCount"`
		Uptime       time.Time      `json:"uptime"`
	}
)

func NewStats() *Stats {
	return &Stats{
		Uptime:   time.Now(),
		Statuses: make(map[string]int),
		ApiStats: make(map[string]*ApiStat),
	}
}

// Process is the middleware function.
func (s *Stats) Process(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := next(c); err != nil {
			c.Error(err)
		}
		s.mutex.Lock()
		defer s.mutex.Unlock()
		status := strconv.Itoa(c.Response().Status)

		if _, ok := s.ApiStats[c.Path()]; !ok {
			apiStat := &ApiStat{
				Uptime:   time.Now(),
				Path:     c.Path(),
				Statuses: make(map[string]int),
			}
			s.ApiStats[c.Path()] = apiStat
		}

		s.ApiStats[c.Path()].Statuses[status]++
		s.ApiStats[c.Path()].RequestCount++

		s.RequestCount++
		s.Statuses[status]++
		return nil
	}
}

// Handle is the endpoint to get stats.
func (s *Stats) Handle(c echo.Context) error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return c.JSON(http.StatusOK, s)
}

func getChickenSoul() *ChickenSoul {
	chickenSoul := &ChickenSoul{}
	db.Raw("select * from chicken_soul order by RANDOM() limit 1").Scan(chickenSoul)

	return chickenSoul
}

func Quotes(c echo.Context) error {
	return c.Render(http.StatusOK, "index", getChickenSoul())
}

func RandomChickenSoul(c echo.Context) error {
	return c.JSON(http.StatusOK, getChickenSoul())
}

func main() {
	e := echo.New()
	e.Renderer = &Template{
		templates: template.Must(template.ParseGlob("views/index.html")),
	}

	e.Static("/", "statics")

	// Stats
	s := NewStats()
	e.Use(s.Process)
	e.GET("/stats", s.Handle) // Endpoint to get stats

	e.GET("/", Quotes)
	e.GET("/random-chicken-soul", RandomChickenSoul)

	e.Logger.Fatal(e.Start(":8081"))
}
