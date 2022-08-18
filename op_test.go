package uni_filter

import (
	"fmt"
	"testing"
)

type Person struct {
	Name string
	Age  int
}

type Location struct {
	Jd      float64
	Wd      float64
	Country string
}

type Car struct {
	Price    float64
	Name     string
	Location Location
}

func (p Person) GetCars() *[]*Car {
	var cars []*Car
	cars = append(cars, &Car{
		Price: 1000,
		Name:  "lanbojini",
		Location: Location{
			Jd:      100,
			Wd:      100,
			Country: "china",
		},
	})

	cars = append(cars, &Car{
		Price: 1500,
		Name:  "lanboniuni",
		Location: Location{
			Jd:      100,
			Wd:      100,
			Country: "usa",
		},
	})

	cars = append(cars, &Car{
		Price: 1800,
		Name:  "mashaladi",
		Location: Location{
			Jd:      102,
			Wd:      103,
			Country: "hk",
		},
	})

	return &cars
}

func TestAll(t *testing.T) {

	//expr, err := Parse("!country=china or (!age__gte=18 and score__lte=9.0) or lover.name=jany or www__ex")
	//expr, err := Parse("pets[].habbits[]__length=39")
	expr, err := Parse("GetCars--call.Location.Country=china")
	//expr, err := Parse("lover.name = jany")
	if err != nil {
		t.Fatal(err)
	}

	p := Person{}
	matched := expr.Match(p)
	if matched {
		fmt.Println("match this person")
	}

	var cases = []struct {
		v       map[string]any
		matched bool
	}{
		{
			map[string]any{
				"name":    "myth",
				"country": "what",
				"age":     10,
				"score":   9.5,
				"lover": map[string]string{
					"name": "jany",
				},
				"pets": []map[string]any{
					{
						"name": "maizi",
						"habbits": []map[string]int{
							{
								"price": 30,
							},
						},
					},
				},
			},
			true,
		},
		{
			map[string]any{
				"name":    "jany",
				"country": "china",
				"age":     18,
				"score":   9.1,
			},
			true,
		},
		{
			map[string]any{
				"name":    "lili",
				"country": "usa",
				"age":     13,
				"score":   9.23,
				"www":     "",
				"pets": []map[string]any{
					{
						"name": "chengzi",
					},
				},
			},

			true,
		},
		{
			map[string]any{
				"name":    "marry",
				"country": "usa",
				"age":     18,
				"score":   9.3,
			},
			true,
		},
	}

	for i, c := range cases {
		matched := expr.Match(c.v)
		if matched != c.matched {
			t.Errorf("test case at index %d failed\n", i)
		}
	}
}
