//go:build grpc_test
// +build grpc_test

package tests

import (
	"context"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
	act_device_api "gitlab.ozon.dev/qa/classroom-4/act-device-api/pkg/act-device-api/gitlab.ozon.dev/qa/classroom-4/act-device-api/pkg/act-device-api"
	"google.golang.org/grpc"
)

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

func Test_grpc_ListDevicesV1(t *testing.T) {
	ctx := context.Background()
	addr := "localhost:8082"
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	require.NoError(t, err)
	defer conn.Close()
	client := act_device_api.NewActDeviceApiServiceClient(conn)

	req := &act_device_api.ListDevicesV1Request{Page: 1, PerPage: 5}
	res, err := client.ListDevicesV1(ctx, req)

	require.NoError(t, err)
	require.NotEmpty(t, res.Items)
}

func Test_grpc_CreateDeviceV1(t *testing.T) {
	ctx := context.Background()
	addr := "localhost:8082"
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	require.NoError(t, err)
	defer conn.Close()
	client := act_device_api.NewActDeviceApiServiceClient(conn)

	platformEx := "Ios"
	const userIdEx uint64 = 131313
	req := &act_device_api.CreateDeviceV1Request{Platform: platformEx, UserId: userIdEx}
	res, err := client.CreateDeviceV1(ctx, req)

	require.NoError(t, err)
	require.NotNil(t, res.DeviceId)

	//Now match the results of the Describe Device to equal contents of the device created
	reqCheck := &act_device_api.DescribeDeviceV1Request{DeviceId: res.DeviceId}
	resCheck, err := client.DescribeDeviceV1(ctx, reqCheck)

	require.NoError(t, err)
	require.NotNil(t, resCheck.Value)
	require.Equal(t, platformEx, resCheck.Value.Platform)
	require.Equal(t, userIdEx, resCheck.Value.UserId)
}

func Test_grpc_DescribeDeviceV1(t *testing.T) {

	ctx := context.Background()
	addr := "localhost:8082"
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	require.NoError(t, err)
	defer conn.Close()
	client := act_device_api.NewActDeviceApiServiceClient(conn)

	req := &act_device_api.ListDevicesV1Request{Page: 1, PerPage: 5}
	res, err := client.ListDevicesV1(ctx, req)
	existentDeviceId := res.Items[0].Id
	reqDesc := &act_device_api.DescribeDeviceV1Request{DeviceId: existentDeviceId}
	resDesc, err := client.DescribeDeviceV1(ctx, reqDesc)

	require.NoError(t, err)
	require.Equal(t, res.Items[0].Id, resDesc.Value.Id)
	require.Equal(t, res.Items[0].Platform, resDesc.Value.Platform)
	require.Equal(t, res.Items[0].UserId, resDesc.Value.UserId)
}

func Test_grpc_RemoveDeviceV1(t *testing.T) {

	ctx := context.Background()
	addr := "localhost:8082"
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	require.NoError(t, err)
	defer conn.Close()
	client := act_device_api.NewActDeviceApiServiceClient(conn)

	// To ensure that we have the needed device each run, we need to get the first existent from a list.
	req := &act_device_api.ListDevicesV1Request{Page: 1, PerPage: 5}
	res, err := client.ListDevicesV1(ctx, req)
	existentDeviceId := res.Items[0].Id
	require.NoError(t, err)
	require.NotEmpty(t, res.Items)

	reqDel := &act_device_api.RemoveDeviceV1Request{DeviceId: existentDeviceId}
	resDel, err := client.RemoveDeviceV1(ctx, reqDel)

	require.NoError(t, err)
	require.Equal(t, true, resDel.Found)
}

func Test_grpc_UpdateDeviceV1(t *testing.T) {

	ctx := context.Background()
	addr := "localhost:8082"
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	require.NoError(t, err)
	defer conn.Close()
	client := act_device_api.NewActDeviceApiServiceClient(conn)

	req := &act_device_api.ListDevicesV1Request{Page: 1, PerPage: 5}
	res, err := client.ListDevicesV1(ctx, req)
	newPlatformEx := "TeslaOX"
	const newUserIdEx uint64 = 666666
	existentDeviceId := res.Items[0].Id
	reqUpd := &act_device_api.UpdateDeviceV1Request{Platform: newPlatformEx, UserId: newUserIdEx, DeviceId: existentDeviceId}
	resUpd, err := client.UpdateDeviceV1(ctx, reqUpd)

	require.NoError(t, err)
	require.Equal(t, true, resUpd.Success)
}

func Test_grpc_SendNotificationV1(t *testing.T) {
	ctx := context.Background()
	addr := "localhost:8082"
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	require.NoError(t, err)
	defer conn.Close()

	clientDevice := act_device_api.NewActDeviceApiServiceClient(conn)
	req := &act_device_api.ListDevicesV1Request{Page: 1, PerPage: 5}
	res, err := clientDevice.ListDevicesV1(ctx, req)
	require.NoError(t, err)
	require.NotEmpty(t, res.Items)
	existentDeviceId := res.Items[0].Id

	clientNotif := act_device_api.NewActNotificationApiServiceClient(conn)
	reqNotif := &act_device_api.SendNotificationV1Request{Notification: &act_device_api.Notification{NotificationId: 0, DeviceId: existentDeviceId, Username: "0", Message: "Hello", Lang: act_device_api.Language_LANG_ENGLISH, NotificationStatus: act_device_api.Status_STATUS_CREATED}}
	resNotif, err := clientNotif.SendNotificationV1(ctx, reqNotif)

	require.NoError(t, err)
	require.NotNil(t, resNotif.NotificationId)

}

func Test_grpc_GetNotificationV1(t *testing.T) {
	ctx := context.Background()
	addr := "localhost:8082"
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	require.NoError(t, err)
	defer conn.Close()

	clientDevice := act_device_api.NewActDeviceApiServiceClient(conn)
	req := &act_device_api.ListDevicesV1Request{Page: 1, PerPage: 5}
	res, err := clientDevice.ListDevicesV1(ctx, req)
	require.NoError(t, err)
	require.NotEmpty(t, res.Items)
	existentDeviceId := res.Items[0].Id

	clientNotif := act_device_api.NewActNotificationApiServiceClient(conn)
	reqNotif := &act_device_api.GetNotificationV1Request{DeviceId: existentDeviceId}
	resNotif, err := clientNotif.GetNotification(ctx, reqNotif)

	require.NoError(t, err)
	require.NotNil(t, resNotif.Notification)
}

func Test_grpc_AckNotificationV1(t *testing.T) {
	ctx := context.Background()
	addr := "localhost:8082"
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	require.NoError(t, err)
	defer conn.Close()

	clientNotif := act_device_api.NewActNotificationApiServiceClient(conn)
	var exNotifId uint64 = 1
	reqNotif := &act_device_api.AckNotificationV1Request{NotificationId: exNotifId}
	resNotif, err := clientNotif.AckNotification(ctx, reqNotif)

	require.NoError(t, err)
	require.Equal(t, true, resNotif.Success)

}

func Test_grpc_SubscribeNotifV1(t *testing.T) {
	ctx := context.Background()
	addr := "localhost:8082"
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	require.NoError(t, err)
	defer conn.Close()

	clientNotif := act_device_api.NewActNotificationApiServiceClient(conn)
	clientDevice := act_device_api.NewActDeviceApiServiceClient(conn)
	req := &act_device_api.ListDevicesV1Request{Page: 1, PerPage: 5}
	res, err := clientDevice.ListDevicesV1(ctx, req)
	require.NoError(t, err)
	require.NotEmpty(t, res.Items)
	existentDeviceId := res.Items[0].Id

	reqNotif := &act_device_api.SubscribeNotificationRequest{DeviceId: existentDeviceId}
	resNotif, err := clientNotif.SubscribeNotification(ctx, reqNotif)

	require.NoError(t, err)
	log.Print(resNotif)
}

// func Test_All_gRPC(t *testing.T) {
// 	t.Run("Should properly provide a list of existing devices", Test_grpc_ListDevicesV1)
// 	t.Run("Should properly create a given device", Test_grpc_CreateDeviceV1)
// 	t.Run("Should properly describe a given device", Test_grpc_DescribeDeviceV1)
// 	t.Run("Should properly delete a given device", Test_grpc_RemoveDeviceV1)
// 	t.Run("Should properly update a given device", Test_grpc_UpdateDeviceV1)
// 	t.Run("Should properly create a notification with given contents", Test_grpc_SendNotificationV1)
// 	t.Run("Should properly get a notification with given DeviceID", Test_grpc_GetNotificationV1)
// 	t.Run("Should properly acknowledge the notification delivery", Test_grpc_AckNotificationV1)
// 	t.Run("Should properly subscribe for a given device", Test_grpc_SubscribeNotifV1)
// }
