package memory

import (
	"context"
	"errors"
	"movieexample.com/pkg/discovery"
	"sync"
	"time"
)

type serviceName string
type instanceID string

type Registry struct {
	sync.RWMutex
	serviceAddrs map[serviceName]map[instanceID]*serviceInstance
}

type serviceInstance struct {
	hostPort   string
	lastActive time.Time
}

func NewRegistry() *Registry {
	return &Registry{serviceAddrs: map[serviceName]map[instanceID]*serviceInstance{}}
}

func (r *Registry) Register(_ context.Context, id string, name string, hostPort string) error {
	r.Lock()
	defer r.Unlock()

	srvName := serviceName(name)

	if _, ok := r.serviceAddrs[srvName]; !ok {
		r.serviceAddrs[srvName] = map[instanceID]*serviceInstance{}
	}
	r.serviceAddrs[srvName][instanceID(id)] = &serviceInstance{hostPort: hostPort, lastActive: time.Now()}
	return nil
}

func (r *Registry) Deregister(_ context.Context, id string, name string) error {
	r.Lock()
	defer r.Unlock()

	srvName := serviceName(name)

	if _, ok := r.serviceAddrs[srvName]; !ok {
		return nil
	}
	delete(r.serviceAddrs[srvName], instanceID(id))
	return nil
}

func (r *Registry) ServiceAddresses(ctx context.Context, name string) ([]string, error) {
	r.RLock()
	defer r.RUnlock()
	if len(r.serviceAddrs[serviceName(name)]) == 0 {
		return nil, discovery.ErrNotFound
	}
	var res []string

	for _, i := range r.serviceAddrs[serviceName(name)] {
		if i.lastActive.Before(time.Now().Add(-5 * time.Second)) {
			continue
		}
		res = append(res, i.hostPort)
	}
	return res, nil
}

func (r *Registry) ReportHealthyState(id string, name string) error {
	r.Lock()
	defer r.Unlock()

	srvName := serviceName(name)
	srvId := instanceID(id)
	if _, ok := r.serviceAddrs[srvName]; !ok {
		return errors.New("service is not registered yet")
	}
	if _, ok := r.serviceAddrs[srvName][srvId]; !ok {
		return errors.New("service instance is not registered yet")
	}

	r.serviceAddrs[srvName][srvId].lastActive = time.Now()

	return nil
}
