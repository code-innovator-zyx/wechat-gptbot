package streamlit_app

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
)

/*
* @Author: zouyx
* @Email: zouyx@knowsec.com
* @Date:   2024/5/28 14:44
* @Package:
 */

func RunStreamlit() {
	if !checkPythonInstallation() {
		// 系统没有安装python，无法启动webui
		fmt.Printf("系统没有安装python，无法启动webui")
		return
	}
	// 创建一个 exec.Command 实例
	cmd := exec.Command("streamlit", "run", "./streamlit_app/home.py")
	// 获取命令的标准输出和标准错误输出管道
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("Error getting stdout pipe: %v\n", err)
		return
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		fmt.Printf("Error getting stderr pipe: %v\n", err)
		return
	}

	// 启动命令
	err = cmd.Start()
	if err != nil {
		fmt.Printf("Error starting command: %v\n", err)
		return
	}

	// 打印标准输出
	go pyOutput(stdoutPipe)
	go pyOutput(stderrPipe)

	// 等待命令完成
	err = cmd.Wait()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Command finished successfully.\n")

}

// 打印输出
func pyOutput(pipe io.ReadCloser) {
	defer pipe.Close()
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading from pipe: %v\n", err)
	}
}

func checkPythonInstallation() bool {
	pythonCmds := []string{"python3", "python"}
	for _, cmd := range pythonCmds {
		if _, err := exec.LookPath(cmd); err == nil {
			return true
		}
	}
	return false
}
