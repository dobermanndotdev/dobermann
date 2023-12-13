package psql

import (
	"context"
	"sync"

	"github.com/flowck/dobermann/backend/internal/domain"
	"github.com/flowck/dobermann/backend/internal/domain/monitor"
)

type MonitorRepositoryMock struct {
	mutex    *sync.RWMutex
	monitors map[domain.ID]*monitor.Monitor
}

func NewMonitorRepositoryMock() *MonitorRepositoryMock {
	return &MonitorRepositoryMock{
		mutex:    &sync.RWMutex{},
		monitors: make(map[domain.ID]*monitor.Monitor),
	}
}

func (m MonitorRepositoryMock) SaveCheckResult(ctx context.Context, monitorID domain.ID, checkResult *monitor.CheckResult) error {
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
	//TODO implement me
	panic("implement me")
}
