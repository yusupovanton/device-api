package tests

import (
	"database/sql"
	"testing"

	// "log"
	"context"

	"github.com/ozontech/allure-go/pkg/framework/suite"

	// "github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	// "github.com/ozontech/allure-go/pkg/framework/runner"
	_ "github.com/lib/pq"
	"github.com/ozontech/allure-go/pkg/framework/runner"
	act_device_api "gitlab.ozon.dev/qa/classroom-4/act-device-api/pkg/act-device-api/gitlab.ozon.dev/qa/classroom-4/act-device-api/pkg/act-device-api"
	"google.golang.org/grpc"
)

var ctx = context.Background()

type BeforeAfterDemoSuite struct {
	suite.Suite
}

func (s *BeforeAfterDemoSuite) BeforeEach(t provider.T) {
	t.NewStep("Before Test Step")

	var connStr = "psql -h localhost -p 5432 -U docker -d act_device_api"
	sql.Open("postgres", connStr)

}

func (s *BeforeAfterDemoSuite) AfterEach(t provider.T) {
	t.NewStep("After Test Step")
}

func (s *BeforeAfterDemoSuite) TestSQL(t provider.T) {

}

func TestBeforesAfters(t *testing.T) {
	t.Parallel()
	suite.RunSuite(t, new(BeforeAfterDemoSuite))
	// В тесте UpdateDevice убедиться что изменения проливаются в БД, а не вернулись из кеша;
	ctx := context.Background()
	addr := "localhost:8082"
	runner.Run(t, "Test_grpc_UpdateDeviceV1", func(t provider.T) {
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
