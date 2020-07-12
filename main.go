package main

import (
	"fmt"
	"hlsParse/ts"
	"io"
	"os"
)

func main() {

	tsFile, err := os.Open("./tmp.ts")
	if err != nil {
		fmt.Printf("os.Open failed, err:%v\n", err)
	}

	buf := make([]byte, 0)
	t := new(ts.Ts)

	nums := 0
	times := 1

	for true {

		tmpBuf := make([]byte, 1024)
		length, errRead := tsFile.Read(tmpBuf)
		if errRead != nil && errRead != io.EOF {
			fmt.Printf("tsFile.Read failed, err:%v\n", errRead)
		}
		if length == 0 {
			fmt.Printf("read end\n")
			break
		}
		buf = append(buf, tmpBuf[:length]...)

		buf, err = t.Parse(buf)
		if err != nil {
			fmt.Printf("t.Parse failed, err:%v\n", err)
			os.Exit(-1)
		}
		if errRead == io.EOF {
			fmt.Printf("already read and deal")
			os.Exit(0)
		}

		nums++
		if nums == times {
			break
		}
	}

	fmt.Println()

}