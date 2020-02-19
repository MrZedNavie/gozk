package persistence

import (
	"bufio"
	"os"

	"gozk/txn"
	"gozk/message"
)

/**
 * @Author: jiajianyun@jd.com
 * @Description:
 * @File:  file_txn_log
 * @Version: 1.0.0
 * @Date: 2020/2/19 1:46 下午
 */

type FileTxnLog struct {
	logFile        *os.File
	logBuf         *bufio.ReadWriter
	streamsToFlush []*bufio.ReadWriter
	lastZxidSeen   int64
	filePandding   *FilePadding
	buf            []byte
}

func (p *FileTxnLog) append(txnHeader *txn.TxnHeader, record interface{}) bool {
	var err error
	if txnHeader == nil {
		return false
	}
	if txnHeader.Zxid <= p.lastZxidSeen {
		//todo
	} else {
		p.lastZxidSeen = txnHeader.Zxid
	}

	if p.logFile == nil {
		p.logFile, err = os.OpenFile("test", os.O_RDWR, 777)
		if err != nil {
			//todo
			return false
		}
		p.logBuf = &bufio.ReadWriter{
			Reader: bufio.NewReader(p.logFile),
			Writer: bufio.NewWriter(p.logFile),
		}
		fileHeader := &FileHeader{
			Magic:   0,
			Version: 0,
			DbId:    0,
		}
		n, err :=message.EncodePacket(p.buf[0:], fileHeader)
		if err != nil {
			//todo
			return false
		}
		if _, err := p.logBuf.Write(p.buf[0:n]); err != nil {
			//todo
			return false
		}
		if err := p.logBuf.Flush(); err != nil {
			//todo
			return false
		}
		position, _ := p.logFile.Seek(0,1)
		p.filePandding.CurrentSize = position
		p.streamsToFlush = append(p.streamsToFlush, p.logBuf)
	}

	p.filePandding.PadFile(p.logFile)
	//todo
	return true
}

func (p *FileTxnLog) rollLog() error {
	if p.logBuf != nil {
		if err := p.logBuf.Flush(); err != nil {
			return err
		}
	}
	p.logBuf = nil
	return nil
}

func (p *FileTxnLog) close() error {
	if p.logBuf != nil {
		if err := p.logBuf.Flush(); err != nil {
			return err
		}
	}
	for _, buf := range p.streamsToFlush {
		if err := buf.Flush(); err != nil {
			return err
		}
	}
	return nil
}