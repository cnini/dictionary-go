package dictionary

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type Dictionary struct {
	RedisClient *redis.Client
}

func NewDictionary(errors chan<- error) *Dictionary {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		errors <- err
	}

	return &Dictionary{
		RedisClient: client,
	}
}

func (d Dictionary) Add(word string, definition string, wg *sync.WaitGroup, errors chan<- error) {
	wg.Add(1)
	defer wg.Done()

	_, err := d.RedisClient.Set(ctx, word, definition, 0).Result()
	if err != nil {
		errors <- err
	}
}

func (d Dictionary) Get(searchTerm string, errors chan<- error) (string, string) {
	// If searchTerm is the word and not the definition
	if result, _ := d.RedisClient.Get(ctx, searchTerm).Result(); result != "" {
		return searchTerm, result
	} else {
		// If searchTerm is the definition
		keys, err := d.RedisClient.Keys(ctx, "*").Result()
		if err != nil {
			errors <- err
		}

		for _, word := range keys {
			definition, err := d.RedisClient.Get(ctx, word).Result()
			if err != nil {
				errors <- err
			}

			if searchTerm == definition {
				return word, searchTerm
			}
		}
	}

	return "", ""
}

func (d Dictionary) Remove(termToRemove string, wg *sync.WaitGroup, errors chan<- error) {
	wg.Add(1)
	defer wg.Done()

	// Get the correct word even if termToRemove is the definition
	word, _ := d.Get(termToRemove, errors)

	err := d.RedisClient.Del(ctx, word).Err()
	if err != nil {
		errors <- err
	}
}

func (d Dictionary) List(errors chan<- error) []string {
	var sortedDictionary []string

	// Get all keys from database
	keys, _ := d.RedisClient.Keys(ctx, "*").Result()

	// Sort all keys
	sort.Strings(keys)

	// For each keys, get the definition and add their formatted form into sortedDictionary
	for _, word := range keys {
		_, definition := d.Get(word, errors)

		if definition != "" {
			sortedDictionary = append(sortedDictionary, fmt.Sprintf("%s: %s", word, definition))
		}
	}

	return sortedDictionary
}
