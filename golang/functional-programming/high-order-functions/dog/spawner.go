package main

import "fmt"

type Dog struct {
	Name   string
	Breed  Breed
	Gender Gender
}

type NameToDogFunc func(string) Dog

const (
	Bulldog  Breed = "bulldog"
	Havanese       = "havanese"
	Cavalier       = "cavalier"
	Poodle         = "poodle"
)

type Gender string
type Breed string

const (
	Male   Gender = "male"
	Female        = "female"
)

func DogSpawner(breed Breed) func(gender Gender) NameToDogFunc {
	return func(gender Gender) NameToDogFunc {
		return func(name string) Dog {
			return Dog{
				Name:   name,
				Breed:  breed,
				Gender: gender,
			}
		}
	}
}

func breedDog(breed Breed, genderFunc func() Gender) NameToDogFunc {
	return func(name string) Dog {
		return Dog{
			Name:   name,
			Gender: genderFunc(),
			Breed:  breed,
		}
	}
}

func main() {
	maleHavaneseSpawner := DogSpawner(Havanese)(Male)
	femalePoodleSpawner := DogSpawner(Poodle)(Female)
	bucky := maleHavaneseSpawner("bucky")
	rocky := maleHavaneseSpawner("rocky")
	tipsy := femalePoodleSpawner("tipsy")

	fmt.Println(bucky)
	fmt.Println(rocky)
	fmt.Println(tipsy)

}
