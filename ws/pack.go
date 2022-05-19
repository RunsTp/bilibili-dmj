package ws

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"errors"
	"fmt"
	"io/ioutil"
)

const (
	OperationHeartbeat          = 2
	OperationHeartbeatReply     = 3
	OperationMessage            = 5
	OperationUserAuthentication = 7
	OperationConnectSuccess     = 8

	PackageHeaderTotalLength = 16

	PackageOffset   = 0
	HeaderOffset    = 4
	VersionOffset   = 6
	OperationOffset = 8
	SequenceOffset  = 12

	BodyProtocolVersionNormal  = 0
	BodyProtocolDeflateVersion = 2
	HeaderDefaultVersion       = 1
	HeaderDefaultOperation     = 1
	HeaderDefaultSequence      = 1
	AuthOk                     = 0
	AuthTokenError             = -101
)

func PackHeader(packLength, operation uint32) Header {
	return Header{
		PackLength: PackageHeaderTotalLength + packLength,
		HeadLength: PackageHeaderTotalLength,
		Version:    HeaderDefaultVersion,
		Operation:  operation,
		Sequence:   HeaderDefaultSequence,
	}
}

func UnpackHeader(head []byte) (Header, error) {
	if len(head) != PackageHeaderTotalLength {
		return Header{}, errors.New(fmt.Sprintf("parse header fail, head length is [%d], expected length is [%d]", len(head), PackageHeaderTotalLength))
	}

	return Header{
		PackLength: binary.BigEndian.Uint32(head[PackageOffset : PackageOffset+4]),
		HeadLength: binary.BigEndian.Uint16(head[HeaderOffset : HeaderOffset+2]),
		Version:    binary.BigEndian.Uint16(head[VersionOffset : VersionOffset+2]),
		Operation:  binary.BigEndian.Uint32(head[OperationOffset : OperationOffset+4]),
		Sequence:   binary.BigEndian.Uint32(head[SequenceOffset : SequenceOffset+4]),
	}, nil
}

// UnpackMessage 解析直播间消息
// 被压缩的消息是一组 Message, 所以解析完成会返回 []Message 而非 Message
func UnpackMessage(raw []byte) ([]Message, error) {
	messages := make([]Message, 0, 4)

	if len(raw) <= PackageHeaderTotalLength {
		return messages, errors.New(fmt.Sprintf("packet defect, raw length [%d]", len(raw)))
	}

	head, err := UnpackHeader(raw[:PackageHeaderTotalLength])
	if err != nil {
		return messages, err
	} else if int(head.PackLength) > len(raw) {
		return messages, errors.New(fmt.Sprintf("packet defect, raw length [%d], expected length is [%d]", len(raw), head.PackLength))
	}

	// unzlib
	if head.Version == BodyProtocolDeflateVersion {
		reader, err := zlib.NewReader(bytes.NewReader(raw[PackageHeaderTotalLength:]))
		if err != nil {
			return messages, fmt.Errorf("unzlib fail, %w", err)
		}
		defer reader.Close()

		if raw, err = ioutil.ReadAll(reader); err != nil {
			return messages, err
		}
	}

	for len(raw) > 0 {
		head, err = UnpackHeader(raw[:PackageHeaderTotalLength])
		if err != nil {
			return messages, err
		} else if int(head.PackLength) > len(raw) {
			return messages, errors.New(fmt.Sprintf("packet defect, raw length [%d], expected length is [%d]", len(raw), head.PackLength))
		}

		// dereference
		payload := make([]byte, 0)
		payload = append(payload, raw[head.HeadLength:head.PackLength]...)

		messages = append(messages, Message{
			header:  head,
			payload: payload,
		})

		raw = raw[head.PackLength:]
	}

	return messages, nil
}

func PackMessage(operation uint32, raw []byte) Message {
	return Message{
		header:  PackHeader(uint32(len(raw)), operation),
		payload: raw,
	}
}

func bigEndianUint32(num uint32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, num)
	return b
}

func bigEndianUint16(num uint16) []byte {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, num)
	return b
}
