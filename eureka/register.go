package eureka

import (
	"fmt"
	slog "github.com/go-eden/slf4go"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd/eureka"
	"github.com/hudl/fargo"
	"net"
	"os"
)

func Register(eurekaHost string, port int, serviceName, homePage string) (*eureka.Registrar, error) {

	ip, err := GetLocalIP()

	if err != nil {
		return nil, err
	}

	hostname, err := os.Hostname()

	if err != nil {
		return nil, err
	}

	logger := log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestamp)

	slog.Info(fmt.Sprintf("IP: %s hostname: %s", ip, hostname))

	var fargoConfig fargo.Config
	// Target Eureka server(s).
	fargoConfig.Eureka.ServiceUrls = []string{eurekaHost}
	// How often the subscriber should poll for updates.
	fargoConfig.Eureka.PollIntervalSeconds = 1

	instanceTest1 := &fargo.Instance{
		HostName:         fmt.Sprintf("%s", hostname),
		InstanceId:       fmt.Sprintf("%s:%s:%d", hostname, serviceName, port),
		Port:             port,
		PortEnabled:      true,
		App:              serviceName,
		IPAddr:           ip,
		VipAddress:       ip,
		SecureVipAddress: ip,
		HealthCheckUrl:   fmt.Sprintf("http://%s:%d/health", ip, port),
		StatusPageUrl:    fmt.Sprintf("http://%s:%d/info", ip, port),
		HomePageUrl:      fmt.Sprintf("http://%s:%d/%s", ip, port, homePage),
		Status:           fargo.UP,
		DataCenterInfo:   fargo.DataCenterInfo{Name: fargo.MyOwn},
		CountryId:        1,
		LeaseInfo:        fargo.LeaseInfo{RenewalIntervalInSecs: 30},
	}

	// Create a Fargo connection and a Eureka registrar.
	fargoConnection := fargo.NewConnFromConfig(fargoConfig)
	registrar1 := eureka.NewRegistrar(&fargoConnection, instanceTest1, log.With(logger, "component", "registrar1"))

	// Register one instance.
	registrar1.Register()

	return registrar1, nil
}

func GetLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", fmt.Errorf("can't get InterfaceAddrs %v", err)
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", fmt.Errorf("can't find any IP")
}
