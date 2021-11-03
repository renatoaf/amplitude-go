// Amplitude unofficial client for Go, inspired in their official SDK for Node.
//
//import (
//	"github.com/renatoaf/amplitude-go/amplitude"
//	"github.com/renatoaf/amplitude-go/amplitude/data"
//	"log"
//)
//
//const ApiKey = "<your-api-key>"
//
//func main() {
//	client := amplitude.NewDefaultClient(ApiKey)
//	client.Start()
//
//	err := client.LogEvent(&data.Event{
//		UserID:    "datamonster@gmail.com",
//		EventType: "test-event",
//		EventProperties: map[string]interface{}{
//			"source": "notification",
//		},
//		UserProperties: map[string]interface{}{
//			"age":    25,
//			"gender": "female",
//		},
//	})
//
//	if err != nil {
//		log.Printf("failed to queue event: %v", err)
//	}
//
//	client.Shutdown()
//}
package amplitude