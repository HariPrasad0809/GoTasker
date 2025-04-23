// cache/cache.go
package cache

import (
	"time"

	"github.com/harip/GoTasker/models"

	"github.com/patrickmn/go-cache"
)

var taskCache *cache.Cache

func InitCache() {
	taskCache = cache.New(5*time.Minute, 10*time.Minute)
}

func GetTasks() ([]models.Task, bool) {
	if tasks, found := taskCache.Get("tasks"); found {
		return tasks.([]models.Task), true
	}
	return nil, false
}

func SetTasks(tasks []models.Task) {
	taskCache.Set("tasks", tasks, cache.DefaultExpiration)
}

func InvalidateCache() {
	taskCache.Delete("tasks")
}
