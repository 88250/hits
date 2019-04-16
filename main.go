// Hits - GitHub repository hits counter.
// Copyright (C) 2019-present, b3log.org
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"flag"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var logger *Logger

var dirPath = "./"

func init() {
	rand.Seed(time.Now().Unix())

	SetLevel("info")
	logger = NewLogger(os.Stdout)
	gin.SetMode(gin.ReleaseMode)

	dir := flag.String("dir", "", "path of data dir directory, for example /opt/hits/data")
	flag.Parse()

	if "" != *dir {
		dirPath = *dir
	} else {
		dirPath = filepath.Join(UserHome(), "hits")
		if !IsExist(dirPath) {
			if err := os.Mkdir(dirPath, 0644); nil != err {
				logger.Fatalf("create data directory [%s] failed [%s]", dirPath, err.Error())
			}
		}
	}
}

func mapRoutes() *gin.Engine {
	ret := gin.New()
	ret.Use(gin.Recovery())

	ret.GET("/:owner/:repo", hit)
	ret.NoRoute(func(c *gin.Context) {
		c.String(http.StatusOK, "The piper will lead us to reason.\n\n欢迎访问黑客与画家的社区 https://hacpai.com")
	})

	return ret
}

var locker = sync.Mutex{}

func hit(c *gin.Context) {
	owner := c.Param("owner")
	repo := c.Param("repo")
	if strings.Contains(owner, "/") || strings.Contains(repo, "/") || !strings.Contains(repo, ".svg") {
		c.Status(404)

		return
	}

	repo = repo[:strings.LastIndex(repo, ".svg")]
	key := owner + "-" + repo

	locker.Lock()
	_, count := writeData(key)
	locker.Unlock()

	svg := `<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" width="88" height="20"><g shape-rendering="crispEdges"><path fill="#555" d="M0 0h37v20H0z"/><path fill="#4c1" d="M37 0h51v20H37z"/></g><g fill="#fff" text-anchor="middle" font-family="DejaVu Sans,Verdana,Geneva,sans-serif" font-size="110">`
	svg += `<text x="195" y="140" transform="scale(.1)">hits</text>`
	svg += `<text x="615" y="140" transform="scale(.1)">` + count + `</text></g></svg>`
	c.Data(200, "image/svg+xml;charset=utf-8", []byte(svg))
}

func writeData(fileName string) (count int, countStr string) {
	count, countStr = 1, "1"
	var err error

	dataFilePath := dirPath + "/" + fileName
	var f *os.File
	if f, err = os.OpenFile(dataFilePath, os.O_CREATE|os.O_RDWR, 0664); nil != err {
		logger.Errorf("open file [%s] failed [%s]", dataFilePath, err.Error())
		return
	}
	if bytes, err := ioutil.ReadAll(f); nil != err {
		logger.Errorf("read file [%s] failed [%s]", dataFilePath, err.Error())
		return
	} else {
		countStr = string(bytes)
	}
	f.Close()

	if "" == countStr {
		countStr = "1"
	}
	countStr = strings.TrimSpace(countStr)
	count = 1
	if count, err = strconv.Atoi(countStr); nil != err {
		logger.Errorf("read count of file [%s] failed  [%s]", dataFilePath, err.Error())
		return
	}

	count++
	countStr = strconv.Itoa(count)
	if err = ioutil.WriteFile(dataFilePath, []byte(countStr), 0644); nil != err {
		logger.Errorf("write count to file [%s] failed [%s]", dataFilePath, err.Error())
		return
	}

	return
}

func main() {
	router := mapRoutes()
	server := &http.Server{
		Addr:    "127.0.0.1:1124",
		Handler: router,
	}

	logger.Infof("hits is running, data directory is [" + dirPath + "]")
	server.ListenAndServe()
}
