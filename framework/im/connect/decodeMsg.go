package connect

import (
	"net"
	"time"
	"errors"
	"gopush/framework/helper"
	"gopush/const"
	"encoding/binary"
)

type Codec struct {
	Conn     net.Conn
	ReadBuf  helper.Buffer // 读缓冲
	WriteBuf []byte // 写缓冲
}

func NewCodec(conn net.Conn) *Codec {
	return &Codec{
		Conn:     conn,
		ReadBuf:  helper.NewBuffer(conn, constdefine.IMBufLen),
		WriteBuf: make([]byte, constdefine.IMBufLen),
	}
}

func (c *Codec) Read() (int, error) {
	return c.ReadBuf.ReadFromReader()
}

// Eecode 编码数据
func (c *Codec) Eecode(pack Package, duration time.Duration) error {
	contentLen := len(pack.Content)
	if contentLen > constdefine.IMContentMaxLen {
		return errors.New(constdefine.GetMsg(constdefine.IM_ERROR_OUT_OF_SIZE))
	}
	binary.BigEndian.PutUint16(c.WriteBuf[0:constdefine.IMTypeLen], uint16(pack.Code))
	binary.BigEndian.PutUint16(c.WriteBuf[constdefine.IMLenLen:constdefine.IMHeadLen], uint16(len(pack.Content)))
	copy(c.WriteBuf[constdefine.IMHeadLen:], pack.Content[:contentLen])
	c.Conn.SetWriteDeadline(time.Now().Add(duration))
	_, err := c.Conn.Write(c.WriteBuf[:constdefine.IMHeadLen+contentLen])
	if err != nil {
		return err
	}
	return nil
}

// Decode 解码数据
func (c *Codec) Decode() (*Package, bool) {
	var err error
	// 读取数据类型
	typeBuf, err := c.ReadBuf.Seek(0, constdefine.IMTypeLen)
	if err != nil {
		return nil, false
	}
	// 读取数据长度
	lenBuf, err := c.ReadBuf.Seek(constdefine.IMTypeLen, constdefine.IMHeadLen)
	if err != nil {
		return nil, false
	}
	// 读取数据内容
	valueType := int(binary.BigEndian.Uint16(typeBuf))
	valueLen := int(binary.BigEndian.Uint16(lenBuf))

	valueBuf, err := c.ReadBuf.Read(constdefine.IMHeadLen, valueLen)
	if err != nil {
		return nil, false
	}
	message := Package{Code: valueType, Content: valueBuf}
	return &message, true
}

