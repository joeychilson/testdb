package gen

import (
	"time"
)

type Person struct {
	FirstName string    `fake:"{firstname}"`
	LastName  string    `fake:"{lastname}"`
	Age       int32     `fake:"{number:18,100}"`
	Salary    int       `fake:"{number:30000,1000000}"`
	StartDate time.Time `fake:"{date}"`
	Languages []string  `fakesize:"2,6"`
}
