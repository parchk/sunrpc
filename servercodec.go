// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package sunrpc

import (
	"bytes"
	"io/ioutil"
	//	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net/rpc"

	"github.com/rasky/go-xdr/xdr2"

	//"github.com/davecgh/go-xdr/xdr2"
)

func init() {

}

type auth_unix struct {
	Stamp       uint32
	Machinename string
	Uid         uint32
	Gid         uint32
	Gids        []uint
}

type CallArags struct {
	Auth_unix auth_unix
	Arags     []byte
}

type serverCodec struct {
	conn         io.ReadWriteCloser
	closed       bool
	notifyClose  chan<- io.ReadWriteCloser
	recordReader io.Reader
	auth_tmp     auth_unix
}

// NewServerCodec returns a new rpc.ServerCodec using Sun RPC on conn.
// If a non-nil channel is passed as second argument, the conn is sent on
// that channel when Close() is called on conn.
func NewServerCodec(conn io.ReadWriteCloser, notifyClose chan<- io.ReadWriteCloser) rpc.ServerCodec {
	return &serverCodec{conn: conn, notifyClose: notifyClose}
}

func (c *serverCodec) ReadRequestHeader(req *rpc.Request) error {
	// NOTE:
	// Errors returned by this function aren't relayed back to the client
	// as WriteResponse() isn't called. The net/rpc package will call
	// c.Close() when this function returns an error.

	// Read entire RPC message from network

	record, err := ReadFullRecord(c.conn)
	if err != nil {
		if err != io.EOF {
			log.Println(err)
		}
		return err
	}

	//fmt.Println("svr read byte :", record)

	//fmt.Printf("%s\n", hexdump(record))

	c.recordReader = bytes.NewReader(record)

	// Unmarshall RPC message
	var call RPCMsgCall
	_, err = xdr.Unmarshal(c.recordReader, &call)
	if err != nil {
		log.Println(err)
		return err
	}

	if call.Header.Type != Call {
		log.Println(ErrInvalidRPCMessageType)
		return ErrInvalidRPCMessageType
	}

	if call.Body.Cred.Flavor == AuthSys {

		auth_buff := bytes.NewBuffer(call.Body.Cred.Body)

		var auth auth_unix

		_, err := xdr.Unmarshal(auth_buff, &auth)

		if err != nil {
			log.Println(err)
			return err
		}

		c.auth_tmp = auth
	}

	// Set req.Seq and req.ServiceMethod
	req.Seq = uint64(call.Header.Xid)
	procedureID := ProcedureID{call.Body.Program, call.Body.Version, call.Body.Procedure}
	procedureName, ok := GetProcedureName(procedureID)
	if ok {
		req.ServiceMethod = procedureName
	} else {
		// Due to our simpler map implementation, we cannot distinguish
		// between ErrProgUnavail and ErrProcUnavail
		log.Printf("%s: %+v\n", ErrProcUnavail, procedureID)
		return ErrProcUnavail
	}

	fmt.Println("ReadRequestHeader end")

	return nil
}

func (c *serverCodec) ReadRequestBody(funcArgs interface{}) error {

	if funcArgs == nil {
		return nil
	}

	callargs := funcArgs.(*CallArags)
	var err error
	callargs.Arags, err = ioutil.ReadAll(c.recordReader)

	if err != nil {
		if err != io.EOF {
			log.Println("ReadRequestBody readall error :", err)
			return err
		}
	}

	callargs.Auth_unix = c.auth_tmp
	/*
		if _, err := xdr.Unmarshal(c.recordReader, &funcArgs); err != nil {
			c.Close()
			return err
		}
	*/
	fmt.Println("ReadRequestBody end")

	return nil
}

func (c *serverCodec) WriteResponse(resp *rpc.Response, result interface{}) error {

	if resp.Error != "" {
		// The remote function returned error (shouldn't really happen)
		log.Println(resp.Error)
	}

	var buf bytes.Buffer

	replyMessage := RPCMsgReply{
		Header: RPCMessageHeader{
			Xid:  uint32(resp.Seq),
			Type: Reply,
		},
		Stat: MsgAccepted,
		Areply: AcceptedReply{
			Stat: Success,
		},
	}

	if _, err := xdr.Marshal(&buf, replyMessage); err != nil {
		c.Close()
		return err
	}

	// Marshal and fill procedure-specific reply into the buffer

	if _, err := xdr.Marshal(&buf, result); err != nil {
		c.Close()
		return err
	}

	//fmt.Println("svr write byte :", buf.Bytes())

	//fmt.Printf("%s\n", hexdump(buf.Bytes()))

	// Write buffer contents to network
	if _, err := WriteFullRecord(c.conn, buf.Bytes()); err != nil {
		c.Close()
		return err
	}

	return nil
}

func (c *serverCodec) Close() error {
	if c.closed {
		return nil
	}

	err := c.conn.Close()
	if err == nil {
		c.closed = true
		if c.notifyClose != nil {
			c.notifyClose <- c.conn
		}
	}

	return err
}
