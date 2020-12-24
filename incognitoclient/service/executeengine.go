package service

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"os/exec"
	"strings"
)

type ExecuteEngine struct {
	PathCmd string
}

func (e *ExecuteEngine) Do(payload *Payload) (interface{}, string, error) {
	payloadData, _ := json.Marshal(payload.Data)
	command := fmt.Sprintf("%v execute -m=%v -p=%v -d='%v'", e.PathCmd, payload.Method, payload.PrivateKey, string(payloadData))
	//fmt.Println(command)

	result, err := e.runCmd(command)
	if err != nil {
		return nil, "", err
	}

	var data interface{}
	if err := json.Unmarshal([]byte(result), &data); err != nil {
		return nil, "", errors.Wrap(err, "json.Unmarshal")
	}

	return data, result, nil
}

func (e *ExecuteEngine) runCmd(command string) (string, error)  {
	parts := strings.Fields(command)
	cmd := exec.Command(parts[0], parts[1:]...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}
	if err := cmd.Start(); err != nil {
		return "", err
	}

	scanner := bufio.NewScanner(stdout)
	res := make([]string, 0)
	for scanner.Scan() {
		line := scanner.Text()
		res = append(res, line)
	}

	if err = scanner.Err(); err != nil {
		return "", err
	}

	if err := cmd.Wait(); err != nil {
		return "", err
	}

	if len(res) > 0 {
		return res[len(res)-1], nil
	}

	return "", errors.New("Result is empty")
}