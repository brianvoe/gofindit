package gofindit

import (
	"fmt"
	"math/rand/v2"
	"time"
)

var TestIndex *Index
var randSource *rand.Rand = rand.New(rand.NewPCG(uint64(time.Now().Nanosecond()), uint64(time.Now().Nanosecond())))

type TestData struct {
	Name      string    `find:"name"`
	Age       int       `find:"age"`
	Hobbies   []string  `find:"hobbies"`
	Pets      []TestPet `find:"pets"`
	Bio       string    `find:"bio"`
	isStudent bool      `find:"isStudent"`
	Birthday  time.Time `find:"birthday"`
}

type TestPet struct {
	Name  string   `find:"name"`
	Age   int      `find:"age"`
	Type  string   `find:"type"`
	Breed string   `find:"breed"`
	Toys  []string `find:"toys"`
}

func init() {
	TestIndex = New()

	// Generate a bunch of test documents
	generateTestDocs()
}

// array of 50 random first names
var firstnames = []string{
	"John", "Jane", "Billy", "Bob", "Sally", "Sue", "Tom", "Tim", "Tina", "Terry",
	"Mike", "Molly", "Megan", "Morgan", "Maddie", "Maggie", "Marge", "Marge", "Marge", "Marge",
	"Steve", "Stacy", "Samantha", "Sam", "Sue", "Sue", "Sue", "Sue", "Sue", "Sue",
	"Chris", "Christine", "Christina", "Christy", "Christina", "Christina", "Christina", "Christina", "Christina", "Christina",
	"David", "Diana", "Diane", "Dina", "Dina", "Dina", "Dina", "Dina", "Dina", "Dina",
}

// array of 50 random last names
var lastnames = []string{
	"Smith", "Johnson", "Williams", "Jones", "Brown", "Davis", "Miller", "Wilson", "Moore", "Taylor",
	"Anderson", "Thomas", "Jackson", "White", "Harris", "Martin", "Thompson", "Garcia", "Martinez", "Robinson",
	"Clark", "Rodriguez", "Lewis", "Lee", "Walker", "Hall", "Allen", "Young", "Hernandez", "King",
	"Wright", "Lopez", "Hill", "Scott", "Green", "Adams", "Baker", "Gonzalez", "Nelson", "Carter",
	"Mitchell", "Perez", "Roberts", "Turner", "Phillips", "Campbell", "Parker", "Evans", "Edwards", "Collins",
}

// array of 50 random hobbies
var hobbies = []string{
	"Rowing", "Skiing", "Soccer", "Football", "Basketball", "Baseball", "Hockey", "Golf", "Tennis", "Rugby",
	"Swimming", "Cycling", "Running", "Jogging", "Walking", "Hiking", "Camping", "Fishing", "Hunting", "Shooting",
	"Archery", "Bowling", "Billiards", "Darts", "Poker", "Chess", "Checkers", "Backgammon", "Go", "Mahjong",
	"Bridge", "Canasta", "Euchre", "Gin", "Hearts", "Pinochle", "Rummy", "Spades", "Whist", "Bingo",
	"Cricket", "Polo", "Racquetball", "Squash", "Badminton", "Table Tennis", "Volleyball", "Handball", "Fencing", "Judo",
}

// array of 50 lorem ipsum words
var loremipsums = []string{
	"et", "in", "sed", "non", "ac", "nec", "nulla", "eu", "orci", "quis",
	"quam", "odio", "sit", "amet", "lorem", "ipsum", "dolor", "vitae", "a", "morbi",
	"leo", "risus", "porta", "ac", "consectetur", "ac", "vestibulum", "at", "eros", "praesent",
	"commodo", "cursus", "magna", "vel", "scelerisque", "nisl", "consectetur", "et", "cras", "mattis",
	"consectetur", "purus", "sit", "amet", "fermentum", "aenean", "lacinia", "bibendum", "nulla", "sed",
}

var petTypes = []string{
	"Dog", "Cat", "Bird", "Fish", "Rabbit", "Hamster", "Guinea Pig", "Turtle", "Lizard", "Snake",
}

var petBreeds = map[string][]string{
	"Dog":        {"Labrador Retriever", "German Shepherd", "Golden Retriever", "French Bulldog", "Bulldog"},
	"Cat":        {"Persian", "Maine Coon", "Siamese", "Ragdoll", "Bengal"},
	"Bird":       {"Parakeet", "Cockatiel", "Canary", "Lovebird", "African Grey Parrot"},
	"Fish":       {"Goldfish", "Betta", "Angelfish", "Guppy", "Oscar"},
	"Rabbit":     {"Holland Lop", "Netherland Dwarf", "Rex", "Mini Lop", "Flemish Giant"},
	"Hamster":    {"Syrian", "Dwarf Campbells Russian", "Dwarf Winter White Russian", "Roborovski", "Chinese"},
	"Guinea Pig": {"American", "Abyssinian", "Peruvian", "Teddy", "Silkie"},
	"Turtle":     {"Red-Eared Slider", "African Aquatic Sideneck", "Painted Turtle", "Eastern Box Turtle", "Mississippi Map Turtle"},
	"Lizard":     {"Bearded Dragon", "Leopard Gecko", "Crested Gecko", "Blue-Tongued Skink", "Green Anole"},
	"Snake":      {"Corn Snake", "Ball Python", "California Kingsnake", "Garter Snake", "Rosy Boa"},
}

var petNames = map[string][]string{
	"Dog":        {"Bark Twain", "Sir Waggington", "Chewbarka", "Bark Obama", "Jimmy Chew"},
	"Cat":        {"Purrman Meowville", "Catrick Swayze", "Kitty Purry", "Chairman Meow", "The Great Catsby"},
	"Bird":       {"Tweety Pie", "Peckachu", "Quack Sparrow", "Beaky Blinder", "Feather Locklear"},
	"Fish":       {"Gill-bert", "Fin Diesel", "Sushi", "Bubbles", "Swim Shady"},
	"Rabbit":     {"Jumpstart", "Bun Bun", "Furrball", "Hoptimist", "Dust Bunny"},
	"Hamster":    {"Hamlet", "Nibbles", "Paws", "Whiskers", "Fuzzy Wuzzy"},
	"Guinea Pig": {"Pigasso", "Wheekly", "Fluffy", "Cuddlebug", "Pudgy"},
	"Turtle":     {"Shell Shock", "Speedy", "Turbo", "Snappy", "Shell-don"},
	"Lizard":     {"Lizzy", "Scaley", "Drago", "Camouflage", "Iggy"},
	"Snake":      {"Hiss-topher", "Slitherin", "Monty Python", "Cobra Winfrey", "Sssam"},
}

var funnyPetToys = []string{
	"Chewbacca",
	"Furrball Frenzy",
	"Bark Vader",
	"Purr-cussion Drum",
	"Meowtain Climber",
	"Squeak-a-Boo",
	"Tail Spinner",
	"Bouncy Burger",
	"Ruff Tuff Tug",
	"Quack Attack",
	"Drool Pool",
	"Wobble Wag Giggle",
	"Fluff-n-Tuff",
	"Snuggle Bug",
	"Pounce n' Play",
	"Jingle Jaws",
	"Fuzzy Frenzy",
	"Slobber Sling",
	"Kitty Kicker",
	"Barky Burger",
	"Giggle Gator",
	"Squeaky Squirrel",
	"Mystery Motion",
	"Treat Teaser",
	"Zoomer Zinger",
	"Whisker Twister",
	"Feather Flinger",
	"Catch-Me-If-You-Can",
	"Hide-n-Squeak",
	"Tug-a-Jug",
}

func generateTestDocs() {
	// Generate a bunch of test documents
	for i := 0; i < 10000; i++ {
		id, doc := generateDoc()

		err := TestIndex.Index(id, doc)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func generateDoc() (string, TestData) {
	id := fmt.Sprintf("%d", randSource.Uint64())
	name := fmt.Sprintf("%s %s", firstnames[randSource.IntN(len(firstnames)-1)], lastnames[randSource.IntN(len(lastnames)-1)])
	age := randSource.IntN(60)
	isStudent := age < 25
	myhobbies := make([]string, randSource.IntN(5)+1)
	for j := 0; j < len(myhobbies); j++ {
		myhobbies[j] = hobbies[randSource.IntN(len(hobbies)-1)]
	}
	bio := ""
	for j := 0; j < randSource.IntN(10)+1; j++ {
		bio += loremipsums[randSource.IntN(len(loremipsums)-1)] + " "
	}
	birthday := time.Now().AddDate(-age, -randSource.IntN(12), -randSource.IntN(28))

	// Generate pets
	pets := make([]TestPet, randSource.IntN(3))
	for j := 0; j < len(pets); j++ {
		// Grab stuff
		petType := petTypes[randSource.IntN(len(pets)-1)]
		petBreed := petBreeds[petType][randSource.IntN(len(petBreeds[petType])-1)]
		petName := petNames[petType][randSource.IntN(len(petNames[petType])-1)]
		petAge := randSource.IntN(15)
		petToys := make([]string, randSource.IntN(3))
		for k := 0; k < len(petToys); k++ {
			petToys[k] = funnyPetToys[randSource.IntN(len(funnyPetToys)-1)]
		}

		// Create pet
		pets[j] = TestPet{
			Name:  petName,
			Age:   petAge,
			Type:  petType,
			Breed: petBreed,
			Toys:  petToys,
		}
	}

	doc := TestData{
		Name:      name,
		Age:       age,
		Hobbies:   myhobbies,
		Bio:       bio,
		isStudent: isStudent,
		Birthday:  birthday,
	}

	return id, doc
}
