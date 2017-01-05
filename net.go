package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

func waitForNetwork() error {
	maxWait := time.Second * 3
	checkInterval := time.Second
	timeStarted := time.Now()

	for {
		interfaces, err := net.Interfaces()
		if err != nil {
			return err
		}

		// pretty basic check ...
		// > 1 as a lo device will already exist
		if len(interfaces) > 1 {
			return nil
		}

		if time.Since(timeStarted) > maxWait {
			return fmt.Errorf("Timeout after %s waiting for network", maxWait)
		}

		time.Sleep(checkInterval)
	}
}

func exitIfNetsetgoNotFound(netsetgoPath string) {
	if _, err := os.Stat(netsetgoPath); os.IsNotExist(err) {
		usefulErrorMsg := fmt.Sprintf(`
Unable to find the netsetgo binary at "%s".
netsetgo is an external binary used to configure networking.
You must download netsetgo, chown it to the root user and apply the setuid bit.
This can be done as follows:

wget "https://github.com/teddyking/netsetgo/releases/download/0.0.1/netsetgo"
sudo mv netsetgo /usr/local/bin/
sudo chown root:root /usr/local/bin/netsetgo
sudo chmod 4755 /usr/local/bin/netsetgo
`, netsetgoPath)

		fmt.Println(usefulErrorMsg)
		os.Exit(1)
	}
}
