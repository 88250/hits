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
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var logger *Logger

func init() {
	rand.Seed(time.Now().Unix())

	SetLevel("info")
	logger = NewLogger(os.Stdout)
	gin.SetMode(gin.ReleaseMode)
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
var hitMap = map[string]int{}

func hit(c *gin.Context) {
	owner := c.Param("owner")
	repo := c.Param("repo")
	key := owner + "/" + repo

	locker.Lock()
	count, ok := 1, false
	if count, ok = hitMap[key]; ok {
		hitMap[key] = count + 1
	} else {
		hitMap[key] = count
	}
	locker.Unlock()

	svg := `<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" width="88" height="20"><g shape-rendering="crispEdges"><path fill="#555" d="M0 0h37v20H0z"/><path fill="#4c1" d="M37 0h51v20H37z"/></g><g fill="#fff" text-anchor="middle" font-family="DejaVu Sans,Verdana,Geneva,sans-serif" font-size="110">`
	svg += `<text x="195" y="140" transform="scale(.1)">hits</text>`
	svg += `<text x="615" y="140" transform="scale(.1)">` + strconv.Itoa(count) + `</text></g></svg>`
	c.Data(200, "image/svg+xml;charset=utf-8", []byte(svg))
}

func main() {
	router := mapRoutes()
	server := &http.Server{
		Addr:    "127.0.0.1:1124",
		Handler: router,
	}

	logger.Infof("hits is running")
	server.ListenAndServe()
}
