package registry

import (
	consul "github.com/hashicorp/consul/api"
	"github.com/justcy/ygo/ygo/registry/iface"
	"github.com/justcy/ygo/ygo/utils"
	"github.com/justcy/ygo/ygo/ylog"
	"sync"
)

type ConsulRegistry struct {
	sync.Mutex
	client *consul.Client
	config       *consul.Config
	QueryOptions *consul.QueryOptions
}

func (c *ConsulRegistry) Init() {
	ylog.Debug("consul Registry init")
	if c.config == nil {
		ylog.Debug(utils.GlobalObject.ConsulAddress)
		c.config = &consul.Config{
			Address: utils.GlobalObject.ConsulAddress,
		}
		ylog.Infof("%v,$v",c.config,c.client)
	}
	if c.client == nil {
		var err error
		c.client, err = consul.NewClient(c.config)
		if err != nil {
			ylog.Info(err)
		}
	}
}

func (c *ConsulRegistry) Register(service iface.Service) {
	c.Init()
	tags := encodeMetadata(service.Metadata)
	tags = append(tags, encodeVersion(service.Version)...)
	registration := consul.AgentServiceRegistration{
		ID:      service.Id,
		Name:    service.Name,
		Port:    service.Port,
		Tags:    tags,
		Address: service.Address,
		Check: &consul.AgentServiceCheck{
			TCP:                            service.Address,
			Timeout:                        "5s",
			Interval:                       "5s",
			DeregisterCriticalServiceAfter: "600s",
		},
	}
	ylog.Infof("%v,$v",c.config,c.client)
	if err := c.client.Agent().ServiceRegister(&registration); err != nil {
		ylog.Error(err)
	}
}

func (c *ConsulRegistry) UnRegister(service iface.Service) {
	c.Init()
	c.client.Agent().ServiceDeregister(service.Id)
}
func (c *ConsulRegistry) UnRegisterById(id string) {
	c.Init()
	c.client.Agent().ServiceDeregister(id)
}

func (c *ConsulRegistry) GetService(name string,opt ...GetOptions) ([]*iface.Service, error){
	c.Init()
	var resp []*consul.ServiceEntry
	var err error
	resp ,_,err = c.client.Health().Service(name,"",false,c.QueryOptions)
	if err != nil {
		return nil, err
	}
	var services []*iface.Service
	for _, s := range resp {
		ylog.Debugf("%s,%s,%s:%d",s.Checks[0].Status,s.Service.Service,s.Service.Address,s.Service.Port)
		if s.Service.Service != name {
			continue
		}
		// address is service address
		address := s.Service.Address

		version, _ := decodeVersion(s.Service.Tags)
		// use node address
		if len(address) == 0 {
			address = s.Node.Address
		}
		service := &iface.Service{
			Metadata: decodeMetadata(s.Service.Tags),
			Name:      s.Service.Service,
			Version:   version,
			Address: s.Service.Address,
			Port: s.Service.Port,
		}
		services = append(services, service)
	}
	return services,nil
}
