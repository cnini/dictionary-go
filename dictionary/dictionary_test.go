package dictionary

import (
	"sync"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

var miniredisServer *miniredis.Miniredis

func initDictionary(miniRedis *miniredis.Miniredis) *Dictionary {
	return &Dictionary{
		RedisClient: redis.NewClient(&redis.Options{
			Addr:     miniRedis.Addr(),
			Password: "",
			DB:       0,
		}),
	}
}

func loadFixtures(miniRedis *miniredis.Miniredis) {
	fixtures := [][2]string{
		{"Chapeau", "Hat"},
		{"Verre", "Glass"},
		{"Actrice", "Actress"},
		{"Acteur", "Actor"},
		{"Or", "Gold"},
		{"Visage", "Face"},
	}

	for _, fixture := range fixtures {
		miniRedis.Set(fixture[0], fixture[1])
	}
}

func TestSuccessAdd(t *testing.T) {
	r := miniredis.RunT(t)
	dictionary := initDictionary(r)

	dictionary.Add("Souris", "Mouse", &sync.WaitGroup{}, make(chan<- error))
	expected, _ := r.Get("Souris")
	assert.Equal(t, expected, "Mouse", "Le mot en français pour \"Souris\" est bien \"Mouse\".")
}

func TestFailAdd(t *testing.T) {
	r := miniredis.RunT(t)
	dictionary := initDictionary(r)

	dictionary.Add("Fille", "Girl", &sync.WaitGroup{}, make(chan<- error))
	expected, _ := r.Get("Fille")
	assert.NotEqual(t, expected, "Boy", "Le mot en français pour \"Fille\" n'est pas \"Boy\".")
}

func TestSuccessGet(t *testing.T) {
	r := miniredis.RunT(t)
	dictionary := initDictionary(r)

	loadFixtures(r)

	_, definition := dictionary.Get("Chapeau", make(chan<- error))
	expectedDefinition, _ := r.Get("Chapeau")
	assert.Equal(t, expectedDefinition, definition, "Le mot en anglais pour \"Chapeau\" est bien \"Hat\".")

	word, _ := dictionary.Get("Glass", make(chan<- error))
	assert.Equal(t, "Verre", word, "Le mot en français pour \"Glass\" est bien \"Verre\".")
}

func TestFailGet(t *testing.T) {
	r := miniredis.RunT(t)
	dictionary := initDictionary(r)

	loadFixtures(r)

	wordEmpty, definitionEmpty := dictionary.Get("", make(chan<- error))
	assert.Empty(t, wordEmpty, "Le mot renvoyé est bien vide.")
	assert.Empty(t, definitionEmpty, "La définition renvoyée est bien vide.")

	_, definitionNil := dictionary.Get("Médecin", make(chan<- error))
	assert.Empty(t, definitionNil, "Le mot \"Médecin\" et sa définition n'existent pas.")
}

func TestSuccessRemove(t *testing.T) {
	r := miniredis.RunT(t)
	dictionary := initDictionary(r)

	loadFixtures(r)

	dictionary.Remove("Acteur", &sync.WaitGroup{}, make(chan<- error))
	expectedNil, _ := r.Get("Acteur")
	assert.Empty(t, expectedNil, "Le mot \"Acteur\" et sa définition ont été supprimé.")

	dictionary.Remove("Face", &sync.WaitGroup{}, make(chan<- error))
	expectedNil, _ = r.Get("Face")
	assert.Empty(t, expectedNil, "Le mot \"Face\" et sa définition ont été supprimé.")
}
