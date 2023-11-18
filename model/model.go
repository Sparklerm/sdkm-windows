package model

// Config 结构定义了配置文件的结构
type Config struct {
	JDKDir  string `json:"JDK_DIR"`
	EnvName string `json:"JDK_ENV_NAME"`
}

type JdkVersion struct {
	Version     string `json:"version"`
	DownloadUrl string `json:"download"`
}

type JdkType struct {
	JdkType  string       `json:"JdkType"`
	Versions []JdkVersion `json:"versions"`
}

type ZuluApiResult struct {
	DownloadUrl string `json:"download_url"`
}
