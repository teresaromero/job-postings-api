package storage

import (
	"job-postings-api/internal/models"
	"sort"
	"time"

	"github.com/go-faker/faker/v4"
)

func (s *Storage) Seed() error {
	for range 100 {
		p, err := fakePost()
		if err != nil {
			return err
		}
		s.postMap[p.ID] = *p
		s.companyIndex[p.Company]++
	}

	return nil
}

func fakePost() (*models.Post, error) {

	now := time.Now()

	randomDay, err := faker.RandomInt(1, 20, 1)
	if err != nil {
		return nil, err
	}

	createdAt := now.Add(time.Hour * -24 * time.Duration(randomDay[0]))

	updatedAt := createdAt.Add(time.Hour * 24)

	idx, err := faker.RandomInt(0, 4, 1)
	if err != nil {
		return nil, err
	}

	companies := map[int]string{
		0: "Google",
		1: "Facebook",
		2: "Amazon",
		3: "Apple",
		4: "Microsoft",
	}

	locations := map[int]string{
		0: "Zaragoza",
		1: "Madrid",
		2: "Los Angeles",
		3: "Remote",
		4: "Barcelona",
	}

	jobtype := map[int]models.JobType{
		0: models.FullTime,
		1: models.PartTime,
		2: models.Contract,
		3: models.Internship,
		4: models.Freelance,
	}

	salary, err := faker.RandomInt(1000, 100000, 2)
	if err != nil {
		return nil, err
	}
	sort.Ints(salary)

	fakePost := models.Post{
		ID:          faker.UUIDHyphenated(),
		Title:       faker.Sentence(),
		Company:     companies[idx[0]],
		Location:    locations[idx[0]],
		Description: faker.Paragraph(),
		Type:        jobtype[idx[0]],
		Salary:      salary,
		Perks:       []string{faker.Word(), faker.Word()},
		Extras:      faker.Paragraph(),
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
	return &fakePost, nil
}
