package account

import "github.com/flowck/doberman/internal/domain"

type Team struct {
	id      domain.ID
	name    string
	members []Member
}

func NewTeam(
	id domain.ID,
	name string,
	members []Member,
) (*Team, error) {
	return &Team{
		id:      id,
		name:    name,
		members: members,
	}, nil
}

func (t *Team) ID() domain.ID {
	return t.id
}

func (t *Team) Name() string {
	return t.name
}

func (t *Team) Members() []Member {
	return t.members
}
