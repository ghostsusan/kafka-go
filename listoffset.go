package kafka

import "bufio"

type listOffsetRequestV1 struct {
	ReplicaID int32
	Topics    []listOffsetRequestTopicV1
}

func (r listOffsetRequestV1) size() int32 {
	return 4 + sizeofArray(len(r.Topics), func(i int) int32 { return r.Topics[i].size() })
}

func (r listOffsetRequestV1) writeTo(wb *writeBuffer) {
	wb.writeInt32(r.ReplicaID)
	wb.writeArray(len(r.Topics), func(i int) { r.Topics[i].writeTo(wb) })
}

type listOffsetRequestTopicV1 struct {
	TopicName  string
	Partitions []listOffsetRequestPartitionV1
}

func (t listOffsetRequestTopicV1) size() int32 {
	return sizeofString(t.TopicName) +
		sizeofArray(len(t.Partitions), func(i int) int32 { return t.Partitions[i].size() })
}

func (t listOffsetRequestTopicV1) writeTo(wb *writeBuffer) {
	wb.writeString(t.TopicName)
	wb.writeArray(len(t.Partitions), func(i int) { t.Partitions[i].writeTo(wb) })
}

type listOffsetRequestPartitionV1 struct {
	Partition int32
	Time      int64
}

func (p listOffsetRequestPartitionV1) size() int32 {
	return 4 + 8
}

func (p listOffsetRequestPartitionV1) writeTo(wb *writeBuffer) {
	wb.writeInt32(p.Partition)
	wb.writeInt64(p.Time)
}

type listOffsetResponseV1 []listOffsetResponseTopicV1

func (r listOffsetResponseV1) size() int32 {
	return sizeofArray(len(r), func(i int) int32 { return r[i].size() })
}

func (r listOffsetResponseV1) writeTo(wb *writeBuffer) {
	wb.writeArray(len(r), func(i int) { r[i].writeTo(wb) })
}

type listOffsetResponseTopicV1 struct {
	TopicName        string
	PartitionOffsets []partitionOffsetV1
}

func (t listOffsetResponseTopicV1) size() int32 {
	return sizeofString(t.TopicName) +
		sizeofArray(len(t.PartitionOffsets), func(i int) int32 { return t.PartitionOffsets[i].size() })
}

func (t listOffsetResponseTopicV1) writeTo(wb *writeBuffer) {
	wb.writeString(t.TopicName)
	wb.writeArray(len(t.PartitionOffsets), func(i int) { t.PartitionOffsets[i].writeTo(wb) })
}

type partitionOffsetV1 struct {
	Partition int32
	ErrorCode int16
	Timestamp int64
	Offset    int64
}

func (p partitionOffsetV1) size() int32 {
	return 4 + 2 + 8 + 8
}

func (p partitionOffsetV1) writeTo(wb *writeBuffer) {
	wb.writeInt32(p.Partition)
	wb.writeInt16(p.ErrorCode)
	wb.writeInt64(p.Timestamp)
	wb.writeInt64(p.Offset)
}

func (p *partitionOffsetV1) readFrom(r *bufio.Reader, sz int) (remain int, err error) {
	if remain, err = readInt32(r, sz, &p.Partition); err != nil {
		return
	}
	if remain, err = readInt16(r, remain, &p.ErrorCode); err != nil {
		return
	}
	if remain, err = readInt64(r, remain, &p.Timestamp); err != nil {
		return
	}
	if remain, err = readInt64(r, remain, &p.Offset); err != nil {
		return
	}
	return
}
