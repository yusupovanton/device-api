package tests

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/runner"
	act_device_api "gitlab.ozon.dev/qa/classroom-4/act-device-api/pkg/act-device-api/gitlab.ozon.dev/qa/classroom-4/act-device-api/pkg/act-device-api"
	"google.golang.org/grpc"
)

func Test_grpc_ListDevicesV1(t *testing.T) {
	ctx := context.Background()
	runner.Run(t, "Test_grpc_ListDevicesV1", func(t provider.T) {
		t.Title("GRPC test for listing devices")
		t.Description("Dials the service, gets the lists of devices")

		t.WithNewStep("Testing ListDevicesV1", func(sCtx provider.StepCtx) {

			sCtx.NewStep("Dialing the service")
			addr := "localhost:8082"
			conn, err := grpc.Dial(addr, grpc.WithInsecure())
			sCtx.Require().NoError(err)
			defer conn.Close()

			sCtx.NewStep("Requesting the list of devices")
			client := act_device_api.NewActDeviceApiServiceClient(conn)
			req := &act_device_api.ListDevicesV1Request{Page: 1, PerPage: 5}
			sCtx.WithNewAttachment("list request", allure.Text, []byte(fmt.Sprintf("%v", req)))
			res, err := client.ListDevicesV1(ctx, req)
			sCtx.WithNewAttachment("list response", allure.Text, []byte(fmt.Sprintf("%v", res)))
			sCtx.Require().NoError(err)
			sCtx.Require().NotEmpty(res.Items)
		})
	})
}

func Test_grpc_CreateDeviceV1(t *testing.T) {
	ctx := context.Background()
	addr := "localhost:8082"
	runner.Run(t, "Test_grpc_CreateDeviceV1", func(t provider.T) {
		t.Title("GRPC test for creating devices")
		t.Description("Dials the service, creates the device, then checks the contents of the device")

		t.WithNewStep("Testing ListDevicesV1", func(sCtx provider.StepCtx) {

			sCtx.NewStep("Dialing the service")
			conn, err := grpc.Dial(addr, grpc.WithInsecure())
			sCtx.Require().NoError(err)
			defer conn.Close()

			sCtx.NewStep("Creating the device")
			client := act_device_api.NewActDeviceApiServiceClient(conn)
			platformEx := "Ios"
			const userIdEx uint64 = 199283
			req := &act_device_api.CreateDeviceV1Request{Platform: platformEx, UserId: userIdEx}
			res, err := client.CreateDeviceV1(ctx, req)
			sCtx.Require().NoError(err)
			sCtx.Require().NotNil(res.DeviceId)

			sCtx.NewStep("Check Contents of the device created")
			//Now match the results of the Describe Device to equal contents of the device created
			reqCheck := &act_device_api.DescribeDeviceV1Request{DeviceId: res.DeviceId}
			resCheck, err := client.DescribeDeviceV1(ctx, reqCheck)
			sCtx.Require().NoError(err)
			sCtx.Require().NotNil(resCheck.Value)
			sCtx.Require().Equal(platformEx, resCheck.Value.Platform)
			sCtx.Require().Equal(userIdEx, resCheck.Value.UserId)
		})
	})
}

func Test_grpc_DescribeDeviceV1(t *testing.T) {
	ctx := context.Background()
	addr := "localhost:8082"
	runner.Run(t, "Test_grpc_DescribeDeviceV1", func(t provider.T) {

		t.Title("GRPC test for describing devices")
		t.Description("Dials the service, creates the device, then checks the contents of the device")

		t.WithNewStep("Testing DescribeDeviceV1", func(sCtx provider.StepCtx) {

			sCtx.NewStep("Dialing the service")
			conn, err := grpc.Dial(addr, grpc.WithInsecure())
			sCtx.Require().NoError(err)
			defer conn.Close()

			sCtx.NewStep("Creating the device")
			client := act_device_api.NewActDeviceApiServiceClient(conn)
			platformEx := "Ios"
			const userIdEx uint64 = 515379
			req := &act_device_api.CreateDeviceV1Request{Platform: platformEx, UserId: userIdEx}
			res, err := client.CreateDeviceV1(ctx, req)
			sCtx.Require().NoError(err)
			sCtx.Require().NotNil(res.DeviceId)

			sCtx.NewStep("Describing the device")
			reqDesc := &act_device_api.DescribeDeviceV1Request{DeviceId: res.DeviceId}
			resDesc, err := client.DescribeDeviceV1(ctx, reqDesc)
			sCtx.Require().NoError(err)
			sCtx.Require().Equal(res.DeviceId, resDesc.Value.Id)
			sCtx.Require().Equal(platformEx, resDesc.Value.Platform)
			sCtx.Require().Equal(userIdEx, resDesc.Value.UserId)
		})
	})
}

func Test_grpc_RemoveDeviceV1(t *testing.T) {
	ctx := context.Background()
	addr := "localhost:8082"
	runner.Run(t, "Test_grpc_RemoveDeviceV1", func(t provider.T) {
		t.Title("GRPC test for removing devices")
		t.Description("Dials the service, creates the device, then deletes it")

		t.WithNewStep("Testing RemoveDeviceV1", func(sCtx provider.StepCtx) {

			sCtx.NewStep("Dialing the service")
			conn, err := grpc.Dial(addr, grpc.WithInsecure())
			sCtx.Require().NoError(err)
			defer conn.Close()

			sCtx.NewStep("Creating the device")
			client := act_device_api.NewActDeviceApiServiceClient(conn)
			platformEx := "TeslaOX"
			const userIdEx uint64 = 171717
			req := &act_device_api.CreateDeviceV1Request{Platform: platformEx, UserId: userIdEx}
			res, err := client.CreateDeviceV1(ctx, req)
			sCtx.Require().NoError(err)
			sCtx.Require().NotNil(res.DeviceId)

			sCtx.NewStep("Deleting the device")
			reqDel := &act_device_api.RemoveDeviceV1Request{DeviceId: res.DeviceId}
			resDel, err := client.RemoveDeviceV1(ctx, reqDel)
			sCtx.Require().NoError(err)
			sCtx.Require().Equal(true, resDel.Found)
		})
	})
}

func Test_grpc_UpdateDeviceV1(t *testing.T) {

	ctx := context.Background()
	addr := "localhost:8082"
	runner.Run(t, "Test_grpc_UpdateDeviceV1", func(t provider.T) {
		t.Title("GRPC test for updating devices")
		t.Description("Dials the service, creates the device, updates it and then checks the contents of the updated device")

		t.WithNewStep("Testing UpdateDeviceV1", func(sCtx provider.StepCtx) {

			sCtx.NewStep("Dialing the service")
			conn, err := grpc.Dial(addr, grpc.WithInsecure())
			sCtx.Require().NoError(err)
			defer conn.Close()

			sCtx.NewStep("Creating the device")
			client := act_device_api.NewActDeviceApiServiceClient(conn)
			platformEx := "MegaOS"
			const userIdEx uint64 = 991991
			req := &act_device_api.CreateDeviceV1Request{Platform: platformEx, UserId: userIdEx}
			res, err := client.CreateDeviceV1(ctx, req)
			sCtx.Require().NoError(err)
			sCtx.Require().NotNil(res.DeviceId)

			sCtx.NewStep("Updating the device")
			newPlatformEx, newUserIdEx := "Ios", uint64(100101)
			reqUpd := &act_device_api.UpdateDeviceV1Request{Platform: newPlatformEx, UserId: newUserIdEx, DeviceId: res.DeviceId}
			resUpd, err := client.UpdateDeviceV1(ctx, reqUpd)
			sCtx.Require().NoError(err)
			sCtx.Require().Equal(true, resUpd.Success)

			sCtx.NewStep("Describing the device")
			reqDesc := &act_device_api.DescribeDeviceV1Request{DeviceId: res.DeviceId}
			resDesc, err := client.DescribeDeviceV1(ctx, reqDesc)
			sCtx.Require().NoError(err)
			sCtx.Require().Equal(newPlatformEx, resDesc.Value.Platform)
			sCtx.Require().Equal(newUserIdEx, resDesc.Value.UserId)
		})
	})
}

func Test_grpc_SendNotificationV1(t *testing.T) {
	ctx := context.Background()
	addr := "localhost:8082"
	runner.Run(t, "Test_grpc_SendNotificationV1", func(t provider.T) {

		t.Title("GRPC test for sending a notification")
		t.Description("Dials the service, creates the device, then sends a notification to it")

		t.WithNewStep("Testing SendNotificationV1", func(sCtx provider.StepCtx) {

			sCtx.NewStep("Dialing the service")
			conn, err := grpc.Dial(addr, grpc.WithInsecure())
			sCtx.Require().NoError(err)
			defer conn.Close()

			sCtx.NewStep("Creating the device")
			client := act_device_api.NewActDeviceApiServiceClient(conn)
			platformEx := "DeltaOS"
			const userIdEx uint64 = 129097
			req := &act_device_api.CreateDeviceV1Request{Platform: platformEx, UserId: userIdEx}
			res, err := client.CreateDeviceV1(ctx, req)
			sCtx.Require().NoError(err)
			sCtx.Require().NotNil(res.DeviceId)

			sCtx.NewStep("Sending the notification to the device")
			clientNotif := act_device_api.NewActNotificationApiServiceClient(conn)
			reqSendNotif := &act_device_api.SendNotificationV1Request{Notification: &act_device_api.Notification{NotificationId: 0, DeviceId: res.DeviceId, Username: "0", Message: "Hello", Lang: act_device_api.Language_LANG_ENGLISH, NotificationStatus: act_device_api.Status_STATUS_CREATED}}
			resSendNotif, err := clientNotif.SendNotificationV1(ctx, reqSendNotif)
			sCtx.Require().NoError(err)
			sCtx.Require().NotNil(resSendNotif.NotificationId)
		})
	})
}

func Test_grpc_GetNotificationV1(t *testing.T) {
	ctx := context.Background()
	addr := "localhost:8082"
	runner.Run(t, "Test_grpc_GetNotificationV1", func(t provider.T) {
		t.Title("GRPC test for sending a notification")
		t.Description("Dials the service, creates the device, then sends a notification to it and checks the notification")

		t.WithNewStep("Testing GetNotificationV1", func(sCtx provider.StepCtx) {

			sCtx.NewStep("Dialing the service")
			conn, err := grpc.Dial(addr, grpc.WithInsecure())
			sCtx.Require().NoError(err)
			defer conn.Close()

			sCtx.NewStep("Creating the device")
			client := act_device_api.NewActDeviceApiServiceClient(conn)
			platformEx := "DeltaOS"
			const userIdEx uint64 = 129097
			req := &act_device_api.CreateDeviceV1Request{Platform: platformEx, UserId: userIdEx}
			res, err := client.CreateDeviceV1(ctx, req)
			sCtx.Require().NoError(err)
			sCtx.Require().NotNil(res.DeviceId)

			sCtx.NewStep("Sending the notification to the device")
			clientNotif := act_device_api.NewActNotificationApiServiceClient(conn)
			reqSendNotif := &act_device_api.SendNotificationV1Request{Notification: &act_device_api.Notification{NotificationId: 0, DeviceId: res.DeviceId, Username: "0", Message: "Thank You", Lang: act_device_api.Language_LANG_ESPANOL, NotificationStatus: act_device_api.Status_STATUS_CREATED}}
			resSendNotif, err := clientNotif.SendNotificationV1(ctx, reqSendNotif)
			sCtx.Require().NoError(err)
			sCtx.Require().NotNil(resSendNotif.NotificationId)

			sCtx.NewStep("Getting the notification from the device")
			reqGetNotif := &act_device_api.GetNotificationV1Request{DeviceId: res.DeviceId}
			resGetNotif, err := clientNotif.GetNotification(ctx, reqGetNotif)
			sCtx.Require().NoError(err)
			sCtx.Require().NotNil(resGetNotif.Notification)
		})
	})
}

func Test_grpc_AckNotificationV1(t *testing.T) {

	ctx := context.Background()
	addr := "localhost:8082"
	runner.Run(t, "Test_grpc_AckNotificationV1", func(t provider.T) {
		t.Title("GRPC test for sending a notification")
		t.Description("Dials the service, creates the device, then sends a notification to it and acknowledges the delivery")
		t.WithNewStep("Testing AckNotificationV1", func(sCtx provider.StepCtx) {

			sCtx.NewStep("Dialing the service")
			conn, err := grpc.Dial(addr, grpc.WithInsecure())
			sCtx.Require().NoError(err)
			defer conn.Close()

			sCtx.NewStep("Creating the device")
			client := act_device_api.NewActDeviceApiServiceClient(conn)
			platformEx := "DeltaOS"
			const userIdEx uint64 = 129097
			req := &act_device_api.CreateDeviceV1Request{Platform: platformEx, UserId: userIdEx}
			res, err := client.CreateDeviceV1(ctx, req)
			sCtx.Require().NoError(err)
			sCtx.Require().NotNil(res.DeviceId)

			sCtx.NewStep("Sending the notification to the device")
			clientNotif := act_device_api.NewActNotificationApiServiceClient(conn)
			reqSendNotif := &act_device_api.SendNotificationV1Request{Notification: &act_device_api.Notification{NotificationId: 0, DeviceId: res.DeviceId, Username: "0", Message: "Hello", Lang: act_device_api.Language_LANG_ITALIAN, NotificationStatus: act_device_api.Status_STATUS_CREATED}}
			resSendNotif, err := clientNotif.SendNotificationV1(ctx, reqSendNotif)
			sCtx.Require().NoError(err)
			sCtx.Require().NotNil(resSendNotif.NotificationId)

			sCtx.NewStep("Acknowledging the notification delivery status")
			reqAckNotif := &act_device_api.AckNotificationV1Request{NotificationId: resSendNotif.NotificationId}
			resAckNotif, err := clientNotif.AckNotification(ctx, reqAckNotif)
			sCtx.Require().NoError(err)
			sCtx.Require().Equal(true, resAckNotif.Success)
		})
	})
}

func Test_grpc_SubscribeNotifV1(t *testing.T) {
	ctx := context.Background()
	addr := "localhost:8082"
	runner.Run(t, "Test_grpc_SubscribeNotifV1", func(t provider.T) {
		t.Title("GRPC test for subscribing to notifications for a given device")
		t.Description("Dials the service, creates the device, then subscribes for delivery")
		t.WithNewStep("Testing SubscribeNotifV1", func(sCtx provider.StepCtx) {

			sCtx.NewStep("Dialing the service")
			conn, err := grpc.Dial(addr, grpc.WithInsecure())
			sCtx.Require().NoError(err)
			defer conn.Close()

			sCtx.NewStep("Creating the device")
			client := act_device_api.NewActDeviceApiServiceClient(conn)
			platformEx := "DeltaOS"
			const userIdEx uint64 = 129097
			req := &act_device_api.CreateDeviceV1Request{Platform: platformEx, UserId: userIdEx}
			res, err := client.CreateDeviceV1(ctx, req)
			sCtx.Require().NoError(err)
			sCtx.Require().NotNil(res.DeviceId)

			sCtx.NewStep("Subscribing to the device")
			clientNotif := act_device_api.NewActNotificationApiServiceClient(conn)
			reqSubNotif := &act_device_api.SubscribeNotificationRequest{DeviceId: res.DeviceId}
			resSubNotif, err := clientNotif.SubscribeNotification(ctx, reqSubNotif)
			sCtx.Require().NoError(err)
			log.Print(resSubNotif)
		})
	})
}
