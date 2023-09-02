package notifications

import (
	"context"
	"feuerwehr-ooe-backend/model"
	"feuerwehr-ooe-backend/notifications/template_id"
	"feuerwehr-ooe-backend/utils"
	"fmt"
	"github.com/OneSignal/onesignal-go-api"
	"log"
	"os"
)

var apiClient *onesignal.APIClient
var appAuth context.Context

func InitializeOneSignal() {
	configuration := onesignal.NewConfiguration()
	apiClient = onesignal.NewAPIClient(configuration)

	appAuth = context.WithValue(context.Background(), onesignal.AppAuth, os.Getenv("ONESIGNAL_REST_API_KEY"))
}

var appId string

func SendOneSignalNotification(einsatz model.Einsatz) {
	variable := utils.LoadEnvVariable("ONESIGNAL_APP_ID", &appId)
	log.Println("Sending notification for", einsatz.ID, "to", variable)
	notification := *onesignal.NewNotification(variable)
	notification.TemplateId = onesignal.PtrString(template_id.EinsatzGlobal)
	notification.IncludedSegments = []string{"Total Subscriptions"}

	notification.Headings = getHeader(einsatz)
	notification.Contents = getContent(einsatz)

	resp, r, err := apiClient.DefaultApi.CreateNotification(appAuth).Notification(notification).Execute()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.CreateNotification``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateNotification`: CreateNotificationSuccessResponse
	fmt.Fprintf(os.Stdout, "Response from `DefaultApi.CreateNotification`: %v\n", resp)
}

func getHeader(einsatz model.Einsatz) onesignal.NullableStringMap {
	headings := onesignal.NewStringMap()
	headings.SetEn("[" + einsatz.Alarmstufe.String() + "] Incident in " + einsatz.Adresse.Ort)
	headings.SetDe("[" + einsatz.Alarmstufe.String() + "] Einsatz in " + einsatz.Adresse.Ort)
	return *onesignal.NewNullableStringMap(headings)
}

func getContent(einsatz model.Einsatz) onesignal.NullableStringMap {
	contents := onesignal.NewStringMap()
	contents.SetEn(einsatz.Einsatztyp.Name + " in " + einsatz.Adresse.PrettyName())
	contents.SetDe(einsatz.Einsatztyp.Name + " in " + einsatz.Adresse.PrettyName())
	return *onesignal.NewNullableStringMap(contents)
}
