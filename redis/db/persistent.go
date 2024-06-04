package db

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc64"
	"os"
	"redis/interface/reply"
	"redis/protocol"

	"github.com/sirupsen/logrus"
)

func init() {
	registerCMD("bgsave", cmdBgsave)
}

// cmdBgsave is the handler for the "bgsave" command.
func cmdBgsave(db *DB, args [][]byte) reply.Reply {
	go func() {
		err := generateRDB(db, "dump.rdb")
		if err != nil {
			logrus.Error("Failed to save RDB: ", err)
		}
	}()
	return protocol.MakeSimpleStr("Background saving started")
}

func generateRDB(db *DB, fileName string) error {
	result := bytes.NewBuffer([]byte{})
	br := bufio.NewWriter(result)

	_, err := br.WriteString("REDIS")
	if err != nil {
		return err
	}

	// TODO 不要直接写二进制
	_, err = br.Write([]byte{0x30, 0x30, 0x30, 0x36})
	if err != nil {
		return err
	}

	// _, err = br.Write([]byte{0xfe})
	// if err != nil {
	// 	return err
	// }

	// eof
	err = br.WriteByte(0xff)
	if err != nil {
		return err
	}

	if err = br.Flush(); err != nil {
		return err
	}

	// print result.Bytes to 16进制
	fmt.Printf("%x\n", result.Bytes())

	// calculate crc64 of writed content
	// write crc64
	checksum := crc64.Checksum(result.Bytes(), crc64.MakeTable(crc64.ISO))
	// write checksum to buffer
	// NOTE
	err = binary.Write(br, binary.BigEndian, checksum)
	if err != nil {
		return err
	}

	if err = br.Flush(); err != nil {
		return err
	}

	tmpRDB := "tmp.rdb"
	file, err := os.Create(tmpRDB)
	if err != nil {
		// TODO if err is file exist
		return err
	}

	logrus.Info("result bytes", result.Bytes())
	// _, err = file.Write(result.Bytes())
	if err != nil {
		return err
	}

	defer func() {
		file.Close()
		os.Remove(tmpRDB)
	}()

	// save bytes to tmp file
	_, err = file.Write(result.Bytes())
	if err != nil {
		return err
	}

	// rename tmp to file
	if _, err = os.Stat(fileName); err != nil {
		_, err = os.Create(fileName)
		if err != nil {
			return err
		}
	}
	err = os.Rename(tmpRDB, fileName)

	return err
}
