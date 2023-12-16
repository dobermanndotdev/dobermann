package psql

import (
	"context"
	"sync"

	"github.com/flowck/dobermann/backend/internal/domain"
	"github.com/flowck/dobermann/backend/internal/domain/account"
	"github.com/flowck/dobermann/backend/internal/domain/monitor"
)

type AccountRepositoryMock struct {
	mutex    *sync.RWMutex
	accounts map[domain.ID]*account.Account
}

func NewAccountRepositoryMock() AccountRepositoryMock {
	return AccountRepositoryMock{
		mutex:    &sync.RWMutex{},
		accounts: make(map[domain.ID]*account.Account),
	}
}

func (p AccountRepositoryMock) Insert(ctx context.Context, acc *account.Account) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.accounts[acc.ID()] = acc

	return nil
}

// Users

func NewUserRepositoryMock() UserRepositoryMock {
	return UserRepositoryMock{}
}

type UserRepositoryMock struct {
	mutex *sync.RWMutex
	users map[domain.ID]*account.User
}

func (p UserRepositoryMock) FindByID(ctx context.Context, id domain.ID) (*account.User, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	user, exists := p.users[id]
	if !exists {
		return nil, account.ErrUserNotFound
	}

	return user, nil
}

func (p UserRepositoryMock) FindByEmail(ctx context.Context, email account.Email) (*account.User, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	return nil, nil
}

func (p UserRepositoryMock) Insert(ctx context.Context, user *account.User) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return nil
}

// Incidents

type IncidentRepositoryMock struct {
	mutex     *sync.RWMutex
	incidents map[domain.ID]*monitor.Incident
}

func NewIncidentRepositoryMock() IncidentRepositoryMock {
	return IncidentRepositoryMock{
		mutex:     &sync.RWMutex{},
		incidents: make(map[domain.ID]*monitor.Incident),
	}
}

func (i IncidentRepositoryMock) FindByID(ctx context.Context, id domain.ID) (*monitor.Incident, error) {
	return nil, nil
}

func (i IncidentRepositoryMock) Create(ctx context.Context, monitorID domain.ID, incident *monitor.Incident) error {
	return nil
}

func (i IncidentRepositoryMock) Update(
	ctx context.Context,
	id, monitorID domain.ID,
	fn func(incident *monitor.Incident) error,
) error {
	incident := i.incidents[id]
	return fn(incident)
}

func (i IncidentRepositoryMock) AppendIncidentAction(
	ctx context.Context,
	incidentID domain.ID,
	action *monitor.IncidentAction,
) error {
	return nil
}

// Monitor

type MonitorRepositoryMock struct {
	mutex    *sync.RWMutex
	monitors map[domain.ID]*monitor.Monitor
}

func NewMonitorRepositoryMock() MonitorRepositoryMock {
	return MonitorRepositoryMock{
		mutex:    &sync.RWMutex{},
		monitors: make(map[domain.ID]*monitor.Monitor),
	}
}

func (m MonitorRepositoryMock) SaveCheckResult(
	ctx context.Context,
	monitorID domain.ID,
	checkResult *monitor.CheckResult,
) error {
	return nil
}

func (m MonitorRepositoryMock) Insert(ctx context.Context, monitor *monitor.Monitor) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.monitors[monitor.ID()] = monitor

	return nil
}

func (m MonitorRepositoryMock) FindByID(ctx context.Context, ID domain.ID) (*monitor.Monitor, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	foundMonitor, exists := m.monitors[ID]
	if !exists {
		return nil, monitor.ErrMonitorNotFound
	}

	return foundMonitor, nil
}

func (m MonitorRepositoryMock) Update(ctx context.Context, ID domain.ID, fn func(monitor *monitor.Monitor) error) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	//TODO implement me
	panic("implement me")
}

func (m MonitorRepositoryMock) UpdateForCheck(ctx context.Context, fn func(foundMonitors []*monitor.Monitor) error) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	monitorList := make([]*monitor.Monitor, len(m.monitors))

	index := 0
	for _, v := range m.monitors {
		monitorList[index] = v
		index++
	}

	err := fn(monitorList)
	if err != nil {
		return err
	}

	return nil
}

func (m MonitorRepositoryMock) Delete(ctx context.Context, ID domain.ID) error {
	delete(m.monitors, ID)
	return nil
}
