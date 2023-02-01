package gen

import (
	"context"
	"fmt"
	"log"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/jackc/pgtype"
	"github.com/joeychilson/testdb/db"
	"github.com/joeychilson/testdb/db/sqlc"
	"github.com/joeychilson/testdb/gen"
	"github.com/spf13/cobra"
)

type PeopleGenCmd struct {
	db *db.Postgres
}

func NewPeopleGen(db *db.Postgres) *PeopleGenCmd {
	return &PeopleGenCmd{db: db}
}

func (g *PeopleGenCmd) PeopleGenCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "people [flags]",
		Short: "Generate realistic fake data for the people schema",
		RunE: func(cmd *cobra.Command, args []string) error {
			people, err := cmd.Flags().GetInt("people")
			if err != nil {
				return err
			}

			return g.handleCmd(cmd.Context(), people)
		},
	}

	cmd.Flags().IntP("people", "p", 1, "The number of people to generate")
	return cmd
}

func (g *PeopleGenCmd) handleCmd(ctx context.Context, people int) error {
	log.Printf("Generating %d people the people schema", people)

	for i := 0; i < people; i++ {
		gofakeit.Seed(0)

		var genPerson *gen.Person
		gofakeit.Struct(&genPerson)

		salary := pgtype.Numeric{}
		err := salary.Set(genPerson.Salary)
		if err != nil {
			log.Fatalf("failed to set salary: %v", err)
		}

		phoneNumber := gofakeit.Phone()

		phone := pgtype.JSON{}
		err = phone.Set(`{"area_code": "` + phoneNumber[0:3] + `","phone": "` + phoneNumber + `"}`)
		if err != nil {
			log.Fatalf("failed to set phone: %v", err)
		}

		fullName := fmt.Sprintf("%s %s", genPerson.FirstName, genPerson.LastName)

		person := sqlc.CreatePersonParams{
			FirstName: genPerson.FirstName,
			LastName:  genPerson.LastName,
			FullName:  fullName,
			Age:       genPerson.Age,
			Salary:    salary,
			StartDate: genPerson.StartDate,
			Phone:     phone,
			Languages: genPerson.Languages,
		}

		err = g.db.CreatePerson(ctx, person)
		if err != nil {
			log.Fatalf("failed to create person: %v", err)
		}

		log.Printf("Created person %s", fullName)
	}
	return nil
}
