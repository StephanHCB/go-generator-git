package docs

import (
	"log"
)

// these do nothing really, but they make tests and their log output way more readable

func Given(s string) {
	log.Print("Given ", s)
}

func When(s string) {
	log.Print("When ", s)
}

func Then(s string) {
	log.Print("Then ", s)
}
