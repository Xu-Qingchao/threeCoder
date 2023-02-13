package ffmpeg

import (
	"fmt"
	"os/exec"
)

func CallCommandRun(cmdName string, args []string) (string, error) {
	cmd := exec.Command(cmdName, args...)
	fmt.Println("CallCommand Run 参数=> ", args)
	fmt.Println("CallCommand Run 执行命令=> ", cmd)
	bytes, err := cmd.Output()
	if err != nil {
		fmt.Println("CallCommand Run 出错了.....", err.Error())
		fmt.Println(err)
		return "", err
	}
	resp := string(bytes)
	fmt.Println(resp)
	fmt.Println("CallCommand Run 调用完成.....")
	return resp, nil
}

// 根据URL调用ffmpeg 获取截图
func GetIpcScreenShot(ffmpegPath string, url string, screenShotPath string) string {
	var params []string
	params = append(params, "-i")
	params = append(params, url)
	params = append(params, "-ss")
	params = append(params, "00:00:01")
	params = append(params, "-vframes")
	params = append(params, "1")
	params = append(params, screenShotPath)

	resp, err := CallCommandRun(ffmpegPath, params)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("获取截图出错，url为--->", url)
	}
	return resp
}
