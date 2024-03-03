package gofindit

import (
	"fmt"
	"math/rand/v2"
	"time"
)

var TestIndex *Index
var randSource *rand.Rand

type TestData struct {
	Name     string    `find:"name"`
	Age      int       `find:"age"`
	Hobbies  []string  `find:"hobbies"`
	Bio      string    `find:"bio"`
	Birthday time.Time `find:"birthday"`
}

func init() {
	TestIndex = New()

	// new rand
	randSource = rand.New(rand.NewPCG(uint64(time.Now().Nanosecond()), uint64(time.Now().Nanosecond())))

	// array of 50 random first names
	firstnames := []string{
		"John", "Jane", "Billy", "Bob", "Sally", "Sue", "Tom", "Tim", "Tina", "Terry",
		"Mike", "Molly", "Megan", "Morgan", "Maddie", "Maggie", "Marge", "Marge", "Marge", "Marge",
		"Steve", "Stacy", "Samantha", "Sam", "Sue", "Sue", "Sue", "Sue", "Sue", "Sue",
		"Chris", "Christine", "Christina", "Christy", "Christina", "Christina", "Christina", "Christina", "Christina", "Christina",
		"David", "Diana", "Diane", "Dina", "Dina", "Dina", "Dina", "Dina", "Dina", "Dina",
	}

	// array of 50 random last names
	lastnames := []string{
		"Smith", "Johnson", "Williams", "Jones", "Brown", "Davis", "Miller", "Wilson", "Moore", "Taylor",
		"Anderson", "Thomas", "Jackson", "White", "Harris", "Martin", "Thompson", "Garcia", "Martinez", "Robinson",
		"Clark", "Rodriguez", "Lewis", "Lee", "Walker", "Hall", "Allen", "Young", "Hernandez", "King",
		"Wright", "Lopez", "Hill", "Scott", "Green", "Adams", "Baker", "Gonzalez", "Nelson", "Carter",
		"Mitchell", "Perez", "Roberts", "Turner", "Phillips", "Campbell", "Parker", "Evans", "Edwards", "Collins",
	}

	// array of 50 random hobbies
	hobbies := []string{
		"Rowing", "Skiing", "Soccer", "Football", "Basketball", "Baseball", "Hockey", "Golf", "Tennis", "Rugby",
		"Swimming", "Cycling", "Running", "Jogging", "Walking", "Hiking", "Camping", "Fishing", "Hunting", "Shooting",
		"Archery", "Bowling", "Billiards", "Darts", "Poker", "Chess", "Checkers", "Backgammon", "Go", "Mahjong",
		"Bridge", "Canasta", "Euchre", "Gin", "Hearts", "Pinochle", "Rummy", "Spades", "Whist", "Bingo",
		"Cricket", "Polo", "Racquetball", "Squash", "Badminton", "Table Tennis", "Volleyball", "Handball", "Fencing", "Judo",
	}

	// array of 50 lorem ipsum words
	loremipsums := []string{
		"et", "in", "sed", "non", "ac", "nec", "nulla", "eu", "orci", "quis",
		"quam", "odio", "sit", "amet", "lorem", "ipsum", "dolor", "vitae", "a", "morbi",
		"leo", "risus", "porta", "ac", "consectetur", "ac", "vestibulum", "at", "eros", "praesent",
		"commodo", "cursus", "magna", "vel", "scelerisque", "nisl", "consectetur", "et", "cras", "mattis",
		"consectetur", "purus", "sit", "amet", "fermentum", "aenean", "lacinia", "bibendum", "nulla", "sed",
	}

	// Generate a bunch of test documents
	for i := 0; i < 10000; i++ {
		id := fmt.Sprintf("%d", randSource.Uint64())
		name := fmt.Sprintf("%s %s", firstnames[randSource.IntN(len(firstnames)-1)], lastnames[randSource.IntN(len(lastnames)-1)])
		age := 20 + randSource.IntN(60)
		myhobbies := make([]string, randSource.IntN(5)+1)
		for j := 0; j < len(myhobbies); j++ {
			myhobbies[j] = hobbies[randSource.IntN(len(hobbies)-1)]
		}
		bio := ""
		for j := 0; j < randSource.IntN(10)+1; j++ {
			bio += loremipsums[randSource.IntN(len(loremipsums)-1)] + " "
		}
		birthday := time.Now().AddDate(-age, -randSource.IntN(12), -randSource.IntN(28))

		doc := TestData{
			Name:     name,
			Age:      age,
			Hobbies:  myhobbies,
			Bio:      bio,
			Birthday: birthday,
		}

		err := TestIndex.Index(id, doc)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
