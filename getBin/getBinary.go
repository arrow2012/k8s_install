package getBin

import (
	"k8s_install/common/config"
	"k8s_install/common/utils"
)

var (
	log = config.Logger

)


//func ImagePull(imageName string, cli *client.Client, ctx context.Context) error {
//	authConfig := types.AuthConfig{
//		//Username: docker_reg_username,
//		//Password: docker_reg_password,
//	}
//	encodedJSON, err := json.Marshal(authConfig)
//	if err != nil {
//		return err
//	}
//	authStr := base64.URLEncoding.EncodeToString(encodedJSON)
//
//	out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{RegistryAuth: authStr})
//	if err != nil {
//		return err
//	}
//	defer out.Close()
//	io.Copy(os.Stdout, out)
//	return nil
//}
//
//func GetBinary()  {
//	cli, err := client.NewEnvClient()
//	utils.CheckErrExit(err)
//	ctx := context.Background()
//	err = ImagePull(kubeImage, cli, ctx)
//	utils.CheckErrExit(err)
//	err = ImagePull(etcdImage, cli, ctx)
//	utils.CheckErrExit(err)
//	err = ImagePull(pauseImage, cli, ctx)
//	utils.CheckErrExit(err)
//}

func UnariveFile()  {
	utils.ExecCmd("/bin/bash down-base.sh all","getBin/",nil)
}

func Task()  {
	UnariveFile()
}