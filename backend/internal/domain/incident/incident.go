package incident

import (
	"net/url"
	"time"

	"github.com/flowck/doberman/internal/domain"
)

type Incident struct {
	id                  domain.ID
	title               string
	description         string
	attachments         []url.URL
	reporterEmail       string
	notificationMethods []interface{}
	peopleToBeNotified  []interface{}
	startedAt           time.Time
	resolvedAt          time.Time
	comments            []Comment
}

func (i *Incident) ID() domain.ID {
	return i.id
}

func (i *Incident) Title() string {
	return i.title
}

func (i *Incident) Description() string {
	return i.description
}

func (i *Incident) Attachments() []url.URL {
	return i.attachments
}

func (i *Incident) ReporterEmail() string {
	return i.reporterEmail
}

func (i *Incident) NotificationMethods() []interface{} {
	return i.notificationMethods
}

func (i *Incident) PeopleToBeNotified() []interface{} {
	return i.peopleToBeNotified
}

func (i *Incident) StartedAt() time.Time {
	return i.startedAt
}

func (i *Incident) ResolvedAt() time.Time {
	return i.resolvedAt
}

func (i *Incident) Comments() []Comment {
	return i.comments
}

type Comment struct {
	id          domain.ID
	authorID    domain.ID
	body        string
	attachments []url.URL
}

func (c *Comment) Id() domain.ID {
	return c.id
}

func (c *Comment) AuthorID() domain.ID {
	return c.authorID
}

func (c *Comment) Body() string {
	return c.body
}

func (c *Comment) Attachments() []url.URL {
	return c.attachments
}
