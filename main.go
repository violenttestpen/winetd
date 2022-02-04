package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/spf13/viper"
	"github.com/violenttestpen/winetd/pkg/log"
	"github.com/violenttestpen/winetd/pkg/windows"
)

const integrityUsage = "Run service with assigned integrity level: %v"

var (
	integrity = "Untrusted"
	verbosity = 0

	configName      string
	integrityLevels []string
)

var logger log.Log

func init() {
	integrityLevels = make([]string, 0, len(windows.SidWinIntegrityLevels))
	for k := range windows.SidWinIntegrityLevels {
		integrityLevels = append(integrityLevels, k)
	}
}

func loadConfig(configName string) map[string]Service {
	viper.SetConfigName(configName)
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		logger.Fatal(err)
	}

	var services map[string]Service
	if err := viper.Unmarshal(&services); err != nil {
		logger.Fatal(err)
	}
	return services
}

func saveDefaultConfig(configName string) {
	viper.SetConfigFile(configName)
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	viper.Set("default", defaultService)
	viper.SafeWriteConfig()
}

func main() {
	flag.StringVar(&configName, "c", "config.json", "Path to services config file")
	flag.StringVar(&integrity, "integrity", integrity, fmt.Sprintf(integrityUsage, integrityLevels))
	flag.IntVar(&verbosity, "verbosity", verbosity, "Verbosity mode (0-2)")
	flag.Parse()

	logger = log.NewLogger(verbosity)
	if flag.Arg(0) == "init" {
		saveDefaultConfig(configName)
		logger.Info(configName, "is updated")
		return
	}
	services := loadConfig(configName)

	var wg sync.WaitGroup
	wg.Add(len(services))
	for name, service := range services {
		if service.Disable {
			continue
		}

		logger.Info("Running server", name)
		go func(service Service) {
			if err := run(service, integrity); err != nil {
				logger.Error(err)
			}
			wg.Done()
		}(service)
	}
	wg.Wait()
}

func run(svc Service, integrity string) error {
	sid, ok := windows.SidWinIntegrityLevels[integrity]
	if !ok {
		return windows.ErrInvalidIntegrityLevel
	}

	token, err := windows.GetIntegrityLevelToken(sid)
	if err != nil {
		return err
	}
	defer token.Close()

	if _, err := os.Stat(svc.Server); err != nil {
		return err
	}

	server, err := net.Listen(svc.Protocol, net.JoinHostPort(svc.Bind, strconv.Itoa(svc.Port)))
	if err != nil {
		return err
	}
	defer server.Close()

	var waitMutex sync.Mutex
	return doListen(server, func(conn net.Conn) {
		if svc.Wait {
			waitMutex.Lock()
			defer waitMutex.Unlock()
		}

		if err := handleConn(conn, svc, syscall.Token(token)); err != nil {
			logger.Error(err)
		}
	})
}

func doListen(server net.Listener, handler func(net.Conn)) error {
	for {
		conn, err := server.Accept()
		if err != nil {
			return err
		}
		logger.Info("New connection from", conn.RemoteAddr().String())
		go handler(conn)
	}
}

func handleConn(conn net.Conn, svc Service, token syscall.Token) error {
	defer conn.Close()

	cmd := exec.Command(svc.Server)
	cmd.SysProcAttr = &syscall.SysProcAttr{Token: token}
	cmd.Stdout = conn

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	job, err := windows.NewJobFromProcess(cmd.Process)
	if err != nil {
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
		return err
	}

	go func() {
		if svc.Timeout > 0 {
			conn.SetDeadline(time.Now().Add(time.Duration(svc.Timeout) * time.Second))
		}
		if _, err := io.Copy(stdin, conn); err != nil {
			logger.Warning(err)
		}

		if err := windows.TerminateJob(job); err != nil {
			logger.Error(err)
		}
	}()

	logger.Info(fmt.Sprintf("Started service \"%s\" with pid %d", svc.Server, cmd.Process.Pid))

	if err := cmd.Wait(); err != nil {
		logger.Warning(err)
	}

	return nil
}
