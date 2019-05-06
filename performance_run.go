/*
 * Copyright (c) 2019 QLC Chain Team
 *
 * This software is released under the MIT License.
 * https://opensource.org/licenses/MIT
 */

package commands

import (
	"bufio"
	"fmt"
	"github.com/abiosoft/ishell"
	"github.com/qlcchain/go-qlc/common/types"
	"github.com/qlcchain/go-qlc/rpc"
	"os"
	"strconv"
	"time"
)

func performanceTest() {
	if interactive {
		c := &ishell.Cmd{
			Name: "perfTest",
			Help: "performance test",
			Func: func(c *ishell.Context) {
				err := performanceTestAction()
				if err != nil {
					Warn(err)
					return
				}
				Info("performance test success!")
			},
		}
		shell.AddCmd(c)
	}
}

func performanceTestAction() error {
	//l := ledger.NewLedger("D:\\goproject\\go-qlc\\lc2\\ledger", nil)

	client, err := rpc.Dial(endpointP)
	if err != nil {
		return err
	}
	defer client.Close()

	loc, _ := time.LoadLocation("Asia/Shanghai")

	f, err := os.Open("D:\\testLedger")
	if err != nil {
		fmt.Println("open test ledger err", err.Error())
		return err
	}

	finfo, _ := f.Stat()

	block := &types.StateBlock{}
	r := bufio.NewReaderSize(f, int(finfo.Size()))
	var buf [1024]byte
	var count int64
	var h types.Hash
	var blockNum int

	fmt.Printf("%s Start\n", time.Unix(0, time.Now().UnixNano()).In(loc).Format("2006-01-02 15:04:05.000000"))
	for {
		f.Seek(count, 0)
		l, err := r.ReadBytes(':')
		if err != nil {
			fmt.Println("read len err:", err.Error())
			break
		}
		count += int64(len(l))
		ll, err := strconv.Atoi(string(l[:len(l) - 1]))
		count += int64(ll)

		rl := buf[:ll]
		_, err = r.Read(rl)
		if err != nil {
			fmt.Println("read bytes err:", err.Error())
			break
		}

		block.UnmarshalMsg(rl)
		err = client.Call(&h, "ledger_processTes", block)
		blockNum++
	}
	fmt.Printf("%s End:send %d blocks\n", time.Unix(0, time.Now().UnixNano()).In(loc).Format("2006-01-02 15:04:05.000000"), blockNum)

	return nil
}
