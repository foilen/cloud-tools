package main

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/dns/armdns"
)

func main() {
	subscriptionID := os.Getenv("AZURE_SUBSCRIPTION_ID")

	// Create an Azure DNS client
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatal("Failed to obtain Azure authentication credentials", err)
	}
	clientFactory, err := armdns.NewClientFactory(subscriptionID, cred, nil)
	dnsClient := clientFactory.NewRecordSetsClient()

	// Get environment
	keepAlive := os.Getenv("KEEP_ALIVE")

	// Get the values from the arguments
	if len(os.Args) < 4 {
		log.Fatal("Missing arguments: resourceGroupName zoneName domain")
	}
	resourceGroupName := os.Args[1]
	zoneName := os.Args[2]
	domain := os.Args[3]

	// Get the subdomain based on the full domain and zone and remove the dot at the end
	subdomain := strings.TrimSuffix(domain, zoneName)
	subdomain = strings.TrimSuffix(subdomain, ".")

	// Obtain the current record
	recordSet, err := dnsClient.Get(context.Background(), resourceGroupName, zoneName, subdomain, armdns.RecordTypeA, nil)
	if err != nil {
		if !strings.Contains(err.Error(), "404") {
			log.Fatal("Failed to obtain the current record ", err)
		}
	}

	// Find the IP address
	previousIP := ""
	if recordSet.Properties != nil && len(recordSet.Properties.ARecords) == 1 {
		previousIP = *recordSet.Properties.ARecords[0].IPv4Address
	}

	// Loop to update DNS record every 10 minutes
	for {

		// Obtain the current public IP address
		currentIP, err := getCurrentIP()
		if err != nil {
			log.Fatal("Failed to obtain the public IP ", err)
		}

		// Check if the IP address needs to be updated
		if previousIP != currentIP {
			log.Println("Updating DNS record", domain, "with new IP address", currentIP)

			recordSet := armdns.RecordSet{
				Properties: &armdns.RecordSetProperties{
					TTL: to.Ptr[int64](300),
					ARecords: []*armdns.ARecord{
						{
							IPv4Address: &currentIP,
						},
					},
				},
			}

			_, err = dnsClient.CreateOrUpdate(context.Background(), resourceGroupName, zoneName, subdomain, armdns.RecordTypeA, recordSet, nil)
			if err != nil {
				log.Fatal("Failed to update the record", err)
			}
			previousIP = currentIP
		} else {
			log.Println("IP address is up-to-date")
		}

		if keepAlive == "true" {
			// Pause for 5 minutes before the next update
			time.Sleep(5 * time.Minute)
		} else {
			break
		}
	}
}

// Get the current public IP address
func getCurrentIP() (string, error) {
	resp, err := http.Get("https://checkip.foilen.com")
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	ip := string(body)

	return ip, nil
}
