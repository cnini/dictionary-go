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

type DictionaryEntry struct {
	Word       string
	Definition string
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

	dictionaryEntry := DictionaryEntry{
		Word:       word,
		Definition: definition,
	}

	_, err := d.RedisClient.Set(ctx, dictionaryEntry.Word, dictionaryEntry.Definition, 0).Result()
	if err != nil {
		errors <- err
	}
}

func (d Dictionary) Get(searchTerm string, errors chan<- error) (string, string) {
	result, err := d.RedisClient.Get(ctx, searchTerm).Result()
	if err != nil {
		errors <- err
		return "", ""
	}

	if result != "" {
		return searchTerm, result
	}

	return "", ""
}

func (d Dictionary) Remove(termToRemove string, wg *sync.WaitGroup, errors chan<- error) {
	wg.Add(1)
	defer wg.Done()

	err := d.RedisClient.Del(ctx, termToRemove).Err()
	if err != nil {
		errors <- err
	}
}

func (d Dictionary) List(errors chan<- error) []string {
	var sortedDictionary []string

	keys, _ := d.RedisClient.Keys(ctx, "*").Result()

	sort.Strings(keys)

	for _, word := range keys {
		definition, err := d.RedisClient.Get(ctx, word).Result()
		if err != nil {
			errors <- err
		}

		sortedDictionary = append(sortedDictionary, fmt.Sprintf("%s: %s", word, definition))
	}

	return sortedDictionary
}
