package tests

import (
	"database/sql"
	"fmt"
	"testing"

	// "log"
	"context"

	"github.com/ozontech/allure-go/pkg/framework/suite"

	// "github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	// "github.com/ozontech/allure-go/pkg/framework/runner"
	_ "github.com/lib/pq"
	act_device_api "gitlab.ozon.dev/qa/classroom-4/act-device-api/pkg/act-device-api/gitlab.ozon.dev/qa/classroom-4/act-device-api/pkg/act-device-api"
	"google.golang.org/grpc"
)

var ctx = context.Background()

const (
	host     = "localhost"
	port     = 5432
	user     = "docker"
	password = "docker"
	dbname   = "postgres"
)

type BeforeAfterDemoSuite struct {
	suite.Suite
}

func (s *BeforeAfterDemoSuite) BeforeEach(t provider.T) {
	t.NewStep("Before Test Step")
	fmt.Println("Test started")

}

func (s *BeforeAfterDemoSuite) AfterEach(t provider.T) {
	t.NewStep("After Test Step")
	fmt.Println("Test completed")
}

func (s *BeforeAfterDemoSuite) TestSQLUpdateDeviceV1(t provider.T) {

	// В тесте UpdateDevice убедиться что изменения проливаются в БД, а не вернулись из кеша;

	t.Title("SQL tests of the update function")
	t.Description("Creates the device, updates it and then looks for the updates in the postgres service")
	ctx := context.Background()
	addr := "localhost:8082"

	t.NewStep("Dialing the service")
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	t.Require().NoError(err)
	defer conn.Close()

	t.NewStep("Creating the device")
	client := act_device_api.NewActDeviceApiServiceClient(conn)
	platformEx := "MegaOS"
	const userIdEx uint64 = 991991
	req := &act_device_api.CreateDeviceV1Request{Platform: platformEx, UserId: userIdEx}
	res, err := client.CreateDeviceV1(ctx, req)
	t.Require().NoError(err)
	t.Require().NotNil(res.DeviceId)

	t.NewStep("Updating the device")
	newPlatformEx, newUserIdEx := "Ios", uint64(100101)
	reqUpd := &act_device_api.UpdateDeviceV1Request{Platform: newPlatformEx, UserId: newUserIdEx, DeviceId: res.DeviceId}
	resUpd, err := client.UpdateDeviceV1(ctx, reqUpd)
	t.Require().NoError(err)
	t.Require().Equal(true, resUpd.Success)

	t.NewStep("Describing the device")
	reqDesc := &act_device_api.DescribeDeviceV1Request{DeviceId: res.DeviceId}
	resDesc, err := client.DescribeDeviceV1(ctx, reqDesc)
	t.Require().NoError(err)
	t.Require().Equal(newPlatformEx, resDesc.Value.Platform)
	t.Require().Equal(newUserIdEx, resDesc.Value.UserId)

	t.NewStep("SQL check")
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	t.Require().NoError(err)
	queryUpdatedDevice := fmt.Sprintf("SELECT user_id, platform FROM devices WHERE id = %d;", res.DeviceId)
	row := db.QueryRow(queryUpdatedDevice)
	device := new(Device)
	err = row.Scan(&device.UserId, &device.Platform)
	t.Require().NoError(err)
	t.Require().Equal(newUserIdEx, device.UserId)
	t.Require().Equal(newPlatformEx, device.Platform)

}

func TestSQLUpdateDevice(t *testing.T) {
	t.Parallel()
	suite.RunSuite(t, new(BeforeAfterDemoSuite))

}
