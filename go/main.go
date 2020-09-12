package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/newrelic/go-agent/v3/integrations/nrecho-v4"
	_ "github.com/newrelic/go-agent/v3/integrations/nrmysql"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/patrickmn/go-cache"
)

const Limit = 20
const NazotteLimit = 50
const cacheSleep = 42

var mySQLConnectionData MySQLConnectionEnv
var chairSearchCondition ChairSearchCondition
var estateSearchCondition EstateSearchCondition

var estateCache *cache.Cache
var chairCache *cache.Cache
var stockCache *cache.Cache

type InitializeResponse struct {
	Language string `json:"language"`
}

type RecordMapper struct {
	Record []string

	offset int
	err    error
}

func (r *RecordMapper) next() (string, error) {
	if r.err != nil {
		return "", r.err
	}
	if r.offset >= len(r.Record) {
		r.err = fmt.Errorf("too many read")
		return "", r.err
	}
	s := r.Record[r.offset]
	r.offset++
	return s, nil
}

func (r *RecordMapper) NextInt() int {
	s, err := r.next()
	if err != nil {
		return 0
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		r.err = err
		return 0
	}
	return i
}

func (r *RecordMapper) NextFloat() float64 {
	s, err := r.next()
	if err != nil {
		return 0
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		r.err = err
		return 0
	}
	return f
}

func (r *RecordMapper) NextString() string {
	s, err := r.next()
	if err != nil {
		return ""
	}
	return s
}

func (r *RecordMapper) Err() error {
	return r.err
}

//ConnectDB isuumoデータベースに接続する
func (mc MySQLConnectionEnv) ConnectDB() (dbType, error) {
	withState, _ := sqlx.Open("nrmysql", fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", mc.withState.User, mc.withState.Password, mc.withState.Host, mc.withState.Port, mc.withState.DBName))
	noState, _ := sqlx.Open("nrmysql", fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", mc.noState.User, mc.noState.Password, mc.noState.Host, mc.noState.Port, mc.noState.DBName))
	return dbType{
		withState: withState,
		noState:   noState,
	}, nil
}

func init() {
	jsonText, err := ioutil.ReadFile("../fixture/chair_condition.json")
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	json.Unmarshal(jsonText, &chairSearchCondition)

	jsonText, err = ioutil.ReadFile("../fixture/estate_condition.json")
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	json.Unmarshal(jsonText, &estateSearchCondition)
}

func main() {
	// Echo instance
	e := echo.New()
	e.Debug = true
	e.Logger.SetLevel(log.DEBUG)

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	app, _ := newrelic.NewApplication(
		newrelic.ConfigAppName("ISUCON10"),
		newrelic.ConfigLicense("d3224d588a43c8ea493456a20f605978471bNRAL"),
		newrelic.ConfigDistributedTracerEnabled(true),
	)
	e.Use(nrecho.Middleware(app))

	// Initialize
	e.POST("/initialize", initialize)

	// Chair Handler
	e.GET("/api/chair/:id", getChairDetail)
	e.POST("/api/chair", postChair)
	e.GET("/api/chair/search", searchChairs)
	e.GET("/api/chair/low_priced", getLowPricedChair)
	e.GET("/api/chair/search/condition", getChairSearchCondition)
	e.POST("/api/chair/buy/:id", buyChair)

	// Estate Handler
	e.GET("/api/estate/:id", getEstateDetail)
	e.POST("/api/estate", postEstate)
	e.GET("/api/estate/search", searchEstates)
	e.GET("/api/estate/low_priced", getLowPricedEstate)
	e.POST("/api/estate/req_doc/:id", postEstateRequestDocument)
	e.POST("/api/estate/nazotte", searchEstateNazotte)
	e.GET("/api/estate/search/condition", getEstateSearchCondition)
	e.GET("/api/recommended_estate/:id", searchRecommendedEstateWithChair)

	mySQLConnectionData = NewMySQLConnectionEnv()

	var err error
	db, err = mySQLConnectionData.ConnectDB()
	if err != nil {
		e.Logger.Fatalf("DB connection failed : %v", err)
	}
	db.withState.SetMaxOpenConns(10)
	db.noState.SetMaxOpenConns(10)
	defer db.withState.Close()
	defer db.noState.Close()

	estateCache = cache.New(5*time.Minute, 10*time.Minute)
	chairCache = cache.New(5*time.Minute, 10*time.Minute)
	stockCache = cache.New(5*time.Minute, 10*time.Minute)

	// Start server
	serverPort := fmt.Sprintf(":%v", "1323")
	e.Logger.Fatal(e.Start(serverPort))
}

func initialize(c echo.Context) error {
	sqlDir := filepath.Join("..", "mysql", "db")
	paths := []string{
		filepath.Join(sqlDir, "0_Schema.sql"),
		filepath.Join(sqlDir, "1_DummyEstateData.sql"),
		filepath.Join(sqlDir, "2_DummyChairData.sql"),
	}
	ch1 := make(chan bool)
	ch2 := make(chan bool)
	go func() {
		for _, p := range paths {
			sqlFile, _ := filepath.Abs(p)
			cmdStr := fmt.Sprintf("mysql -h %v -u %v -p%v -P %v %v < %v",
				mySQLConnectionData.withState.Host,
				mySQLConnectionData.withState.User,
				mySQLConnectionData.withState.Password,
				mySQLConnectionData.withState.Port,
				mySQLConnectionData.withState.DBName,
				sqlFile,
			)
			if err := exec.Command("bash", "-c", cmdStr).Run(); err != nil {
				c.Logger().Errorf("Initialize script error : %v", err)
				panic(err)
			}
		}
		ch1 <- true
	}()
	go func() {
		for _, p := range paths {
			sqlFile, _ := filepath.Abs(p)
			cmdStr := fmt.Sprintf("mysql -h %v -u %v -p%v -P %v %v < %v",
				mySQLConnectionData.noState.Host,
				mySQLConnectionData.noState.User,
				mySQLConnectionData.noState.Password,
				mySQLConnectionData.noState.Port,
				mySQLConnectionData.noState.DBName,
				sqlFile,
			)
			if err := exec.Command("bash", "-c", cmdStr).Run(); err != nil {
				c.Logger().Errorf("Initialize script error : %v", err)
				panic(err)
			}
		}
		ch2 <- true
	}()
	<-ch1
	<-ch2

	estateCache.Flush()
	chairCache.Flush()
	stockCache.Flush()

	return c.JSON(http.StatusOK, InitializeResponse{
		Language: "go",
	})
}
