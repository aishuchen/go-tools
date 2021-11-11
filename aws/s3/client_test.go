package s3

import (
	"context"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/spf13/viper"
	"gitlab.hypers.com/server-go/tools/aws/common"
	"gitlab.hypers.com/server-go/tools/config"
	"gitlab.hypers.com/server-go/tools/internal"
)

var ctx = context.Background()

func testNewFromViper(t *testing.T) *Client {
	configFilePath := internal.GetTestConfigFile()
	if err := config.SetGlobalConfig(configFilePath); err != nil {
		t.Fatal(err)
	}
	v := viper.GetViper()
	awsCfg, err := common.NewAWSCfgFromViper(v)
	if err != nil {
		t.Fatal(err)
	}
	clnt, err := NewFromViper(viper.GetViper(), awsCfg)
	if err != nil {
		t.Fatal(err)
	}
	return clnt
}

func genKey() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}

func TestNewFromFile(t *testing.T) {
	testNewFromViper(t)
}

func TestClient_Listdir(t *testing.T) {
	clnt := testNewFromViper(t)
	files, err := clnt.Listdir(ctx, "")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(files)
}

// 此测试用例会执行 Client_Put
func TestClient_Upload(t *testing.T) {
	clnt := testNewFromViper(t)
	local := "/tmp/for_upload.txt"
	f, err := os.Create(local)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = f.WriteString("test s3_storage for put"); err != nil {
		t.Fatal(err)
	}
	if err := f.Sync(); err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
	key := genKey()
	if err := clnt.Upload(ctx, local, key); err != nil {
		t.Fatal(err)
	}
}

// 此测试用例会测试 Client_Open
// 执行此测试用例之前请确保已手动上传过文件
func TestClient_Download(t *testing.T) {
	clnt := testNewFromViper(t)
	local := "/tmp/for_open.txt"
	err := clnt.Download(ctx, "for_open.txt", local)
	if err != nil {
		t.Fatal(err)
	}
	_, err = os.Open(local) // 确认文件是否下载成功
	if err != nil {
		t.Fatalf(`file does not download to "%s"`, local)
	}
}

func TestClient_Delete(t *testing.T) {
	// 执行此测试用例之前请确保已手动上传过文件
	clnt := testNewFromViper(t)
	err := clnt.Delete(ctx, "for_delete.txt")
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_Exists(t *testing.T) {
	// 执行此测试用例之前请确保已手动上传过文件
	clnt := testNewFromViper(t)
	exists, err := clnt.Exists(ctx, "for_open.txt")
	if !exists {
		t.Fatalf("file already exists, but returned false, err: %v", err)
	}

	exists, err = clnt.Exists(ctx, "non-exists.txt")
	if exists {
		t.Fatalf("file dose not exists, but returned ture, err: %v", err)
	}
}
