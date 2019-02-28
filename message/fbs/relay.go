// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package fbs

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type Relay struct {
	_tab flatbuffers.Table
}

func GetRootAsRelay(buf []byte, offset flatbuffers.UOffsetT) *Relay {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &Relay{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *Relay) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Relay) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *Relay) Blockchain() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Relay) Netid() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Relay) Version() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Relay) Data() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Relay) DevID() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func RelayStart(builder *flatbuffers.Builder) {
	builder.StartObject(5)
}
func RelayAddBlockchain(builder *flatbuffers.Builder, blockchain flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(blockchain), 0)
}
func RelayAddNetid(builder *flatbuffers.Builder, netid flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(netid), 0)
}
func RelayAddVersion(builder *flatbuffers.Builder, version flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(version), 0)
}
func RelayAddData(builder *flatbuffers.Builder, data flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(3, flatbuffers.UOffsetT(data), 0)
}
func RelayAddDevID(builder *flatbuffers.Builder, devID flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(4, flatbuffers.UOffsetT(devID), 0)
}
func RelayEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
