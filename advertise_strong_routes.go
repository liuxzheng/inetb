package main

import (
	"github.com/codeout/inetb/client"
	"log"
	"time"
)

func advertiseStrongRoutes(client1 *client.Client, client2 *client.Client) error {
	log.Print("Start benchmarking - Advertise strong routes from client2")

	if err := client2.NexthopSelf(); err != nil {
		log.Fatal(err)
	}

	if err := client2.AcceptExport(); err != nil {
		log.Fatal(err)
	}

	reports := make([]*Report, 0)
	sent := 0
	received := 0
	timeout := 10

	for tick := 0; tick < timeout; tick++ {
		func() {
			for {
				select {
				case update := <-client2.Updates:
					if routerId, _ := client2.RouterId(); update.Nexthop == routerId {
						sent += len(update.Raw.NLRI)
					}
					tick = 0
				default:
					return
				}
			}
		}()

		func() {
			for {
				select {
				case update := <-client1.Updates:
					if routerId, _ := client1.RouterId(); update.Nexthop != routerId {
						received += len(update.Raw.NLRI)
					}
					tick = 0
				default:
					return
				}
			}
		}()

		report := &Report{
			Time:     time.Now().Format("15:04:05"),
			Sent:     sent,
			Received: received,
		}

		log.Print(report.String())
		reports = append(reports, report)

		if sent == 0 && received == 0 {
			tick = 0
		}

		time.Sleep(time.Second)
	}

	log.Print("Stop benchmarking - Advertise new routes from client2")

	return WriteReport("advertise_strong_routes.json", reports)
}
