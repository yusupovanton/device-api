//go:build http_test
// +build http_test

package tests

import (
	"io"
	"net/http"

	// "net/url"
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

const baseurl string = "http://localhost:8080"

type AutoGenerated struct {
	Notification *Notification `json:"notification"`
}

type Notification struct {
	NotificationID     string `json:"notificationId"`
	DeviceID           string `json:"deviceId"`
	Username           string `json:"username"`
	Message            string `json:"message"`
	Lang               string `json:"lang"`
	NotificationStatus string `json:"notificationStatus"`
}

func Test_http_ListDevicesV1(t *testing.T) {
	resp, err := http.Get(baseurl + "/api/v1/devices?page=1&perPage=15")
	assert.NoError(t, err)

	if resp.StatusCode != http.StatusOK {
		respBody, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		t.Fatalf("Invalid status code. Expected %d, got %d. Response: %s",
			http.StatusOK, resp.StatusCode, respBody)
	}
}

func Test_http_CreateDeviceV1(t *testing.T) {
	values := map[string]string{"platform": "Ios", "userId": "123123"}
	json_data, err := json.Marshal(values)
	assert.NoError(t, err)

	resp, err := http.Post(baseurl+"/api/v1/devices", "application/json",
		bytes.NewBuffer(json_data))

	if resp.StatusCode != http.StatusOK {
		respBody, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		t.Fatalf("Invalid status code. Expected %d, got %d. Response: %s",
			http.StatusOK, resp.StatusCode, respBody)
	}
}

func Test_http_DescribeDeviceV1(t *testing.T) {
	deviceId := "7"
	resp, err := http.Get(baseurl + "/api/v1/devices/" + deviceId)
	assert.NoError(t, err)

	if resp.StatusCode != http.StatusOK {
		respBody, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)

		t.Fatalf("Invalid status code. Expected %d, got %d for %s. Response: %s",
			http.StatusOK, resp.StatusCode, deviceId, respBody)
	}
}

func Test_http_DeleteDeviceV1(t *testing.T) {
	deviceId := "4"
	url := baseurl + "/api/v1/devices/" + deviceId
	req, err := http.NewRequest("DELETE", url, nil)
	assert.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	if resp.StatusCode != http.StatusOK {
		respBody, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		t.Fatalf("Invalid status code. Expected %d, got %d. Response: %s",
			http.StatusOK, resp.StatusCode, respBody)
	}
}

func Test_http_UpdateDeviceV1(t *testing.T) {
	values := map[string]string{"platform": "Ios", "userId": "123123"}
	json_data, err := json.Marshal(values)

	assert.NoError(t, err)

	deviceId := "2"
	url := baseurl + "/api/v1/devices/" + deviceId
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(json_data))
	assert.NoError(t, err)
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	if resp.StatusCode != http.StatusOK {
		respBody, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		t.Fatalf("Invalid status code. Expected %d, got %d. Response: %s",
			http.StatusOK, resp.StatusCode, respBody)
	}
}

func Test_http_GetNotifications(t *testing.T) {
	deviceId := "11"
	resp, err := http.Get(baseurl + "/api/v1/notification?deviceId=" + deviceId)
	assert.NoError(t, err)

	if resp.StatusCode != http.StatusOK {
		respBody, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		t.Fatalf("Invalid status code. Expected %d, got %d for %s. Response: %s",
			http.StatusOK, resp.StatusCode, deviceId, respBody)
	}

}

func Test_http_SendNotifV1(t *testing.T) {

	values := &AutoGenerated{&Notification{NotificationID: "0", DeviceID: "11", Username: "0", Message: "Hello", Lang: "LANG_ENGLISH", NotificationStatus: "STATUS_CREATED"}}
	json_data, err := json.Marshal(values)
	assert.NoError(t, err)
	resp, err := http.Post(baseurl+"/api/v1/notification", "application/json",
		bytes.NewBuffer(json_data))
	if resp.StatusCode != http.StatusOK {
		respBody, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		t.Fatalf("Invalid status code. Expected %d, got %d. Response: %s",
			http.StatusOK, resp.StatusCode, respBody)
	}
}
func Test_http_AckNotifV1(t *testing.T) {
	notifId := "1"
	url := baseurl + "/api/v1/notification/ack/" + notifId
	req, err := http.NewRequest("PUT", url, nil)
	assert.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	if resp.StatusCode != http.StatusOK {
		respBody, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		t.Fatalf("Invalid status code. Expected %d, got %d. Response: %s",
			http.StatusOK, resp.StatusCode, respBody)
	}
}

// func Test_http_SubscribeNotifV1(t *testing.T) {
// 	deviceId := "5"
// 	url := baseurl + "/api/v1/notification/subscribe/" + deviceId
// 	resp, err := http.Get(url)
// 	assert.NoError(t, err)

// 	if resp.StatusCode != http.StatusOK {
// 		respBody, err := io.ReadAll(resp.Body)
// 		assert.NoError(t, err)
// 		t.Fatalf("Invalid status code. Expected %d, got %d for %s. Response: %s",
// 			http.StatusOK, resp.StatusCode, deviceId, respBody)
// 	}
// }

// func Test_All_HTTP(t *testing.T) {
// 	t.Run("Should list devices", Test_http_ListDevicesV1)
// 	t.Run("Should create devices", Test_http_CreateDeviceV1)
// 	t.Run("Should describe devices", Test_http_DescribeDeviceV1)
// 	t.Run("Should update devices", Test_http_UpdateDeviceV1)
// 	t.Run("Should delete devices", Test_http_DeleteDeviceV1)
// 	t.Run("Should send out notifications for a given device", Test_http_SendNotifV1)
// 	t.Run("Should show notifications for a given device", Test_http_GetNotifications)
// 	t.Run("Should acknowledge delivery for a given notification", Test_http_AckNotifV1)
// 	// t.Run("Should subscribe for delivery for a given device", Test_http_SubscribeNotifV1)

// }
